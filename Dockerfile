FROM golang:1.18.4 as build-env

WORKDIR /go/src/learn_docker
ADD . /go/src/learn_docker

RUN go get -d -v ./...
RUN go build -o /go/bin/learn_docker .

USER 1000

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/learn_docker /

CMD ["./learn_docker"]

# I have barely any idea what's happening here...
