FROM golang:1.24-alpine
WORKDIR /app
# Install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Install TZdata
RUN apk add --no-cache tzdata
# Copy the source code
COPY . .
RUN go build -o main .
# Expose the port the app runs on
CMD ["./main"]
