FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /project cmd/server/main.go
RUN chmod +x /project

FROM scratch

COPY --from=builder /project /project

ENTRYPOINT ["/project"]
