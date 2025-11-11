FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./ 
COPY go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app

FROM debian:bullseye-slim 
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/main /bin/main
ENTRYPOINT ["/bin/main"]