FROM golang:latest AS app

ENV GO111MODULE=on

RUN apk add git

WORKDIR /app
COPY . .

RUN go get -d -v
RUN go build -o ./boxyd_poc .

EXPOSE 8080
CMD ["./boxyd_poc"]
