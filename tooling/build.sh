#!/bin/sh

BASEDIR=$(dirname "$0")
cd $BASEDIR/native-runtime
docker build -t oreilly/rac-native-runtime .

cd $BASEDIR/native-devel
docker build -t oreilly/rac-native-devel .