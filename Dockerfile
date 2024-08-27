FROM golang:alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/kube-agent

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Install gcc and g++
RUN apk --update upgrade && \
    apk add gcc && \
    apk add g++ && \
    rm -rf /var/cache/apk/*

#Build the Go app
RUN go build -o ./out/kube-agent .

#Sets the env variable to 8080 so our server can use that port
ENV PORT 8080

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
RUN chmod +x ./out/kube-agent
CMD ["./out/kube-agent"]
