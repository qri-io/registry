FROM golang:1.10
LABEL maintainer="sparkle_pony_2000@qri.io"

ADD . /go/src/github.com/qri-io/registry
RUN cd /go/src/github.com/qri-io/registry
RUN go get ./...
RUN go install github.com/qri-io/registry
RUN cd /go/src/github.com/qri-io/registry/regserver

# set default port
# ENV PORT=3001
# EXPOSE 3001

# Set binary as entrypoint, initalizing ipfs repo if none is mounted
CMD ["regserver"]