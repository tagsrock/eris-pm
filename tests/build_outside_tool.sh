#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the build script for epm. It will build the tool
# into docker containers in a reliable and predicatable
# manner.

# ----------------------------------------------------------
# REQUIREMENTS

# docker installed locally

# ----------------------------------------------------------
# USAGE

# build_tool.sh

# ----------------------------------------------------------
# Set defaults
set -e
if [ "$CIRCLE_BRANCH" ]
then
  repo=`pwd`
else
  repo=$GOPATH/src/github.com/eris-ltd/eris-pm
fi
branch=${CIRCLE_BRANCH:=master}
branch=${branch/-/_}
testimage=${testimage:="quay.io/eris/epm"}

release_min=$(cat $repo/version/version.go | tail -n 1 | cut -d \  -f 4 | tr -d '"')
release_maj=$(echo $release_min | cut -d . -f 1-2)

# ---------------------------------------------------------------------------
# Go!

# This is a temporary solution to avoid gliding inside the docker container
# rather building it outside of the docker

# rm -rf $repo/vendor &>/dev/null
cd $repo
# glide install
go build --ldflags '-extldflags "-static"' -o ./epm_built ./cmd/epm

if [[ "$branch" = "master" ]]
then
  docker build -f ./Dockerfile_outsidebuild -t $testimage:latest $repo
  docker tag $testimage:latest $testimage:$release_maj
  docker tag $testimage:latest $testimage:$release_min
else
  docker build -f ./Dockerfile_outsidebuild -t $testimage:$release_min $repo
fi