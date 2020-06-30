#!/usr/bin/env bash

main() {

  make geth

  geth_cmd="$1"
  conf="$2"
  test_target_dir="$4"

  [[ ! -z $test_target_dir ]] && [[ ! -d $test_target_dir ]] && mkdir -p "$test_target_dir"

  echo "Geth: $geth_cmd  Config: $conf --> Test target dir: $test_target_dir"

  cp "$conf" "$test_target_dir/"
  $geth_cmd version >"$test_target_dir/geth_version.txt"

  if [[ $3 == reset ]]; then
    echo "reset: clearing and re-initing ./data"
    rm -rf ./data && ./$geth_cmd --datadir ./data init "$conf"
  fi

  (
    sleep 5
    echo '{"jsonrpc":"2.0","method":"miner_stop","params":[],"id":1}' | nc -U -W1 /tmp/geth.ipc
    echo '{"jsonrpc":"2.0","method":"miner_start","params":[1],"id":1}' | nc -U -W1 /tmp/geth.ipc
  ) &

  (
  "$geth_cmd" 2>geth.log --syncmode=full --gcmode=archive --datadir ./data --ipcpath /tmp/geth.ipc \
    --http --http.corsdomain='*' \
    --mine --miner.gasprice=1 --miner.gastarget=8000000 --miner.gaslimit=10000000 --miner.etherbase=0x25b7955e43adf9c2a01a9475908702cce67f302a --miner.recommit=2s \
    --nodiscover --debug
  )&
  geth_pid=$!

  ./runtest-squashing.sh &
  test_pid=$!

  trap "kill -2 $geth_pid $test_pid" SIGINT EXIT
  wait $test_pid

  # ./"$geth_cmd" attach /tmp/geth.ipc
}
main $*

# run tests:
# while :; do ./runtest-classicsquashed.sh; done
