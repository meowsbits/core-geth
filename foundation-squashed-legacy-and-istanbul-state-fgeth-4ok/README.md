
this is a datadir and export created with ethereum/go-ethereum

included is the configuration

i created the blockchain using 

- ./runminer.sh reset
- ./runtest-classicsquashed.sh

you need to hardcode and swap the desired client (./build/bin/geth, .builds/fgeth) and
desired configuration to use the former.

the _runtest_ program depends on some changes at `tests/` which are a little unpolished.


## important

the export.rlp.gz file included here imports successfully against core-geth and go-ethereum.

it was generated with go-ethereum.



