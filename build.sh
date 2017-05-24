#!/bin/bash

docker system prune -f
docker build . -f Dockerfile -t targence/commits
