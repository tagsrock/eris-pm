[![Circle CI](https://circleci.com/gh/eris-ltd/epm-go/tree/master.svg?style=svg)](https://circleci.com/gh/eris-ltd/eris-pm/tree/master)

[![GoDoc](https://godoc.org/github.com/eris-pm?status.png)](https://godoc.org/github.com/eris-ltd/eris-pm)

# Eris Package Manager: The Smart Contract Package Manager

This repository continues the work which began with the [EPM](https://github.com/eris-ltd/epm-go) version of EPM.


## Install

1. [Install go](https://golang.org/doc/install)
2. Ensure you have gmp installed (sudo apt-get install libgmp3-dev || brew install gmp)
3. `go get github.com/eris-ltd/eris-pm/cmd/epm`


## Formatting

Ethereum input data and storage deals strictly in 32-byte segments or words, most conveniently represented as 64 hex characters. When representing data, strings are right padded while ints/hex are left padded.

EPM accepts integers, strings, and explicitly hexidecimal strings (ie. "0x45"). If your string is strictly hex characters but missing a "0x", it will still be treated as a normal string, so add the "0x". Addresses should be prefixed with "0x" whenever possible. Integers in base-10 will be handled, hopefully ok.

Values stored as EPM variables will be immediately converted to the proper hex representation. That is, if you store "dog", you will find it later as `0x0000000000000000000000000000000000000000000000000000646f67`.

## Directory

As part of the larger suite of Eris libraries, epm works out of the core directory in `~/.eris`.

## Command Line Interface

Assuming your `go bin` is on your path, the cli is accessible as `epm`. EPM provides a git-like interface for managing smart contract libraries. For more details, see `epm --help` or `epm [command] --help`.