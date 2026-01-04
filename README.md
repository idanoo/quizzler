# Quizzler

Simple self-hosted quiz/flashcard app with a Go API and Svelte frontend.
Supports importing quizlet exports.

## Requirements
- Docker

## Quick Start with docker-compose
```bash
wget -O docker-compose.yml https://raw.githubusercontent.com/idanoo/quizzler/refs/heads/main/docker-compose.prod.yml
wget -O .env https://raw.githubusercontent.com/idanoo/quizzler/refs/heads/main/.env.example
# Generate JWT_SECRET with: openssl rand -base64 32
docker compose up -d
```

## Development
```bash
git clone https://github.com/idanoo/quizzler.git
cd quizzler
make
```

## Access
Frontend: http://127.0.0.1:5172
API: http://127.0.0.1:8132
First user created will be an admin

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make` | Start containers in background |
| `make build` | Start containers with rebuild |
| `make down` | Stop containers |
| `make logs` | Follow container logs |
| `make clean` | Stop containers, remove volumes |
