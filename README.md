# Paas
A mini Platform-as-a-Service in Go. Deploy an app from a directory to a running 
container with one HTTP call.

## What works today
- image building: tar over HTTP → POST /apps/{name}/deploy → tagged image
- build-stream inspection: in-stream build failures returns real HTTP errors
- container lifecycle: create/start with Docker-allocated host ports
- state tracking: container ID, port, status in SQLite
- redeploy: old container stopped + removed, replaced atomically
- app metadata CRUD: create, get

## In progress
 
- reverse proxy: <app>.localhost routing
- CLI: platform deploy <dir>

## Stack
- Go · go-chi/chi · SQLite (modernc.org/sqlite) · Docker SDK (moby/moby/client)

## Docs
- Design decisions and the costs they accepted: [`docs/decisions.md`](docs/decisions.md)
- Bugs that taught something: [`docs/bugs.md`](docs/bugs.md)