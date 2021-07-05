#!/bin/bash


regex=$1

go test -v -run $1 github.com/Gravity-Tech/gravity-node-data-extractor/v2/extractors/susy/bridge