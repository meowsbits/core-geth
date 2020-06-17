#!/usr/bin/env bash

# Show that all package tests succeed, unchanged by the patch reverts.
make test

# Show that the consensus eq. test fails.
make test-coregeth-consensus
