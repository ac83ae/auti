#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o closc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/closc_ce_merkle_proof_verify.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20; do
  #   for j in {1..10}; do
  #  echo "No: $j" >>$LOG_FILE_DIR
  ./closc.out -phase ce -process merkle_proof_verify -num $i -numIter 10 | tee -a $LOG_FILE_DIR
  sleep 1
  #   done
done

rm closc.out
