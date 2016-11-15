#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
[ -d vendor ] && rm -r vendor
. 'tools/.vendor-helpers.sh'

clone git gopkg.in/iconv.v1 v1.1.1
clone git github.com/Sirupsen/logrus v0.9.0 
clone git github.com/codegangsta/cli 6086d7927ec35315964d9fea46df6c04e6d697c1
clone git github.com/mattn/go-sqlite3 v1.1.0

clean && mv vendor/src/* vendor && rmdir vendor/src

