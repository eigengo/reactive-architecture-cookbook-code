#!/bin/sh

BASEDIR="$(cd "$(dirname "$0")" && pwd)"

cd $BASEDIR/native-devel
docker build -t oreilly/rac-native-devel .

cd $BASEDIR/native-runtime
docker build -t oreilly/rac-native-runtime .

BASEDIR_BIN=$BASEDIR/.bin
DOCKER_SQUASH=$BASEDIR_BIN/docker-squash
if [ ! -f $DOCKER_SQUASH ]; then
  DOCKER_SQUASH_URL=""
  if [[ "$OSTYPE" == "linux-gnu" ]]; then
    DOCKER_SQUASH_URL="https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-linux-amd64-v0.2.0.tar.gz"
  elif [[ "$OSTYPE" == "darwin"* ]]; then
    DOCKER_SQUASH_URL="https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-darwin-amd64-v0.2.0.tar.gz"
  else
    echo "Don't know how to stip outside mac OS or Linux."
  fi
  mkdir -p $BASEDIR_BIN
  curl $DOCKER_SQUASH_URL -L -o $BASEDIR_BIN/x.tgz && tar xzf x.tgz && rm x.tgz
fi

docker save `docker images oreilly/rac-native-runtime:latest -qa` | $BASEDIR_BIN/docker-squash -t squash -verbose | docker load
