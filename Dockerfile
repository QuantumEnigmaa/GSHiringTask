##################################
# Building go binary
##################################
FROM golang:1.18-alpine3.15 as builder
RUN mkdir /build
COPY . /build/
WORKDIR /build
RUN go build -buildvcs=false . 

##################################
# Creating final image from binary
##################################
FROM alpine:3.15
COPY --from=builder /build/beeant .
RUN addgroup -g 1000 beeant && adduser -u 1000 -G beeant -D beeant
USER beeant
EXPOSE 8088

ENTRYPOINT ["./beeant"]