FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/isotopsweden/docker-farmer

# Install docker farmer.
RUN go install github.com/isotopsweden/docker-farmer

# Run the docker farmer command when the container starts.
ENTRYPOINT /go/bin/docker-farmer

# http server listens on port 8080.
EXPOSE 8080