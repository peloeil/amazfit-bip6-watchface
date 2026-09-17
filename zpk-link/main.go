// ZPK LINK: Windows helper. Network access is controlled by Tailscale.
package main

import (
	"archive/zip"
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const packageName = "watchface.zpk"
const maxPackage = 8 * 1024 * 1024

//go:embed web/*
var assets embed.FS

type status struct {
	Phase       string `json:"phase"`
	Message     string `json:"message"`
	SetupURL    string `json:"setupURL,omitempty"`
	InstallURI  string `json:"installURI,omitempty"`
	DownloadURL string `json:"downloadURL,omitempty"`
	LastRequest string `json:"lastRequest,omitempty"`
	Selected    bool   `json:"selected"`
}
type application struct {
	mu                    sync.Mutex
	ctx                   context.Context
	state                 status
	remoteHost, localHost string
	payload               []byte
	stop                  context.CancelFunc
	preparing             bool
}

func newApplication(stop context.CancelFunc) *application {
	ctx, cancel := context.WithCancel(context.Background())
	return &application{ctx: ctx, stop: func() { cancel(); stop() }, state: status{Phase: "starting", Message: "【PC】Tailscaleを確認しています。"}}
}
func (a *application) snapshot() status { a.mu.Lock(); defer a.mu.Unlock(); return a.state }
func (a *application) message(phase, msg string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.Phase = phase
	a.state.Message = msg
}
func (a *application) urls(host string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.remoteHost = host
	// This endpoint serves raw ZPK bytes. The watchface:// flow expects dial
	// details (including url/preview) on Zepp 10.8.1, not a raw package.
	a.state.InstallURI = "zpkd1://" + host + "/" + packageName
	a.state.DownloadURL = "https://" + host + "/" + packageName
}

// HTTPS's default port and DNS letter case do not change the destination.
func sameHTTPSHost(got, want string) bool {
	return want != "" && strings.EqualFold(strings.TrimSuffix(got, ":443"), strings.TrimSuffix(want, ":443"))
}

type responseResult struct {
	http.ResponseWriter
	code, bytes int
}

func (w *responseResult) WriteHeader(code int) {
	if w.code != 0 {
		return
	}
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}
func (w *responseResult) Write(b []byte) (int, error) {
	if w.code == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}
func (w *responseResult) Unwrap() http.ResponseWriter { return w.ResponseWriter }

func (a *application) handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'none'")
	a.mu.Lock()
	remote := a.remoteHost
	a.mu.Unlock()
	// Host/Origin checks protect the browser interface, not Tailscale identity.
	// No Tailscale user headers, owner checks, credentials or expiring tokens.
	if r.URL.Path == "/"+packageName {
		result := &responseResult{ResponseWriter: w}
		w = result
		defer func() {
			code := result.code
			if code == 0 {
				code = http.StatusOK
			}
			a.mu.Lock()
			a.state.LastRequest = fmt.Sprintf("%s · %s → HTTP %d · %d バイト応答", time.Now().Format("15:04:05"), r.Method, code, result.bytes)
			a.mu.Unlock()
		}()
		if !sameHTTPSHost(r.Host, remote) {
			http.Error(w, "Invalid host", 403)
			return
		}
		if r.Method != "GET" && r.Method != "HEAD" {
			http.Error(w, "Method not allowed", 405)
			return
		}
		a.mu.Lock()
		b := a.payload
		a.mu.Unlock()
		if len(b) == 0 {
			http.Error(w, "【PC】ZPKを選んでください。", 409)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, packageName, time.Time{}, bytes.NewReader(b))
		return
	}
	// Only the file download is exposed through Serve. The UI stays on this PC.
	if !strings.HasPrefix(r.URL.Path, "/pc/") {
		http.NotFound(w, r)
		return
	}
	if r.Host != a.localHost {
		http.Error(w, "Invalid host", 403)
		return
	}
	route := strings.TrimPrefix(r.URL.Path, "/pc/")
	if strings.Contains(route, "/") || r.URL.RawPath != "" {
		http.NotFound(w, r)
		return
	}
	origin := "http://" + a.localHost
	if r.Method == "POST" {
		if r.Header.Get("Origin") != origin || r.Header.Get("X-ZPK-Action") != "1" {
			http.Error(w, "Invalid browser request", 403)
			return
		}
		switch route {
		case "select":
			b, e := io.ReadAll(http.MaxBytesReader(w, r.Body, maxPackage))
			if e != nil {
				http.Error(w, "ZPKは8MB以内で選んでください。", 413)
				return
			}
			if e = a.setPackage(b); e != nil {
				http.Error(w, e.Error(), 422)
				return
			}
		case "retry":
			go a.prepare()
		case "stop":
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"ok":true}`)
			go func() { time.Sleep(150 * time.Millisecond); a.stop() }()
			return
		default:
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"ok":true}`)
		return
	}
	if r.Method != "GET" && r.Method != "HEAD" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if route == "state" {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(a.snapshot())
		return
	}
	if route == "" {
		route = "index.html"
	}
	types := map[string]string{"index.html": "text/html; charset=utf-8", "style.css": "text/css; charset=utf-8", "app.js": "text/javascript; charset=utf-8", "qr.js": "text/javascript; charset=utf-8"}
	typ, ok := types[route]
	if !ok {
		http.NotFound(w, r)
		return
	}
	b, e := assets.ReadFile("web/" + route)
	if e != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", typ)
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	if r.Method != "HEAD" {
		w.Write(b)
	}
}

