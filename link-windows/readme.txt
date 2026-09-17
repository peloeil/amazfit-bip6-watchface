LUMEN ORBIT LINK 1.0.5 — Windows x64

初回準備:
【PCとAndroid】Tailscaleをインストールして接続する。同じアカウントを使えばよい。
https://tailscale.com/download
【Android・Zepp】開発者モードを有効にする。
プロフィール → 設定 → バージョン情報 → Zeppロゴを7回タップ。
https://docs.zepp.com/docs/guides/tools/zepp-app/

最短の操作:
1.【Windows PC】ZIPを展開し、手元のZPKを lumen-orbit-link.exe のアイコンへ
   ドラッグ＆ドロップする。PCのブラウザに導入用QRが表示される。
   初回にHTTPS設定を求められたら、PC画面のリンクから有効にする。
2.【Android・Zepp】Bip 6 → 一般 → 開発者モード → ＋ → スキャンを開き、
   PC画面のQRを読み取る。
3.【Windows PC】Zeppと時計で転送完了を確認したら、PC画面の「終了」を押す。

exeを普通に開いた場合も、PCのブラウザにある「ZPKを選択」でファイルを選べる。
Androidのブラウザ・ファイル選択・中継用QRは使わない。

PC上のZPKをそのままZeppへ渡す。ファイルは書き換えない。
アクセス制御はTailscaleの設定に従う。独自のユーザー認証・有効期限はない。
PCのサーバーは127.0.0.1だけで待ち受け、Tailscale ServeでZPKを配信する。
一般公開用のFunnelは起動しない。
https://tailscale.com/docs/features/tailscale-serve

接続できない場合:
「Fail to get watch face details」が出た場合:
【PC】QRの下の「ファイル要求」と、【Android・Zepp】のバージョンを確認する。
・未到達: この配信処理にZPKの要求が届いていない。接続またはQR解釈を調べる。
・HTTP 200: PCがZPKの応答を返した。Zeppの解析や時計への導入成功は意味しない。
・HTTP 403など: PC側で要求を拒否した。表示をそのまま報告する。
【PC】「エラーが出た場合」を開くと、AndroidのChromeで接続を確認するURLがある。
【PC・Android両方】Tailscaleの接続とアクセス設定を確認。
【PC】転送が終わるまでスリープさせない。
【Android】ZeppをTailscaleのVPN除外対象にしない。

検証:
PCファイルの読込と完全一致、QRの別実装によるデコード、
模擬Tailscaleでの起動・終了、PC画面の操作を確認。
Windowsでの実行・Tailscale実網・Zepp/Bip 6実機の導入は未検証。
画面操作はDOMシミュレーションで検証し、実描画は未確認。

1.0.3では同じHTTPSホストの「:443」付き表記・大文字表記を拒否する不具合を修正。
これはコード上で再現・修正した不具合であり、報告されたZeppエラーの原因は未確定。

1.0.4では、ZPK本体を直接返す配信内容に合わせてQRをzpkd1://形式へ変更。
Zepp 10.8.1でwatchface://を読んだ際の「Fail to get watch face details」に対する修正。
公式Zeus CLIの時計盤フローは配布URLとpreviewの情報を読み取る。
直接ZPKを返す今回の処理にwatchface://を使っていたため、形式が一致していなかった。
2026-09-13に公開配布サイトのQRでもzpkd1://を確認した。
ZPKやpackageInfoは変更していない。10.8.1実機での改善確認は利用者の再試行待ち。

Goソースとwebディレクトリは確認・再ビルド用。通常はexeだけを使う。
文字盤は別配布の lumen-orbit-bip6.zpk を使う。
文字盤1.9.1と導入ツール1.0.5は、それぞれ独立したバージョン。
Go/Node.js/Pythonのインストールは不要。

1.0.5では、ファイル選択画面の旧版ファイル名を汎用的な案内に変更。
テストはその場で作るZIPを使い、旧版の文字盤と固定ハッシュ・容量への依存を解消。

再ビルド（Go 1.24以降）:
作業ツリー・配布ZIPとも、link-windows ディレクトリで実行。
  go test ./...
Windows PowerShell:
  go build -trimpath -ldflags="-s -w" -o lumen-orbit-link.exe .
Linux/macOSからWindows x64向け:
  GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o lumen-orbit-link.exe .
