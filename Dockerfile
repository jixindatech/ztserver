FROM node:16.8.0 as builder
ARG VUE=/usr/src/vue
COPY ./dashboard $VUE
WORKDIR $VUE
RUN npm install --registry=https://registry.npm.taobao.org
RUN npm run build:prod

FROM golang:alpine AS development
WORKDIR $GOPATH/src
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
ENV CGO_ENABLED=1
COPY ./server .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add build-base
RUN go mod tidy & go mod vendor
RUN go build -a -ldflags '-extldflags "-static"' -o ./bin/ztserver  ./cmd/ztserver

FROM alpine:latest AS production
WORKDIR /opt/ztserver
COPY --from=development /go/src/bin/ztserver .
RUN mkdir etc
COPY --from=development /go/src/etc/config.yaml etc/
RUN mkdir -p dashboard/dist
COPY --from=builder /usr/src/vue/dist dashboard/dist

EXPOSE 8000
ENTRYPOINT ["./ztserver", "-config", "etc/config.yaml"]


