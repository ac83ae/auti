#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o closc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/closc_tr_local_chain_submit.log"
if [ -f $LOG_FILE_DIR ]; then
  rm $LOG_FILE_DIR
fi
touch $LOG_FILE_DIR

# install fablo if not installed
[ -f ./fablo ] ||
  curl -Lf https://github.com/hyperledger-labs/fablo/releases/download/1.1.0/fablo.sh -o ./fablo && chmod +x ./fablo

export AUTI_LOCAL_CHAIN_DIR=${PWD}

FABLO_LOCAL_CHAIN_CONFIG="local-chain-config.yaml"

# 1k test
clean_up
#./fablo generate $FABLO_ORG_CHAIN_CONFIG
#./script/replace_port.sh ./fablo-target/fabric-docker/docker-compose.yaml
./fablo up $FABLO_LOCAL_CHAIN_CONFIG
docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
sleep 5
./closc.out -phase tr -process local_submit -num 1000 -numIter 10 | tee -a $LOG_FILE_DIR
sleep 5

for i in 10000 100000 1000000; do
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    clean_up
    #    ./fablo generate $FABLO_ORG_CHAIN_CONFIG
    #    ./script/replace_port.sh ./fablo-target/fabric-docker/docker-compose.yaml
    ./fablo up $FABLO_LOCAL_CHAIN_CONFIG
    docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
    docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
    sleep 5
    ./closc.out -phase tr -process local_submit -num $i -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
  done
done

clean_up
rm closc.out
