#!/usr/bin/env bash




function IsCommitInBranch() {
  COMMIT_SHA=$1
  REQUIRED_BRANCH=$2
  for FBRANCH in $(git for-each-ref --format "%(refname)" refs/heads); do
    RL=`git rev-list ${FBRANCH}`
    echo "$RL" | grep -q ${COMMIT_SHA} || true
    if $(echo "${RL}" | grep -q ${COMMIT_SHA}); then
      if [ "$(echo ${FBRANCH} | sed "s/refs\/heads\///ig")" == "$REQUIRED_BRANCH" ]; then
        echo "$REQUIRED_BRANCH"
        return
      fi
    fi
  done
}



function PrepareXGO() {
  export NODE_PATH=/usr/lib/node_modules
  cp -rf /go_src/* /go/src/
  mkdir -p /go/src/gitlab.768bit.com
  rm -rf /go/src/gitlab.768bit.com/*
}



function QuitTagged() {
  APP_VERSION="$(git tag --contains ${1})"
  if [ ! "$APP_VERSION" ]; then
    echo "Commit is untagged - do build";
  else
    echo "Commit is tagged {APP_VERSION} - ignore build";
    exit 0;
  fi
}
