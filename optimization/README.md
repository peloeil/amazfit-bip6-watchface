# LUMEN ORBIT 1.9.2 — ZPK最適化

標準ビルドの `.zab` / `.zpk` を、対応するソースPNGから最適化する。
[uv](https://docs.astral.sh/uv/getting-started/installation/) が必要。
Pythonと依存パッケージは `optimize-zpk.py` のインラインメタデータに従って用意される。
初回実行時はインターネット接続が必要。

## 実行手順

[watchfaceをビルド](../watchface/README.md)し、`optimization/` で実行する。
ZIPで作業する場合は、ソースZIPと最適化用ZIPを同じ親フォルダへ展開する。

```sh
uv run optimize-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6.zpk --source ../watchface --fixed
```

出力先の `publish/` は自動で作成される。
`dist/` にZABが複数ある場合は、対象の `.zab` を直接指定する。`.zpk` も入力できる。
`--source` には、そのビルドに対応するソースディレクトリを指定する。

## 処理と検証

- 元PNGの色・透明度を端末用画像に反映し、重複色・未使用パレット枠を整理する。
- 画像ヘッダー・寸法・保存幅・行末の余白と、設定・署名・実行コードなど画像以外のデータを維持する。
- `--fixed` ではファイル順を維持し、内側ZIPを無圧縮、外側をDEFLATEで圧縮する。
  省略すると、パレット順・ファイル順・内側の圧縮方法も比較する。

実行時にZIPのCRC・重複エントリーを検査し、Pillowで元PNGとの色・透明度の一致を検証する。
変更対象外のデータの一致も検証する。想定外の画像形式や、ソースとパッケージの
アプリID・名称・バージョンの不一致はエラーにする。

開発・配布のルールは [AGENTS.md](../AGENTS.md) を参照。

スクリプトは [MIT License](../LICENSE) で公開する。
