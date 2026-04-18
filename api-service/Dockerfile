FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api .

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /api /api
EXPOSE 8081
CMD ["/api"]