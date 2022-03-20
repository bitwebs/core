#!/bin/bash

VERSION="${1:-v0.5.11-oracle}"

pushd .. 

git checkout $VERSION
docker build -t bitwebs/iq-core:$VERSION .
git checkout -

popd

docker build --build-arg version=$VERSION --build-arg chainid=swartz-1 -t bitwebs/iq-core-node:$VERSION .
docker build --build-arg version=$VERSION --build-arg chainid=mcafee-1 -t bitwebs/iq-core-node:$VERSION-testnet .