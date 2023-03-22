#!/bin/bash
echo "Building Running server"
mkdir -p build
mkdir -p bin
make all
./bin/server
