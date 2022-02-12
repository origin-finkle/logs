#!/bin/bash

version=$(git describe --tags)
go install -ldflags "-X 'github.com/origin-finkle/logs/internal/version.Version=$version'" ./...
