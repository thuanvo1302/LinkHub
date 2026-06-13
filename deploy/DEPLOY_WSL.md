# Deploy on WSL Ubuntu

## Recommended GitHub flow

Use 2 GitHub Actions workflows:

- `CI`
  - runs automatically on every push to `main`
  - builds backend and frontend
- `Deploy`
  - started manually with **Run workflow**
  - runs on a **self-hosted runner** installed in your WSL Ubuntu or VPS
  - executes Docker deploy locally on that machine

This is the best fit for WSL because GitHub-hosted runners cannot deploy directly into a private local WSL instance unless you expose SSH publicly.

## Target architecture

- WSL2 Ubuntu hosts Docker Engine
- `docker compose` runs:
  - `postgres`
  - `redis`
  - `backend`
  - `frontend`
  - `nginx`
- Nginx is the only public entrypoint on port `80`

## Routing note

To avoid conflict between:

- public profile frontend route: `/{username}`
- short-link backend route: `/{code}`

short links are exposed under:

```text
/s/{code}
```

Examples:

```text
http://localhost/s/abc123
http://your-domain.com/s/abc123
```

## Prepare WSL

From Ubuntu WSL:

```bash
cd /mnt/c/Link-in-bio/deploy
chmod +x scripts/setup-wsl-ubuntu.sh
./scripts/setup-wsl-ubuntu.sh
```

## Prepare deploy env on runner machine

```bash
cd /mnt/c/Link-in-bio/deploy
cp .env.prod.example .env.prod
```

Edit `deploy/.env.prod` with real values before first deploy.

## Install GitHub self-hosted runner on WSL or VPS

1. Open GitHub repository
2. Go to `Settings` -> `Actions` -> `Runners`
3. Click `New self-hosted runner`
4. Choose `Linux` and `x64`
5. Run the provided commands inside your Ubuntu WSL
6. When configuring labels, add:

```text
linkhub
```

The deploy workflow in this repo expects these labels:

```text
self-hosted, linux, x64, linkhub
```

## Deploy manually from GitHub

After the runner is online:

1. Push code to `main`
2. Wait for workflow `CI` to finish
3. Open workflow `Deploy`
4. Click `Run workflow`
5. Choose the ref, usually `main`

The deploy job will run:

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml up -d --build
```

## Manual deploy from WSL

You can still deploy directly:

```bash
cd /mnt/c/Link-in-bio
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml up -d --build
```

## Check services

```bash
docker compose --env-file deploy/.env.prod -f deploy/docker-compose.prod.yml ps
curl http://localhost/healthz
```

## Notes

- Current backend still uses in-memory app data for auth/profile/links/short-links.
- PostgreSQL and Redis are provisioned now so the deployment shape is ready for the next migration step.
- If you want HTTPS later, add Certbot or place this stack behind a cloud/VPS reverse proxy.

