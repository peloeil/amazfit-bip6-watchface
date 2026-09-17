package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// The helper transfers opaque ZIP bytes; it does not depend on a watchface release.
func testPackage(t *testing.T) string {
	t.Helper()
	var b bytes.Buffer
	z := zip.NewWriter(&b)
	w, err := z.Create("app.json")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = w.Write([]byte(`{"app":{"appName":"test watchface"}}`)); err != nil {
		t.Fatal(err)
	}
	if err = z.Close(); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "watchface.zpk")
	if err = os.WriteFile(path, b.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRawPackageQRProtocol(t *testing.T) {
	a := fixture(t)
	for _, host := range []string{"pc.example.ts.net", "pc.example.ts.net:8443"} {
		a.urls(host)
		s := a.snapshot()
		if s.InstallURI != "zpkd1://"+host+"/watchface.zpk" || s.DownloadURL != "https://"+host+"/watchface.zpk" {
			t.Fatal("raw package must use the direct ZPK download scheme", s)
		}
	}
}

func fixture(t *testing.T) *application {
	t.Helper()
	a := newApplication(func() {})
	a.localHost = "127.0.0.1:12345"
	a.urls("pc.example.ts.net")
	t.Cleanup(a.stop)
	return a
}
func call(a *application, method, path, host, identity, origin string, body []byte) *httptest.ResponseRecorder {
	r := httptest.NewRequest(method, "http://"+host+path, bytes.NewReader(body))
	r.Host = host
	if identity != "" {
		r.Header.Set("Tailscale-User-Login", identity)
	}
	if origin != "" {
		r.Header.Set("Origin", origin)
		r.Header.Set("X-Lumen-Action", "1")
	}
	w := httptest.NewRecorder()
	a.handler(w, r)
	return w
}
func TestAccessDelegatedToTailscale(t *testing.T) {
	a := fixture(t)
	if e := a.loadPackage(testPackage(t)); e != nil {
		t.Fatal(e)
	}
	// Serve already enforces network access. The application does not filter identities.
	for _, identity := range []string{"", "owner@example.com", "another-member@example.com"} {
		w := call(a, "GET", "/"+packageName, a.remoteHost, identity, "", nil)
		if w.Code != 200 {
			t.Fatalf("identity %q rejected: %d", identity, w.Code)
		}
	}
	if w := call(a, "POST", "/pc/select", a.localHost, "", "https://other-site.example", nil); w.Code != 403 {
		t.Fatal("cross-site browser mutation accepted")
	}
	if w := call(a, "GET", "/pc/", a.localHost, "", "", nil); w.Code != 200 {
		t.Fatal("PC page")
	}
	if w := call(a, "GET", "/phone/../main.go", a.remoteHost, "", "", nil); w.Code != 404 {
		t.Fatal("arbitrary file exposed")
	}
}
func TestUploadedZPKIsUnchanged(t *testing.T) {
	a := fixture(t)
	original, e := os.ReadFile(testPackage(t))
	if e != nil {
		t.Fatal(e)
	}
	if w := call(a, "GET", "/"+packageName, a.remoteHost, "", "", nil); w.Code != 409 {
		t.Fatal("file before selection")
	}
	if w := call(a, "POST", "/pc/select", a.localHost, "", "http://"+a.localHost, []byte("invalid")); w.Code != 422 {
		t.Fatal("invalid file accepted")
	}
	if w := call(a, "POST", "/pc/select", a.localHost, "", "http://"+a.localHost, original); w.Code != 200 {
		t.Fatal(w.Code, w.Body.String())
	}
	w := call(a, "GET", "/"+packageName, a.remoteHost, "", "", nil)
	if w.Code != 200 || !bytes.Equal(w.Body.Bytes(), original) {
		t.Fatal("package modified")
	}
	if w.Header().Get("Content-Type") != "application/octet-stream" {
		t.Fatal("MIME")
	}
	w = call(a, "HEAD", "/"+packageName, a.remoteHost, "", "", nil)
	if w.Code != 200 || w.Body.Len() != 0 || w.Header().Get("Content-Length") != strconv.Itoa(len(original)) {
		t.Fatal("HEAD")
	}
	r := httptest.NewRequest("GET", "https://"+a.remoteHost+"/"+packageName, nil)
	r.Header.Set("Range", "bytes=10-19")
	w = httptest.NewRecorder()
	a.handler(w, r)
	if w.Code != 206 || !bytes.Equal(w.Body.Bytes(), original[10:20]) {
		t.Fatal("Range")
	}
}
func TestTailscaleConfiguration(t *testing.T) {
	m := map[string]any{"BackendState": "Running", "Self": map[string]any{"DNSName": "pc.example.ts.net."}}
	host, e := nodeDNS(m)
	if e != nil || host != "pc.example.ts.net" {
		t.Fatal("DNS must work without user identity", host, e)
	}
	cfg := map[string]any{"TCP": map[string]any{"443": nil}, "Foreground": map[string]any{"session": map[string]any{"TCP": map[string]any{"8443": nil}}}}
	if p, e := selectPort(cfg); e != nil || p != 10000 {
		t.Fatal("existing ports overwritten")
	}
	target := "http://127.0.0.1:3210"
	own := map[string]any{"Web": map[string]any{host + ":443": map[string]any{"Handlers": map[string]any{"/": map[string]any{"Proxy": target}}}}}
	if !configOwns(own, host, 443, target) {
		t.Fatal("serve config")
	}
	own["AllowFunnel"] = map[string]any{host + ":443": true}
	if configOwns(own, host, 443, target) {
		t.Fatal("Funnel accepted as private Serve")
	}
}

func TestPhoneUploadInterfaceRemoved(t *testing.T) {
	a := fixture(t)
	for _, path := range []string{"/phone/", "/phone/select", "/phone/state"} {
		w := call(a, "GET", path, a.remoteHost, "", "", nil)
		if w.Code != 404 {
			t.Fatal(path, w.Code)
		}
	}
}

func TestHTTPSDefaultPortAndHostCase(t *testing.T) {
	a := fixture(t)
	if e := a.loadPackage(testPackage(t)); e != nil {
		t.Fatal(e)
	}
	for _, host := range []string{"pc.example.ts.net", "pc.example.ts.net:443", "PC.EXAMPLE.TS.NET:443"} {
		w := call(a, "GET", "/"+packageName, host, "", "", nil)
		if w.Code != 200 || w.Body.Len() != len(a.payload) {
			t.Errorf("equivalent HTTPS host %q: HTTP %d", host, w.Code)
		}
	}
	for _, host := range []string{"pc.example.ts.net:8443", "other.example.ts.net", a.localHost} {
		if w := call(a, "GET", "/"+packageName, host, "", "", nil); w.Code != 403 {
			t.Errorf("unexpected host accepted: %q", host)
		}
	}
	a.urls("pc.example.ts.net:8443")
	if w := call(a, "GET", "/"+packageName, "PC.EXAMPLE.TS.NET:8443", "", "", nil); w.Code != 200 {
		t.Fatal("explicit alternate port rejected")
	}
	if w := call(a, "GET", "/"+packageName, "pc.example.ts.net:443", "", "", nil); w.Code != 403 {
		t.Fatal("wrong port accepted")
	}
}

func TestRequestDiagnostics(t *testing.T) {
	a := fixture(t)
	if a.snapshot().LastRequest != "" {
		t.Fatal("request exists before access")
	}
	if e := a.loadPackage(testPackage(t)); e != nil {
		t.Fatal(e)
	}
	for _, tc := range []struct{ method, host, want string }{
		{"GET", "pc.example.ts.net:443", fmt.Sprintf("GET → HTTP 200 · %d バイト応答", len(a.payload))},
		{"HEAD", "pc.example.ts.net", "HEAD → HTTP 200 · 0 バイト応答"},
		{"GET", "wrong.example.ts.net", "GET → HTTP 403"},
	} {
		call(a, tc.method, "/"+packageName, tc.host, "", "", nil)
		if got := a.snapshot().LastRequest; !strings.Contains(got, tc.want) {
			t.Fatal(got)
		}
	}
	if w := call(a, "GET", "/pc/state", a.remoteHost, "", "", nil); w.Code != 403 {
		t.Fatal("PC diagnostics exposed remotely")
	}
}
