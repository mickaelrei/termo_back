# Termo Backend

Backend for Termo web application

## Configuration

The file `config_example.json` specifies the template for the necessary configurations for running the server.
Create a copy of this file named `config.json` with your configurations to be read when the server is executed.

The project uses PASETO authentication. The key pair, stored in the `auth.private_key` and `auth.public_key` fields can
be generated in the following way:

```go
package main

import (
	"fmt"
	"aidanwoods.dev/go-paseto"
)

func main() {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()

	fmt.Println("Private:", privateKey.ExportHex())
	fmt.Println("Public:", publicKey.ExportHex())
}
```

## Deploy

The script `deploy_example.ps1` is a PowerShell script containing a template for building and deploying the server on
the cloud. Change the parameters at the start of the script to your needs.

The script builds the project, sends it to the cloud via SSH and attempts to call a `restart_backend.sh` script on the
cloud instance. Here's the template for the restart script:

```bash
#!/bin/bash
set -e

echo "Restarting Go backend via systemd..."

# Stop the systemd service if running
sudo systemctl stop backend || echo "Service not running yet"

# Copy uploaded build files from /tmp to target location
mkdir -p ~/backend
cp /tmp/backend/* ~/backend/

# Ensure it's executable
chmod +x ~/backend/main

# Start the service again
sudo systemctl start backend

# Show status
sudo systemctl status backend --no-pager

echo "Backend restarted using systemd."
```

The restart script expects a working systemd daemon for the backend. Here's a template for `backend.service` file:

```
[Unit]
Description=Go Backend Service
After=network.target

[Service]
Type=simple
User=ubuntu
ExecStart=/home/ubuntu/backend/main
WorkingDirectory=/home/ubuntu/backend
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```