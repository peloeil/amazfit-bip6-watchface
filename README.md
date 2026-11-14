# LUMEN ORBIT

Amazfit Bip 6用watchface **2.0.3**。
アナログ時計に、デジタル時刻・日付・日本語の曜日・電池残量を組み合わせる。

![LUMEN ORBITの表示例](watchface/assets/bip6/preview.png)

## 必要な環境

- Amazfit Bip 6と、Zeppで時計に接続したスマートフォン。
- Windows x64のPC。
- PCとスマートフォンの両方で接続したTailscale。
- Zeppの開発者モード。

初回設定・詳しい操作・接続できない場合は [導入ツールの手順](zpk-link/README.txt) を参照。

## ダウンロード・導入

[GitHub ReleasesのLUMEN ORBIT 2.0.3](https://github.com/peloeil/amazfit-bip6-watchface/releases/tag/watchface-v2.0.2)から次の2ファイルを取得する。

| ファイル | 用途 |
| --- | --- |
| `lumen-orbit-bip6.zpk` | 時計へ導入するwatchface |
| `zpk-link-windows.zip` | ZPK配信用ツール ZPK LINK 1.0.6（Windows） |

1. Windows導入ツールのZIPを展開し、中の `zpk-link/` フォルダを開く。
2. watchfaceのZPKを、そのフォルダの `zpk-link.exe` にドラッグ＆ドロップする。
3. PCに表示されたQRを、スマートフォンのZeppの開発者モードで読み取る。
4. Zeppと時計で転送完了を確認し、PC画面の「終了」を押す。

## 改変・ビルド

```text
watchface/ のソース・画像
  → Zeusでビルド → watchface/dist/*.zab
  → packaging/prepare-zpk.py で画像補正・検証・圧縮
  → publish/lumen-orbit-bip6.zpk → GitHub Releasesで配布
```

| パス | 役割 |
| --- | --- |
| `watchface/src/index.js` | 時計の表示を作るコード |
| `watchface/assets/bip6/` | 元画像 |
| `watchface/app.json`・`watchface/build.js` | 機種・バージョン設定とビルド |
| `packaging/` | 元画像の色・透明度を反映し、配布用ZPKを作成・検証する処理 |
| `zpk-link/` | Windows用の導入ツール |
| `publish/` | 配布物の出力先（Git管理対象外） |

- [watchfaceのビルド](watchface/README.md)
- [配布用ZPKの作成](packaging/README.md)
- [導入ツールのテスト・再ビルド](zpk-link/README.txt)
- [開発方針・配布手順](AGENTS.md)

### ソースZIPから作業する場合

同じリリースの `lumen-orbit-bip6-source.zip` と `lumen-orbit-bip6-packaging.zip` を取得し、
それぞれの中身を1つの作業フォルダへ展開する。ZIP名のフォルダが作られた場合は、
その中身を移して、次の配置に揃える。

```text
作業フォルダ/
├── AGENTS.md
├── LICENSE
├── watchface/
│   ├── src/index.js
│   ├── assets/bip6/
│   ├── app.json
│   ├── build.js
│   └── README.md
└── packaging/
    ├── prepare-zpk.py
    └── README.md
```

`watchface/` と `packaging/` が同じ階層にあることを確認する。
重複する `AGENTS.md` と `LICENSE` は同じリリースなら同一内容なので、上書きしてよい。

## ライセンス

自作コード・watchface画像は [MIT License](LICENSE) で公開する。
数字・記号・曜日の画像は、[Inter](https://github.com/rsms/inter) と
[Noto Sans JP](https://github.com/notofonts/noto-cjk) で描画した。
どちらのフォントも [SIL Open Font License 1.1](https://openfontlicense.org) で提供されている。
第三者コードには、それぞれのライセンスが適用される。

- [QRコードライブラリ](zpk-link/QR-LICENSE.txt)
- [Goランタイム](zpk-link/GO-LICENSE.txt)