func (a *application) setPackage(b []byte) error {
	if len(b) > maxPackage {
		return errors.New("ZPKは8MB以内で選んでください。")
	}
	if _, e := zip.NewReader(bytes.NewReader(b), int64(len(b))); e != nil {
		return errors.New("ZPKファイルを選んでください。")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.payload = b
	a.state.Selected = true
	a.state.LastRequest = ""
	return nil
}
func (a *application) loadPackage(path string) error {
	f, e := os.Open(path)
	if e != nil {
		return e
	}
	defer f.Close()
	b, e := io.ReadAll(io.LimitReader(f, maxPackage+1))
	if e != nil {
		return e
	}
	return a.setPackage(b)
}

func tailscalePath() (string, error) {
	if p, e := exec.LookPath("tailscale"); e == nil {
		return p, nil
	}
	var candidates []string
	for _, env := range []string{"ProgramFiles", "ProgramW6432"} {
		if p := os.Getenv(env); p != "" {
			candidates = append(candidates, filepath.Join(p, "Tailscale", "tailscale.exe"))
		}
	}
	for _, p := range candidates {
		if s, e := os.Stat(p); e == nil && !s.IsDir() {
			return p, nil
		}
	}
	return "", errors.New("【PC・スマートフォン両方】Tailscaleを接続してください。その後、このPC画面で「再確認」を押してください。")
}
func cliJSON(path string, args ...string) (map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	b, err := exec.CommandContext(ctx, path, args...).Output()
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(b)) == 0 || string(bytes.TrimSpace(b)) == "null" {
		return map[string]any{}, nil
	}
	var m map[string]any
	err = json.Unmarshal(b, &m)
	return m, err
}
func nodeDNS(m map[string]any) (string, error) {
	if m["BackendState"] != "Running" {
		return "", errors.New("【PC】Tailscaleを起動して接続をONにし、「再確認」を押してください。")
	}
	self, _ := m["Self"].(map[string]any)
	dns, _ := self["DNSName"].(string)
	dns = strings.TrimSuffix(dns, ".")
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]*\.ts\.net$`).MatchString(dns) {
		return "", errors.New("【PC】TailscaleのMagicDNS設定を確認してください。")
	}
	return dns, nil
}

// Conservatively reserve ports from TCP, Web and AllowFunnel, including foreground sessions.
func usedPorts(config map[string]any) map[int]bool {
	used := map[int]bool{}
	var visit func(map[string]any)
	visit = func(m map[string]any) {
		for key, val := range m {
			if key == "TCP" || key == "Web" || key == "AllowFunnel" {
				if ports, ok := val.(map[string]any); ok {
					for k := range ports {
						p := k
						if i := strings.LastIndex(k, ":"); i >= 0 {
							p = k[i+1:]
						}
						if n, e := strconv.Atoi(p); e == nil {
							used[n] = true
						}
					}
				}
			}
			if sub, ok := val.(map[string]any); ok {
				visit(sub)
			}
		}
	}
	visit(config)
	return used
}
func selectPort(config map[string]any) (int, error) {
	used := usedPorts(config)
	for _, p := range []int{443, 8443, 10000, 18443, 28443} {
		if !used[p] {
			return p, nil
		}
	}
	return 0, errors.New("使用できるHTTPSポートがありません。既存の配信設定は変更していません。")
}
func configOwns(config map[string]any, host string, port int, target string) bool {
	key := host + ":" + strconv.Itoa(port)
	found := false
	public := false
	var visit func(map[string]any)
	visit = func(m map[string]any) {
		if funnel, ok := m["AllowFunnel"].(map[string]any); ok {
			if enabled, _ := funnel[key].(bool); enabled {
				public = true
			}
		}
		if web, ok := m["Web"].(map[string]any); ok {
			if hostConfig, ok := web[key].(map[string]any); ok {
				if handlers, ok := hostConfig["Handlers"].(map[string]any); ok {
					if root, ok := handlers["/"].(map[string]any); ok {
						if root["Proxy"] == target {
							found = true
						}
					}
				}
			}
		}
		for _, v := range m {
			if sub, ok := v.(map[string]any); ok {
				visit(sub)
			}
		}
	}
	visit(config)
	return found && !public
}

type serveLog struct {
	a    *application
	mu   sync.Mutex
	tail string
}

var setupPattern = regexp.MustCompile(`https://login\.tailscale\.com/[^\s<>"']+`)

