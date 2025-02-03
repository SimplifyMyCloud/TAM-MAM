#!/bin/bash
set -e

apt-get update && apt-get upgrade -y
apt-get install -y nginx

# Setup frontend
rm -rf /var/www/html/*
mv /tmp/build/* /var/www/html/

# Configure nginx
cat > /etc/nginx/sites-available/default << 'EOF'
server {
   listen 80 default_server;
   root /var/www/html;
   index index.html;

   location / {
       try_files $uri $uri/ /index.html;
   }

   location /api/ {
       proxy_pass http://backend_ip:8080/;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
   }
}
EOF

# Replace the template variable after creating the file
# sed -i "s|\${backend_url}|${backend_url}|g" /etc/nginx/sites-available/default

# Enable and start nginx
# systemctl enable nginx || true
# systemctl restart nginx