FROM golang:1.18 as builder
WORKDIR /app
COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Set compilation target variables in case someone is building this on a Mac.
RUN go build -o rental github.com/shidenkai0/rental/cmd/api


FROM alpine:3.15 as runtime
RUN apk add --update --no-cache ca-certificates
COPY --from=builder /app/rental /

# Ah, non-root user...
USER nobody
ENTRYPOINT ["/rental", "--logtostderr=true", "--cert-dir=/tmp"]
