FROM golang:1.18-alpine3.16 AS deps

RUN apk update && apk add gcc musl-dev gcompat libc-dev linux-headers
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

FROM deps as builder
COPY . .

RUN go build -tags=jsoniter -ldflags "-linkmode external -extldflags -static" -o backend

FROM alpine:3.16
WORKDIR /app
EXPOSE 8080

COPY --from=builder /app/backend /app/backend

RUN chmod +x /app/backend

ENTRYPOINT [ "/app/backend" ]
