#!/usr/bin/env bash

CURR_DIR=${PWD}

cd "$CURR_DIR"/2_transaction_record || exit
./1_bench_local_chain_submit.sh
sleep 5
./2_bench_local_chain_read.sh
sleep 5
./3_bench_org_chain_submit.sh
sleep 5
./4_bench_org_chain_read.sh
sleep 5

cd "$CURR_DIR"/3_consistency_examination || exit
./7_bench_aud_chain_submit.sh
sleep 5
./8_bench_aud_chain_read.sh

cd "$CURR_DIR" || exit