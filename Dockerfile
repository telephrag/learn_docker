
# creating a build stage (environment) where we gonna do stuff using golang
FROM golang:1.18.4 as build-env

# setting workdirectory where operations shall be performed
WORKDIR /go/src/learn_docker
# copying everything in Dockerfile's directory to the said workdir
ADD . /go/src/learn_docker

# inside the said workdir pulling dependencies and building executable
RUN go get -d -v ./...
RUN go build -o /go/bin/learn_docker .

# it's a rootless container so we have no need for root access
USER 1000

# importing distroless linux to container as operating system
FROM gcr.io/distroless/base

# from build environment copying executable into container's root directory
COPY --from=build-env /go/bin/learn_docker /

# launching the said executable
CMD ["./learn_docker"]
