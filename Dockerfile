FROM golang:1.22 AS builder

ENV APP_HOME /code/newbee_api
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -v -o newbee_api main.go && chmod +x newbee_api

FROM alpine:latest

COPY --from=builder /code/newbee_api/newbee_api /app/newbee_api
COPY --from=builder /code/newbee_api/config/newbee.env.yaml /app/config/newbee.env.yaml

WORKDIR /app

EXPOSE 8081

CMD ["./newbee_api"]