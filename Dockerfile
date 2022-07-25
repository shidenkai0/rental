FROM golang:1.18 as builder
WORKDIR /app
COPY . .

# Install go-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# Set compilation target variables in case someone is building this on a Mac.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rental github.com/shidenkai0/rental/cmd/api


FROM alpine:3.15 as runtime
RUN apk add --update --no-cache ca-certificates
COPY --from=builder /app/rental /

# Add migrate binary to the runtime image so we can run migrations in production
COPY --from=builder /go/bin/migrate /bin/migrate

# Copy migrations directory for running migrations in production
COPY db/migrations migrations

# Ah, non-root user...
USER nobody
ENTRYPOINT ["/rental", "--logtostderr=true", "--cert-dir=/tmp"]
