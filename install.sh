#!/bin/bash

# Define the installation directory and binary name
directory="/usr/local/iptablelb4"
binary_name="iptablelb4"
service_name="iptablelb4.service"

# Start installation
echo -e "\033[1;34mStarting installation of $binary_name...\033[0m"
echo "Creating installation directory at $directory..."

mkdir -p $directory
if [ $? -ne 0 ]; then
    echo -e "\033[1;31mError: Failed to create directory. Please check permissions.\033[0m"
    exit 1
fi

cd $directory || { echo -e "\033[1;31mError: Failed to navigate to the directory.\033[0m"; exit 1; }

# Downloading binary
echo -e "\033[1;32mDownloading $binary_name...\033[0m"
wget https://github.com/Sithukyaw666/iptablelb4/raw/refs/heads/main/iptablelb4 -q --show-progress
if [ "$?" -eq "0" ]; then
    echo -e "\033[1;32m - Downloaded successfully!\033[0m"
else
    echo -e "\033[1;31m - Download failed! Please check your internet connection and try again.\033[0m"
    exit 1
fi

echo -e "\033[1;34mSetting executable permissions...\033[0m"
chmod 755 "${directory}/${binary_name}"

# Creating systemd service file
echo -e "\033[1;34mGenerating systemd service file...\033[0m"
cat > /lib/systemd/system/$service_name << EOT
[Unit]
Description=$binary_name

[Service]
Type=simple
Restart=always
RestartSec=5s
ExecStart=$directory/$binary_name

[Install]
WantedBy=multi-user.target
EOT

if [ $? -eq 0 ]; then
    echo -e "\033[1;32m - Systemd service file generated successfully!\033[0m"
else
    echo -e "\033[1;31m - Error: Failed to generate service file.\033[0m"
    exit 1
fi

# Starting the service
echo -e "\033[1;34mStarting the $binary_name service...\033[0m"
systemctl daemon-reload &>/dev/null
systemctl enable $service_name &>/dev/null
systemctl start $service_name &>/dev/null

if [ $? -eq 0 ]; then
    echo -e "\033[1;32m - $binary_name service started successfully!\033[0m"
else
    echo -e "\033[1;31m - Error: Failed to start the service.\033[0m"
    exit 1
fi

echo -e "\033[1;34mInstallation completed successfully!\033[0m"

