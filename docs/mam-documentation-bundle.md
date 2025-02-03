# MAM System Documentation Bundle

## Table of Contents
1. [Backend Deployment Guide](#backend-deployment)
2. [Frontend Deployment Guide](#frontend-deployment)
3. [TAMS Integration Guide](#tams-integration)
4. [System Monitoring Guide](#system-monitoring)
5. [Backup & Recovery Guide](#backup-recovery)

<a name="frontend-deployment"></a>
## Frontend Deployment Guide

### System Requirements
- Node.js 18.x or higher
- Nginx
- 2GB RAM minimum
- 20GB storage minimum

### Installation Steps
```bash
# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
apt install -y nodejs

# Setup application
mkdir -p /opt/mam-frontend
cd /opt/mam-frontend
npm install pm2 -g

# Build application
npm install
npm run build

# Configure PM2
pm2 start npm --name "mam-frontend" -- start
pm2 startup
pm2 save

# Configure Nginx
cat > /etc/nginx/sites-available/mam-frontend << EOF
server {
    listen 80;
    server_name mam-ui.yourdomain.com;
    root /opt/mam-frontend/build;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    location /api {
        proxy_pass http://backend-server:8080;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF
```

<a name="tams-integration"></a>
## TAMS Integration Guide

### Configuration
```yaml
tams:
  base_url: http://tams-server:4010
  api_version: v5.1
  auth:
    type: bearer
    token: your_tams_token
  storage:
    type: http_object_store
    bucket: mam-assets
  retry:
    max_attempts: 3
    backoff: 1s
```

### Integration Points
1. Source Creation
2. Flow Management
3. Segment Upload
4. Content Retrieval

### Error Handling
```go
// Example error handling configuration
{
    "retryable_errors": [
        "connection_timeout",
        "rate_limit_exceeded"
    ],
    "fatal_errors": [
        "invalid_credentials",
        "insufficient_permissions"
    ]
}
```

<a name="system-monitoring"></a>
## System Monitoring Guide

### Prometheus Setup
```bash
# Install Prometheus
apt install -y prometheus

# Configure monitoring
cat > /etc/prometheus/prometheus.yml << EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'mam_backend'
    static_configs:
      - targets: ['localhost:8080']
  - job_name: 'node_exporter'
    static_configs:
      - targets: ['localhost:9100']
EOF
```

### Grafana Dashboard
```bash
# Install Grafana
apt install -y grafana

# Start services
systemctl enable prometheus grafana-server
systemctl start prometheus grafana-server
```

### Alert Configuration
```yaml
alerts:
  disk_space:
    threshold: 90%
    notification: email
  memory_usage:
    threshold: 85%
    notification: slack
  api_errors:
    threshold: 50/minute
    notification: pagerduty
```

<a name="backup-recovery"></a>
## Backup & Recovery Guide

### Daily Backup Script
```bash
#!/bin/bash
# /opt/mam/backup/daily-backup.sh

# Variables
BACKUP_DIR="/var/backups/mam"
DATE=$(date +%Y%m%d)
DB_NAME="mam_db"
MEDIA_DIR="/var/lib/mam/media"

# Create backup directory
mkdir -p $BACKUP_DIR/$DATE

# Database backup
pg_dump $DB_NAME | gzip > $BACKUP_DIR/$DATE/database.sql.gz

# Configuration backup
tar czf $BACKUP_DIR/$DATE/config.tar.gz /opt/mam/config.yaml

# Media files backup
rsync -az $MEDIA_DIR $BACKUP_DIR/$DATE/media

# Cleanup old backups (keep last 7 days)
find $BACKUP_DIR -type d -mtime +7 -exec rm -rf {} +
```

### Recovery Procedures

1. Database Recovery
```bash
# Restore database
gunzip < backup.sql.gz | psql mam_db
```

2. Configuration Recovery
```bash
# Restore config
tar xzf config.tar.gz -C /
```

3. Media Files Recovery
```bash
# Restore media files
rsync -az backup/media/ /var/lib/mam/media/
```

### Monitoring Backup Status
```bash
# Add to crontab
0 1 * * * /opt/mam/backup/daily-backup.sh
0 2 * * * /opt/mam/backup/check-backup.sh

# Monitoring script
cat > /opt/mam/backup/check-backup.sh << EOF
#!/bin/bash
if [ ! -f /var/backups/mam/\$(date +%Y%m%d)/database.sql.gz ]; then
    echo "Backup failed" | mail -s "MAM Backup Alert" admin@yourdomain.com
fi
EOF
```

### Disaster Recovery Checklist

1. [ ] Stop MAM services
2. [ ] Restore database
3. [ ] Restore configurations
4. [ ] Restore media files
5. [ ] Verify file permissions
6. [ ] Start services
7. [ ] Verify system functionality
8. [ ] Check TAMS integration
9. [ ] Validate media access

### Emergency Contacts
```yaml
contacts:
  primary:
    name: System Administrator
    email: sysadmin@yourdomain.com
    phone: +1-555-0123

  backup:
    name: DevOps Lead
    email: devops@yourdomain.com
    phone: +1-555-0124

  vendor:
    name: TAMS Support
    email: support@tams.com
    phone: +1-555-0125
```
