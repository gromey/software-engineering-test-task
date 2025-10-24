FROM golang:1.25-alpine AS builder

WORKDIR /app
ENV CGO_ENABLED=0
COPY . .

RUN go mod download
RUN go build -o /usr/local/bin/app ./cmd/
RUN apk add -U --no-cache ca-certificates

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/local/bin/app /usr/local/bin/app

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/app"]