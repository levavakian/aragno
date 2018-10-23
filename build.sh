#!/bin/sh
SOURCE_STAGE="git"
VALUE=${1:-$SOURCE_STAGE}
echo "Building using $VALUE stage (git or local)"
docker build \
	--build-arg SSH_KNOWN_HOSTS="$(cat ~/.ssh/known_hosts)"  \
	--build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" \
	--build-arg SOURCE_STAGE=$VALUE \
	--build-arg GIT_RANDOMIZER="$(date|md5sum)" \
	-t aragno .
