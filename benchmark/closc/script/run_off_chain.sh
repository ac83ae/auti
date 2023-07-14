#!/usr/bin/env bash

CURR_DIR=$(pwd)

cd "$CURR_DIR"/1_initialization || exit
./1_bench_default.sh
sleep 5
./2_bench_rand_gen.sh
sleep 5

cd "$CURR_DIR"/2_transaction_record || exit
./1_bench_commitment.sh
sleep 5
./2_bench_merkle_proof_gen.sh
sleep 5

cd "$CURR_DIR"/3_consistency_examination || exit
./1_bench_merkle_proof_verify.sh
sleep 5
./2_bench_merkle_proof_merge.sh
sleep 5
./3_bench_summarize_proof_result.sh
sleep 5
./4_bench_verify_commitment.sh
sleep 5
./5_bench_accumulate_commitments.sh
sleep 5

cd "$CURR_DIR"/4_result_verification || exit
./1_bench_merkle_batch_proof.sh
sleep 5
./2_bench_summarize_batch_proof.sh
sleep 5
./3_bench_verify_commit.sh

cd "$CURR_DIR" || exit
