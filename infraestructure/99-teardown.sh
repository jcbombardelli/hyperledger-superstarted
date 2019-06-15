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