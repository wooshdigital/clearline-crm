# ClearLine CRM

A lightweight self-hosted CRM built with an Astro frontend and a Go REST API backend. Manage contacts, track deals through pipelines, view activity timelines, and enforce role-based access control — all in a single Docker Compose stack.

## Features

- **Contact Management** — Create, update, and search contacts with full details
- **Deal Pipelines** — Kanban-style pipeline with customizable stages
- **Activity Timeline** — Log calls, emails, meetings, and notes per contact or deal
- **Role-Based Access Control** — Admin, Manager, and Rep roles with JWT auth
- **Self-Hosted** — Runs entirely on your own infrastructure with PostgreSQL

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | Astro + TypeScript |
| Backend | Go (net/http) |
| Database | PostgreSQL 15 |
| Auth | JWT (HS256) |
| Deployment | Docker + Docker Compose |

## Quick Start

bash
git clone https://github.com/yourorg/clearline-crm
cd clearline-crm
docker compose up --build


Then open [http://localhost:4321](http://localhost:4321).

**Demo credentials:**

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@clearline.local | admin123 |
| Manager | manager@clearline.local | manager123 |
| Rep | rep@clearline.local | rep123 |

## Development Setup

### Backend (Go)

bash
cd backend
go mod download
go run ./cmd/server


### Frontend (Astro)

bash
cd frontend
npm install
npm run dev


### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://crm:crm@localhost:5432/clearline` | PostgreSQL DSN |
| `JWT_SECRET` | — | Secret key for signing JWTs |
| `PORT` | `8080` | Backend HTTP port |
| `FRONTEND_URL` | `http://localhost:4321` | CORS origin |

## API Reference


POST   /api/auth/login
GET    /api/contacts
POST   /api/contacts
GET    /api/contacts/:id
PUT    /api/contacts/:id
DELETE /api/contacts/:id
GET    /api/deals
POST   /api/deals
GET    /api/deals/:id
PUT    /api/deals/:id
GET    /api/activities
POST   /api/activities
GET    /api/users
GET    /api/health


## License

MIT
