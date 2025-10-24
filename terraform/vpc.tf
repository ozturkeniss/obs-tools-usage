# VPC Module
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
  version = "~> 5.0"

  name = local.name_prefix
  cidr = var.vpc_cidr

  azs             = var.availability_zones != [] ? var.availability_zones : slice(data.aws_availability_zones.available.names, 0, 3)
  private_subnets = [for i in range(3) : cidrsubnet(var.vpc_cidr, 8, i)]
  public_subnets  = [for i in range(3) : cidrsubnet(var.vpc_cidr, 8, i + 10)]
  intra_subnets   = [for i in range(3) : cidrsubnet(var.vpc_cidr, 8, i + 20)]

  enable_nat_gateway     = var.enable_nat_gateway
  enable_vpn_gateway     = var.enable_vpn_gateway
  enable_dns_hostnames   = true
  enable_dns_support     = true

  # VPC Flow Logs
  enable_flow_log                      = true
  create_flow_log_cloudwatch_iam_role  = true
  create_flow_log_cloudwatch_log_group = true

  # Tags
  tags = local.common_tags

  public_subnet_tags = {
    "kubernetes.io/role/elb" = "1"
    "kubernetes.io/cluster/${local.name_prefix}-eks" = "shared"
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = "1"
    "kubernetes.io/cluster/${local.name_prefix}-eks" = "shared"
  }
}

# Security Groups
resource "aws_security_group" "eks_cluster" {
  name_prefix = "${local.name_prefix}-eks-cluster-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = var.allowed_cidr_blocks
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-eks-cluster-sg"
  })
}

resource "aws_security_group" "eks_nodes" {
  name_prefix = "${local.name_prefix}-eks-nodes-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "Node to node communication"
    from_port   = 0
    to_port     = 65535
    protocol    = "tcp"
    self        = true
  }

  ingress {
    description     = "Cluster API to node"
    from_port       = 0
    to_port         = 65535
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_cluster.id]
  }

  ingress {
    description = "Node to node communication (UDP)"
    from_port   = 0
    to_port     = 65535
    protocol    = "udp"
    self        = true
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-eks-nodes-sg"
  })
}

resource "aws_security_group" "rds" {
  name_prefix = "${local.name_prefix}-rds-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "PostgreSQL from EKS"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_nodes.id]
  }

  ingress {
    description     = "MariaDB from EKS"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_nodes.id]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-rds-sg"
  })
}

resource "aws_security_group" "elasticache" {
  name_prefix = "${local.name_prefix}-elasticache-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "Redis from EKS"
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_nodes.id]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-elasticache-sg"
  })
}

resource "aws_security_group" "msk" {
  name_prefix = "${local.name_prefix}-msk-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description     = "Kafka from EKS"
    from_port       = 9092
    to_port         = 9092
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_nodes.id]
  }

  ingress {
    description     = "Kafka JMX from EKS"
    from_port       = 11001
    to_port         = 11001
    protocol        = "tcp"
    security_groups = [aws_security_group.eks_nodes.id]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name_prefix}-msk-sg"
  })
}
