# This Dockerfile is used by GoReleaser to build
# a minimal Docker image for Smash.
FROM alpine:3.22

RUN apk add --no-cache ca-certificates
RUN adduser -D -H -s /bin/false smash

# Create output directory and set ownership
RUN mkdir -p /output && chown smash:smash /output

COPY smash /usr/local/bin/smash
RUN chmod +x /usr/local/bin/smash

USER smash
WORKDIR /data

ENTRYPOINT ["/usr/local/bin/smash"]
