# LUMEN ORBIT

Amazfit Bip 6用文字盤 **1.9.1** と、作成・最適化・導入に使うツール。
Windows導入ツール **LUMEN ORBIT LINK 1.0.5** のバージョンは文字盤と独立して管理する。

## 構成

ファイル名・ディレクトリ名は小文字のハイフン区切りに統一する。
Goがテストを認識するために必要な `main_test.go` など、ツールが要求する名前は維持する。

```text
watchface/       文字盤のソース・画像・ビルド設定
optimization/    ZPK最適化スクリプト
link-windows/    Windows導入ツールのGoソース・画面・ライセンス
publish/        配布するZPKとZIP
```

実行ファイルはWindows用ZIPに収録する。配布ZIPの先頭ディレクトリも
`watchface/`、`optimization/`、`link-windows/` とし、作業ツリーと揃える。
依存パッケージの `node_modules/` やビルド途中の `dist/` はZIPに含めない。

## 時計へ導入

1. [Windows導入ツール](publish/lumen-orbit-link-windows.zip)を展開する。
2. [文字盤1.9.1](publish/lumen-orbit-bip6-1.9.1.zpk)を `lumen-orbit-link.exe` にドラッグ＆ドロップする。
3. PCに表示されたQRをAndroidのZeppで読み取る。

PC・AndroidのTailscale接続とZeppの開発者モードが必要。
準備・操作・接続確認は [導入ツールのREADME](link-windows/readme.txt) を参照。
Windows・Tailscale実網・Zepp/Bip 6実機での動作は未確認。

## 再ビルド・最適化

文字盤のビルド（Node.js・npmが必要）:

```sh
cd watchface
npm ci
node build.js build
```

続いて最適化（uvが必要）:

```sh
cd ../optimization
uv run optimize-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6-1.9.1.zpk --source ../watchface --fixed
```

詳細は [文字盤のREADME](watchface/readme.md) と
[最適化のREADME](optimization/readme.md) を参照。
Windows導入ツールのテスト・再ビルド手順は [こちら](link-windows/readme.txt)。

## 配布物

文字盤の次の3ファイルは、同じバージョンのセットで返却する。
ソースZIPと最適化ZIPは分け、各READMEの返却・引き継ぎルールも維持する。

1. [lumen-orbit-bip6-1.9.1.zpk](publish/lumen-orbit-bip6-1.9.1.zpk)
2. [lumen-orbit-bip6-1.9.1-rebuild-source.zip](publish/lumen-orbit-bip6-1.9.1-rebuild-source.zip)
3. [lumen-orbit-bip6-1.9.1-optimization.zip](publish/lumen-orbit-bip6-1.9.1-optimization.zip)

導入ツールは [lumen-orbit-link-windows.zip](publish/lumen-orbit-link-windows.zip)。
文字盤のZPKを同梱せず、利用者が選択したファイルをそのまま配信する。
