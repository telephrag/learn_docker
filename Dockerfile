FROM golang:1.18.4

WORKDIR learn_docker/
COPY . .
RUN go get -d -v ./...
RUN go build -o learn_docker .
CMD ["./learn_docker"]
