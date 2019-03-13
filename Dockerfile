FROM golang:alpine

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

EXPOSE 9090
EXPOSE 9999/udp
ENV PORT=9090
ENV ISDOCKER=true
# Declare volumes to mount
VOLUME ["/go-file-server/files"]

CMD ["go", "run", "main.go", "pinger.go"]