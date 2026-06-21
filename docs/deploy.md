# デプロイ手順書 — ExpenseTracker（EC2 + RDS / 都度起動）

本書は ExpenseTracker を **AWS EC2（アプリ）+ RDS（MySQL）** にデプロイする手順をまとめたもの。
学習用課題のため **コスト最小化**を重視し、**EC2 は使うときだけ起動（都度起動）**する運用とする。

- 構成図・設計の背景は [アーキテクチャ設計書](architecture.md) を参照。
- 本番では **Docker は使わない**（ローカル開発専用）。Go はバイナリ＋systemd、Nuxt は静的SPA を nginx で配信する。

---

## 0. 全体像

```
Browser ──HTTP──▶ nginx (EC2:80)
                   ├─ /         → Nuxt 静的SPA（/var/www/expense-tracker）
                   └─ /api/...  → reverse proxy → Go API (127.0.0.1:8080, systemd)
                                                       │ 3306
                                                       ▼
                                                  RDS MySQL 8（プライベート, EC2のSGからのみ許可）
```

- フロントは `ssr: false` の静的SPA。`NUXT_PUBLIC_API_BASE=/api`（相対パス）で配信するため、**EC2 のパブリック IP が起動毎に変わってもビルドし直し不要**。
- フロントと API は同一オリジン（nginx）になるので **CORS 設定は不要**。

リポジトリ内の関連ファイル:

| ファイル | 用途 |
|----------|------|
| [deploy/expense-tracker-api.service](../deploy/expense-tracker-api.service) | Go API の systemd ユニット |
| [deploy/nginx/expense-tracker.conf](../deploy/nginx/expense-tracker.conf) | nginx 設定（静的配信 + `/api` プロキシ） |
| [deploy/api.env.example](../deploy/api.env.example) | 本番 API 環境変数のサンプル |
| [deploy/deploy.sh](../deploy/deploy.sh) | 2回目以降の再デプロイ補助スクリプト |

---

## 1. コスト方針（重要）

- **EC2**: 停止すればコンピュート課金はゼロ（EBS の保存料のみ。8GB なら月 $1 未満）。**使うときだけ起動**する。
- **RDS**: EC2 とは**独立して課金される**。ただしアカウント作成から **12か月間は無料利用枠**があり、`db.t3.micro`（または `db.t4g.micro`）1台・ストレージ20GBまでなら実質無料。
  - RDS も停止できるが、AWS の仕様で **停止後7日経つと自動的に再起動**する。長期間使わない／課金を完全に止めたい場合は、最終的に **スナップショットを取得してから削除**する（§9 参照）。
- **パブリック IPv4**: 起動中の EC2 にも料金が発生する（少額）。Elastic IP は付けない方針（停止中も課金され、IP 変化は相対 `/api` で吸収できるため不要）。

> まとめ: ふだんは **EC2 を停止**。デモ・採点のときだけ起動する。RDS は初年度なら付けっぱなしでも無料枠内。

---

## 2. RDS（MySQL）を作成する

AWS コンソール → RDS → 「データベースの作成」。

| 項目 | 値 |
|------|----|
| エンジン | MySQL 8.0 |
| テンプレート | 無料利用枠 |
| インスタンスクラス | `db.t3.micro`（無料枠対象） |
| ストレージ | 20 GB（gp3）。**ストレージ自動拡張はオフ**（想定外課金の防止） |
| マルチAZ | 無効（単一インスタンス） |
| パブリックアクセス | **なし**（EC2 経由でのみ接続） |
| VPC | EC2 と同じ VPC |
| DB 認証情報 | マスターユーザー名・強力なパスワードを控える |
| 初期データベース名 | 空でよい（後で `expense_tracker` を作る） |

- セキュリティグループは新規作成（例: `expense-tracker-rds-sg`）。インバウンドは後述の EC2 用 SG からの **3306 のみ**許可する（§3 で設定）。
- 文字コードは接続文字列側で `utf8mb4` を指定するため、パラメータグループの変更は必須ではない。タイムゾーンを揃えたい場合はパラメータグループで `time_zone = Asia/Tokyo` を設定してもよい。
- 作成後、**エンドポイント**（`xxxx.ap-northeast-1.rds.amazonaws.com`）を控える。

