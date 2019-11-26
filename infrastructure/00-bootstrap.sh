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

# for use version 1.3.0, uncomment this line above and comment line 21 
# curl -sSL http://bit.ly/2ysbOFE | bash -s 1.3.0

# curl -sSL http://bit.ly/2ysbOFE | bash -s -- <fabric_version> <fabric-ca_version> <thirdparty_version>
curl -sSL http://bit.ly/2ysbOFE | bash -s -- 1.4.4 1.4.4 0.4.18

cp -r ./fabric-samples/bin ../
mv ./fabric-samples ../../

cd ..
echo export HYPERLEDGER=$PWD/bin >> ~/.bashrc
source ~/.bashrc
export PATH=$HYPERLEDGER:$PATH