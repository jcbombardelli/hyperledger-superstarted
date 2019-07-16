#!/bin/sh

set -e
# delete previous creds
rm -rf ./hfc-key-store
rm -rf ../hfc-key-store
rm -rf ~/.hfc-key-store/*

# remove previous crypto material and config transactions
rm -rf ./config/*
rm -rf ../config/*
rm -rf ./crypto-config/*
rm -rf ../crypto-config/*

mkdir -p ./.hfc-key-store

curl -sSL http://bit.ly/2ysbOFE | bash -s 1.3.0

cp -r ./fabric-samples/bin ../
mv ./fabric-samples ../../