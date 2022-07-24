FROM golang:1.16 as builder
ARG Version
ARG CommitVersion
ARG BuildTime
LABEL version=$Version comshbuimit=$CommitVersion create_time=$BuildTime
ENV GOPROXY https://goproxy.cn
ADD . /tape
WORKDIR /tape
RUN  go version && go env && gcc -v && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build \
     -v -o tape cmd/tape/main.go

#生成中间镜像后,将build之后的可执行文件考到新的镜像中
FROM  alpine:3.14 as tape
ARG Version
ARG CommitVersion
ARG BuildTime
LABEL version=$Version commit=$CommitVersion create_time=$BuildTime
COPY --from=builder  /tape/tape /usr/local/bin
COPY --from=builder  /tape/deployments/config/policy /tape/
COPY --from=builder  /tape/deployments/config/tape-config.yaml /etc/tape/
# 切换软件源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
WORKDIR /tape
#容器内部开放端口
ENTRYPOINT ["tape"]
