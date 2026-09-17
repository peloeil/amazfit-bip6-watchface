# LUMEN ORBIT

Amazfit Bip 6用の文字盤。アナログ時計に、デジタル時刻・日付・日本語の曜日・電池残量を組み合わせる。
現在のバージョンは **1.9.2**。

![LUMEN ORBITの表示例](watchface/assets/bip6/preview.png)

画像は旧版のプレビュー。1.9.2のデジタル表示は「時:分」で、アナログ秒針は維持している。

## 特徴

- 黒背景のアナログ時計。時針・分針・秒針を表示する。
- デジタル時刻は時・分を表示する。
- 月日、日本語の曜日、電池残量のパーセントを表示する。

## 必要な環境

このリポジトリで案内する導入方法には、次の環境が必要。

- Amazfit Bip 6と、Zeppで時計に接続したAndroid端末。
- Windows x64のPC。
- PCとAndroidの両方で接続したTailscale。
- Zeppの開発者モード。

文字盤の導入にNode.js・Python・Goのインストールは不要。
初回設定は [導入ツールの手順](link-windows/README.txt) を参照。

## ダウンロード・導入

配布ファイルはGitHub Releasesで公開予定。
公開後は、このリポジトリのReleasesを開き、同じリリースのAssetsから次の2ファイルを取得する。

| ファイル | 用途 |
| --- | --- |
| `lumen-orbit-bip6.zpk` | 時計へ導入する文字盤1.9.2 |
| `lumen-orbit-link-windows.zip` | Windows導入ツール LUMEN ORBIT LINK 1.0.6 |

1. Windows導入ツールのZIPを展開する。
2. 文字盤のZPKを `lumen-orbit-link.exe` にドラッグ＆ドロップする。
3. PCのブラウザに表示されたQRを、AndroidのZeppの開発者モードで読み取る。
4. Zeppと時計で転送完了を確認し、PC画面の「終了」を押す。

詳しい操作や接続できない場合の確認方法は [導入ツールの手順](link-windows/README.txt) を参照。

## 改変・ビルド

- [文字盤のビルド](watchface/README.md)
- [ビルドしたZPKの最適化](optimization/README.md)
- [導入ツールのテスト・再ビルド](link-windows/README.txt)
- [構成・開発方針・リリース手順](AGENTS.md)

再ビルド用のソースZIPと最適化用ZIPは、開発者向けの配布物としてReleasesに添付する。

## ライセンス

自作コード・文字盤画像のライセンスは未設定。
第三者由来のコードのライセンスは、次のファイルを参照。

- [QRコードライブラリ](link-windows/QR-LICENSE.txt)
- [Goランタイム](link-windows/GO-LICENSE.txt)
