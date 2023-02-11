FROM golang:1.19 AS build

ENV GOPATH /go
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/aliyun-ddns main.go

RUN strip /go/bin/aliyun-ddns
RUN test -e /go/bin/aliyun-ddns

FROM alpine:latest

LABEL org.opencontainers.image.source=https://github.com/chyroc/aliyun-ddns
LABEL org.opencontainers.image.description="DDNS Tool, Automatically Update Your Public IP to Aliyun DNS."
LABEL org.opencontainers.image.licenses="Apache-2.0"

ENV ALIYUN_ACCESS_KEY_ID=""
ENV ALIYUN_ACCESS_KEY_SECRET=""
ENV DOMAIN=""
ENV RR=""
ENV IP_TYPE="ipv6"

COPY --from=build /go/bin/aliyun-ddns /bin/aliyun-ddns

CMD /bin/aliyun-ddns auto-update -domain=$DOMAIN -rr=$RR -ip-type=$IP_TYPE