---

## 3. EC2 を作成する

AWS コンソール → EC2 → 「インスタンスを起動」。

| 項目 | 値 |
|------|----|
| OS | Ubuntu Server 24.04 LTS |
| インスタンスタイプ | `t3.micro`（または無料枠の `t2.micro`） |
| キーペア | SSH 用に作成・ダウンロード |
| ストレージ | 8〜16 GB（gp3） |
| VPC | RDS と同じ VPC |

**セキュリティグループ（例: `expense-tracker-ec2-sg`）**

| タイプ | ポート | ソース | 用途 |
|--------|--------|--------|------|
| SSH | 22 | マイIP | 管理用 |
| HTTP | 80 | 0.0.0.0/0 | アプリ公開 |

**RDS 側 SG にルールを追加**: `expense-tracker-rds-sg` のインバウンドに「MySQL/Aurora (3306)、ソース = `expense-tracker-ec2-sg`」を追加する（EC2 からのみ DB に到達可能にする）。

SSH 接続:

```bash
ssh -i <キーペア>.pem ubuntu@<EC2のパブリックDNS>
```

---

## 4. ミドルウェアを導入する（EC2 上）

```bash
sudo apt update && sudo apt upgrade -y

# nginx / git / MySQL クライアント（RDS 操作用）
sudo apt install -y nginx git mysql-client rsync

# Node.js 20（NodeSource）
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Go 1.22（tarball）
curl -fsSL https://go.dev/dl/go1.22.12.linux-amd64.tar.gz -o /tmp/go.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf /tmp/go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh
source /etc/profile.d/go.sh
go version

# golang-migrate（マイグレーション CLI）
curl -fsSL https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz
tar -C /tmp -xzf /tmp/migrate.tar.gz migrate
sudo install -m 0755 /tmp/migrate /usr/local/bin/migrate
migrate -version
```

> **t2.micro（メモリ1GB）の場合**: `nuxt generate` がメモリ不足で失敗することがある。先に swap を 2GB 追加しておく。
>
> ```bash
> sudo fallocate -l 2G /swapfile && sudo chmod 600 /swapfile
> sudo mkswap /swapfile && sudo swapon /swapfile
> echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab   # 再起動後も有効化
> ```

---

## 5. ソースを取得する

```bash
sudo mkdir -p /opt/expense-tracker && sudo chown ubuntu:ubuntu /opt/expense-tracker
git clone https://github.com/koma236/expense-tracker.git /opt/expense-tracker
cd /opt/expense-tracker
```

---

## 6. データベースを準備する（RDS へ）

EC2 から RDS に接続して、アプリ用 DB・ユーザーを作成する。

```bash
# マスターユーザーで接続（パスワードはプロンプト入力）
mysql -h <RDSエンドポイント> -u <マスターユーザー> -p
```

```sql
CREATE DATABASE IF NOT EXISTS expense_tracker
  CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE USER 'app'@'%' IDENTIFIED BY '<アプリ用の強力なパスワード>';
GRANT ALL PRIVILEGES ON expense_tracker.* TO 'app'@'%';
FLUSH PRIVILEGES;
EXIT;
```

マイグレーションとシード投入:

```bash
cd /opt/expense-tracker/backend

# マイグレーション
migrate -path db/migrations \
  -database "mysql://app:<アプリ用パスワード>@tcp(<RDSエンドポイント>:3306)/expense_tracker" up

# 初期カテゴリ（日本語の文字化け防止に utf8mb4 を指定）
mysql -h <RDSエンドポイント> -u app -p \
  --default-character-set=utf8mb4 expense_tracker < db/seeds/0001_categories.sql
```

---

## 7. バックエンド（Go API）をデプロイする

