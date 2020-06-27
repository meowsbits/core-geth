#!/usr/bin/env bash

set -e

make all

./build/bin/geth --datadir ./data init classic.ref.json


