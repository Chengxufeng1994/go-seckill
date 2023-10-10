# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /usr/app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o ./user-svc

EXPOSE 9009 9008

CMD ["/usr/app/user-svc"]