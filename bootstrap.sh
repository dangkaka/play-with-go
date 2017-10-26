#!/usr/bin/env bash
#
# bootstrap zfeed
# Param $1 => [environment dev/staging/live]
# Param $2 => [number of kafka nodes 1/2/3]
#

if test z$1 = 'zstaging'; then
	ENV='staging'
else
    ENV='dev'
    if [[ -z "${DOCKER_HOST_IP-}" ]]; then
      docker_host_ip=$(docker run --rm --net host alpine ip address show eth0 | awk '$1=="inet" {print $2}' | cut -f1 -d'/')
      # Work around Docker for Mac 1.12.0-rc2-beta16 (build: 9493)
      if [[ $docker_host_ip = '192.168.65.2' ]]; then
        docker_host_ip=$(/sbin/ifconfig | grep -v '127.0.0.1' | awk '$1=="inet" {print $2}' | cut -f1 -d'/' | head -n 1)
      fi
      export DOCKER_HOST_IP=$docker_host_ip
    fi
fi

if test z$2 = 'z'; then
    brokers=3
else
    brokers=$2
fi
echo 'brokers: ' $brokers

export APP_PATH=$(pwd)/app
echo 'app path: ' $APP_PATH

docker_compose_file='env/'$ENV'/docker-compose.yml'
echo "docker-compose path" $docker_compose_file

docker-compose -f $docker_compose_file build --pull
docker-compose -f $docker_compose_file -p zfeed up -d --scale kafka=$brokers