# CloudFront ディストリビューション。
#   /*     → S3（フロント静的SPA, キャッシュ有）
#   /api/* → EC2（API, キャッシュ無・全ヘッダ転送）
# 公開URLは CloudFront の *.cloudfront.net（HTTPS）。

# S3 オリジン用 OAC
resource "aws_cloudfront_origin_access_control" "frontend" {
  name                              = "${var.project}-frontend-oac"
  description                       = "OAC for ${var.project} frontend bucket"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

# SPA ルーティング用 CloudFront Function。
# 拡張子を持たないパス（/dashboard, /transactions/123 等）は SPA シェル /index.html に書き換える。
# 静的アセット（/_nuxt/*.js 等、拡張子あり）はそのまま通す。
# ※ default ビヘイビア(S3)にのみ関連付けるため /api/* には影響しない。
resource "aws_cloudfront_function" "spa_rewrite" {
  name    = "${var.project}-spa-rewrite"
  runtime = "cloudfront-js-2.0"
  comment = "Rewrite extensionless paths to /index.html for SPA"
  publish = true
  code    = <<-EOT
    function handler(event) {
      var request = event.request;
      var uri = request.uri;
      if (!uri.includes('.')) {
        request.uri = '/index.html';
      }
      return request;
    }
  EOT
}

resource "aws_cloudfront_distribution" "main" {
  enabled             = true
  comment             = "${var.project} distribution"
  default_root_object = "index.html"
  price_class         = "PriceClass_200"

  # --- オリジン: S3（フロント）---
  origin {
    origin_id                = "s3-frontend"
    domain_name              = aws_s3_bucket.frontend.bucket_regional_domain_name
    origin_access_control_id = aws_cloudfront_origin_access_control.frontend.id
  }

  # --- オリジン: EC2（API, HTTP）---
  origin {
    origin_id   = "ec2-api"
    domain_name = aws_eip.app.public_dns

    custom_origin_config {
      http_port              = 80
      https_port             = 443
      origin_protocol_policy = "http-only"
      origin_ssl_protocols   = ["TLSv1.2"]
    }
  }

  # --- デフォルト: フロント(S3) ---
  default_cache_behavior {
    target_origin_id       = "s3-frontend"
    viewer_protocol_policy = "redirect-to-https"
    allowed_methods        = ["GET", "HEAD"]
    cached_methods         = ["GET", "HEAD"]
    cache_policy_id        = data.aws_cloudfront_cache_policy.optimized.id
    compress               = true

    function_association {
      event_type   = "viewer-request"
      function_arn = aws_cloudfront_function.spa_rewrite.arn
    }
  }

  # --- /api/*: API(EC2) ---
  ordered_cache_behavior {
    path_pattern             = "/api/*"
    target_origin_id         = "ec2-api"
    viewer_protocol_policy   = "redirect-to-https"
    allowed_methods          = ["GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"]
    cached_methods           = ["GET", "HEAD"]
    cache_policy_id          = data.aws_cloudfront_cache_policy.disabled.id
    origin_request_policy_id = data.aws_cloudfront_origin_request_policy.all_viewer_except_host.id
    compress                 = true
  }

  restrictions {
    geo_restriction {
      restriction_type = "none"
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }

  tags = { Name = "${var.project}-cf" }
}
