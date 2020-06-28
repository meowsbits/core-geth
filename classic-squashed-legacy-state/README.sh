./build/bin/geth --datadir ./tmpdd init classic-squashed-legacy-state/classic.ref.json
./build/bin/geth --datadir ./tmpdd import classic-squashed-legacy-state/export.rlp.gz
