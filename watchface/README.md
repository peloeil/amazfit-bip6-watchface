# LUMEN ORBIT 1.9.2 — 再ビルド

Node.js・npm・インターネット接続が必要。
ビルドツールは `package-lock.json` で固定している（Zeus CLI 1.9.3）。

## ビルド

リポジトリ、またはソースZIPの `watchface/` で実行する。

```sh
npm ci
node build.js build
```

成果物は `dist/` の `.zab` ファイル。ZIP形式で、端末向けパッケージの `.zpk` を含む。
ビルド時にZeppの機種情報を取得するため、取得に失敗した場合はネットワークを確認して再実行する。

## 配布用ZPKの作成

[最適化の手順](../optimization/README.md) に従い、ビルドした画像をソースPNGに合わせたZPKを作成する。
ZIPで作業する場合は、別配布の最適化用ZIPをソースZIPと同じ親フォルダへ展開する。

開発・配布のルールは [AGENTS.md](../AGENTS.md) を参照。
