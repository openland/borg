FROM openland/lalalinux:latest

# Copying sources
RUN mkdir -p /root/.go/src/github.com/statecrafthq/borg/
COPY . /root/.go/src/github.com/statecrafthq/borg/

# Go Dependencies
RUN go get \
    gopkg.in/kyokomi/emoji.v1 \
    github.com/urfave/cli \
    github.com/twpayne/go-geom \
    github.com/buger/jsonparser \
    gopkg.in/cheggaaa/pb.v1 \
    cloud.google.com/go/storage \
    golang.org/x/sync/semaphore \
    github.com/stretchr/testify \
    github.com/umahmood/haversine \
    github.com/aws/aws-sdk-go/aws/..

# Building Go
RUN cd /root/.go/src/github.com/statecrafthq/borg/ && go test ./... && go build && mv borg /usr/bin/

ENTRYPOINT ["/usr/bin/borg"]