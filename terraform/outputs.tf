output "ec2_instance_id" {
  description = "EC2 インスタンス ID（start/stop に使用）"
  value       = aws_instance.app.id
}

output "ec2_public_dns" {
  description = "EC2 のパブリック DNS（Elastic IP のため停止/起動でも不変。SSH 用）"
  value       = aws_eip.app.public_dns
}

output "ec2_public_ip" {
  description = "EC2 の Elastic IP（停止/起動でも不変）"
  value       = aws_eip.app.public_ip
}

output "cloudfront_url" {
  description = "公開URL（CloudFront, HTTPS）"
  value       = "https://${aws_cloudfront_distribution.main.domain_name}"
}

output "cloudfront_distribution_id" {
  description = "CloudFront ディストリビューション ID（キャッシュ invalidation 用）"
  value       = aws_cloudfront_distribution.main.id
}

output "s3_frontend_bucket" {
  description = "フロント配信用 S3 バケット名（aws s3 sync 先）"
  value       = aws_s3_bucket.frontend.id
}

output "rds_endpoint" {
  description = "RDS エンドポイント（api.env の DB_HOST に設定）"
  value       = aws_db_instance.main.address
}

output "rds_port" {
  description = "RDS ポート"
  value       = aws_db_instance.main.port
}

output "ssh_command" {
  description = "EC2 への SSH コマンド例"
  value       = "ssh -i <my-aws-key>.pem ubuntu@${aws_eip.app.public_dns}"
}
