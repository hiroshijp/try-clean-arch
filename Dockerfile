# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21.4 as builder
# Create and change to the app directory.
WORKDIR /app
# Retrieve application dependencies using go modules.
# Allows container builds to reuse downloaded dependencies.
COPY ./app/go.* ./
RUN go mod download
# Copy local code to the container image.
COPY ./app ./
# Build the binary.
# -mod=readonly ensures immutable go.mod and go.sum in container builds.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o server

# Build the runtime container image from scratch, copying what is needed from the previous stage.  
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM scratch
# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
ENTRYPOINT ["/server"]