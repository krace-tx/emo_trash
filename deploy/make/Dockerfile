# 使用轻量级Linux发行版Alpine作为基础镜像
FROM alpine:3.20

# 设置环境变量
ENV GO_VERSION=1.24.5 \
    GOPATH=/go \
    GOROOT=/usr/local/go \
    PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# 更换Alpine仓库为国内镜像（解决网络访问问题）
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

# 安装基础工具链和依赖
RUN apk add --no-cache \
    bash \
    make \
    git \
    wget \
    curl \
    gcc \
    musl-dev \
    && rm -rf /var/cache/apk/*

# 安装Go环境
RUN wget -q https://mirrors.aliyun.com/golang/go${GO_VERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz \
    && rm go${GO_VERSION}.linux-amd64.tar.gz \
    && mkdir -p $GOPATH/src $GOPATH/bin

# 设置工作区
WORKDIR /app

# 容器启动命令
CMD ["bash"]