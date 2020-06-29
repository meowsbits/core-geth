#!/usr/bin/env bash



# make geth

if [[ $1 == reset ]]; then
  rm -rf ./data && ./build/bin/geth --datadir ./data init foundation.conf.json
fi

(
sleep 5
echo '{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":1}' | nc -U -W1 /tmp/geth.ipc
echo '{"jsonrpc":"2.0","method":"miner_start","params":[1],"id":1}' | nc -U -W1 /tmp/geth.ipc
)&

2>geth.log ./build/bin/geth --syncmode=full --gcmode=archive --datadir ./data --ipcpath /tmp/geth.ipc \
  --http --http.corsdomain='*' \
  --mine --miner.gasprice=1 --miner.gastarget=10000000 --miner.etherbase=0x25b7955e43adf9c2a01a9475908702cce67f302a --miner.recommit=3s \
  --nodiscover --debug console 

# run tests:
# while :; do ./runtest-classicsquashed.sh; done
