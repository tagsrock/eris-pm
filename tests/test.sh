#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the test manager for epm. It will run the testing
# sequence for epm using docker.

# ----------------------------------------------------------
# REQUIREMENTS

# eris installed locally

# ----------------------------------------------------------
# USAGE

# test.sh

# ----------------------------------------------------------
# Set defaults

# Where are the Things
base=github.com/eris-ltd/eris-pm
if [ "$CIRCLE_BRANCH" ]
then
  repo=${GOPATH%%:*}/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
  circle=true
else
  repo=$GOPATH/src/$base
  circle=false
fi
branch=${CIRCLE_BRANCH:=master}
branch=${branch/-/_}

# Define now the tool tests within the Docker container will be booted from docker run
entrypoint="/home/eris/test_tool.sh"
testimage=quay.io/eris/epm
testuser=eris
was_running=0

# ---------------------------------------------------------------------------
# Needed functionality

is_it_running(){
  if [[ "$(eris services ls | grep $1 | awk '{print $2}')" == "No" ]]
  then
    return 0
  else
    was_running=1
    return 1
  fi
}

# ---------------------------------------------------------------------------
# Get the things build and dependencies turned on

echo "Hello! I'm the testing suite for epm."
echo ""
echo "Building epm in a docker container."
strt=`pwd`
cd $repo
export testimage
export repo
# suppressed by default as too chatty
tests/build_tool.sh > /dev/null
# tests/build_tool.sh
if [ $? -ne 0 ]
then
  echo "Could not build epm. Debug via by directly running [`pwd`/tests/build_tool.sh]"
  exit 1
fi

# ---------------------------------------------------------------------------
# Go!

uuid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)
is_it_running keys
if [ $? -eq 0 ]
then
  eris services start keys
fi
eris chains new epm-tests-$uuid

test_exit=0

# ---------------------------------------------------------------------------
# Cleaning up

eris chains stop epm-tests-$uuid
eris chains rm -xf epm-tests-$uuid
if [ "$was_running" -eq 0 ]
then
  eris services stop keys
fi
echo ""
echo ""
echo "Done. Exiting with code: $test_exit"
cd $strt
exit $test_exit