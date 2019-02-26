#!/usr/bin/env bash

function CloneRepo() {
  REPO_IN=$1
  JOB_TOKEN=$2
  echo "CLONING REPO: $REPO_IN"
  git clone https://gitlab-ci-token:${JOB_TOKEN}@gitlab.768bit.com/${REPO_IN}.git /go/src/gitlab.768bit.com/${REPO_IN} || true
}
