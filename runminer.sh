#!/usr/bin/env bash

if [[ $1 == reset ]]; then
  rm -rf ./data && ./build/bin/geth --datadir ./data init foundation.conf.json
fi

(
sleep 5
echo '{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":1}' | nc -U -W1 /tmp/geth.ipc
echo '{"jsonrpc":"2.0","method":"miner_start","params":[1],"id":1}' | nc -U -W1 /tmp/geth.ipc
)&

./build/bin/geth --syncmode=full --gcmode=archive --datadir ./data --ipcpath /tmp/geth.ipc --mine --miner.gasprice=1 --keystore ./keys --nodiscover console |& tee geth.log

# run tests:
# while :; do ./runtest-classicsquashed.sh; done