func (l *serveLog) Write(b []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tail += string(b)
	if len(l.tail) > 4000 {
		l.tail = l.tail[len(l.tail)-4000:]
	}
	if u := setupPattern.FindString(l.tail); u != "" {
		l.a.mu.Lock()
		l.a.state.SetupURL = u
		l.a.state.Message = "初回のみ、TailscaleのHTTPS設定が必要です。「HTTPSを有効にする」を開いてください。"
		l.a.mu.Unlock()
	}
	return len(b), nil
}
func (l *serveLog) last() string { l.mu.Lock(); defer l.mu.Unlock(); return l.tail }
func (a *application) prepare() {
	a.mu.Lock()
	if a.preparing || a.state.Phase == "ready" {
		a.mu.Unlock()
		return
	}
	a.preparing = true
	a.state.SetupURL = ""
	a.mu.Unlock()
	defer func() { a.mu.Lock(); a.preparing = false; a.mu.Unlock() }()
	a.message("starting", "Tailscaleと非公開のHTTPS接続を確認しています。")
	ts, err := tailscalePath()
	if err != nil {
		a.message("setup", err.Error())
		return
	}
	m, err := cliJSON(ts, "status", "--json")
	if err != nil {
		a.message("setup", "Tailscaleに接続できません。PCのTailscaleを起動してから再確認してください。")
		return
	}
	dns, err := nodeDNS(m)
	if err != nil {
		a.message("setup", err.Error())
		return
	}
	config, err := cliJSON(ts, "serve", "status", "--json")
	if err != nil {
		a.message("error", "Tailscale Serveの既存設定を取得できません。Tailscaleを更新してください。設定は変更していません。")
		return
	}
	port, err := selectPort(config)
	if err != nil {
		a.message("error", err.Error())
		return
	}
	host := dns
	if port != 443 {
		host += ":" + strconv.Itoa(port)
	}
	a.urls(host)
	ctx, cancel := context.WithCancel(a.ctx)
	target := "http://" + a.localHost
	cmd := exec.CommandContext(ctx, ts, "serve", "--https="+strconv.Itoa(port), target)
	logger := &serveLog{a: a}
	cmd.Stdout = logger
	cmd.Stderr = logger
	if err = cmd.Start(); err != nil {
		cancel()
		a.message("error", "Tailscale Serveを起動できませんでした。")
		return
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	deadline := time.NewTimer(3 * time.Minute)
	defer deadline.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-done:
			msg := "非公開配信が終了しました。再確認で起動し直せます。"
			if err != nil {
				msg = "Tailscale Serveを開始できませんでした。PCのTailscaleとHTTPS設定を確認してください。"
				if strings.Contains(logger.last(), "Access denied") {
					msg += " TailscaleからServeを実行する権限が拒否されています。"
				}
			}
			cancel()
			a.message("error", msg)
			return
		case <-deadline.C:
			cancel()
			a.message("setup", "HTTPSの設定待ちが終了しました。設定を済ませてから再確認してください。")
			return
		case <-ticker.C:
			c, e := cliJSON(ts, "serve", "status", "--json")
			if e == nil && configOwns(c, dns, port, target) {
				a.message("ready", "【PC】接続できました。ZPKを選ぶとZepp用のQRを表示します。")
				go func() {
					select {
					case <-done:
						a.message("error", "非公開配信が終了しました。PCで再確認してください。")
					case <-ctx.Done():
					}
				}()
				return
			}
		}
	}
}
func main() {
	noBrowser := flag.Bool("no-browser", false, "do not open the Windows browser")
	flag.Parse()
	if flag.NArg() > 1 {
		log.Fatal("ZPKは1ファイルだけ指定してください。")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	a := newApplication(cancel)
	if flag.NArg() == 1 {
		if e := a.loadPackage(flag.Arg(0)); e != nil {
			log.Fatal(e)
		}
	}
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}
	a.localHost = listener.Addr().String()
	server := &http.Server{Handler: http.HandlerFunc(a.handler), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 15 * time.Second, MaxHeaderBytes: 16 * 1024}
	go func() {
		if e := server.Serve(listener); e != nil && !errors.Is(e, http.ErrServerClosed) {
			cancel()
		}
	}()
	localURL := "http://" + a.localHost + "/pc/"
	fmt.Println("ZPK LINK 1.0.6 / Windows\n【PC】ブラウザで開いてください：\n" + localURL + "\n終了：PC画面の終了ボタン、または Ctrl+C。")
	if !*noBrowser {
		cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", localURL)
		if e := cmd.Start(); e == nil {
			go cmd.Wait()
		}
	}
	go a.prepare()
	<-ctx.Done()
	a.stop()
	shutdown, c := context.WithTimeout(context.Background(), 2*time.Second)
	defer c()
	server.Shutdown(shutdown)
	fmt.Println("終了しました。")
}
