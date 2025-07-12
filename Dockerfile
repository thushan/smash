# This Dockerfile is used by GoReleaser to build
# a minimal Docker image for Smash.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates
RUN adduser -D -H -s /bin/false smash
COPY smash /usr/local/bin/smash
RUN chmod +x /usr/local/bin/smash
USER smash

ENTRYPOINT ["/usr/local/bin/smash"]
