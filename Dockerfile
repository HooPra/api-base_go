# # Start from a Debian image with the latest version of Go installed
# # and a workspace (GOPATH) configured at /go.
# FROM golang

# # Copy the local package files to the container's workspace.
# ADD . /go/src/github.com/hoopra/api-base_go

# # Build the outyet command inside the container.
# # (You may fetch or manage dependencies here,
# # either manually or with a tool like "godep".)
# RUN go install github.com/hoopra/api-base_go

# # Run the outlet command by default when the container starts.
# ENTRYPOINT /go/bin/api-base_go

# # Document that the service listens on port 8080.
# EXPOSE 8080

FROM golang:alpine
RUN mkdir -p /go/src/github.com/hoopra/api-base_go
WORKDIR /go/src/github.com/hoopra/api-base_go
RUN apk add git && adduser -S -D -H -h /app appuser
# RUN export GOPATH=$(pwd)
ADD . ./
RUN go get ./ && go build -o main .
USER appuser
CMD ["./main"]
