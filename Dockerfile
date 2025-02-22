FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /openweather-influxdb-writer

FROM alpine:3.19

RUN apk add --no-cache supercronic
COPY crontab /crontab

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

WORKDIR /

COPY --from=builder /openweather-influxdb-writer /openweather-influxdb-writer

RUN chown nonroot:nonroot /openweather-influxdb-writer

USER nonroot:nonroot

CMD ["supercronic", "/crontab"]
