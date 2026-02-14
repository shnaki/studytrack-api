# CLAUDE.md

## 目的

本プロジェクトは「学習進捗管理Webアプリ」のバックエンドREST APIである。  
設計は **ドメイン駆動設計（Domain-Driven Design, DDD）** と  
**クリーンアーキテクチャ（Clean Architecture, CA）** に従う。

HTTP層は **Huma（huma）** を利用し、OpenAPI仕様の生成もHumaに委ねる。

DBは PostgreSQL を第一候補とし、DBアクセスは **pgx** と **sqlc** を利用する。

---

# アーキテクチャ原則

## 1. 依存方向の原則（最重要）

依存は常に **外側 → 内側** に向かう。

interface (HTTP/Huma)  
↓  
application (usecase)  
↓  
domain  

infrastructure は application の port(interface) を実装する。

### 禁止事項

- domain が infrastructure に依存してはならない
- domain が Huma / HTTP / DB ライブラリを import してはならない
- application が具体的なDB実装に依存してはならない

---

# レイヤ責務

## domain

純粋なビジネスルールのみを持つ。

### 含むもの

- Entity
- Value Object
- Domain Service
- Domain Error

### 禁止

- DBアクセス
- HTTP依存
- 外部ライブラリ依存（可能な限り）

---

## application（usecase）

ユースケースを実装する層。

### 含むもの

- Usecase struct
- Repository interface（port）
- トランザクション境界（ユースケース単位）
- ドメインルールの組み合わせ

### 禁止

- HTTPレスポンス生成
- JSONタグ付きDTO
- DB具体実装

---

## infrastructure

外部I/Oの実装。

### 含むもの

- Repository実装（PostgreSQL）
- DB接続（pgx）
- sqlc 生成コードの利用
- 設定管理
- Logger
- マイグレーション適用

### 役割

application層のRepository interfaceを実装する。

---

## interface（delivery）

HumaによるHTTP層。

### 含むもの

- Huma handler
- Request/Response DTO
- 入力バリデーション
- エラー変換（domain/application → HTTP）

### 原則

- domain entity を直接返さない
- DTO変換を明示する

---

# DB/Repository 実装方針（pgx + sqlc）

## 採用方針

- PostgreSQL のドライバ/接続管理は **pgx** を使用する
- クエリは SQL を明示的に書き、コード生成は **sqlc** に任せる
- 生成コードは **infrastructure 層でのみ利用**し、domain/application に漏らさない

## SQL配置ルール

- `internal/infrastructure/db/query/*.sql` にクエリを置く（例）
- `internal/infrastructure/db/migrations` にマイグレーションを置く（例）
- `sqlc.yaml`（または `sqlc.json`）をルートに置く
- `make sqlc` で生成できるようにする

## Repository 実装ルール

- application層に `XxxRepository` の interface（port）を定義する
- infrastructure層で `XxxRepository` を実装する
- sqlc の生成型（Params/Row/Model）は **infrastructure 内に閉じ込める**
  - application/domain に返す型は domain entity または application用の構造体に変換する
- SQLのエラー（unique違反、not found 等）は infrastructure で吸収し、
  domain/application のエラーに変換して返す

---

# トランザクション方針（pgx）

- 原則 usecase 単位でトランザクションを管理する
- application層では Tx の開始/コミット/ロールバックを「抽象（port）」として扱えるようにする
- infrastructure層で pgx の Tx（`pgx.Tx`）を利用して実装する

---

# エラーハンドリング規約

## domain error

domain内で定義する。

例:

- ErrSubjectNotFound
- ErrInvalidStudyDuration

---

## application error

domain error をラップして返す。

---

## HTTP error変換規則

| domain/application error | HTTP |
|--------------------------|------|
| NotFound                 | 404  |
| Validation               | 400  |
| Conflict                 | 409  |
| その他                   | 500  |

Humaのエラーレスポンス形式に統一する。

---

# DTO規約

- Request/Responseは interface 層に定義
- JSONタグは interface 層のみに置く
- domain構造体にJSONタグを付けない
- domainとDTOは明示的に変換関数を作る

---

# Huma利用方針

- ルーティングは Huma を使用
- OpenAPIはHumaに自動生成させる
- /openapi.json または /docs を公開
- バリデーションはHumaの型定義を活用

---

# テスト方針

## domain

- 純粋な単体テストを書く
- DB不要

## usecase

- repositoryをモックしてテスト

## infrastructure

- 必要に応じて統合テスト（PostgreSQLを使う場合はテスト用DBを用意）

---

# 新機能追加ルール

新しいAPIを追加する場合:

1. domainに必要なEntity/ValueObjectを追加
2. applicationにUsecase追加
3. repository interface定義（port）
4. SQL追加（query/*.sql）→ `make sqlc` で生成
5. infrastructureにRepository実装（sqlc生成コード利用、型変換）
6. interface(Huma)でDTOとhandler作成
7. エラー変換追加
8. テスト追加

この順序を守る。

---

# リファクタリング時の注意

- 依存方向が壊れていないか必ず確認
- domainに外部依存が入り込んでいないか確認
- DTOとdomainが混ざっていないか確認
- sqlc生成型がapplication/domainに漏れていないか確認

---

# コーディング規約

- gofmt必須
- golangci-lintを通す
- 明確な命名（UserRepository など）
- 曖昧なパッケージ名を避ける（util, commonは禁止）

---

# このプロジェクトでのアンチパターン

- handlerから直接DB呼び出し
- domainにjsonタグ
- infrastructureからdomainを書き換える
- usecaseを通さずrepository呼び出し
- sqlc生成型をapplication/domainへ返す（漏らす）

---

# 最重要原則

このプロジェクトは「動くコード」よりも設計を守ることを優先する。
