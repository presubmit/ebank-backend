# eBank backend 🚀

Go backend that powers eBank.

## Dependencies

- [go](https://golang.org/doc/install) - primary language
- [docker](https://www.docker.com/) - used to run containerized services

## How to run?

- `yarn db`
  - Runs a PostgresDB (see config in .env)
  - Runs `yarn db:migrate up` to apply outstanding migrations
- `yarn watch`

  - Runs all services on ports :80 (http) and :8079 (grpc)
  - Runs a Traefik proxy that redirects requests to the appropiate service
  - Runs `yarn proto` on every proto change

- `yarn proto` - Generate protocol buffers
- `yarn db:migrate up [N]` - Apply all or N up migrations
- `yarn db:migrate down [N]` - Apply all or N down migrations
- `yarn db:migration NAME` - Create a set of timestamped up/down migrations titled NAME in db/migrations.
