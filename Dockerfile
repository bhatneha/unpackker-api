FROM golang:alpine as builder
COPY . /unpackker-api-bin
WORKDIR /unpackker-api-bin
ENV GO111MODULE=on
RUN go mod vendor
RUN CGO_ENABLED=0 GOOS=linux go build -o unpackker-api
RUN go get github.com/nikhilsbhat/unpackker


FROM golang:alpine
WORKDIR /root/
COPY --from=builder /go/bin/unpackker /usr/bin/unpackker
COPY --from=builder unpackker-api-bin/unpackker-api unpackker-api
ENTRYPOINT ["./unpackker-api"]