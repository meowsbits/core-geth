#!/usr/bin/env bash

geth_cmd="$1"
conf="$2"

"$geth_cmd" --datadir ./data export export.rlp.gz
rm -rf ./tmpdd
"$geth_cmd" --datadir ./tmpdd init "$conf"
"$geth_cmd" --datadir tmpdd import export.rlp.gz

