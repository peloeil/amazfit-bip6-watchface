ZPK LINK 1.0.6 — Windows x64

watchfaceのZPKをPCから配信し、Zeppで読み取るQRを表示する。
時計に対応するZPKを用意する。

初回準備:
【PCとスマートフォン】Tailscaleをインストールして接続する。同じアカウントを使えばよい。
https://tailscale.com/download
【スマートフォン・Zepp】開発者モードを有効にする。
プロフィール → 設定 → バージョン情報 → Zeppロゴを7回タップ。
https://docs.zepp.com/docs/guides/tools/zepp-app/

最短の操作:
1.【Windows PC】ZIPを展開し、手元のZPKを zpk-link.exe のアイコンへ
   ドラッグ＆ドロップする。PCのブラウザに導入用QRが表示される。
   初回にHTTPS設定を求められたら、PC画面のリンクから有効にする。
2.【スマートフォン・Zepp】接続した時計 → 一般 → 開発者モード → ＋ → スキャンを開き、
   PC画面のQRを読み取る。
3.【Windows PC】Zeppと時計で転送完了を確認したら、PC画面の「終了」を押す。

exeを普通に開いた場合も、PCのブラウザにある「ZPKを選択」でファイルを選べる。

ZPKはTailscale内で配信する。アクセスできる端末はTailscaleの設定に従う。

接続できない場合:
「Fail to get watch face details」が出た場合:
【PC】QRの下の「ファイル要求」と、【スマートフォン・Zepp】のバージョンを確認する。
・未到達: この配信処理にZPKの要求が届いていない。接続またはQR解釈を調べる。
・HTTP 200: PCがZPKの応答を返した。Zeppの解析や時計への導入成功は意味しない。
・HTTP 403など: PC側で要求を拒否した。表示をそのまま報告する。
【PC】「エラーが出た場合」を開くと、スマートフォンのブラウザで接続を確認するURLがある。
【PC・スマートフォン両方】Tailscaleの接続とアクセス設定を確認。
【PC】転送が終わるまでスリープさせない。
【Android】Tailscaleの「App-based split tunneling」でZeppを除外している場合は、Zeppのチェックを外す。

開発・配布のルールは、親フォルダの AGENTS.md を参照。

再ビルド（Go 1.24以降）:
作業ツリー・配布ZIPとも、zpk-link ディレクトリで実行。
  go test ./...
画面更新のテスト（Node.jsが必要）:
  node --test test-app.js
Windows PowerShell:
  go build -trimpath -ldflags="-s -w" -o zpk-link.exe .
Linux/macOSからWindows x64向け:
  GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o zpk-link.exe .
