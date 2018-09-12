# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.9.2 AS builder
RUN go version

COPY . /go/src/github.com/Jason-Bai/go-clean-arch/
WORKDIR /go/src/github.com/Jason-Bai/go-clean-arch/
RUN set -x && \
    go get -u github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o go-clean-arch .

# Stage 2 (to create a downsized "container executable", ~7MB)

# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:
#
#   FROM alpine:3.7
#   RUN apk --no-cache add ca-certificates
FROM scratch
WORKDIR /root/
# Copy the go-clean-arch binary file into /root/
COPY --from=builder /go/src/github.com/Jason-Bai/go-clean-arch/go-clean-arch .
# Copy the config.yaml file into /root/conf/
COPY --from=builder /go/src/github.com/Jason-Bai/go-clean-arch/conf/config.yaml .
# http server listens on port 8080
EXPOSE 8080
# Run the entryport main program
CMD ["./go-clean-arch", "-c", "./config.yaml"]
