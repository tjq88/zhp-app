FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/zhp-app ./cmd/app

FROM alpine:3.20

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /out/zhp-app /app/zhp-app

EXPOSE 8080

USER app

ENTRYPOINT ["/app/zhp-app"]
