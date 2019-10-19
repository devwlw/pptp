FROM anapsix/alpine-java
MAINTAINER kev <noreply@easypi.pro>
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
#RUN apk add --no-cache musl-dev go
#RUN apk add git
#RUN mkdir -p /go/src/mail/deliver/go
#RUN mkdir -p /go/src/github.com/urfave
#ENV GOPATH=/go
#RUN go version
#RUN go get github.com/urfave/cli

#COPY go/ /go/src/mail/deliver/go
RUN apk add --no-cache \
            --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing/ \
            pptpclient
RUN apk add curl
COPY entrypoint.sh /entrypoint.sh
COPY tick.sh /tick.sh
COPY pptp_start.sh /pptp_start.sh
COPY pptp_stop.sh /pptp_stop.sh
COPY iem-server-163.jar /iem-server-163.jar
#为了方便,直接将编译好的程序copy到docker镜像里面去
COPY deliver /deliver
#RUN cd /go/src/mail/deliver/go && \
#       go build -o /app
ENTRYPOINT ["/entrypoint.sh"]