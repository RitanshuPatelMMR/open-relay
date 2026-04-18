FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ingestion .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /ingestion /ingestion
EXPOSE 8080
CMD ["/ingestion"]