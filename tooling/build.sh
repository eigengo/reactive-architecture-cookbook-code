#!/bin/sh

BASEDIR="$(cd "$(dirname "$0")" && pwd)"

cd $BASEDIR/native-runtime
docker build -t oreilly/rac-native-runtime .

cd $BASEDIR/native-devel
docker build -t oreilly/rac-native-devel .
