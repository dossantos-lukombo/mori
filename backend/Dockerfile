# backend/Dockerfile (Go Backend)
FROM golang:1.20-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

EXPOSE 8081
CMD ["go", "run", "."]