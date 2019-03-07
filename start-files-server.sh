IP=$(ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1')
PORT=9090
DIR=$(echo $PWD)
mkdir -p $DIR/files
docker run -it --rm -p $IP:$PORT:$PORT -v $DIR/files:/go-file-server/files file-server

 