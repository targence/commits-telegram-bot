FROM golang:1.8.0-alpine

RUN apk add --update bash git curl ca-certificates openssl && update-ca-certificates
WORKDIR /go/src/app
COPY . $WORKDIR
RUN curl https://glide.sh/get | sh
RUN glide install
EXPOSE 3000
CMD ["go", "run", "*.go"]

