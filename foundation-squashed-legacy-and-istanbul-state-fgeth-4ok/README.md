
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


---


core-geth.export.rlp.gz created as:

core-geth import export.rlp.gz (go-ethereum created)
core-geth export core-geth.export.rlp.gz

This file (core-geth.export.rlp.gz) is importable OK with core-geth and go-ethereum.

This eliminates the possibility that the issue is with core-geth RLP encoding at the `export` stage,
since the RLP is consumable by both go-ethereum and core-geth.



So, what do we know?

- core-geth, when using the foundation.conf.json provided and manipulated with the ./runtest-classicsquashed.sh script
 (creating lots of transactions taken from the GeneralStateTests of ethereum/tests), generates a chain which, when exported,
  is _not_ importable, neither by core-geth itself nor go-ethereum. 
