# FROM ferry:latest
FROM golang:1.22-alpine3.20

# 设置工作目录
WORKDIR /app

# 复制 Go 源码到容器
COPY sms_gateway_demo /app

# 编译 Go 程序
RUN go build -o sms_gateway_demo sms_gateway_demo.go

# 暴露端口
EXPOSE 8080

# 运行程序
CMD ["./sms_gateway_demo"]
