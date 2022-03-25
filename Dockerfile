# syntax=docker/dockerfile:1.2

FROM golang:1.17
RUN mkdir /app
WORKDIR /app
ADD . /app

RUN apt-get install git

COPY go.mod ./
COPY go.sum ./

COPY *.go ./

ARG PAT
ENV PAT=$PAT

RUN git config --global url."https://${PAT}:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN git config --global http.sslVerify false

RUN go mod download

RUN go build -o bookstore .
CMD [ "/app/bookstore" ]