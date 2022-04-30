FROM golang:1.15.7 AS build_base
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod tidy

# RUN go mod vendor
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GO11MODULE=on  go build -mod=readonly -v -a -o mac .

FROM alpine:latest
# FROM centurylink/ca-certs
# FROM debian:buster-slim
ENV MNT_DIR /asset

WORKDIR /
RUN mkdir -p $MNT_DIR
COPY --from=build_base /app/mac /
COPY --from=build_base /usr/local/go/lib/time/zoneinfo.zip /
ENV TZ=Asia/Hong_Kong
ENV ENV=ksc
ENV ZONEINFO=/zoneinfo.zip
COPY /asset $MNT_DIR

CMD ["/mac"]
