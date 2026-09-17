# LUMEN ORBIT 1.9.1 — ZPK最適化

## 成果物の返却・引き継ぎルール

配布ファイル名は固定し、履歴はGitのコミットで管理します。
アプリ内部のバージョン番号は、端末の更新判定のため維持します。

このウォッチフェイスの成果物を返す際は、変更の有無にかかわらず、最終回答に
対応する最新版の次の3ファイルを、必ずこの順でまとめて提示してください。

1. `lumen-orbit-bip6.zpk`：最適化済みの端末向けパッケージ。
2. `lumen-orbit-bip6-rebuild-source.zip`：再ビルドに必要なソース・画像・設定・依存関係・手順。
3. `lumen-orbit-bip6-optimization.zip`：最適化に必要なスクリプト・依存関係・実行手順。

3ファイルは同じリリースに対応させてください。source ZIPと最適化用ZIPは分離し、
不要な照合結果・作業記録は含めません。別セッションでもこの3ファイルを使うため、
この返却ルールを両ZIPのREADMEに引き継いでください。

標準ビルドの .zab / .zpk を、対応するソースPNGから最適化します。
[uv](https://docs.astral.sh/uv/getting-started/installation/) を使用します。
Python 3.10以降と依存パッケージ（Pillow・zopfli）の固定バージョンは、
optimize-zpk.py 冒頭の [インラインメタデータ](https://docs.astral.sh/uv/guides/scripts/#declaring-script-dependencies)
で管理します。uv が必要なPythonと依存パッケージを自動で用意するため、
初回実行時はインターネット接続が必要です。

## 実行手順

再ビルド用source ZIPとこのZIPを同じ親フォルダへ展開します。
source側のREADMEに従ってビルドした後、
optimization ディレクトリで実行してください。

```sh
uv run optimize-zpk.py ../watchface/dist ./lumen-orbit-bip6.zpk --source ../watchface --fixed
```

出力は、このディレクトリの lumen-orbit-bip6.zpk です。
入力のdistにZABが複数ある場合は、対象の .zab を直接指定してください。
標準ビルドから取り出した .zpk も入力できます。
--source には、そのビルドに対応するソースディレクトリを指定します。

## 適用する処理

1. 元PNGのRGBAを端末用パレット画像に直接反映し、ビルド時の減色による差を解消。
2. 重複色・未使用パレット枠を整理。画像ヘッダー、寸法、16px単位の保存幅と行末の余白は維持。
3. 内側ZIPをSTORE（無圧縮）にし、外側のみ標準DEFLATEで圧縮。
4. --fixed では元のファイル順を使用し、Zopfli 120 iterationsとzlib level 9を比較して短い方を採用。

--fixed を省略すると、パレット順・ZIP内の並び順・内側の圧縮方法も比較します。

## 自動検証

Pillowによる独立した復号で、28枚の端末用画像の表示領域が元PNGの色・透明度と
完全一致することを確認します。画像ヘッダーと行末の余白、および設定・有効期限・
署名・実行コードなど画像以外のデータは、入力パッケージと同一であることを検証します。
想定外の画像形式や、ビルドとソースのアプリID・名称・バージョンの不一致はエラーにします。

今回の右寄せ版の出力は16,276Bです。再ビルド時の日時・一時パスなどによって、
圧縮後の容量は多少変わり得ます。実機でのインストールと通知表示は未確認です。
