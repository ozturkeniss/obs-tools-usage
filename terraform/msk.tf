# MSK Configuration
resource "aws_msk_configuration" "main" {
  kafka_versions = ["3.5.1"]
  name           = "${local.name_prefix}-kafka-config"

  server_properties = <<PROPERTIES
auto.create.topics.enable=true
default.replication.factor=3
min.insync.replicas=2
num.partitions=3
log.retention.hours=168
log.segment.bytes=1073741824
log.cleanup.policy=delete
compression.type=snappy
PROPERTIES

  tags = local.common_tags
}

# MSK Cluster
resource "aws_msk_cluster" "main" {
  cluster_name           = "${local.name_prefix}-kafka"
  kafka_version          = "3.5.1"
  number_of_broker_nodes = var.msk_number_of_broker_nodes

  # Broker configuration
  broker_node_group_info {
    instance_type   = var.msk_instance_type
    ebs_volume_size = 20
    client_subnets  = module.vpc.private_subnets
    security_groups = [aws_security_group.msk.id]
  }

  # Configuration
  configuration_info {
    arn      = aws_msk_configuration.main.arn
    revision = aws_msk_configuration.main.latest_revision
  }

  # Encryption
  encryption_info {
    encryption_at_rest_kms_key_id = aws_kms_key.msk.arn
    encryption_in_transit {
      client_broker = "TLS"
      in_cluster    = true
    }
  }

  # Logging
  logging_info {
    broker_logs {
      cloudwatch_logs {
        enabled   = true
        log_group = aws_cloudwatch_log_group.kafka.name
      }
      firehose {
        enabled = false
      }
      s3 {
        enabled = false
      }
    }
  }

  # Open monitoring
  open_monitoring {
    prometheus {
      jmx_exporter {
        enabled_in_broker = true
      }
      node_exporter {
        enabled_in_broker = true
      }
    }
  }

  tags = local.common_tags
}

# KMS Key for MSK
resource "aws_kms_key" "msk" {
  description             = "KMS key for MSK cluster ${local.name_prefix}"
  deletion_window_in_days = 7

  tags = local.common_tags
}

resource "aws_kms_alias" "msk" {
  name          = "alias/${local.name_prefix}-msk"
  target_key_id = aws_kms_key.msk.key_id
}

# CloudWatch Log Group for Kafka
resource "aws_cloudwatch_log_group" "kafka" {
  name              = "/aws/msk/${local.name_prefix}"
  retention_in_days = 7

  tags = local.common_tags
}
