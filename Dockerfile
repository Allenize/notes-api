FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /notes-api ./cmd/notes-api

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /notes-api /app/notes-api

EXPOSE 9001
ENTRYPOINT ["/app/notes-api"]
