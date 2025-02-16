FROM golang:1.23.0


WORKDIR /server

COPY . .

RUN go mod download

RUN go build -o /build ./cmd/main.go 
    # && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]