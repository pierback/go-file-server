

FROM golang:alpine
# Build Args
ARG APP_NAME=go-file-server
ARG FILE_DIR=/${APP_NAME}/files

RUN apk update && \
    apk upgrade && \
    apk add git

# Create Log Directory
RUN mkdir -p ${FILE_DIR}

# Environment Variables
ENV LOG_FILE_LOCATION=${FILE_DIR}/app.log 


ADD ./server /go/src/app
WORKDIR /go/src/app

ENV GO111MODULE=on

COPY ./server/go.mod . 
COPY ./server/go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Download dependencies
# RUN go get -d -v ./...

# # Install the package
# RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 9090
EXPOSE 9999/udp
ENV PORT=9090
# Declare volumes to mount
VOLUME ["/go-file-server/files"]

CMD ["go", "run", "main.go", "pinger.go"]