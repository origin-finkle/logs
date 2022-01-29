#!/bin/bash
set -e

if [ -z "$CONFIG_FOLDER" ]
then
    CONFIG_FOLDER=$GOPATH/src/github.com/origin-finkle/wcl-origin/data/config
fi

CONFIG_FOLDER=$CONFIG_FOLDER go test -v ./... $@