#!/bin/sh

BASEDIR="$(cd "$(dirname "$0")" && pwd)"

cd $BASEDIR/native-devel
docker build -t oreilly/rac-native-devel .

cd $BASEDIR/native-runtime
docker build -t oreilly/rac-native-runtime .

exit 0

BASEDIR_BIN=$BASEDIR/.bin
DOCKER_SQUASH=$BASEDIR_BIN/docker-squash
if [[ "$OSTYPE" == "linux-gnu" ]]; then
  DOCKER_SQUASH_URL="https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-linux-amd64-v0.2.0.tar.gz"
elif [[ "$OSTYPE" == "darwin"* ]]; then
  PATH="/usr/local/opt/gnu-tar/libexec/gnubin":$PATH
  if [[ `tar --version` == "bsdtar"* ]]; then
    echo "This script needs gnu-tar to work."
  fi
  DOCKER_SQUASH_URL="https://github.com/jwilder/docker-squash/releases/download/v0.2.0/docker-squash-darwin-amd64-v0.2.0.tar.gz"
else
  echo "Don't know how to stip outside mac OS or Linux."
fi
if [ ! -f $DOCKER_SQUASH ]; then
  mkdir -p $BASEDIR_BIN
  curl $DOCKER_SQUASH_URL -L -o $BASEDIR_BIN/x.tgz && tar xzf $BASEDIR_BIN/x.tgz && rm $BASEDIR_BIN/x.tgz
fi
docker save `docker images oreilly/rac-native-runtime:latest -qa` | sudo $BASEDIR_BIN/docker-squash -t squash -verbose | docker load
