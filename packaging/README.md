# LUMEN ORBIT 2.0.1 — 配布用ZPKの作成

標準ビルドの `.zab` / `.zpk` に元PNGの色・透明度を反映し、検証・圧縮した配布用ZPKを作成する。
[uv](https://docs.astral.sh/uv/getting-started/installation/) が必要。
Pythonと依存パッケージは `prepare-zpk.py` のインラインメタデータに従って用意される。
初回実行時はインターネット接続が必要。

処理の流れ: watchfaceのソース・画像 → ZeusでZABを作成 → このスクリプトで画像補正・検証・圧縮 → 配布用ZPK。

## ZIPの展開先

同じリリースの `lumen-orbit-bip6-source.zip` と `lumen-orbit-bip6-packaging.zip` の中身を、
1つの作業フォルダへ展開する。ZIP名のフォルダが作られた場合は、その中身を移して次の配置にする。

```text
作業フォルダ/
├── AGENTS.md
├── LICENSE
├── watchface/       ← src/、assets/、app.json、build.js
└── packaging/      ← このREADME、prepare-zpk.py
```

`watchface/` と `packaging/` は同じ階層に置く。
重複する `AGENTS.md` と `LICENSE` は同じリリースなら同一内容なので、上書きしてよい。

## 実行手順

[watchfaceをビルド](../watchface/README.md)し、作業フォルダの `packaging/` で実行する。

```sh
uv run prepare-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6.zpk --source ../watchface --fixed
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
