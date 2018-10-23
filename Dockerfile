FROM golang:latest as src_download

ARG GIT_RANDOMIZER

# Authorize ssh private key
RUN mkdir -p /root/.ssh/
ARG SSH_PRIVATE_KEY
RUN echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_rsa
ARG SSH_KNOWN_HOSTS
RUN echo "$SSH_KNOWN_HOSTS" > /root/.ssh/known_hosts
RUN chmod 400 /root/.ssh/id_rsa

git clone git@github.com:levavakian/aragno.git

FROM golang:latest

COPY --from=src_download /aragno /go/src/aragno
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/aragno
RUN dep ensure
