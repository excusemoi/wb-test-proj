FROM golang:1.13.8

RUN mkdir opt

WORKDIR /opt

COPY ./ /opt

EXPOSE 8080

RUN go mod download && cd cmd && go build

CMD ["cmd/cmd", "config/config.json"]