#!/usr/bin/env bash

env AM=on AM_CLIENT=/tmp/geth.ipc AM_KEYSTORE=$(pwd)/keys AM_ADDRESS=0xb1355e69c1ba7b401e38c95cdb936727ae88be76 AM_PASSWORDFILE=$(pwd)/keys/pass.txt go test -timeout 99999s -count 1 -run TestState -v ./tests |& tee gethsquashed.out
