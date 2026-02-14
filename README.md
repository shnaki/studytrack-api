# StudyTrack API

学習進捗管理Webアプリのバックエンド REST API サーバ。

## 技術スタック

- **言語**: Go 1.23+
- **HTTP/OpenAPI**: [Huma v2](https://huma.rocks/)
- **ルーター**: Chi v5
- **DB**: PostgreSQL 16 + pgx v5
- **マイグレーション**: golang-migrate
- **アーキテクチャ**: DDD × Clean Architecture

## アーキテクチャ

```
cmd/api/main.go            # エントリポイント・DI配線
internal/
  domain/                  # エンティティ・ドメインエラー（外部依存なし）
  application/
    port/                  # リポジトリインタフェース（port）
    usecase/               # ユースケース実装
  infrastructure/
    config/                # 環境変数ベースの設定
    postgres/              # DB接続・マイグレーション・リポジトリ実装
  handler/                 # Huma ハンドラー・DTO・エラー変換・ルーター
db/
  migrations/              # SQLマイグレーション
  query/                   # sqlc用クエリSQL
```

依存方向: `handler → usecase → domain ← infrastructure`

## セットアップ

### 前提条件

- Go 1.23+
- PostgreSQL 16+ （またはDocker）
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)（任意: Makefile経由でマイグレーション）

### 1. PostgreSQL起動

```bash
docker compose up -d
```

### 2. 環境変数

```bash
cp .env.example .env
# 必要に応じて .env を編集
```

### 3. APIサーバ起動

```bash
make dev
```

起動後、マイグレーションが自動適用されます。

## 開発コマンド

| コマンド | 説明 |
|---|---|
| `make dev` | APIサーバ起動（ホットリロードなし） |
| `make build` | バイナリビルド |
| `make test` | テスト実行 |
| `make test-cover` | カバレッジ付きテスト |
| `make lint` | golangci-lint 実行 |
| `make fmt` | フォーマット |
| `make migrate-up` | マイグレーション適用 |
| `make migrate-down` | マイグレーション1つロールバック |
| `make migrate-create` | 新規マイグレーション作成 |
| `make sqlc` | sqlcコード生成 |
| `make docker-build` | Dockerイメージビルド |
| `make docker-up` | docker compose up |

## OpenAPI

APIサーバ起動後、以下で確認できます:

- **ドキュメント**: http://localhost:8080/docs
- **OpenAPI JSON**: http://localhost:8080/openapi.json
- **OpenAPI YAML**: http://localhost:8080/openapi.yaml

## API エンドポイント

### Users
- `POST /users` - ユーザー作成
- `GET /users/{id}` - ユーザー取得

### Subjects
- `POST /users/{userId}/subjects` - 学習分野作成
- `GET /users/{userId}/subjects` - 学習分野一覧
- `PUT /subjects/{id}` - 学習分野更新
- `DELETE /subjects/{id}` - 学習分野削除

### StudyLogs
- `POST /users/{userId}/study-logs` - 学習記録作成
- `GET /users/{userId}/study-logs?from=&to=&subjectId=` - 学習記録一覧
- `DELETE /study-logs/{id}` - 学習記録削除

### Goals
- `PUT /users/{userId}/goals/{subjectId}` - 目標設定（upsert）
- `GET /users/{userId}/goals` - 目標一覧

### Stats
- `GET /users/{userId}/stats/weekly?weekStart=YYYY-MM-DD` - 週次統計

## 環境変数

| 変数 | デフォルト | 説明 |
|---|---|---|
| `PORT` | `8080` | サーバポート |
| `DB_URL` | `postgres://studytrack:studytrack@localhost:5432/studytrack?sslmode=disable` | PostgreSQL接続URL |
| `CORS_ORIGINS` | `http://localhost:5173` | CORS許可オリジン（カンマ区切り） |
| `LOG_LEVEL` | `debug` | ログレベル |

## エラーレスポンス

Huma の RFC 9457 (Problem Details) 形式で統一:

```json
{
  "status": 404,
  "title": "Not Found",
  "detail": "user not found"
}
```

| ドメインエラー | HTTPステータス |
|---|---|
| NotFound | 404 |
| Validation | 400 |
| Conflict | 409 |
| その他 | 500 |
