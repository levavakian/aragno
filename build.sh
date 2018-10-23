docker build --build-arg GIT_RANDOMIZER="$(date|md5sum)" --build-arg SSH_KNOWN_HOSTS="$(cat ~/.ssh/known_hosts)"  --build-arg SSH_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" -t aragno .
