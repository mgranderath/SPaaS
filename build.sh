#!/usr/bin/env bash
make server_linux
cd frontend
./node_modules/.bin/webpack --env prod
cd ..
docker build -t mgranderath/spaas:latest .