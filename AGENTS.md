# 開発・リリース

## 性能最適化

- 表示・配置・回転中心・透明度・機能を保ち、不要な処理・描画・メモリ使用を減らす。
  標準ウィジェットで済む処理を自作しない。
- 容量だけを減らす変更や、容量のために復号・描画処理を増やす変更は行わない。
- 削減できる処理・メモリを説明できなければ変更しない。無駄がなければ完了とする。

## 作業別の手順

- watchfaceのビルド: `watchface/README.md`
- 配布用ZPKの作成・画像とパッケージの検証: `packaging/README.md`
- 導入ツールのテスト・ビルド: `zpk-link/README.txt`

## 配布物の作成・引き継ぎ

配布物はGit管理対象外の `publish/` に作成し、固定のファイル名でGitHub Releasesへ添付する。
watchfaceと導入ツールのバージョンは独立して管理し、端末の更新判定に使う内部バージョン番号を省略しない。

watchfaceの配布物を引き継ぐ場合は、同じリリースの次の3ファイルをこの順でまとめて提示する。

1. `lumen-orbit-bip6.zpk`: 最適化済みのwatchface。
2. `lumen-orbit-bip6-source.zip`: `watchface/` のソース・画像・設定・ビルド手順。
3. `lumen-orbit-bip6-packaging.zip`: `packaging/` のスクリプト・実行手順。

導入ツールは `zpk-link-windows.zip` に、`zpk-link/` のWindows x64実行ファイル・
ソース・操作手順・第三者ライセンスを収録する。
各ZIPのルートに `AGENTS.md` と、プロジェクトのMITライセンス文書 `LICENSE` を含める。

## リリース

- 対象ソースからビルド・検証した配布物を揃え、ZIPの内容とバージョンがソースに対応することを確認する。
  対応するコミットへ `watchface-v1.9.2`・`link-v1.0.6` のように対象を区別したタグを付ける。
- watchfaceのReleaseには上記4ファイル、導入ツール単独のReleaseにはWindows用ZIPを添付する。
  本文には各バージョンと実機確認の範囲を記載し、利用者向けのZPK・導入ツールを先に案内する。
- READMEからwatchfaceのReleaseへ案内する。導入ツール単独のReleaseでLatestを更新しない。
