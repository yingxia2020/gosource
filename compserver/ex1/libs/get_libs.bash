#!/bin/bash

cd `dirname $0`

mkdir -p usr/lib
mkdir -p usr/include

USR_DIR="`pwd`/usr"

#isa-l
if [ -d ./isa-l/.git ]; then
    pushd isa-l
    git pull
    popd
else
    git clone https://github.com/01org/isa-l.git
fi

pushd isa-l
./autogen.sh
./configure --prefix=$USR_DIR --libdir=$USR_DIR/lib
make
make install
popd

#benchmark
if [ -d ./benchmark/.git ]; then
    pushd benchmark
    git pull
    popd
else
    git clone https://github.com/google/benchmark.git
fi

pushd benchmark
rm -f CMakeCache.txt
cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=$USR_DIR
make
make install
popd

#zlib
if [ -d ./zlib/.git ]; then
    pushd zlib
    git pull
    popd
else
    git clone https://github.com/madler/zlib.git
fi

pushd zlib
cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX:PATH=$USR_DIR
make
make install
popd
