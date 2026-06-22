#!/usr/bin/env bash
# EC2 初回起動時のブートストラップ（cloud-init / 1回のみ実行）。
# 基本ミドルウェアと swap を導入するだけ。秘密情報やアプリのビルドは含めない。
# アプリ配置・DB セットアップは docs/deploy.md の手順で実施する。
set -euxo pipefail

export DEBIAN_FRONTEND=noninteractive

# 2GB swap（t2/t3.micro で nuxt generate の OOM を防ぐ）
if [ ! -f /swapfile ]; then
  fallocate -l 2G /swapfile
  chmod 600 /swapfile
  mkswap /swapfile
  swapon /swapfile
  echo '/swapfile none swap sw 0 0' >> /etc/fstab
fi

apt-get update -y
apt-get install -y nginx git mysql-client rsync

# Node.js 20
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs

# Go 1.22
curl -fsSL https://go.dev/dl/go1.22.12.linux-amd64.tar.gz -o /tmp/go.tar.gz
rm -rf /usr/local/go
tar -C /usr/local -xzf /tmp/go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' > /etc/profile.d/go.sh

# golang-migrate
curl -fsSL https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz -o /tmp/migrate.tar.gz
tar -C /tmp -xzf /tmp/migrate.tar.gz migrate
install -m 0755 /tmp/migrate /usr/local/bin/migrate

# アプリ配置用ディレクトリ
mkdir -p /opt/expense-tracker
chown ubuntu:ubuntu /opt/expense-tracker
