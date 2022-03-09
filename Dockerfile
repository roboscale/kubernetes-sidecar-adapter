FROM golang:1.17 as adapter
SHELL [ "/bin/bash", "-c" ]
RUN mkdir -p /go/src/github.com/roboscale/sidecar-adapter
COPY . /go/src/github.com/roboscale/sidecar-adapter/
WORKDIR /go/src/github.com/roboscale/sidecar-adapter/
RUN go build -o /adapter main.go

FROM ubuntu:focal as traverser
SHELL [ "/bin/bash", "-c" ]
RUN apt update && apt install -y g++
COPY ./helpers/folder_traverser.cpp .
RUN g++ -std=c++17 -o traverser folder_traverser.cpp

FROM ubuntu:focal
RUN apt update && apt install -y debootstrap
COPY --from=adapter ./adapter .
COPY --from=traverser ./traverser .