#!/usr/bin/env bash
set -euo pipefail

echo "[1/5] Updating apt package index"
sudo apt-get update

echo "[2/5] Installing prerequisites"
sudo apt-get install -y ca-certificates curl gnupg nginx

if ! command -v docker >/dev/null 2>&1; then
  echo "[3/5] Installing Docker Engine"
  sudo install -m 0755 -d /etc/apt/keyrings
  curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  sudo chmod a+r /etc/apt/keyrings/docker.gpg
  . /etc/os-release
  echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    ${VERSION_CODENAME} stable" | \
    sudo tee /etc/apt/sources.list.d/docker.list >/dev/null
  sudo apt-get update
  sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
fi

echo "[4/5] Enabling Docker for current user"
sudo usermod -aG docker "$USER" || true

echo "[5/5] Setup complete"
echo "Re-open your Ubuntu WSL shell, then deploy with:"
echo "  cd /mnt/c/Link-in-bio/deploy"
echo "  cp .env.prod.example .env.prod"
echo "  docker compose --env-file .env.prod -f docker-compose.prod.yml up -d --build"

