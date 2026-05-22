# 機能
- MySQL によるファイルパス保存
- host(server) 側のディレクトリ・ファイル管理
- client 側からのファイル操作リクエスト
- 各言語対応のライブラリおよびシェル用コマンド
- 疑似 raid 1 (ミラーリング) 対応
- prometheus exporter 対応
- *LAN内運用を想定

# 仕様
### 管理
- hardlayer, applayer 二層化によるファイル名衝突の回避
- hard, app, name クエリパラメータ
- MySQL 側の初期設定およびユーザー管理が必須
- 致命的なエラーによりサービス停止(クエリエラーでは停止しない)
- 設定ファイル /etc/fileserver/fileserver.conf
- whitelistファイル /etc/fileserver/whitelist.conf
- logファイル(パス変更可) /var/log/fileserver/fileserver.log
- デフォルトポート 50080 (変更可)
### エンドポイント
- http://host:port ...
- /upload?params=values アップロード(上書き禁止)
- /download?params=values ダウンロード
- /delete?params=values 削除
- /overwrite?params=values 上書き
- *無効なパスおよびパラメーターに対しhttpエラーを出力

to do<br>
・ディレクトリ操作<br>
・openAPIによるclient側の整備基盤作成(openapi.yml)<br>
・ファイル送受信<br>
・mysqlユーザー, パスワードおよびhostポート番号管理機能<br>
・delete時 ID scanエラー<br>
・疑似 raid 1 有効の場合push/overwriteにより2つのディレクトリに書き込む<br>
・送信前にtar.gz圧縮を実行<br>
・whitelist有効時クライアントIPアドレスの照合<br>
・クライアントIPアドレスおよび処理内容をログファイルに書き込み<br>