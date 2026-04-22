# 機能
- MySQL によるファイルパス保存
- host 側のディレクトリ・ファイル管理
- crient 側からのファイル操作リクエスト
- 各言語対応のライブラリおよびシェル用コマンド
- LAN内運用を想定
- 疑似 raid 0 対応
- prometheus exporter 対応

# 仕様
### 管理
- _hardlayer, applayer 二層化によるファイル名衝突の回避
- _hard, app, name クエリパラメータにより管理
- MySQL 側の初期設定およびユーザー管理が必須
- _インジェクション対策
- 送信前にtar.gz圧縮を実行
- 疑似 raid 0 有効の場合push/overwriteにより2つのディレクトリに書き込む
- whitelist有効時クライアントIPアドレスの照合
- クライアントIPアドレスおよび処理内容をログファイルに書き込み
### リクエスト
- http://host:port ...
- /upload?params=values によりファイルアップロード(上書き禁止)
- /download?params=values によりファイルダウンロード
- /delete?params=values によりファイル削除
- /overwrite?params=values によりファイル上書き
- _無効なパスおよびパラメーターに対しhttpエラーを出力

to do<br>
・ディレクトリ操作<br>
・openAPIによるclient側の整備基盤作成(openapi.yml)<br>
・ファイル送受信<br>
・mysqlユーザー, パスワードおよびhostポート番号管理機能<br>
・delete時 ID scanエラー<br>

備考<br>
実装済み項目は _付き