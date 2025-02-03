# MAM Backend Deployment Guide

## System Requirements

- Ubuntu 20.04 LTS or higher
- 4 CPU cores minimum
- 8GB RAM minimum
- 100GB storage minimum

## Prerequisites

```bash
# Update system
apt update && apt upgrade -y

# Install required packages
apt install -y postgresql-14 redis-server ffmpeg golang-1.20 nginx
```

## Database Setup

```bash
# Configure PostgreSQL
sudo -u postgres psql

# Create database and user
CREATE DATABASE mam_db;
CREATE USER mam_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE mam_db TO mam_user;

# Enable required extensions
\c mam_db
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
```

## Redis Configuration

```bash
# Edit Redis configuration
vim /etc/redis/redis.conf

# Set password and bind address
requirepass your_redis_password
bind 127.0.0.1

# Restart Redis
systemctl restart redis
```

## Application Deployment

```bash
# Create application directory
mkdir -p /opt/mam
cd /opt/mam

# Copy application files
cp mam-backend .
cp config.yaml .

# Set permissions
chown -R mam:mam /opt/mam
chmod 755 /opt/mam/mam-backend

# Create systemd service
cat > /etc/systemd/system/mam.service << EOF
[Unit]
Description=MAM Backend Service
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=mam
WorkingDirectory=/opt/mam
ExecStart=/opt/mam/mam-backend
Restart=always
Environment=MAM_CONFIG=/opt/mam/config.yaml

[Install]
WantedBy=multi-user.target
EOF

# Start service
systemctl enable mam
systemctl start mam
```

## Environment Configuration

Create `/opt/mam/config.yaml`:
```yaml
database:
  host: localhost
  port: 5432
  name: mam_db
  user: mam_user
  password: your_secure_password

redis:
  host: localhost
  port: 6379
  password: your_redis_password

tams:
  url: http://tams-server:4010
  auth_token: your_tams_token

storage:
  path: /var/lib/mam/media
  temp_path: /var/lib/mam/temp

server:
  port: 8080
  max_upload_size: 1GB
```

## Media Storage Setup

```bash
# Create media directories
mkdir -p /var/lib/mam/{media,temp}
chown -R mam:mam /var/lib/mam
chmod 755 /var/lib/mam
```

## Nginx Reverse Proxy

```bash
# Create Nginx configuration
cat > /etc/nginx/sites-available/mam << EOF
server {
    listen 80;
    server_name mam.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        client_max_body_size 1G;
    }
}
EOF

# Enable site and restart Nginx
ln -s /etc/nginx/sites-available/mam /etc/nginx/sites-enabled/
systemctl restart nginx
```

## Health Check

```bash
# Check service status
systemctl status mam

# View logs
journalctl -u mam -f

# Test API endpoint
curl http://localhost:8080/api/v1/health
```

## Backup Configuration

```bash
# Database backup script
cat > /opt/mam/backup.sh << EOF
#!/bin/bash
BACKUP_DIR=/var/backups/mam
mkdir -p \$BACKUP_DIR
pg_dump mam_db | gzip > \$BACKUP_DIR/mam_db_\$(date +%Y%m%d).sql.gz
EOF

# Add to crontab
echo "0 1 * * * /opt/mam/backup.sh" | crontab -
```

## Security Considerations

1. Configure firewall:
```bash
ufw allow ssh
ufw allow http
ufw allow https
ufw enable
```

2. Set secure file permissions
3. Use strong passwords
4. Keep system updated
5. Monitor logs regularly

## Troubleshooting

Common issues and solutions:

1. Service won't start:
   - Check logs: `journalctl -u mam -f`
   - Verify config file permissions
   - Ensure PostgreSQL is running

2. Database connection issues:
   - Check PostgreSQL status
   - Verify connection settings
   - Ensure firewall allows connection

3. Media processing fails:
   - Check FFmpeg installation
   - Verify storage permissions
   - Monitor disk space
