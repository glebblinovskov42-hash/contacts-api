FROM golang:1.26.5-bookworm

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/exe ./cmd/api

CMD ["/app/exe"]