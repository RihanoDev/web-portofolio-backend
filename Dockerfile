########################
# Build stage
########################
FROM golang:1.22-alpine AS builder
WORKDIR /src
RUN apk add --no-cache git build-base
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/web-porto-backend ./

########################
# Run stage
########################
FROM alpine:3.20
WORKDIR /app
RUN adduser -D -H app && apk add --no-cache ca-certificates tzdata && update-ca-certificates
COPY --from=builder /out/web-porto-backend /app/web-porto-backend
# Expect config.json to be mounted into /app/config.json
ENV GIN_MODE=release
EXPOSE 8080
USER app
ENTRYPOINT ["/app/web-porto-backend"]
