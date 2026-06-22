output "ec2_instance_id" {
  description = "EC2 インスタンス ID（start/stop に使用）"
  value       = aws_instance.app.id
}

output "ec2_public_dns" {
  description = "EC2 のパブリック DNS（起動毎に変わる。SSH/ブラウザ用）"
  value       = aws_instance.app.public_dns
}

output "ec2_public_ip" {
  description = "EC2 のパブリック IP（起動毎に変わる）"
  value       = aws_instance.app.public_ip
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
  value       = "ssh -i <my-aws-key>.pem ubuntu@${aws_instance.app.public_dns}"
}
