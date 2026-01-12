# セキュリティチェックシート抽出・ナレッジ化Webアプリ

SaaS事業者向けのセキュリティチェックシート抽出・ナレッジ管理システムです。
顧客から受領したバラバラなExcelフォーマットのチェックシートから、質問・回答・担当部門を効率的に抽出し、永続的なナレッジとして蓄積できます。

## 技術スタック

### フロントエンド
- **フレームワーク**: React 18+ with TypeScript
- **ビルドツール**: Vite
- **UIライブラリ**: TailwindCSS
- **テスト**: Jest + React Testing Library

### バックエンド
- **メインAPI**: Go 1.21+ with Gin
  - 案件管理、ナレッジDB操作、検索、CSVエクスポート
- **Excel処理**: Python 3.11+ with FastAPI
  - Excelファイル解析、シート情報取得、Q/A抽出

### データベース
- **RDBMS**: PostgreSQL 15
- **ORM**: Go側はGinで直接SQL、Python側は必要に応じてSQLAlchemy

### インフラ
- **コンテナ**: Docker + Docker Compose
- **ストレージ**: ローカルファイルシステム（MVP）

## プロジェクト構造

```
security-checksheets/
├── backend/                # Go APIサーバー
│   ├── cmd/api/            # エントリーポイント
│   ├── internal/           # アプリケーションコード
│   │   ├── domain/         # ドメインモデル
│   │   ├── usecase/        # ビジネスロジック
│   │   ├── interface/      # ハンドラー・ミドルウェア
│   │   └── infrastructure/ # DB・外部API実装
│   └── tests/              # テスト
├── excel-service/          # Python Excel処理サービス
│   ├── app/
│   │   ├── api/            # FastAPIエンドポイント
│   │   ├── services/       # ビジネスロジック
│   │   └── utils/          # ユーティリティ
│   └── tests/              # pytest
├── frontend/               # Reactアプリ
│   └── src/
│       ├── components/     # UIコンポーネント（Atomic Design）
│       ├── pages/          # ページコンポーネント
│       ├── hooks/          # カスタムフック
│       ├── services/       # API通信
│       └── types/          # TypeScript型定義
├── database/               # データベース初期化スクリプト
├── uploads/                # アップロードファイル保存先
└── docker-compose.yml      # 開発環境定義
```

## セットアップ

### 必要な環境
- Docker Desktop
- Git

### 環境構築手順

1. リポジトリのクローン
```bash
git clone <repository-url>
cd security_checksheets
```

2. 開発環境の起動
```bash
docker-compose up
```

初回起動時は、各コンテナのビルドに数分かかります。

3. アクセス確認
- **フロントエンド**: http://localhost:3000
- **Go API**: http://localhost:8080
- **Python Excel Service**: http://localhost:8000
- **PostgreSQL**: localhost:5432

### 各サービスの確認

#### ヘルスチェック
```bash
# Go API
curl http://localhost:8080/health

# Python Excel Service
curl http://localhost:8000/health
```

## 開発ガイド

### TDD（テスト駆動開発）必須
このプロジェクトでは、すべての開発にTDDを適用します。

**Red-Green-Refactorサイクル**:
1. **Red**: 失敗するテストを先に書く
2. **Green**: テストを通す最小限の実装
3. **Refactor**: 動作を変えずに改善

### バックエンド開発（Go）

```bash
# コンテナ内でテスト実行
docker-compose exec backend go test -v ./...

# 特定のパッケージのテスト
docker-compose exec backend go test -v ./internal/domain

# カバレッジ確認
docker-compose exec backend go test -cover ./...
```

### Excel Service開発（Python）

```bash
# コンテナ内でテスト実行
docker-compose exec excel-service pytest -v

# カバレッジ付きテスト
docker-compose exec excel-service pytest --cov=app tests/
```

### フロントエンド開発（React）

```bash
# コンテナ内でテスト実行
docker-compose exec frontend npm test

# ウォッチモード
docker-compose exec frontend npm run test:watch
```

## MVP実装フェーズ

### Phase 1: 環境構築 ✅（完了）
- Docker Compose設定
- PostgreSQLスキーマ
- Go/Python/React基本構造
- ヘルスチェックエンドポイント

### Phase 2: 案件・ファイル管理（Week 2）
- 案件CRUD機能
- ファイルアップロード機能

### Phase 3: Excel解析基盤（Week 3）
- Python API実装（openpyxl使用）
- Go→Python API連携

### Phase 4: プレビュー機能（Week 4）
- シートプレビュー生成
- Excel風テーブル表示UI

### Phase 5: Q/A抽出機能（Week 5-6）
- 抽出ロジック実装
- 分割/統合編集UI
- ナレッジ保存

### Phase 6: ナレッジ検索（Week 7）
- 全文検索実装
- フィルタ機能

### Phase 7: エクスポート機能（Week 8）
- CSV生成・ダウンロード

## データベース

### スキーマ
- `projects`: 案件管理
- `uploaded_files`: アップロードファイル管理
- `departments`: 部門マスタ
- `knowledge_items`: ナレッジQ/A
- `extraction_sessions`: 抽出セッション（将来用）

### マイグレーション
初期スキーマは`database/init.sql`で自動実行されます。

## API仕様

### Go API（メイン）
- `GET /api/projects` - 案件一覧
- `POST /api/projects` - 案件作成
- `GET /api/knowledge` - ナレッジ検索
- `POST /api/knowledge` - ナレッジ保存
- `GET /api/departments` - 部門一覧

### Python API（Excel処理）
- `POST /excel/parse` - Excelファイル解析
- `POST /excel/sheets/{sheet_name}/preview` - シートプレビュー
- `POST /excel/extract` - Q/A抽出

詳細は`docs/api/`のOpenAPI仕様書を参照してください（今後追加予定）。

## トラブルシューティング

### コンテナが起動しない
```bash
# ログ確認
docker-compose logs

# 特定サービスのログ
docker-compose logs backend
docker-compose logs excel-service
docker-compose logs frontend

# コンテナ再ビルド
docker-compose down
docker-compose build --no-cache
docker-compose up
```

### ポートが既に使用されている
他のアプリケーションが3000, 8080, 8000, 5432ポートを使用している場合は、
`docker-compose.yml`でポート番号を変更してください。

### データベースの初期化
```bash
# データベースボリュームを削除して再作成
docker-compose down -v
docker-compose up
```

## ライセンス

（ライセンス情報を記載）

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずIssueで議論してください。

## サポート

質問やバグ報告は、GitHubのIssuesセクションにお願いします。
