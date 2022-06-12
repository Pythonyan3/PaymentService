FROM golang:1.18

# additional tools for db migration and postgres waiting script
RUN curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | bash
RUN apt-get update
RUN apt-get install -y migrate
RUN apt-get -y install postgresql-client

COPY . /go/src/app/

WORKDIR /go/src/app/

# build go app
RUN go build -o ./cmd/payment/main ./cmd/payment/main.go

ENTRYPOINT [ "./cmd/payment/main" ]