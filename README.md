# hyperledger-fabric
Model's project for getting started Hyperledger Fabric

# Prerequisites

## Linux

- Install make
- Install build-essential

# Build

## Java

Execute: gradle clean build shadowJar

## Node API

Execute npm install em: /infrastructure/js

## Test on Peer

### Create MisteryBox
peer chaincode invoke -n misterybox-java -c '{"Args":["invoke", "a123b", "XXL", "jcbombardelli"]}' -C jewelchannel

### QueryHistoryMisteryBox

peer chaincode invoke -n misterybox-java -c '{"Args":["queryHistory", "a123b"]}' -C jewelchannel

# Troubleshootings

#### API error (400): OCI runtime create failed: container_linux.go:345: starting container process caused "exec: \"/root/chaincode-java/start\": stat /root/chaincode-java/start: no such file or directory": unknown

Isso pode acontecer caso voce esteja tentando criar imagens diferentes no docker com o mesmo nome