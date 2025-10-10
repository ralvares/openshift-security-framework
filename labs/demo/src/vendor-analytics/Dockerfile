FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 go build -o /out/vendor-analytics ./cmd/app

FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /app
# Running as root intentionally (privileged port) - document risk
USER 0
COPY --from=build /out/vendor-analytics /app/vendor-analytics
EXPOSE 443
ENTRYPOINT ["/app/vendor-analytics"]
