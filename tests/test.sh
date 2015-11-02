#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the test manager for epm. It will run the testing
# sequence for epm using docker.

# ----------------------------------------------------------
# REQUIREMENTS

# eris installed locally
# eris-keys installed locally

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
branch=${CIRCLE_BRANCH:=master}
branch=${branch/-/_}

# Define now the tool tests within the Docker container will be booted from docker run
entrypoint="/home/eris/test_tool.sh"
testimage=quay.io/eris/epm
testuser=eris

# Other needed variables
uuid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)
key1="RqsMlojbh9RwUUVTfIXvLgBq4PzYxG3NNeD1eUWCUxyC2trhU0BkEoM/gJhDntS13MJsC0PCd17/6LkRtJ0Bgg=="
key1=`echo -n "$key1" | base64 -d | hexdump -ve '1/1 "%.2X"'`
key2="cyAlnxHeURCJgNSPNea3aTjmSId0MfKykAR/iSRN19132APZNKM1FETmvHV83w60ds4PVvl153a+4dtqCC4q+Q=="
key2=`echo -n "$key2" | base64 -d | hexdump -ve '1/1 "%.2X"'`
# key1="46AB0C9688DB87D4705145537C85EF2E006AE0FCD8C46DCD35E0F5794582531C82DADAE153406412833F8098439ED4B5DCC26C0B43C2775EFFE8B911B49D0182"
# key2="7320259F11DE51108980D48F35E6B76938E648877431F2B290047F89244DD7DD77D803D934A3351444E6BC757CDF0EB476CE0F56F975E776BEE1DB6A082E2AF9"
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
  eris init --yes --skip-pull 1>/dev/null
  # by default the keys daemon does not export its port to the host
  # for this sequencing to work properly it needs to be exported.
  # this is a hack.
  echo 'ports = [ "4767:4767" ]' >> ~/.eris/services/keys.toml
fi

is_it_running keys
if [ $? -eq 0 ]
then
  eris services start keys 1>/dev/null
  sleep 3 # boot time
fi

keysHost=$(eris services inspect keys NetworkSettings.IPAddress)
eris-keys import "$key1" --no-pass --host $keysHost 1>/dev/null
eris-keys import "$key2" --no-pass --host $keysHost 1>/dev/null

eris chains new epm-tests-$uuid --dir tests/fixtures/chaindata 1>/dev/null
sleep 5 # boot time
echo "Setup complete"

echo ""
if [ "$circle" = true ]
then
  if [[ "$branch" = "master" ]]
  then
    branch="latest"
  fi
  docker run --link eris_service_keys_1:keys --link eris_chain_epm-tests-"$uuid"_1:chain --entrypoint $entrypoint --user $testuser --env CHAINID=epm-tests-$uuid $testimage:$branch
else
  docker run --rm --link eris_service_keys_1:keys --link eris_chain_epm-tests-"$uuid"_1:chain --entrypoint $entrypoint --user $testuser --env CHAINID=epm-tests-$uuid $testimage
fi
test_exit=$?

# ---------------------------------------------------------------------------
# Cleaning up

echo ""
if [ "$circle" = false ]
then
  eris chains stop -xf epm-tests-$uuid
  eris chains rm -xf epm-tests-$uuid
  rm -rf ~/.eris/data/epm-tests-*
  if [ "$was_running" -eq 0 ]
  then
    eris services stop -rx keys
  fi
fi
echo ""
if [ "$test_exit" -eq 0 ]
then
  echo "Tests complete! Tests are Green. :)"
else
  echo "Tests complete. Tests are Red. :("
fi
cd $start
exit $test_exit