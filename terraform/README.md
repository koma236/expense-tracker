# Terraform — ExpenseTracker インフラ（EC2 + RDS）

ExpenseTracker の本番インフラ（EC2 アプリサーバ + RDS MySQL）を Terraform で構築する。
**インフラ作成までが Terraform の責務**で、アプリのビルド・配置・DB セットアップは [../docs/deploy.md](../docs/deploy.md) の手順で行う。

## 作成されるもの

- セキュリティグループ 2つ（EC2 用 / RDS 用）
- EC2（Ubuntu 24.04, 既定 `t3.micro`、初回起動で基本ミドルウェア＋swap を導入）＋ Elastic IP
- RDS（MySQL 8.0, 既定 `db.t3.micro`、パブリックアクセス無効、EC2 SG からのみ 3306 許可）
- **S3**（フロント配信用, 完全非公開・CloudFront OAC 経由のみ）
- **CloudFront**（公開URL=`*.cloudfront.net` / HTTPS。`/*`→S3, `/api/*`→EC2）
- 既存 VPC（デフォルト VPC）と既存キーペア `my-aws-key` を利用

### 配信構成（CloudFront + S3）

```
Browser ──HTTPS──▶ CloudFront (*.cloudfront.net)
                    ├─ /*     → S3（フロント静的SPA, 非公開/OAC, キャッシュ有）
                    └─ /api/* → EC2:80(nginx) → Go:8080 → RDS（キャッシュ無）
```

- EC2 の 80 番は CloudFront の origin-facing プレフィックスリストのみ許可（直接アクセス不可）。
- CloudFront オリジンは **Elastic IP の DNS** を使うため、EC2 を停止/起動してもオリジンが安定する。
- SPA ルーティングは CloudFront Function（拡張子なしパス→`/index.html`）で対応。

## 前提

- AWS CLI が認証済み（`aws sts get-caller-identity` が通る）
- Terraform 1.5+
- SSH 用の既存キーペア `my-aws-key` の秘密鍵(.pem)を手元に持っていること

## 使い方

```bash
cd terraform
cp terraform.tfvars.example terraform.tfvars
# terraform.tfvars を編集:
#   ssh_cidr    … 自分のIP/32（curl -s https://checkip.amazonaws.com で確認）
#   db_password … RDS マスターパスワード

terraform init
terraform plan      # 作成内容の確認（リソースは作られない）
terraform apply     # ← ここで EC2/RDS が作成され課金が発生する
```

`apply` 後、出力に CloudFront URL・EC2 の DNS・RDS エンドポイント・S3 バケット名が表示される。

```bash
terraform output cloudfront_url        # 公開URL（HTTPS）
terraform output s3_frontend_bucket    # フロント配置先バケット
```

続けて [../docs/deploy.md](../docs/deploy.md) の「アプリ配置」（EC2 へのバックエンド配置）と、
**フロントの S3 配信**（リポジトリルートで `./deploy/deploy-frontend.sh`）を実施する。
CloudFront の作成・反映には数分〜十数分かかる。

## コスト・都度起動

- EC2 は使うときだけ起動する運用（`aws ec2 start-instances` / `stop-instances`）。停止中は EBS 保存料のみ。
- RDS は常時稼働になりがち。初年度の無料利用枠（`db.t3.micro` 750h/月）でカバーする想定。
- **後片付け（課金を完全に止める）:**
  ```bash
  terraform destroy   # 作成した EC2/RDS/SG をまとめて削除
  ```

## 注意

- `terraform.tfvars`・state ファイル（`terraform.tfstate*`）・`.terraform/` はコミットしない（リポジトリの `.gitignore` で除外済み）。
- state にはDBパスワード等が平文で含まれるため、ローカル state を共有・公開しないこと。
