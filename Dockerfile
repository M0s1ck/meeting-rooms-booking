# ---- build stage: compile binaries from the codebase ----
FROM golang:1.25.0-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/app ./cmd/app \
 && CGO_ENABLED=0 go build -trimpath -o /out/seed ./cmd/seed

# ---- migrate CLI, pinned version ----
FROM migrate/migrate:v4.18.3 AS migrate

# ---- release image: app + admin tools (migrate, seed) from the same build ----
FROM alpine:3.23 AS release
WORKDIR /app

COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=build /out/app /out/seed ./
COPY internal/infra/postgres/migrations ./migrations
COPY scripts/migrate.sh ./migrate.sh

RUN adduser -D -H app
USER app

EXPOSE 8080
CMD ["./app"]
