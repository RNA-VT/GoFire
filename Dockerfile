FROM golang:1.13.5 AS gofirebuilder

RUN mkdir -p /go/src/GoFire/

ADD /src /go/src/GoFire/

WORKDIR /go/src/GoFire/

RUN go build

FROM alpine:latest

RUN export GOFIRE_MASTER_HOST=`/sbin/ip route|awk '/default/ { print $3 }'` && export GOFIRE_MASTER=true

COPY --from=gofirebuilder /go/src/GoFire/firecontroller .

RUN chmod +x firecontroller

CMD ["firecontroller"]