FROM golang:1.20.8-alpine3.18

# Run the air command in the directory where our code will live
WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy