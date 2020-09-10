FROM hub.xxx.com/library/golang:alpine AS builder

# 设置GOPATH
ENV GOPATH=/opt/gopath

# 设置项目名称
ENV PROJECT_NAME=community-blogger

# 设置项目地址路径
ENV PROJECT_PATH=${GOPATH}/src/${PROJECT_NAME}
ENV GO111MODULE=on
ENV GOFLAGS="-mod=vendor"

# 设置用于存放生成的二进制文件的目录
ADD . ${PROJECT_PATH}
WORKDIR ${PROJECT_PATH}

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
&& apk add --update make \
&& rm -rf /var/cache/apk/* \
&& make build

# Release docker image
FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
&& apk update \
&& apk add ca-certificates \
&& rm -rf /var/cache/apk/*

ENV GOPATH=/opt/gopath
ENV PROJECT_NAME=community-blogger
ENV PROJECT_PATH=${GOPATH}/src/${PROJECT_NAME}
WORKDIR ${PROJECT_PATH}
COPY --from=builder ${PROJECT_PATH} ${PROJECT_NAME}
COPY --from=builder ${PROJECT_PATH}/configs configs

ENV SERVICE=/service
WORKDIR ${SERVICE}

ARG TARGET_HOME
ARG TARGET_ARTICLE
ARG TARGET_USERRPC
ARG TARGET_USER
RUN echo ${TARGET_HOME}
RUN echo ${TARGET_ARTICLE}
RUN echo ${TARGET_USERRPC}
RUN echo ${TARGET_USER}
COPY --from=builder ${PROJECT_PATH}/${TARGET_HOME} ${TARGET_HOME}
COPY --from=builder ${PROJECT_PATH}/${TARGET_ARTICLE} ${TARGET_ARTICLE}
COPY --from=builder ${PROJECT_PATH}/${TARGET_USERRPC} ${TARGET_USERRPC}
COPY --from=builder ${PROJECT_PATH}/${TARGET_USER} ${TARGET_USER}
COPY --from=builder ${PROJECT_PATH}/configs configs
RUN echo "/${SERVICE}/${TARGET_HOME}" > /start.bash; \
    echo "/${SERVICE}/${TARGET_ARTICLE}" > /start.bash; \
    echo "/${SERVICE}/${TARGET_USERRPC}" > /start.bash; \
    echo "/${SERVICE}/${TARGET_USER}" > /start.bash; \
    chmod +x /start.bash

CMD ["/bin/sh", "/start.bash"]