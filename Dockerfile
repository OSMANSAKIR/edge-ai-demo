FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod .
COPY main.go . 
RUN go build -o edge-ai-demo .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/edge-ai-demo .
EXPOSE 8080
CMD ["./edge-ai-demo"]