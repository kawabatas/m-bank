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

# サーバ起動
make serve

# テスト
make test
```
