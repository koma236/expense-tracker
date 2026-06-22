# ---- 既存リソースの参照（データソース）-------------------------------------

# デフォルト VPC とそのサブネットを利用する（学習用のため新規 VPC は作らない）
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# Ubuntu 24.04 LTS (amd64) の最新 AMI（Canonical 公式）
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd-gp3/ubuntu-noble-24.04-amd64-server-*"]
  }

  filter {
    name   = "state"
    values = ["available"]
  }
}

# ---- セキュリティグループ ---------------------------------------------------

# EC2 用: SSH(自分のIPのみ) と HTTP(全許可)
resource "aws_security_group" "ec2" {
  name        = "${var.project}-ec2-sg"
  description = "ExpenseTracker EC2 (app)"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [var.ssh_cidr]
  }

  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.project}-ec2-sg" }
}

# RDS 用: 3306 を EC2 の SG からのみ許可
resource "aws_security_group" "rds" {
  name        = "${var.project}-rds-sg"
  description = "ExpenseTracker RDS (MySQL)"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "MySQL from EC2"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.ec2.id]
  }

  egress {
    description = "all outbound"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = { Name = "${var.project}-rds-sg" }
}

# ---- EC2（アプリ）-----------------------------------------------------------

resource "aws_instance" "app" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = var.instance_type
  subnet_id                   = data.aws_subnets.default.ids[0]
  vpc_security_group_ids      = [aws_security_group.ec2.id]
  key_name                    = var.key_name
  associate_public_ip_address = true

  # 初回起動時に基本ミドルウェアと swap を導入（秘密情報は含めない）。
  # アプリのビルド・配置・DB セットアップは docs/deploy.md の手順で実施する。
  user_data = file("${path.module}/user_data.sh")

  root_block_device {
    volume_type = "gp3"
    volume_size = var.ec2_volume_size
  }

  tags = { Name = "${var.project}-ec2" }
}

# ---- RDS（MySQL）-----------------------------------------------------------

resource "aws_db_subnet_group" "main" {
  name       = "${var.project}-db-subnet-group"
  subnet_ids = data.aws_subnets.default.ids

  tags = { Name = "${var.project}-db-subnet-group" }
}

resource "aws_db_instance" "main" {
  identifier     = "${var.project}-db"
  engine         = "mysql"
  engine_version = "8.0"
  instance_class = var.db_instance_class

  allocated_storage     = var.db_allocated_storage
  max_allocated_storage = 0 # ストレージ自動拡張を無効化（想定外課金の防止）
  storage_type          = "gp3"

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]
  publicly_accessible    = false
  multi_az               = false

  # 学習用のため運用簡素化
  skip_final_snapshot = true
  apply_immediately   = true

  tags = { Name = "${var.project}-db" }
}
