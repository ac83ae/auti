#!/usr/bin/env bash

CURR_DIR=$(pwd)

cd "$CURR_DIR"/2_transaction_record || exit
./3_bench_local_chain_submit.sh
sleep 5
./4_bench_local_chain_read.sh
sleep 5
./5_bench_local_chain_commit_submit.sh
sleep 5
./6_bench_local_chain_commit_read.sh
sleep 5
./7_bench_org_chain_submit.sh
sleep 5
./8_bench_org_chain_read.sh
sleep 5

cd "$CURR_DIR"/3_consistency_examination || exit
./6_bench_aud_chain_submit.sh
sleep 5
./7_bench_aud_chain_read.sh

cd "$CURR_DIR" || exit
