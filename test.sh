#!/bin/bash
set -e

CONFIG_FOLDER=$GOPATH/src/github.com/origin-finkle/wcl-origin/data/config go test -v ./... $@