# Eris-ABI
Eris ABI tool

A simple tool for constructing transaction call data for ABI-enabled contracts.

Features:
- command-line interface
- thats it (its early)

Planned:
- http interface?
- automated abi access in eris directory

##WARNING: This is completely untested software. It should not yet be presumed operational.

#CLI

```
ebi pack --file example1.abi set hello
```

#Commands    

Pack:
```
NAME:
   pack - generate a transaction

USAGE:
   command pack [command options] [arguments...]

OPTIONS:
   --file 	Specify the ABI file (Containing JSON ABI) to use
 ```
