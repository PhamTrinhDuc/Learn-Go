FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /demo-run main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /demo-run .
EXPOSE 8080
CMD ["./demo-run"]

# # Build image
# docker build -t demo-go .

# # Xem image size — nhỏ hơn Python rất nhiều
# docker images demo-go

# # Chạy container
# docker run -p 8080:8080 demo-go