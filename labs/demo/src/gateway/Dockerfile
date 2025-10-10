FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api-gateway ./cmd/api

FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /app
COPY --from=build /out/api-gateway /app/api-gateway
EXPOSE 8080
USER 1001
ENTRYPOINT ["/app/api-gateway"]
