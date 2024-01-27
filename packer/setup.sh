#!/bin/bash

# 创建你的应用的 systemd 服务文件
cat > /tmp/csye6225.service <<EOF
[Unit]
Description=CSYE6225 Web Application
After=cloud-final.service

[Service]
ExecStart=/opt/csye6225/myapp
WorkingDirectory=/opt/csye6225/
Restart=always
RestartSec=5
User=csye6225
Group=csye6225
Environment=PATH=/usr/bin:/usr/local/bin

[Install]
WantedBy=cloud-init.target
EOF

