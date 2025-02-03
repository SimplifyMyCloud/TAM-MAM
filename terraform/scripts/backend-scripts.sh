# scripts/backend-init.sh
#!/bin/bash
set -e

apt-get update && apt-get upgrade -y
apt-get install -y postgresql-client redis-tools ffmpeg golang-1.20 nginx

# Setup application
mkdir -p /opt/mam
mv /tmp/mam-backend /opt/mam/
chmod +x /opt/mam/mam-backend

cat > /opt/mam/config.yaml << EOF
database:
  host: ${db_host}
  name: ${db_name}
  user: ${db_user}
  password: ${db_password}
redis:
  host: ${redis_host}
  port: 6379
server:
  port: 8080
EOF

# Create service
cat > /etc/systemd/system/mam.service << EOF
[Unit]
Description=MAM Backend Service
After=network.target

[Service]
Type=simple
ExecStart=/opt/mam/mam-backend
Restart=always
Environment=MAM_CONFIG=/opt/mam/config.yaml

[Install]
WantedBy=multi-user.target
EOF

systemctl enable mam
systemctl start mam
