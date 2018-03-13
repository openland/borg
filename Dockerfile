FROM alpine:3.5
ENTRYPOINT ["/bin/borg"]

COPY borg /bin/