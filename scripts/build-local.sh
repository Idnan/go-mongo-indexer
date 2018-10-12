#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export CGO_ENABLED=0
my_dir="$(dirname "$0")"

go build -o ${my_dir}/../bin/index cmd/index/*.go
