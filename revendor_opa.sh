#!/usr/bin/env bash

usage() {
    echo "$0 <commit-ish> (e.g., $0 v0.9.1)"
}

if [ $# -eq 0 ]; then
    usage
    exit 1
fi

# Copied from https://unix.stackexchange.com/questions/92895/how-can-i-achieve-portability-with-sed-i-in-place-editing
case $(sed --help 2>&1) in
  *GNU*) sed_i () { sed -i "$@"; };;
  *) sed_i () { sed -i '' "$@"; };;
esac

sed_i "/name = \"github.com\/open-policy-agent\/opa\"/{N;s/version = .*/version = \"$1\"/;}" Gopkg.toml

git status | grep Gopkg.toml

if [ $? -eq 0 ]; then
    dep ensure
    git add .
fi
