# m-bank

## 事前準備

環境変数の読み込みに [direnv](https://github.com/direnv/direnv) を使用しています

```bash
cp .envrc-example .envrc
```

```bash
# DB 起動
docker-compose up -d

# データベース作成、初期データ投入
make create
```

## 動作確認

サーバーを起動

```bash
make serve
```

curl の例

```bash
# ユーザの残高を確認
curl http://127.0.0.1:3000/balances/1

# 残高の加減算（仮押さえ）
curl --request POST \
  --url http://127.0.0.1:3000/payments/try \
  --header 'content-type: application/json' \
  --data '{
  "idempotency_key":"foobar",
  "user_id":1,
  "amount":100
}'

# 残高の加減算（確定）
curl --request POST \
  --url http://127.0.0.1:3000/payments/confirm \
  --header 'content-type: application/json' \
  --data '{
  "idempotency_key":"foobar",
  "user_id":1,
  "amount":100
}'

# 残高の加減算（キャンセル）
curl --request POST \
  --url http://127.0.0.1:3000/payments/cancel \
  --header 'content-type: application/json' \
  --data '{
  "idempotency_key":"foobar",
  "user_id":1,
  "amount":100
}'
```
