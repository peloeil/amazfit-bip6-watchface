# LUMEN ORBIT 1.9.2 — 再ビルド

対象はAmazfit Bip 6。Node.js・npm・インターネット接続が必要。
ビルドツールは `package-lock.json` で固定している（Zeus CLI 1.9.3）。

## ビルド

リポジトリ、または再ビルド用ZIPに含まれる `watchface/` で実行する。

```sh
npm ci
node build.js build
```

成果物は `dist/` の `.zab` ファイル。
`.zab` はZIP形式で、その中の `.zpk` が端末向けパッケージ。

ビルド時にZeppの機種情報を取得する。取得に失敗した場合は、
ネットワークを確認してビルドを再実行する。

## 配布用ZPKの作成

配布用ZPKは、標準ビルドの画像をソースPNGに合わせて最適化して作成する。
リポジトリの `optimization/`、または別配布の最適化用ZIPを使う。
ZIPで作業する場合は、両方のZIPを同じ親フォルダへ展開する。

`uv` を用意し、`watchface/` から次を実行する。

```sh
cd ../optimization
uv run optimize-zpk.py ../watchface/dist ../publish/lumen-orbit-bip6.zpk --source ../watchface --fixed
```

詳細は [最適化の手順](../optimization/README.md) を参照。
開発方針・配布物の作成と引き継ぎは [開発者向け文書](../AGENTS.md) にまとめている。
