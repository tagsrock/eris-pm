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
#INSTALL_BASE="/usr/local/bin"
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

# build the static binary in a clean build environment
# echo "Building static binary in a clean build environment"
# docker build -t $testimage:build $repo

# if [ "$ERIS_PM_LEAN" = true ]
# then
  # copy out the built executable to local host
#   echo "Recovering artifact to local host"
#  docker run --rm --entrypoint cat $testimage:build $INSTALL_BASE/epm > $repo/epm_artifact
  # move the built artefact into a clean docker image
  # echo "Creating fresh deployment image with built artifact"
#   docker build -f $repo/ -t $testimage:deploy $repo
# else
  # rename build image as deployment; skipping repackaging step
#  docker tag $testimage:build $testimage:deploy
# fi

if [[ "$branch" = "master" ]]
then
  # retag original deploy image as latest image
  #docker tag $testimage:deploy $testimage:latest
  docker build -t $testimage:latest $repo
  docker tag $testimage:latest $testimage:$release_maj
  docker tag $testimage:latest $testimage:$release_min
else
  # retag original deploy image as minor release image
  #docker tag $testimage:deploy $testimage:$release_min
  docker build -t $testimage:$release_min $repo
fi

# clean up
# docker rmi $testimage:build
# docker rmi $testimage:deploy