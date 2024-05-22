FROM golang:1.22.2 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /service-b cmd/main.go
CMD [ "/service-b" ]