FROM golang:1.24.0

WORKDIR /build
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 go build -o main ./cmd/main.go
EXPOSE 50053

CMD ["./main"]