#!/usr/bin/env bash

set -e

head -1 "$1"
cat "$1" | while read -r line; do
  if ! grep -q --line-buffered 'p=' <<< "$line"; then continue; fi
  p=$(echo "$line" | cut -d'=' -f2 | cut -d' ' -f1)
  lt=$(echo "$p < 0.05"|bc)
  if [[ $lt == 1 ]]; then
    echo "$line"
    # echo "p => $p"
  fi

done | sort -k5 -h -r
