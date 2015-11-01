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

# Where are the Things?
base=github.com/eris-ltd/eris-pm
if [ "$CIRCLE_BRANCH" ]
then
  repo=${GOPATH%%:*}/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
  circle=true
else
  repo=$GOPATH/src/$base
  circle=false
fi

# Define now the tool tests within the Docker container will be booted from docker run
entrypoint="/home/eris/test_tool.sh"
testimage=quay.io/eris/epm
testuser=eris

# Other needed variables
uuid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)
was_running=0
key1="RqsMlojbh9RwUUVTfIXvLgBq4PzYxG3NNeD1eUWCUxyC2trhU0BkEoM/gJhDntS13MJsC0PCd17/6LkRtJ0Bgg=="
key2="cyAlnxHeURCJgNSPNea3aTjmSId0MfKykAR/iSRN19132APZNKM1FETmvHV83w60ds4PVvl153a+4dtqCC4q+Q=="

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
start=`pwd`
cd $repo
tests/build_tool.sh > /dev/null
if [ $? -ne 0 ]
then
  echo "Could not build epm. Debug via by directly running [`pwd`/tests/build_tool.sh]"
  exit 1
fi
echo "Build complete."
echo ""

# ---------------------------------------------------------------------------
# Go!

echo "Getting Setup"
if [ "$circle" = true ]
then
  export ERIS_PULL_APPROVE="true"
  eris init
fi

is_it_running keys
# mkdir ~/.eris/data/keys/data
# cp $GOPATH/src/github.com/eris-ltd/eris-pm/tests/fixtures/keys/* ~/.eris/data/keys/keys/data/.
# eris data import keys
if [ $? -eq 0 ]
then
  eris services start keys
fi
eris chains new epm-tests-$uuid --dir tests/fixtures/chaindata
sleep 5
echo "Setup complete"
echo ""
docker run --rm --link eris_service_keys_1:keys --link eris_chain_epm-tests-"$uuid"_1:chain --entrypoint $entrypoint --user $testuser --env CHAINID=epm-tests-$uuid $testimage
test_exit=$?

# ---------------------------------------------------------------------------
# Cleaning up

echo ""
eris chains stop -xf epm-tests-$uuid
eris chains rm -xof epm-tests-$uuid
rm -rf ~/.eris/data/epm-tests-*
if [ "$was_running" -eq 0 ]
then
  eris services stop keys
fi
if [ "$test_exit" -eq 0 ]
then
  echo "Tests complete! Tests are Green"
else
  echo "Tests complete. Tests are red :("
fi
cd $start
exit $test_exit