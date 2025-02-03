provider "aws" {
  region = var.aws_region
}

# VPC with single public subnet
resource "aws_vpc" "main" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "mam-vpc"
  }
}

resource "aws_internet_gateway" "main" {
  vpc_id = aws_vpc.main.id
}

resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.main.id
  cidr_block              = "10.0.1.0/24"
  map_public_ip_on_launch = true
  availability_zone       = "${var.aws_region}a"
}

resource "aws_route_table" "public" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.main.id
  }
}

resource "aws_route_table_association" "public" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# Security Groups
resource "aws_security_group" "db_sg" {
  name   = "mam-db-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}

resource "aws_security_group" "redis_sg" {
  name   = "mam-redis-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = ["10.0.0.0/16"]
  }
}

resource "aws_security_group" "backend_sg" {
  name   = "mam-backend-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "frontend_sg" {
  name   = "mam-frontend-sg"
  vpc_id = aws_vpc.main.id

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# RDS Instance
resource "aws_db_subnet_group" "mam" {
  name       = "mam-db-subnet"
  subnet_ids = [aws_subnet.public.id]
}

resource "aws_db_instance" "mam_db" {
  identifier        = "mam-postgres"
  engine            = "postgres"
  engine_version    = "14.7"
  instance_class    = "db.t3.micro"
  allocated_storage = 20

  db_name  = var.db_name
  username = var.db_username
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.db_sg.id]
  db_subnet_group_name   = aws_db_subnet_group.mam.id
  publicly_accessible    = true
  skip_final_snapshot    = true
}

# Redis Instance
resource "aws_elasticache_subnet_group" "mam" {
  name       = "mam-cache-subnet"
  subnet_ids = [aws_subnet.public.id]
}

resource "aws_elasticache_cluster" "mam_redis" {
  cluster_id           = "mam-redis"
  engine              = "redis"
  node_type           = "cache.t3.micro"
  num_cache_nodes     = 1
  parameter_group_name = "default.redis6.x"
  port                = 6379
  security_group_ids  = [aws_security_group.redis_sg.id]
  subnet_group_name   = aws_elasticache_subnet_group.mam.name
}

# EC2 Instances
resource "aws_instance" "mam_backend" {
  ami           = var.ami_id
  instance_type = "t3.micro"

  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.backend_sg.id]

  user_data = templatefile("${path.module}/scripts/backend-init.sh", {
    db_host     = aws_db_instance.mam_db.endpoint
    db_name     = var.db_name
    db_user     = var.db_username
    db_password = var.db_password
    redis_host  = aws_elasticache_cluster.mam_redis.cache_nodes[0].address
  })

  tags = {
    Name = "mam-backend"
  }
}

resource "aws_instance" "mam_frontend" {
  ami           = var.ami_id
  instance_type = "t3.micro"

  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.frontend_sg.id]

  user_data = templatefile("${path.module}/scripts/frontend-init.sh", {
    backend_url = "http://${aws_instance.mam_backend.public_ip}:8080"
  })

  tags = {
    Name = "mam-frontend"
  }
}

# Output values
output "frontend_url" {
  value = "http://${aws_instance.mam_frontend.public_ip}"
}

output "backend_url" {
  value = "http://${aws_instance.mam_backend.public_ip}:8080"
}