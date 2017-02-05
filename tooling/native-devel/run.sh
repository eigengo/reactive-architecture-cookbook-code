#!/bin/sh

cd /var/src/module
mkdir -p target
cd target
rm -rf *
cmake ..
make -j8
ctest -v
