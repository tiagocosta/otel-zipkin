FROM golang:1.22

WORKDIR /go/src

COPY . .

RUN go mod tidy

ENTRYPOINT ["tail", "-f", "/dev/null"]