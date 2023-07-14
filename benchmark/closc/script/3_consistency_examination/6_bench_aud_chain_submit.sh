#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o closc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/closc_ce_aud_chain_submit.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

# install fablo if not installed
[ -f ./fablo ] ||
  curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_AUD_CHAIN_DIR=${PWD}

FABLO_AUD_CHAIN_CONFIG="aud-chain-config.yaml"
rm -f $FABLO_AUD_CHAIN_CONFIG
python3 config_gen.py --output_filename $FABLO_AUD_CHAIN_CONFIG --chaincode_name auti-aud-chain --chaincode_dir contract/aud_chain --num_orderers 3 --num_orgs 0 --num_auditors 8


# 256 test
clean_up
./fablo up $FABLO_AUD_CHAIN_CONFIG
docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
sleep 5
./closc.out -phase ce -process aud_submit -num 256 -numIter 10 | tee -a $LOG_FILE_DIR
sleep 5

# 1k test
clean_up
./fablo up $FABLO_AUD_CHAIN_CONFIG
docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
sleep 5
./closc.out -phase ce -process aud_submit -num 1000 -numIter 10 | tee -a $LOG_FILE_DIR
sleep 5

for i in 10000 100000 1000000; do
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    clean_up
    ./fablo up $FABLO_AUD_CHAIN_CONFIG
    docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
    docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
    sleep 5
    ./closc.out -phase ce -process aud_submit -num $i -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
  done
done

clean_up
rm -f $FABLO_AUD_CHAIN_CONFIG
rm closc.out
