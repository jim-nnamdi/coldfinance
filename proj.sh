#!/bin/bash

# initialize this from the makefile
function runProj{
    echo "starting the coldfinance binary ..."
    go run main.go 
    echo "running on port 9900 ..."
}
runProj
projx=$(runProj)
echo $(projx)