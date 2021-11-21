#!/bin/bash
set -e

CONFIG_FOLDER=`pwd`/data/config go test -v ./... $@