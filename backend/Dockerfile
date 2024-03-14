# Step 1: Build executable binary
FROM golang@sha256:dd8888bb7f1b0b05e1e625aa29483f50f38a9b64073a4db00b04076cec52b71c as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/myapp

# Step 2: Build a small image
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group


# Copy our static executable
COPY --from=builder /go/bin/myapp /go/bin/myapp

USER appuser:appuser

# Expose port if your app is a web server
EXPOSE 8080

ENTRYPOINT ["/go/bin/myapp"]

