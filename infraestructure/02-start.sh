#!/bin/bash
set -e

export MSYS_NO_PATHCONV=1

docker-compose -f docker-compose.yml down
docker-compose -f docker-compose.yml up -d ca.emerald.com orderer.emerald.com peer0.ruby.emerald.com couchdb

sleep 10

docker exec -e "CORE_PEER_LOCALMSPID=RubyMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@ruby.emerald.com/msp" peer0.ruby.emerald.com peer channel create -o orderer.emerald.com:7050 -c jewelchannel -f /etc/hyperledger/configtx/channel.tx
docker exec -e "CORE_PEER_LOCALMSPID=RubyMSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@ruby.emerald.com/msp" peer0.ruby.emerald.com peer channel join -b jewelchannel.block