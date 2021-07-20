FROM centos

MAINTAINER liuzezhong<972568589@qq.com>

WORKDIR /app

ADD main .
ADD ./configs/application.yaml ./configs/application.yaml

EXPOSE 80

# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /app/cmd/ark-iot-gateway/test1.go

ENTRYPOINT ["./main"]
