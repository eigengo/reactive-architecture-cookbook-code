#!/bin/sh
cd native-devel
docker build -t oreilly/rac-native-devel .

cd ../native-runtime
docker build -t oreilly/rac-native-runtime .
