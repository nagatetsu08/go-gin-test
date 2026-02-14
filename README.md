# gin-freemarket

Go (Gin) で構築したフリーマーケット風の REST API です。
ユーザー認証（JWT）付きで、商品の CRUD 操作を提供します。

## 技術スタック

| カテゴリ | 技術 |
|---------|------|
| 言語 | Go 1.24 |
| Web フレームワーク | Gin |
| ORM | GORM |
| データベース | PostgreSQL 16 (本番) / SQLite (テスト) |
| 認証 | JWT (golang-jwt) |
| コンテナ | Docker Compose |
| ホットリロード | Air |
| DB 管理 UI | pgAdmin 4 |

## プロジェクト構成

```
.
├── controllers/    # リクエストハンドラ
├── services/       # ビジネスロジック
├── repositories/   # データアクセス層
├── models/         # GORM モデル定義 (Item, User)
├── dto/            # リクエスト/レスポンス用 DTO
├── middlewares/    # JWT 認証ミドルウェア
├── infra/          # DB 接続・環境変数の初期化
├── migrations/     # GORM AutoMigrate
├── docker/         # Docker 関連ファイル
└── main.go         # エントリーポイント
```

## API エンドポイント

| メソッド | パス | 認証 | 説明 |
|---------|------|------|------|
| GET | `/items/` | 不要 | 商品一覧取得 |
| GET | `/items/:id` | 必要 | 商品詳細取得 |
| POST | `/items/` | 必要 | 商品登録 |
| PUT | `/items/:id` | 必要 | 商品更新 |
| DELETE | `/items/:id` | 必要 | 商品削除 |
| POST | `/auth/signup` | 不要 | ユーザー登録 |
| POST | `/auth/login` | 不要 | ログイン (JWT 発行) |

## セットアップ

### 前提条件

- Go 1.24+
- Docker / Docker Compose

### 1. リポジトリのクローン

```bash
git clone <repository-url>
cd gin-freemarket
```

### 2. 環境変数の設定

`.env` ファイルをプロジェクトルートに作成します。

```env
ENV=prod
DB_HOST=localhost
DB_USER=ginuser
DB_PASSWORD=ginpassword
DB_NAME=fleamarket
DB_PORT=5432
SECRET_KEY=<任意のシークレットキー>
```

### 3. データベースの起動

```bash
docker-compose up -d
```

PostgreSQL が `localhost:5432`、pgAdmin が `localhost:81` で起動します。

### 4. マイグレーション

```bash
go run migrations/migration.go
```

### 5. サーバー起動

```bash
# 通常起動
go run main.go

# ホットリロード (Air)
air
```

サーバーは `localhost:8080` で起動します。

## テスト

テスト時は SQLite（インメモリ）が自動的に使用されます。

```bash
go test ./...
```