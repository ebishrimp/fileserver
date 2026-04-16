# 機能
- MySQL によるファイルパス保存
- host 側のディレクトリ・ファイル管理
- crient 側からのファイル操作リクエスト
- 各言語対応のライブラリおよびシェル用コマンド
- LAN内運用を想定

# 仕様
### 管理
- hardlayer, applayer 二層化によるファイル名衝突の回避
- hard, app, name クエリパラメータにより管理
- MySQL 側の初期設定およびユーザー管理が必須
- インジェクション対策
- 送信前にtarアーカイブ化を実行
### リクエスト
- http://host:port/push?params=values によりファイルアップロード(上書き禁止)
- http://host:port/pull?params=values によりファイルダウンロード
- http://host:port/delete?params=values によりファイル削除
- http://host:port/overwrite?params=values によりファイル上書き
- 無効なパスおよびパラメーターに対しhttpエラーを出力

to do<br>
・ディレクトリ操作<br>
・openAPIによるclient側の整備基盤作成(openapi.yml)
・ファイル送受信
・mysqlユーザー, パスワードおよびhostポート番号管理機能