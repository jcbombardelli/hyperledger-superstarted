#!/bin/bash
set -e

docker-compose -f docker-compose.yml kill && docker-compose -f docker-compose.yml down

docker rm -f $(docker ps -aq)
docker rmi -f $(docker images dev-* -q)
