#!/usr/bin/env bash
set -e

cd "$(dirname "$BASH_SOURCE")/.."
[ -d vendor ] && rm -r vendor
. 'tools/.vendor-helpers.sh'

clone git gopkg.in/iconv.v1 v1.1.1
clone git github.com/Sirupsen/logrus v0.9.0 
clone git github.com/urfave/cli v1.18.1
clone git github.com/mattn/go-sqlite3 v1.1.0
clone git gopkg.in/mgo.v2 3569c88678d88179dcbd68d02ab081cbca3cd4d0

clean && mv vendor/src/* vendor && rmdir vendor/src

