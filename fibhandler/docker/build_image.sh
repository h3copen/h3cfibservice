#!/bin/sh
exe="fibhandler"
build_path=$(readlink -f $(dirname "$0"))
pushd "$build_path"

if [ ! -f ../$exe ];
then
    pushd ../
    go build
    popd
fi

cp ../$exe .

docker build -t h3cfibservice:latest .

rm fibhandler

popd
