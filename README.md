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

# 残高の加減算（仮登録）
curl --request POST \
  --url http://127.0.0.1:3000/payments/try \
  --header 'content-type: application/json' \
  --data '{
  "idempotency_key":"foobar",
  "user_id":1,
  "amount":100
}'

# 残高の加減算（本実行）
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

## 説明

### 機能要件

#### 1. 顧客(ユーザー)の保有する残高を加減算する仕組み

- 残高の加算/減算は周辺のマイクロサービスから要求される
- エラーや障害によって再送される可能性がある

呼び出し元のマイクロサービス側で、エラーや障害によるリトライを実装してもらう想定で、TCC(Try-Confirm/Cancel)パターンで REST API を用意しました。

#### 2. すべての顧客の残高に一斉に残高を加算する仕組み

数千数万ユーザずつ、バッチで処理（バッチが状態を保存する）されることを想定した REST API を用意しました。

なお、REST API の詳細ドキュメントは [swagger.yml](https://github.com/kawabatas/m-bank/swagger.yml) をご覧ください。
