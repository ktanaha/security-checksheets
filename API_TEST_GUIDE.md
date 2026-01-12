# API動作確認ガイド

## アクセス情報

- **Go Backend API**: http://localhost:8080
- **Python Excel Service**: http://localhost:8000
- **React Frontend**: http://localhost:3000
- **PostgreSQL**: localhost:5432

---

## 1. 案件管理API

### 1.1 案件作成（POST /api/projects）

```bash
curl -X POST http://localhost:8080/api/projects \
  -H 'Content-Type: application/json' \
  -d '{
    "customer_name": "株式会社サンプル",
    "description": "セキュリティチェックシート案件",
    "owner": "山田太郎"
  }' | jq .
```

**期待されるレスポンス例**:
```json
{
  "id": 1,
  "customer_name": "株式会社サンプル",
  "description": "セキュリティチェックシート案件",
  "owner": "山田太郎",
  "status": "active",
  "created_at": "2026-01-11T06:00:00Z",
  "updated_at": "2026-01-11T06:00:00Z"
}
```

### 1.2 案件一覧取得（GET /api/projects）

```bash
curl http://localhost:8080/api/projects | jq .
```

**期待されるレスポンス例**:
```json
[
  {
    "id": 1,
    "customer_name": "株式会社サンプル",
    "description": "セキュリティチェックシート案件",
    "owner": "山田太郎",
    "status": "active",
    "created_at": "2026-01-11T06:00:00Z",
    "updated_at": "2026-01-11T06:00:00Z"
  }
]
```

### 1.3 案件詳細取得（GET /api/projects/:id）

```bash
# 案件ID=1の詳細を取得
curl http://localhost:8080/api/projects/1 | jq .
```

### 1.4 案件更新（PUT /api/projects/:id）

```bash
curl -X PUT http://localhost:8080/api/projects/1 \
  -H 'Content-Type: application/json' \
  -d '{
    "customer_name": "株式会社サンプル（更新）",
    "description": "更新されたセキュリティチェックシート案件",
    "owner": "佐藤花子",
    "status": "in_progress"
  }' | jq .
```

### 1.5 案件削除（DELETE /api/projects/:id）

```bash
# 注意: 削除すると復元できません
curl -X DELETE http://localhost:8080/api/projects/1 -w "\nHTTP Status: %{http_code}\n"
```

**期待されるレスポンス**: HTTP 204 (No Content)

---

## 2. ファイル管理API

### 2.1 ファイルアップロード（POST /api/projects/:id/files）

まず、テスト用ファイルを作成します：

```bash
# テストファイルを作成
echo "テスト用のExcelファイルです" > /tmp/sample.txt
```

ファイルをアップロード：

```bash
# 案件ID=1にファイルをアップロード
curl -X POST http://localhost:8080/api/projects/1/files \
  -F "file=@/tmp/sample.txt" \
  -F "uploaded_by=山田太郎" | jq .
```

**期待されるレスポンス例**:
```json
{
  "id": 1,
  "project_id": 1,
  "file_name": "sample.txt",
  "file_path": "/app/uploads/project_1/20260111_060000_sample.txt",
  "file_size": 42,
  "uploaded_by": "山田太郎",
  "uploaded_at": "2026-01-11T06:00:00Z"
}
```

### 2.2 案件のファイル一覧取得（GET /api/projects/:id/files）

```bash
# 案件ID=1のファイル一覧を取得
curl http://localhost:8080/api/projects/1/files | jq .
```

**期待されるレスポンス例**:
```json
[
  {
    "id": 1,
    "project_id": 1,
    "file_name": "sample.txt",
    "file_path": "/app/uploads/project_1/20260111_060000_sample.txt",
    "file_size": 42,
    "uploaded_by": "山田太郎",
    "uploaded_at": "2026-01-11T06:00:00Z"
  }
]
```

### 2.3 ファイル詳細取得（GET /api/files/:id）

```bash
# ファイルID=1の詳細を取得
curl http://localhost:8080/api/files/1 | jq .
```

### 2.4 ファイル削除（DELETE /api/files/:id）

```bash
# 注意: 削除すると物理ファイルも削除されます
curl -X DELETE http://localhost:8080/api/files/1 -w "\nHTTP Status: %{http_code}\n"
```

