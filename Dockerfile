FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bridge-taiga-matrix .

FROM alpine:latest
RUN apk add --no-cache dcron
WORKDIR /app
COPY --from=builder /app/bridge-taiga-matrix .
COPY locales/ locales/
COPY settings.json .

RUN echo "0 9 * * * /app/bridge-taiga-matrix -config /app/settings.json >> /var/log/cron.log 2>&1" > /etc/templates/crontab
RUN crontab /etc/templates/crontab

CMD ["crond", "-f"]
