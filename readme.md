# LUMEN ORBIT

Amazfit Bip 6用文字盤 **1.9.2** と、作成・最適化・導入に使うツール。
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

## 改善方針

目標は、現在の表示・機能を維持し、実機で正しく軽く動くこと。
最適化の採否は、実行時の処理とメモリの無駄を減らせるかで決める。

- 容量だけ減らす変更は採用しない。ZIP圧縮率の追求は打ち切る。
- 不要な更新・重複処理・描画領域を削る。標準ウィジェットで済む処理を自作しない。
- 見た目や機能を削らない。容量のために復号・描画処理を増やさない。
- 何の処理・メモリが減るか説明できない変更は行わない。無駄がなければ完了とする。
- 改善の前提として利用者に電池持ちの比較測定を要求しない。未測定の効果を成果として報告しない。

文字盤1.9.1は利用者が実機で数日間使用し、正常動作を確認済み。
これを基準に、変更で影響する表示・機能を確認する。

1.9.2ではデジタル表示を時・分のみに変更。アナログ秒針は維持する。

## 時計へ導入

1. [Windows導入ツール](publish/lumen-orbit-link-windows.zip)を展開する。
2. [文字盤1.9.2](publish/lumen-orbit-bip6.zpk)を `lumen-orbit-link.exe` にドラッグ＆ドロップする。
3. PCに表示されたQRをAndroidのZeppで読み取る。

PC・AndroidのTailscale接続とZeppの開発者モードが必要。
準備・操作・接続確認は [導入ツールのREADME](link-windows/readme.txt) を参照。
Windows導入ツール経由でBip 6実機へ導入できることを利用者が確認済み。

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
uv run optimize-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6.zpk --source ../watchface --fixed
```

詳細は [文字盤のREADME](watchface/readme.md) と
[最適化のREADME](optimization/readme.md) を参照。
Windows導入ツールのテスト・再ビルド手順は [こちら](link-windows/readme.txt)。

## 配布物

配布ファイル名はバージョン番号を含まない固定名とし、履歴はGitで管理する。
過去の配布物は、その時点のコミットから取得する。
端末の更新判定に使うアプリ内部のバージョン番号は維持する。

文字盤の次の3ファイルは、同じバージョンのセットで返却する。
ソースZIPと最適化ZIPは分け、各READMEの返却・引き継ぎルールも維持する。

1. [lumen-orbit-bip6.zpk](publish/lumen-orbit-bip6.zpk)
2. [lumen-orbit-bip6-rebuild-source.zip](publish/lumen-orbit-bip6-rebuild-source.zip)
3. [lumen-orbit-bip6-optimization.zip](publish/lumen-orbit-bip6-optimization.zip)

導入ツールは [lumen-orbit-link-windows.zip](publish/lumen-orbit-link-windows.zip)。
文字盤のZPKを同梱せず、利用者が選択したファイルをそのまま配信する。
