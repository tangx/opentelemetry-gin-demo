
# language=golang
## IMAGE ARGS

ARG RUNTIME_IMAGE=golang:1.20.5
FROM ${RUNTIME_IMAGE} as builder

ENV GOPROXY=https://goproxy.cn,direct 

WORKDIR /go/src
ADD . .
RUN make install

FROM alpine:3.16 as runtime
EXPOSE 3000
COPY --from=builder /go/bin/webapp /usr/bin/webapp
ENTRYPOINT [ "/usr/bin/webapp" ]