**期待されるレスポンス**: HTTP 204 (No Content)

### 2.5 物理ファイルの確認（Docker内）

```bash
# アップロードされた物理ファイルを確認
docker-compose exec backend ls -la /app/uploads/

# 特定案件のファイルを確認
docker-compose exec backend ls -la /app/uploads/project_1/

# ファイル内容を確認
docker-compose exec backend cat /app/uploads/project_1/20260111_060000_sample.txt
```

---

## 3. 統合動作確認シナリオ

### シナリオ: 案件作成からファイルアップロードまで

```bash
# Step 1: 案件を作成
PROJECT_ID=$(curl -s -X POST http://localhost:8080/api/projects \
  -H 'Content-Type: application/json' \
  -d '{
    "customer_name": "テスト株式会社",
    "description": "動作確認用案件",
    "owner": "テスト太郎"
  }' | jq -r '.id')

echo "作成された案件ID: $PROJECT_ID"

# Step 2: テストファイルを作成
echo "これはテスト用のファイルです" > /tmp/test_file.txt

# Step 3: ファイルをアップロード
FILE_ID=$(curl -s -X POST http://localhost:8080/api/projects/$PROJECT_ID/files \
  -F "file=@/tmp/test_file.txt" \
  -F "uploaded_by=テスト太郎" | jq -r '.id')

echo "アップロードされたファイルID: $FILE_ID"

# Step 4: 案件詳細を確認
echo -e "\n=== 案件詳細 ==="
curl -s http://localhost:8080/api/projects/$PROJECT_ID | jq .

# Step 5: ファイル一覧を確認
echo -e "\n=== ファイル一覧 ==="
curl -s http://localhost:8080/api/projects/$PROJECT_ID/files | jq .

# Step 6: 物理ファイルを確認
echo -e "\n=== 物理ファイル内容 ==="
docker-compose exec backend cat /app/uploads/project_$PROJECT_ID/*.txt
```

---

## 4. エラーケースの確認

### 4.1 バリデーションエラー（顧客名なし）

```bash
curl -X POST http://localhost:8080/api/projects \
  -H 'Content-Type: application/json' \
  -d '{
    "description": "顧客名がない案件",
    "owner": "山田太郎"
  }' | jq .
```

**期待されるレスポンス**:
```json
{
  "error": "無効なリクエストです"
}
```

### 4.2 存在しない案件の取得

```bash
curl http://localhost:8080/api/projects/99999 | jq .
```

**期待されるレスポンス**:
```json
{
  "error": "案件が見つかりません"
}
```

### 4.3 ファイルなしでアップロード

```bash
curl -X POST http://localhost:8080/api/projects/1/files \
  -F "uploaded_by=山田太郎" | jq .
```

**期待されるレスポンス**:
```json
{
  "error": "ファイルが指定されていません"
}
```

---

## 5. データベース直接確認

### PostgreSQLに接続

```bash
docker-compose exec postgres psql -U admin -d security_checksheets
```

### SQLクエリ例

```sql
-- 案件一覧
SELECT * FROM projects;

-- ファイル一覧
SELECT * FROM uploaded_files;

-- 案件とファイルの結合
SELECT
  p.id as project_id,
  p.customer_name,
  f.id as file_id,
  f.file_name,
  f.file_size
FROM projects p
LEFT JOIN uploaded_files f ON p.id = f.project_id;

-- 終了
\q
```

---

## 6. ブラウザでの確認

### フロントエンド

ブラウザで以下にアクセス:
```
http://localhost:3000
```

現在はヘルスチェック表示のみですが、Phase 3以降で案件管理UIを実装予定です。

---

## 7. ログの確認

### バックエンドログ

```bash
docker-compose logs -f backend
```

### データベースログ

```bash
docker-compose logs -f postgres
```

### すべてのログ

```bash
docker-compose logs -f
```

---

## トラブルシューティング

### サービスが起動しない場合

```bash
# すべてのサービスを再起動
docker-compose restart

# 特定のサービスのみ再起動
docker-compose restart backend
```

### データベースをリセットしたい場合

```bash
# 注意: すべてのデータが削除されます
docker-compose down -v
docker-compose up -d
```

