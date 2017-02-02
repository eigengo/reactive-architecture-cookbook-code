#!/bin/sh

cd /var/src/module
mkdir -p target
cd target
rm -rf *
cmake ..
cmake --build .
ctest -v
