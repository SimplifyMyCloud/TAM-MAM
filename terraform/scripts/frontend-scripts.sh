
# scripts/frontend-init.sh
#!/bin/bash
set -e

apt-get update && apt-get upgrade -y
apt-get install -y nginx

# Setup frontend
mv /tmp/build /var/www/html

# Configure nginx
cat > /etc/nginx/sites-available/default << EOF
server {
    listen 80 default_server;
    root /var/www/html;
    index index.html;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    location /api {
        proxy_pass ${backend_url};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
    }
}
EOF

systemctl enable nginx
systemctl restart nginx