# VPC Outputs
output "vpc_id" {
  description = "ID of the VPC"
  value       = module.vpc.vpc_id
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = module.vpc.vpc_cidr_block
}

output "private_subnets" {
  description = "List of IDs of private subnets"
  value       = module.vpc.private_subnets
}

output "public_subnets" {
  description = "List of IDs of public subnets"
  value       = module.vpc.public_subnets
}

# EKS Outputs
output "cluster_id" {
  description = "EKS cluster ID"
  value       = module.eks.cluster_id
}

output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = module.eks.cluster_arn
}

output "cluster_endpoint" {
  description = "Endpoint for EKS control plane"
  value       = module.eks.cluster_endpoint
}

output "cluster_security_group_id" {
  description = "Security group ids attached to the cluster control plane"
  value       = module.eks.cluster_security_group_id
}

output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = module.eks.cluster_certificate_authority_data
}

output "cluster_oidc_issuer_url" {
  description = "The URL on the EKS cluster for the OpenID Connect identity provider"
  value       = module.eks.cluster_oidc_issuer_url
}

output "cluster_primary_security_group_id" {
  description = "The cluster primary security group ID created by the EKS service"
  value       = module.eks.cluster_primary_security_group_id
}

# RDS Outputs
output "rds_postgresql_endpoint" {
  description = "RDS PostgreSQL instance endpoint"
  value       = aws_db_instance.postgresql.endpoint
  sensitive   = true
}

output "rds_postgresql_port" {
  description = "RDS PostgreSQL instance port"
  value       = aws_db_instance.postgresql.port
}

output "rds_mariadb_endpoint" {
  description = "RDS MariaDB instance endpoint"
  value       = aws_db_instance.mariadb.endpoint
  sensitive   = true
}

output "rds_mariadb_port" {
  description = "RDS MariaDB instance port"
  value       = aws_db_instance.mariadb.port
}

# ElastiCache Outputs
output "elasticache_redis_endpoint" {
  description = "ElastiCache Redis cluster endpoint"
  value       = aws_elasticache_replication_group.redis.configuration_endpoint_address
}

output "elasticache_redis_port" {
  description = "ElastiCache Redis cluster port"
  value       = aws_elasticache_replication_group.redis.port
}

# MSK Outputs
output "msk_cluster_arn" {
  description = "MSK cluster ARN"
  value       = aws_msk_cluster.main.arn
}

output "msk_cluster_bootstrap_brokers" {
  description = "MSK cluster bootstrap brokers"
  value       = aws_msk_cluster.main.bootstrap_brokers
  sensitive   = true
}

output "msk_cluster_bootstrap_brokers_tls" {
  description = "MSK cluster bootstrap brokers TLS"
  value       = aws_msk_cluster.main.bootstrap_brokers_tls
  sensitive   = true
}

# Security Group Outputs
output "eks_cluster_security_group_id" {
  description = "EKS cluster security group ID"
  value       = aws_security_group.eks_cluster.id
}

output "eks_nodes_security_group_id" {
  description = "EKS nodes security group ID"
  value       = aws_security_group.eks_nodes.id
}

output "rds_security_group_id" {
  description = "RDS security group ID"
  value       = aws_security_group.rds.id
}

output "elasticache_security_group_id" {
  description = "ElastiCache security group ID"
  value       = aws_security_group.elasticache.id
}

output "msk_security_group_id" {
  description = "MSK security group ID"
  value       = aws_security_group.msk.id
}

# IAM Role Outputs
output "aws_load_balancer_controller_role_arn" {
  description = "AWS Load Balancer Controller IAM role ARN"
  value       = var.enable_aws_load_balancer_controller ? aws_iam_role.aws_load_balancer_controller[0].arn : null
}

output "external_dns_role_arn" {
  description = "External DNS IAM role ARN"
  value       = var.enable_external_dns ? aws_iam_role.external_dns[0].arn : null
}

# Application Configuration
output "app_config" {
  description = "Application configuration for deployment"
  value = {
    namespace = var.app_namespace
    replicas  = var.app_replicas
    domain    = var.domain_name
    certificate_arn = var.certificate_arn
  }
}

# Connection Information
output "connection_info" {
  description = "Connection information for the infrastructure"
  value = {
    kubectl_config = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
    cluster_name   = module.eks.cluster_name
    region         = var.aws_region
  }
}
