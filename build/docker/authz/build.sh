#!/usr/bin/env bash

cd $(dirname "$0")
cp -rv authz.yaml "$DST_DIR"
cp -rv authz.data.json "$DST_DIR"
cp -rv authz.rego "$DST_DIR"
cd -

