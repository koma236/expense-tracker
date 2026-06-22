variable "aws_region" {
  description = "デプロイ先リージョン"
  type        = string
  default     = "ap-northeast-1"
}

variable "project" {
  description = "リソースの Project タグ・名前のプレフィックス"
  type        = string
  default     = "expense-tracker"
}

variable "key_name" {
  description = "EC2 に割り当てる既存のキーペア名（SSH 用）"
  type        = string
  default     = "my-aws-key"
}

variable "ssh_cidr" {
  description = "SSH(22) を許可する CIDR。自分のグローバルIP/32 を推奨（例: 60.67.41.124/32）"
  type        = string
}

variable "instance_type" {
  description = "EC2 インスタンスタイプ（無料枠は t2.micro / t3.micro）"
  type        = string
  default     = "t3.micro"
}

variable "ec2_volume_size" {
  description = "EC2 ルート EBS サイズ（GB）"
  type        = number
  default     = 16
}

variable "db_instance_class" {
  description = "RDS インスタンスクラス（無料枠は db.t3.micro / db.t4g.micro）"
  type        = string
  default     = "db.t3.micro"
}

variable "db_allocated_storage" {
  description = "RDS ストレージ（GB）。無料枠は 20GB まで"
  type        = number
  default     = 20
}

variable "db_name" {
  description = "初期作成するデータベース名"
  type        = string
  default     = "expense_tracker"
}

variable "db_username" {
  description = "RDS マスターユーザー名"
  type        = string
  default     = "app"
}

variable "db_password" {
  description = "RDS マスターパスワード（terraform.tfvars で設定。コミット禁止）"
  type        = string
  sensitive   = true
}
