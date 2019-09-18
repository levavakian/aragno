ARG SOURCE_STAGE=git

FROM golang:latest as src_git
ARG GIT_RANDOMIZER
# Authorize ssh private key
RUN mkdir -p /root/.ssh/
ARG SSH_PRIVATE_KEY
RUN echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_rsa
ARG SSH_KNOWN_HOSTS
RUN echo "$SSH_KNOWN_HOSTS" > /root/.ssh/known_hosts
RUN chmod 400 /root/.ssh/id_rsa

WORKDIR /
RUN git clone git@github.com:levavakian/aragno.git

FROM golang:latest as src_local
COPY . /aragno

FROM src_$SOURCE_STAGE AS stageforcopy

FROM golang:latest

COPY --from=stageforcopy /aragno /go/src/aragno
RUN go get -u github.com/golang/dep/cmd/dep
RUN curl -sL https://deb.nodesource.com/setup_10.x | bash -
RUN DEBIAN_FRONTEND=noninteractive  apt-get update -y && apt-get install nodejs build-essential libgl1-mesa-dev xorg-dev -y

WORKDIR /go/src/aragno
RUN cd ui && npm install
RUN dep ensure
