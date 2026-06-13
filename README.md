# LinkHub MVP

Scaffold monorepo cho du an Link-in-bio / URL Shortener SaaS theo file `link_in_bio_url_shortener_plan.md`.

## Local defaults

- Backend: `http://localhost:8081`
- Frontend: `http://localhost:3002`

## Chay backend local

```bash
cd backend
go run ./cmd/api
```

## Chay frontend local

```bash
cd frontend
npm install
npm run dev
```

Tao file `frontend/.env.local` neu can:

```env
NEXT_PUBLIC_API_URL=http://localhost:8081
```

## Docker Compose

```bash
docker compose up --build
```

Frontend se map ra `3002`, backend map ra `8081`.

## Ghi chu

- Backend CORS da dong bo theo `FRONTEND_URL=http://localhost:3002`
- Trong `development`, backend se echo lai `Origin` hop le de tranh loi CORS khi dev browser
- Rate limit dang bat cho auth, create short link, public profile va redirect
