#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o closc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/closc_tr_local_chain_sc.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

# install fablo if not installed
[ -f ./fablo ] ||
  curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_LOCAL_CHAIN_DIR=${PWD}

FABLO_LOCAL_CHAIN_CONFIG="local-chain-config-sc.yaml"

function print_size() {
  echo "Blockchain size of peer0.org1.example.com:" >>$LOG_FILE_DIR
  docker exec -it peer0.org1.example.com bash -c "ls -lh /var/hyperledger/production/ledgersData/chains/chains/mychannel" | tee -a $LOG_FILE_DIR
  echo "Blockchain size of peer0.aud1.example.com:" >>$LOG_FILE_DIR
  docker exec -it peer0.aud1.example.com bash -c "ls -lh /var/hyperledger/production/ledgersData/chains/chains/mychannel" | tee -a $LOG_FILE_DIR
}

for i in 1 2 4 8 16 32 64 128 256 512 1024; do
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    clean_up

    rm $FABLO_LOCAL_CHAIN_CONFIG
    python generate_local_chain_config.py $i

    start_time=$(($(date +%s%N)/1000000))
    ./fablo up $FABLO_LOCAL_CHAIN_CONFIG
    end_time=$(($(date +%s%N)/1000000))
    echo "Fablo up time: {$((end_time - start_time))} ms" | tee -a $LOG_FILE_DIR

    docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
    docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
    sleep 5

    print_size
    ./closc.out -phase tr -process local_sc -num $i -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
    print_size
  done
done

clean_up
rm closc.out
