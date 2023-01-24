FROM golang:1.19.4-alpine

WORKDIR /app
COPY . .
RUN go build -o main ./main.go
CMD ["/app/main"]
