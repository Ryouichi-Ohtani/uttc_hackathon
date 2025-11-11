# デプロイガイド / Deployment Guide

## 前提条件 / Prerequisites

### Google Cloud Platform
- GCPプロジェクトの作成
- Cloud Run, Cloud SQL, Cloud Buildの有効化
- gcloud CLIのインストールと認証

### Vercel
- Vercelアカウントの作成
- Vercel CLIのインストール（オプション）

## バックエンドデプロイ (Cloud Run)

### 1. Cloud SQLの設定

```bash
# Cloud SQLインスタンスの作成
gcloud sql instances create ecomate-db \
  --database-version=POSTGRES_15 \
  --tier=db-f1-micro \
  --region=asia-northeast1

# データベースの作成
gcloud sql databases create ecomate_db \
  --instance=ecomate-db

# ユーザーの作成
gcloud sql users create ecomate \
  --instance=ecomate-db \
  --password=YOUR_SECURE_PASSWORD
```

### 2. Secret Managerの設定

```bash
# データベースパスワードをSecret Managerに保存
echo -n "YOUR_SECURE_PASSWORD" | gcloud secrets create ecomate-db-password \
  --replication-policy="automatic" \
  --data-file=-

# JWT秘密鍵を保存
echo -n "YOUR_JWT_SECRET" | gcloud secrets create ecomate-jwt-secret \
  --replication-policy="automatic" \
  --data-file=-

# Gemini API キーを保存
echo -n "YOUR_GEMINI_API_KEY" | gcloud secrets create ecomate-gemini-key \
  --replication-policy="automatic" \
  --data-file=-
```

### 3. Cloud Buildでのデプロイ

```bash
# Cloud Buildの実行
gcloud builds submit \
  --config=cloudbuild.yaml \
  --substitutions=\
_DATABASE_HOST=<CLOUD_SQL_PRIVATE_IP>,\
_DATABASE_PORT=5432,\
_DATABASE_USER=ecomate,\
_DATABASE_NAME=ecomate_db,\
_JWT_SECRET=<YOUR_JWT_SECRET>,\
_GEMINI_API_KEY=<YOUR_GEMINI_API_KEY>
```

### 4. 手動デプロイ（オプション）

```bash
# Dockerイメージのビルド
docker build -t gcr.io/$PROJECT_ID/ecomate-backend:latest -f backend/Dockerfile .

# イメージのプッシュ
docker push gcr.io/$PROJECT_ID/ecomate-backend:latest

# Cloud Runへデプロイ
gcloud run deploy ecomate-backend \
  --image=gcr.io/$PROJECT_ID/ecomate-backend:latest \
  --region=asia-northeast1 \
  --platform=managed \
  --allow-unauthenticated \
  --set-env-vars=DATABASE_HOST=<HOST>,DATABASE_PORT=5432,... \
  --set-secrets=DATABASE_PASSWORD=ecomate-db-password:latest
```

## フロントエンドデプロイ (Vercel)

### 1. Vercel CLIでのデプロイ

```bash
# Vercel CLIのインストール
npm install -g vercel

# プロジェクトのルートでデプロイ
vercel

# 本番環境へのデプロイ
vercel --prod
```

### 2. Vercel Dashboard でのデプロイ

1. https://vercel.com にアクセス
2. GitHubリポジトリを接続
3. プロジェクト設定:
   - **Framework Preset**: Vite
   - **Build Command**: `cd frontend && npm run build`
   - **Output Directory**: `frontend/dist`
   - **Install Command**: `cd frontend && npm install`

4. 環境変数の設定:
   - `VITE_API_URL`: Cloud RunのバックエンドURL

5. デプロイボタンをクリック

### 3. 環境変数の設定

```bash
# .env.production (フロントエンド)
VITE_API_URL=https://ecomate-backend-xxxxx-an.a.run.app
```

## CI/CDパイプライン

### GitHub Actions (既存)

`.github/workflows/backend-ci.yml` と `.github/workflows/frontend-ci.yml` が設定済み。

- プッシュ時に自動テスト実行
- main/developブランチへのマージでビルド確認

### Cloud Buildとの連携（オプション）

```bash
# GitHub リポジトリとCloud Buildを連携
gcloud builds triggers create github \
  --repo-name=uttc_hackathon \
  --repo-owner=YOUR_GITHUB_USERNAME \
  --branch-pattern="^main$" \
  --build-config=cloudbuild.yaml
```

## データベースマイグレーション

```bash
# Cloud Run経由でマイグレーション実行
gcloud run services update ecomate-backend \
  --command="/root/main,migrate"

# または、Cloud SQLに直接接続
gcloud sql connect ecomate-db --user=ecomate
# psql内でマイグレーション実行
```

## モニタリングとログ

### Cloud Logging

```bash
# バックエンドのログを確認
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=ecomate-backend" --limit=50
```

### Cloud Monitoring

- Cloud Consoleでメトリクスを確認
- アラートポリシーの設定

## トラブルシューティング

### バックエンドが起動しない

1. ログを確認: `gcloud logging read ...`
2. 環境変数が正しく設定されているか確認
3. Cloud SQLへの接続権限を確認

### フロントエンドがバックエンドに接続できない

1. CORS設定を確認（backend/cmd/api/main.go）
2. `VITE_API_URL`が正しいか確認
3. Cloud RunのURLが`--allow-unauthenticated`になっているか確認

## コスト最適化

- Cloud Run: 最小インスタンス数を0に設定（コールドスタート許容時）
- Cloud SQL: 開発環境では`db-f1-micro`を使用
- 本番環境では適切なスケーリング設定

## セキュリティ

- [ ] Secret Managerで機密情報を管理
- [ ] Cloud SQLはプライベートIPで接続
- [ ] Cloud Runは必要に応じて認証を有効化
- [ ] HTTPS通信の強制
- [ ] CORS設定の適切な制限

---

## デモアカウント

開発・デモ用アカウント:

- **管理者**: admin@automate.com / password123
- **販売者**: demo.seller@automate.com / password123
- **購入者**: demo.buyer@automate.com / password123
