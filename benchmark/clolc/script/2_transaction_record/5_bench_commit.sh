#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o clolc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_tr_commit.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

for i in 1000 10000 100000 1000000; do
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    ./clolc.out -phase tr -process commit -numTXs $i -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 1
  done
done

rm clolc.out
