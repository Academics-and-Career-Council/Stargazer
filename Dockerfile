FROM golang:1.17.2-alpine AS builder
RUN apk update && apk add --no-cache git build-base
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
RUN apk add --no-cache zeromq-dev musl-dev pkgconfig alpine-sdk libsodium-dev libzmq-static libsodium-static
RUN mkdir /stargazer
WORKDIR /stargazer
COPY ./ ./
RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/stargazer
RUN CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium" \
  CGO_ENABLED=1 \
  GOOS=linux \
  go build -v -a --ldflags '-extldflags "-static" -v' \
  -o build/stargazer
  


FROM scratch
LABEL MAINTAINER="Akshit Verma"
LABEL VERSION="0.0.1"
COPY --from=builder /stargazer/build/stargazer /go/bin/stargazer
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/stargazer"]