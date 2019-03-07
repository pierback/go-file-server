

FROM golang:alpine
# Build Args
ARG APP_NAME=go-file-server
ARG FILE_DIR=/${APP_NAME}/files

# Create Log Directory
RUN mkdir -p ${FILE_DIR}

# Environment Variables
ENV LOG_FILE_LOCATION=${FILE_DIR}/app.log 


ADD ./server /go/src/app
WORKDIR /go/src/app

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 9090
ENV PORT=9090
# Declare volumes to mount
VOLUME ["/go-file-server/files"]

CMD ["go", "run", "main.go"]