#!/usr/bin/env bash
# ----------------------------------------------------------
# PURPOSE

# This is the test tool for epm. It should be ran (typically)
# from *inside* a docker container with chains and keys
# turned on and managed by test.sh

# ----------------------------------------------------------
# REQUIREMENTS

# epm installed locally
# assumes chain is running with chainID testChain
# assumes key server is running

# ----------------------------------------------------------
# USAGE

# test_tool.sh [local]

# ----------------------------------------------------------
# Set defaults

start=`pwd`
base=github.com/eris-ltd/eris-pm
if [ "$CIRCLE_BRANCH" ]
then
  repo=${GOPATH%%:*}/src/github.com/${CIRCLE_PROJECT_USERNAME}/${CIRCLE_PROJECT_REPONAME}
  circle=true
else
  repo=$GOPATH/src/$base
  circle=false
fi

# ---------------------------------------------------------
# Setup

echo "Hello! I'm the marmot that tests the epm tooling."
cd $repo/tests/fixtures
apps=(app*/)
for app in "${apps[@]}"
do
  echo ""
  echo -e "Testing EPM using fixture =>\t${app%/}"

  # Run the epm deploy
  cd $app
  if [[ "$1" == "local" ]]
  then
    epm test
  else
    epm --chain chain:46657 --sign keys:4767 test
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

# ---------------------------------------------------------
# Cleanup

if [ $test_exit -ne 0 ]
then
  echo "EPM Log on Failed Test."
  cat $failing_dir/epm.log
fi
echo ""
echo "Done. Exiting with code: $test_exit"
cd $start
exit $test_exit