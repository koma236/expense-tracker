# CloudFront/S3 で利用する既存（AWS マネージド）リソースの参照

data "aws_caller_identity" "current" {}

# CloudFront オリジン向けのマネージドプレフィックスリスト
# （EC2 の 80 番を CloudFront からのアクセスのみに絞るために使用）
data "aws_ec2_managed_prefix_list" "cloudfront" {
  name = "com.amazonaws.global.cloudfront.origin-facing"
}

# CloudFront マネージドポリシー
data "aws_cloudfront_cache_policy" "optimized" {
  name = "Managed-CachingOptimized"
}

data "aws_cloudfront_cache_policy" "disabled" {
  name = "Managed-CachingDisabled"
}

data "aws_cloudfront_origin_request_policy" "all_viewer_except_host" {
  name = "Managed-AllViewerExceptHostHeader"
}
