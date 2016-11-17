#!/bin/bash

if [ "x$ARTIFACTS_DIR" = "x" ] ; then 
  echo "ERROR: the environment variable ARTIFIACTS_DIR must be set"
  exit 1
fi

# should be the standard for every script
OLD_DIR=$( pwd )
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd $DIR

# now our stuff
os=$(go env GOOS)
arch=$(go env GOARCH)
echo "Building ws now..."

# Get the git commit
#GIT_COMMIT="$(git rev-parse HEAD)"
#GIT_DIRTY="$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

CGO_ENABLED=0 go build -v \
  -o $ARTIFACTS_DIR/ws-$os-$arch \
  github.com/mheese/ws

cd $OLD_DIR
