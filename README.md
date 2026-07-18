# Paas
A mini Platform-as-a-Service in Go. Deploy an app from a directory to a running 
container with one HTTP call.

## Quick start

Requires Go and Docker Desktop running.

    # terminal 1 — the platform
    go run ./cmd/server

    # terminal 2 — create the app, deploy a directory (must contain a Dockerfile)
    curl -X POST -d '{"name":"myapp"}' localhost:8080/apps
    go run ./cmd/platform deploy ./myapp

    # your app is live
    curl myapp.localhost:8080

## What works today
- Image building: tar over HTTP → POST /apps/{name}/deploy → tagged image
- Build-stream inspection: in-stream build failures return real HTTP errors
- Container lifecycle: create/start with Docker-allocated host ports
- State tracking: container ID, port, status in SQLite
- Redeploy: old container stopped + removed, replaced atomically
- App metadata CRUD: create, get
- Reverse proxy: <app>.localhost:8080 routes to the app's container — nobody types a port
- CLI: platform deploy <dir> — tars the directory, deploys it, prints the URL

## In progress
- Deployment to a real server (DigitalOcean)
- Zero-downtime deploys: new container up and healthy before the old one dies — redeploys with no dropped requests

## Stack
- Go · go-chi/chi · SQLite (modernc.org/sqlite) · Docker SDK (moby/moby/client)

## Docs
- Design decisions and the costs they accepted: [`docs/decisions.md`](docs/decisions.md)
- Bugs that taught something: [`docs/bugs.md`](docs/bugs.md)