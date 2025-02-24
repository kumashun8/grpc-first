FROM golang:1.23.5-alpine3.21 AS build

ENV ROOT=/go/src/project
WORKDIR $ROOT

COPY . $ROOT

RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:3.21.3

ENV ROOT=/go/src/project
WORKDIR $ROOT

RUN addgroup -S dockergroup && adduser -S docker -G dockergroup
USER docker

COPY --from=build ${ROOT}/server $ROOT

EXPOSE 8080
CMD ["./server"]
