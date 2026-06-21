# Terraform — ExpenseTracker インフラ（EC2 + RDS）

ExpenseTracker の本番インフラ（EC2 アプリサーバ + RDS MySQL）を Terraform で構築する。
**インフラ作成までが Terraform の責務**で、アプリのビルド・配置・DB セットアップは [../docs/deploy.md](../docs/deploy.md) の手順で行う。

## 作成されるもの

- セキュリティグループ 2つ（EC2 用 / RDS 用）
- EC2（Ubuntu 24.04, 既定 `t3.micro`、初回起動で基本ミドルウェア＋swap を導入）
- RDS（MySQL 8.0, 既定 `db.t3.micro`、パブリックアクセス無効、EC2 SG からのみ 3306 許可）
- 既存 VPC（デフォルト VPC）と既存キーペア `my-aws-key` を利用

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

`apply` 後、出力に EC2 のパブリック DNS と RDS エンドポイントが表示される。
続けて [../docs/deploy.md](../docs/deploy.md) の「アプリ配置」以降を実施する。

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
