#!/usr/bin/env bash

./build/bin/geth --datadir ./data export export.rlp.gz
rm -rf ./tmpdd
./build/bin/geth --datadir ./tmpdd init foundation.conf.json
./build/bin/geth --datadir tmpdd import export.rlp.gz

