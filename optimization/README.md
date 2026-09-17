# LUMEN ORBIT 1.9.2 — ZPK最適化

標準ビルドの `.zab` / `.zpk` を、対応するソースPNGから最適化する。
[uv](https://docs.astral.sh/uv/getting-started/installation/) が必要。
Python 3.10以降と依存パッケージ（Pillow・zopfli）の固定バージョンは、
`optimize-zpk.py` 冒頭の [インラインメタデータ](https://docs.astral.sh/uv/guides/scripts/#declaring-script-dependencies)
で管理している。uvが必要なPythonと依存パッケージを用意するため、
初回実行時はインターネット接続が必要。

## 実行手順

リポジトリ、または再ビルド用ソースZIPの [ビルド手順](../watchface/README.md) に従ってビルドする。
ZIPで作業する場合は、ソースZIPと最適化用ZIPを同じ親フォルダへ展開する。
続いて `optimization/` で実行する。

```sh
uv run optimize-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6.zpk --source ../watchface --fixed
```

出力は `publish/lumen-orbit-bip6.zpk`。`publish/` は自動で作成される。
入力の `dist/` にZABが複数ある場合は、対象の `.zab` を直接指定する。
標準ビルドから取り出した `.zpk` も入力できる。
`--source` には、そのビルドに対応するソースディレクトリを指定する。

## 適用する処理

1. 元PNGのRGBAを端末用パレット画像に直接反映し、ビルド時の減色による差を解消する。
2. 重複色・未使用パレット枠を整理する。画像ヘッダー、寸法、16px単位の保存幅と行末の余白は維持する。
3. 内側ZIPをSTORE（無圧縮）にし、外側のみ標準DEFLATEで圧縮する。
4. `--fixed` では元のファイル順を使用し、Zopfli 120 iterationsとzlib level 9を比較して短い方を採用する。

`--fixed` を省略すると、パレット順・ZIP内の並び順・内側の圧縮方法も比較する。

## 自動検証

ZIPの読み込み時にCRCを検証し、重複エントリーも拒否する。
各エントリーは1回で読み込む。

Pillowによる独立した復号で、28枚の端末用画像の表示領域が元PNGの色・透明度と
完全一致することを確認する。画像ヘッダーと行末の余白、および設定・有効期限・
署名・実行コードなど画像以外のデータは、入力パッケージと同一であることを検証する。
想定外の画像形式や、ビルドとソースのアプリID・名称・バージョンの不一致はエラーにする。

開発方針・配布物の作成と引き継ぎは [開発者向け文書](../AGENTS.md) を参照。
