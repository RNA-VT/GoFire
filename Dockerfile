# ------- backend build ------- #
FROM golang:1.13.5 AS gofirebuilder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /app/GoFire

COPY src/go.mod .
COPY src/go.sum .

RUN go mod download

COPY ./src .

RUN go build -o ./gofire .

# ------- frontend build ------- #
FROM node:13.12.0-alpine AS frontendbuilder

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY frontend/package.json ./
COPY frontend/package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent

# add app
COPY ./frontend ./

RUN npm run build

# ------- executable build ------- #
FROM alpine:3.9
RUN apk add ca-certificates

RUN export GOFIRE_MASTER_HOST=`/sbin/ip route|awk '/default/ { print $3 }'` && export GOFIRE_MASTER=true

RUN mkdir -p /frontend/build

COPY --from=frontendbuilder /app/build/* /frontend/build/

COPY --from=gofirebuilder /app/GoFire/gofire /app/
COPY --from=gofirebuilder /app/GoFire/config.yaml /
COPY --from=gofirebuilder /app/GoFire/app/config/ /app/config/

RUN chmod +x /app/gofire

CMD ["/app/gofire"]