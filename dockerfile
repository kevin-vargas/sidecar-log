# syntax=docker/dockerfile:1

FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN apk --no-cache add ca-certificates

RUN go get -d -v

RUN CGO_ENABLED=0 go build -o /bin/app

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /bin/app /app

CMD ["/app"]