```bash
# ビルド
cd /opt/expense-tracker/backend
go build -o /tmp/expense-tracker-api ./cmd/server
sudo install -m 0755 /tmp/expense-tracker-api /usr/local/bin/expense-tracker-api

# 環境変数ファイル（RDS 接続情報）
sudo mkdir -p /etc/expense-tracker
sudo cp /opt/expense-tracker/deploy/api.env.example /etc/expense-tracker/api.env
sudo nano /etc/expense-tracker/api.env     # DB_HOST=RDSエンドポイント, DB_PASSWORD などを設定
sudo chmod 640 /etc/expense-tracker/api.env

# systemd ユニットを配置・有効化（EC2 起動時に自動起動）
sudo cp /opt/expense-tracker/deploy/expense-tracker-api.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable --now expense-tracker-api
sudo systemctl status expense-tracker-api

# 動作確認（DB 接続まで含めた疎通）
curl -fsS http://127.0.0.1:8080/api/health
```

---

## 8. フロントエンド（Nuxt 静的SPA）と nginx

```bash
# 静的SPA をビルド（/api 相対で API を叩く）
cd /opt/expense-tracker/frontend
export NUXT_PUBLIC_API_BASE=/api
npm ci
npm run generate

# 配信先へ配置
sudo mkdir -p /var/www/expense-tracker
sudo rsync -a --delete .output/public/ /var/www/expense-tracker/

# nginx 設定を有効化
sudo cp /opt/expense-tracker/deploy/nginx/expense-tracker.conf /etc/nginx/sites-available/
sudo ln -sf /etc/nginx/sites-available/expense-tracker.conf /etc/nginx/sites-enabled/expense-tracker.conf
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx
```

ブラウザで `http://<EC2のパブリックDNS>/` を開き、各画面が表示されれば成功。

---

## 9. 都度起動の運用（コスト節約）

### 起動・停止（AWS コンソール）
EC2 → インスタンス選択 → 「インスタンスの状態」→ 開始 / 停止。
systemd（API）と nginx は `enable` 済みなので、**EC2 を起動すれば自動で全サービスが復帰**する。

### 起動・停止（ローカルから AWS CLI）

```bash
# 起動
aws ec2 start-instances --instance-ids <インスタンスID>

# 現在のパブリック DNS を確認（起動毎に変わる）
aws ec2 describe-instances --instance-ids <インスタンスID> \
  --query 'Reservations[0].Instances[0].PublicDnsName' --output text

# 使い終わったら停止（課金を止める）
aws ec2 stop-instances --instance-ids <インスタンスID>
```

> EC2 のパブリック DNS/IP は起動毎に変わるが、フロントは相対パス `/api` を使うため**アプリの再ビルドは不要**。アクセスする URL だけが変わる。

### RDS の扱い
- 初年度（無料枠）なら起動したままでも無料。EC2 だけ都度起動すればよい。
- 課金を完全に止めたい場合: RDS を停止（**7日後に自動再起動する点に注意**）。長期間使わないなら、スナップショットを取得してから RDS を削除し、再開時にスナップショットから復元する。

---

## 10. 再デプロイ（コード更新時）

初回セットアップ後にコードを更新したら、補助スクリプトで一括反映できる。

```bash
cd /opt/expense-tracker
./deploy/deploy.sh
```

`git pull` → バックエンド再ビルド & 再起動 → フロント再ビルド & 配置 → nginx リロード を行う。

---

## 11. トラブルシュート

| 症状 | 確認・対処 |
|------|-----------|
| 画面は出るがデータが出ない | `curl http://<DNS>/api/health` を確認。失敗なら API かDB接続の問題 |
| `/api/health` が失敗 | `journalctl -u expense-tracker-api -f` でログ確認。`api.env` の DB 接続情報を見直す |
| DB につながらない | RDS の SG に EC2 の SG からの 3306 が許可されているか、エンドポイント/認証情報が正しいか |
| 日本語が文字化け | シード投入時に `--default-character-set=utf8mb4` を付けたか確認 |
| `nuxt generate` が落ちる | メモリ不足の可能性。§4 の swap 追加を実施 |
| 起動毎に URL が変わって不便 | 恒久 URL が必要なら Elastic IP を割当（停止中も少額課金される点に留意） |
| nginx が 404/既定ページ | `sites-enabled/default` を削除したか、`nginx -t` が通るか確認 |
