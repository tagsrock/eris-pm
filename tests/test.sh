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

# test.sh [setup]

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

# Key variables
key1_addr="1040E6521541DAB4E7EE57F21226DD17CE9F0FB7"
key2_addr="58FD1799AA32DED3F6EAC096A1DC77834A446B9C"

# Other variables
uuid=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)
was_running=0

# ---------------------------------------------------------------------------
# Needed functionality

ensure_running(){
  if [[ "$(eris services ls | grep $1 | awk '{print $2}')" == "No" ]]
  then
    eris services start $1 1>/dev/null
    sleep 3 # boot time
  else
    was_running=1
  fi
}

early_exit(){
  if [ $? -eq 0 ]
  then
    return 0
  fi

  echo "There was an error duing setup; keys were not properly imported. Exiting."
  if [ "$was_running" -eq 0 ]
  then
    if [ "$circle" = true ]
    then
      eris services stop keys
    else
      eris services stop -rx keys
    fi
  fi
  exit 1
}

test_setup(){
  echo "Getting Setup"
  if [ "$circle" = true ]
  then
    export ERIS_PULL_APPROVE="true"
    eris init --yes --skip-pull 1>/dev/null
  fi

  ensure_running keys
  eris keys import $key1_addr --src tests/fixtures/keys/$key1_addr 1>/dev/null
  eris keys import $key2_addr --src tests/fixtures/keys/$key2_addr 1>/dev/null

  # check keys were properly imported
  eris keys pub $key1_addr 1>/dev/null
  early_exit
  eris keys pub $key2_addr 1>/dev/null
  early_exit

  eris chains new epm-tests-$uuid --dir tests/fixtures/chaindata 1>/dev/null
  sleep 5 # boot time
  echo "Setup complete"
}

perform_tests(){
  echo ""
  cd $repo/tests/fixtures
  apps=(app*/)
  for app in "${apps[@]}"
  do
    echo ""
    echo -e "Testing EPM using fixture =>\t${app%/}"

    # Run the epm deploy
    cd $app
    if [ "$circle" = false ]
    then
      eris contracts test --chain "epm-tests-$uuid" -d
      # eris contracts test --chain "epm-tests-$uuid"
    else
      eris contracts test --chain "epm-tests-$uuid" --rm
    fi

    # Set exit code properly
    test_exit=$?
    if [ $test_exit -ne 0 ]
    then
      failing_dir=`pwd`
      break
    fi

    # Reset for next run
    cd ..
  done
}

test_teardown(){
  if [ $test_exit -ne 0 ]
  then
    echo ""
    echo "EPM Log on Failed Test."
    cat $failing_dir/epm.csv
  fi
  if [ "$circle" = false ]
  then
    echo ""
    eris chains stop -rxf epm-tests-$uuid 1>/dev/null
    eris chains rm -f epm-tests-$uuid 1>/dev/null
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
}

# ---------------------------------------------------------------------------
# Get the things build and dependencies turned on

echo "Hello! I'm the marmot that tests the epm tooling."
start=`pwd`
cd $repo
echo ""
echo "Building epm in a docker container."
set -e
tests/build_tool.sh 1>/dev/null
set +e
if [ $? -ne 0 ]
then
  echo "Could not build epm. Debug via by directly running [`pwd`/tests/build_tool.sh]"
  exit 1
fi
echo "Build complete."
echo ""

# ---------------------------------------------------------------------------
# Setup

test_setup

# ---------------------------------------------------------------------------
# Go!

if [[ "$1" != "setup" ]]
then
  perform_tests
fi

# ---------------------------------------------------------------------------
# Cleaning up

if [[ "$1" != "setup" ]]
then
  test_teardown
fi