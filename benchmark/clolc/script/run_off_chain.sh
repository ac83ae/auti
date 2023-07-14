#!/usr/bin/env bash

CURR_DIR=$(pwd)

cd "$CURR_DIR"/1_initialization || exit
./1_bench_default.sh
sleep 5

cd "$CURR_DIR"/2_transaction_record || exit
./5_bench_commit.sh
sleep 5
./6_bench_accumulate.sh
sleep 5

cd "$CURR_DIR"/3_consistency_examination || exit
./1_bench_acc_commit.sh
sleep 5
./2_bench_cal_b.sh
sleep 5
./3_bench_cal_c.sh
sleep 5
./4_bench_cal_d.sh
sleep 5
./5_bench_encrypt.sh
sleep 5
./6_bench_decrypt.sh
sleep 5
./9_bench_check.sh

cd "$CURR_DIR"/4_result_verification || exit
./1_bench_org_and_aud.sh
sleep 5
./2_bench_aud_pair.sh

cd "$CURR_DIR" || exit
