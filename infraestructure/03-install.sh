#!/bin/bash
# name of your project with chaincode default
PROJECT_NAME="misterybox"
LANGUAGE="golang"

export MSYS_NO_PATHCONV=1

CC_SRC_PATH=i3tech.com/emerald/misterybox/chaincode

docker-compose -f ./docker-compose.yml up -d cli

echo "Installing chaincode"
docker exec -e "CORE_PEER_LOCALMSPID=RubyMSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/ruby.emerald.com/users/Admin@ruby.emerald.com/msp" cli peer chaincode install -n $PROJECT_NAME -v 1.0 -p "$CC_SRC_PATH" -l "$LANGUAGE"
sleep 10