#!/bin/bash

SERVER="go-clean-arch"
BASE_DIR=$PWD
INTERVAL=2

# 命令行参数，需要手动指定
ARGS=""

# image name, tag
IMAGE_NAME="jason-bai/go-clean-arch"
IMAGE_TAG="latest"

# container
LOCAL_PORT=8080
DB_ADDR=192.168.7.252:3306
DOCKER_DB_ADDR=192.168.7.252:3306
NAME=go-clean-arch

function build()
{
  echo "Building image '$IMAGE_NAME:$IMAGE_TAG'..."
  docker build -t $IMAGE_NAME:$IMAGE_TAG .
}

function start() 
{
  echo "Start a container from '$IMAGE_NAME:$IMAGE_TAG' image..."
  docker run --name go-clean-arch -p $LOCAL_PORT:8080 --env APISERVER_DB_ADDR=$DB_ADDR --env APISERVER_DOCKER_DB_ADDR=$DOCKER_DB_ADDR -d $IMAGE_NAME:$IMAGE_TAG
}

function stop() {
  echo "Stop current running container..."
  docker rm -f $NAME
}

case "$1" in
	'build')
  build
	;;  
	'start')
  start
	;;  
	'stop')
  stop
	;;  
	*)  
	echo "usage: $0 {build|start|stop}"
	exit 1
	;;  
esac
