FROM golang:1.23.3 as builder

ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . /app
COPY ./config /app/config
COPY ./static /app/static

RUN go build -o app .

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/config /app/config
COPY --from=builder /app/static /app/static

CMD ["./app"]