### コンテナを完全に再ビルドしたい場合

```bash
docker-compose down
docker-compose build --no-cache
docker-compose up -d
```

---

---

## 8. Excel解析API（Phase 3）

### 8.1 Excel解析（POST /excel/parse）

```bash
curl -X POST http://localhost:8000/excel/parse \
  -H 'Content-Type: application/json' \
  -d '{"file_path": "/app/uploads/test_security_check.xlsx"}' | jq .
```

**期待されるレスポンス例**:
```json
{
  "file_name": "test_security_check.xlsx",
  "sheets": [
    {
      "name": "セキュリティチェック",
      "index": 0,
      "row_count": 4,
      "column_count": 3
    }
  ],
  "total_sheets": 1
}
```

### 8.2 シートプレビュー（POST /excel/preview）

```bash
curl -X POST http://localhost:8000/excel/preview \
  -H 'Content-Type: application/json' \
  -d '{
    "file_path": "/app/uploads/test_security_check.xlsx",
    "sheet_name": "セキュリティチェック"
  }' | jq .
```

**期待されるレスポンス例**:
```json
{
  "sheet_name": "セキュリティチェック",
  "cells": [
    {
      "row": 1,
      "column": 1,
      "value": "部門",
      "formatted_value": "部門",
      "is_merged": false
    },
    {
      "row": 1,
      "column": 2,
      "value": "質問",
      "formatted_value": "質問",
      "is_merged": false
    }
  ],
  "total_rows": 4,
  "total_columns": 3
}
```

### 8.3 範囲指定プレビュー

```bash
curl -X POST http://localhost:8000/excel/preview \
  -H 'Content-Type: application/json' \
  -d '{
    "file_path": "/app/uploads/test_security_check.xlsx",
    "sheet_name": "セキュリティチェック",
    "start_row": 1,
    "end_row": 2,
    "start_column": 1,
    "end_column": 2
  }' | jq .
```

---

## 9. Excelプレビュー機能（Phase 4）

### 9.1 ブラウザでのアクセス

```
http://localhost:3000/excel-preview
```

### 9.2 操作手順

1. **ファイルパス入力**: `/app/uploads/test_security_check.xlsx`
2. **解析ボタンをクリック**
3. **シート選択**: 自動的に最初のシートが選択される
4. **プレビュー表示を確認**:
   - 行ヘッダー（1, 2, 3, ...）
   - 列ヘッダー（A, B, C, ...）
   - セルデータ
5. **セルクリック**: セルが青くハイライトされる
6. **範囲選択**: Shiftキーを押しながら別のセルをクリック
7. **選択範囲の表示**: 青い背景で表示される

### 9.3 テスト用Excelファイル作成

```bash
# Pythonコンテナに入る
docker-compose exec excel-service python

# Pythonインタラクティブシェルで実行
from openpyxl import Workbook
wb = Workbook()
ws = wb.active
ws.title = "セキュリティチェック"

# ヘッダー
ws['A1'] = "部門"
ws['B1'] = "質問"
ws['C1'] = "回答"

# データ
ws['A2'] = "情報システム"
ws['B2'] = "アクセス権限の管理方法は?"
ws['C2'] = "役割ベースのアクセス制御を実施"

ws['A3'] = "セキュリティ"
ws['B3'] = "脆弱性スキャンの頻度は?"
ws['C3'] = "月次で実施"

ws['A4'] = "法務"
ws['B4'] = "個人情報の保管期間は?"
ws['C4'] = "法令に基づき5年間保管"

wb.save('/app/uploads/test_security_check.xlsx')
exit()
```

---

## 次のステップ

### Phase 4完了により、以下が動作確認可能です：
- ✅ 案件CRUD操作
- ✅ ファイルアップロード/削除
- ✅ PostgreSQLデータ永続化
- ✅ Clean Architecture実装
- ✅ Excel解析API（Python）
- ✅ Excelプレビュー画面（React）
- ✅ セル選択・範囲選択UI
- ✅ セル結合対応

### Phase 5では以下を実装予定：
- ファイルアップロードUIの統合
- 案件管理画面からのExcelプレビュー
- 範囲選択からのナレッジ抽出機能
- 抽出データのDB保存
