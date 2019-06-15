#!/bin/bash
# name of your project with chaincode default
PROJECT_NAME="misterybox"
LANGUAGE="golang"
VERSION=$1

export MSYS_NO_PATHCONV=1

CC_SRC_PATH=i3tech.com/emerald/misterybox/chaincode

echo "upgrade chaincode"

docker exec cli peer chaincode install -n $PROJECT_NAME -v $VERSION -p "$CC_SRC_PATH" -l "$LANGUAGE"
sleep 5
docker exec cli peer chaincode upgrade -o orderer.emerald.com:7050 -C jewelchannel -n $PROJECT_NAME -v $VERSION -c '{"Args":[""]}' -P "OR ('RubyMSP.member')"

