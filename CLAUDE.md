# 開発ガイドライン

## 目的

本プロジェクトは「学習進捗管理Webアプリ」のバックエンドREST APIである。  
設計は **ドメイン駆動設計（Domain-Driven Design, DDD）** と  
**クリーンアーキテクチャ（Clean Architecture, CA）** に従う。

HTTP層は **Huma（huma）** を利用し、OpenAPI仕様の生成もHumaに委ねる。

DBは PostgreSQL を第一候補とし、DBアクセスは **pgx** と **sqlc** を利用する。

## アーキテクチャ原則

### 依存方向の原則（最重要）

依存は常に **外側 → 内側** に向かう。

repository -> controller -> usecase -> domain

## レイヤ責務

### domain

エンティティ層は、ビジネスルールやオブジェクトの集合を表す。

#### 含むもの

- Entity
- Value Object
- Domain Service
- Domain Error

#### 禁止

- DBアクセス
- HTTP依存
- 外部ライブラリ依存（可能な限り）

### usecase

ユースケース層では、アプリケーション固有のビジネスルールを定義する。

#### 含むもの

- Usecase struct
- Repository interface（`internal/usecase/port`）
- ドメインルールの組み合わせ

#### 禁止

- HTTPレスポンス生成
- JSONタグ付きDTO
- DB具体実装

### controller

HumaによるHTTP層。

#### 含むもの

- Huma handler
- Request/Response DTO (`internal/controller/dto`)
- 入力バリデーション
- エラー変換（domain/application → HTTP）

#### 原則

- domain entity を直接返さない
- DTO変換を明示する

### repository (postgres)

外部I/Oの実装。

#### 含むもの

- Repository実装（PostgreSQL）
- DB接続（pgx）
- sqlc 生成コードの利用
- 設定管理
- マイグレーション適用

#### 役割

usecase層（port）のRepository interfaceを実装する。

## DB/Repository 実装方針（pgx + sqlc）

### 採用方針

- PostgreSQL のドライバ/接続管理は **pgx** を使用する
- クエリは SQL を明示的に書き、コード生成は **sqlc** に任せる
- 生成コードは **repository 層でのみ利用**し、domain/usecase に漏らさない

### SQL配置ルール

- `db/query/*.sql` にクエリを置く
- `db/migrations/*.sql` にマイグレーションを置く
- `sqlc.yaml` をルートに置く
- `make sqlc` で生成できるようにする

### Repository 実装ルール

- usecase層の `internal/usecase/port` に Repository interface を定義する
- `internal/repository/postgres` で interface を実装する
- sqlc の生成型（Params/Row/Model）は **repository 内に閉じ込める**
    - usecase/domain に返す型は domain entity に変換する
- SQLのエラー（unique違反、not found 等）は repository で吸収し、
  domain のエラーに変換して返す

## エラーハンドリング規約

### domain error

domain内で定義する。

例:

- `ErrNotFound`
- `ErrInvalidArgument`

### HTTP error変換規則

`internal/controller/error.go` で変換を行う。

| domain error | HTTP |
|--------------|------|
| NotFound     | 404  |
| Invalid      | 400  |
| Conflict     | 409  |
| その他          | 500  |

Humaのエラーレスポンス形式に統一する。

## DTO規約

- Request/Responseは `internal/controller/dto` に定義
- JSONタグは DTO のみに置く
- domain構造体にJSONタグを付けない
- domainとDTOは明示的に変換関数を作る（`ToXxxResponse` 等）

### Huma利用方針

- ルーティングは Huma を使用（`internal/controller/router.go`）
- OpenAPIはHumaに自動生成させる
- バリデーションはHumaの型定義（struct tag）を活用

## テスト方針

### domain

- 純粋な単体テストを書く
- DB不要

### usecase

- repositoryをモックしてテスト

### repository

- 必要に応じて統合テスト（PostgreSQLを使う場合はテスト用DBを用意）

## 新機能追加ルール

新しいAPIを追加する場合:

1. domainに必要なEntity/ValueObjectを追加
2. usecaseの `port` に repository interface定義を追加
3. usecaseにUsecase追加
4. SQL追加（`db/query/*.sql`）→ `make sqlc` で生成
5. repositoryにRepository実装（sqlc生成コード利用、型変換）
6. controllerの `dto` に型定義を追加
7. controller(Huma)でhandler作成し `router.go` で登録
8. テスト追加

この順序を守る。

## リファクタリング時の注意

- 依存方向が壊れていないか必ず確認
- domainに外部依存が入り込んでいないか確認
- DTOとdomainが混ざっていないか確認
- sqlc生成型がusecase/domainに漏れていないか確認

## コーディング規約

- `make fmt` 必須
- `make lint` を通す
- 明確な命名（`UserRepository` など）
- 曖昧なパッケージ名を避ける（`util`, `common`は禁止）

## このプロジェクトでのアンチパターン

- handlerから直接DB呼び出し
- domainにjsonタグ
- repositoryからdomainを書き換える
- usecaseを通さずrepository呼び出し
- sqlc生成型をusecase/domainへ返す（漏らす）

## 最重要原則

このプロジェクトは「動くコード」よりも設計を守ることを優先する。
