#!/usr/bin/env bash
docker-compose build --pull
docker-compose up -d
docker-compose scale kafka=3