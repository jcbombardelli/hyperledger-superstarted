#!/bin/bash

set -e

docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

rm -f ~/.hfc-key-store/*

docker rm -f $(docker ps -aq)
docker rmi -f $(docker images dev-* -q)
