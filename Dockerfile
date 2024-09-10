FROM alpine:latest
RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata
COPY newbee_api /app/newbee_api
COPY config/newbee.env.yaml /app/config/newbee.env.yaml

WORKDIR /app

EXPOSE 28081

CMD ["./newbee_api"]

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o newbee_api ./main.go