FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/data-syncer ./cmd/sync

FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /app
COPY --from=build /out/data-syncer /app/data-syncer
USER 1001
ENTRYPOINT ["/app/data-syncer"]
