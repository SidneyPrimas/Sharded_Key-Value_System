#!/usr/bin/env bash

mkdir -p Custom_test
for i in {1..100}
do
  go test -run TestSnapshotRPC >> ./Custom_test/TestSnapshotRPC.txt
done