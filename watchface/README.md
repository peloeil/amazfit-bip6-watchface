# LUMEN ORBIT 2.0.0 — 再ビルド

Node.js・npm・インターネット接続が必要。
ビルドツールは `package-lock.json` で固定している（Zeus CLI 1.9.3）。

`src/index.js` が時計の表示を作るコード、`assets/bip6/` が元画像、
`app.json` が機種・バージョン・実行するソースの設定。

ビルドから配布まで: ソース・画像 → ZeusでZABを作成 → 元画像の色・透明度を反映し検証 → 配布用ZPK。

## ZIPの展開先

同じリリースの `lumen-orbit-bip6-source.zip` と `lumen-orbit-bip6-packaging.zip` の中身を、
1つの作業フォルダへ展開する。ZIP名のフォルダが作られた場合は、その中身を移して次の配置にする。

```text
作業フォルダ/
├── AGENTS.md
├── LICENSE
├── watchface/       ← このREADME、src/、assets/、app.json、build.js
└── packaging/      ← prepare-zpk.py、README.md
```

`watchface/` と `packaging/` は同じ階層に置く。
重複する `AGENTS.md` と `LICENSE` は同じリリースなら同一内容なので、上書きしてよい。

## ビルド

リポジトリ、またはソースZIPの `watchface/` で実行する。

```sh
npm ci
node build.js build
```

成果物は `dist/` の `.zab` ファイル。ZIP形式で、端末向けパッケージの `.zpk` を含む。
ビルド時にZeppの機種情報を取得するため、取得に失敗した場合はネットワークを確認して再実行する。

## 配布用ZPKの作成

標準ビルドでは画像の色・透明度が元PNGと異なる場合があるため、
[配布用ZPKの作成手順](../packaging/README.md) で元PNGの値を反映し、検証したZPKを作成する。
成果物は作業フォルダの `publish/lumen-orbit-bip6.zpk`。

開発・配布のルールは [AGENTS.md](../AGENTS.md) を参照。

ソースコード・画像は [MIT License](../LICENSE) で公開する。
