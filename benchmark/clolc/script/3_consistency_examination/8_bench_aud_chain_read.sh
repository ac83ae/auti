#!/usr/bin/env bash

HOME_DIR="../.."
cd $HOME_DIR || exit

source ./script/clean_up.sh

go build -o clolc.out

LOG_DIR="logs"
if [ ! -d $LOG_DIR ]; then
  mkdir $LOG_DIR
fi
LOG_FILE_DIR="${LOG_DIR}/clolc_ce_aud_chain_read.log"
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

clean_up
#./fablo generate $FABLO_ORG_CHAIN_CONFIG
#./script/replace_port.sh ./fablo-target/fabric-docker/docker-compose.yaml
./fablo up $FABLO_AUD_CHAIN_CONFIG
docker ps -a --format '{{.Names}}' | grep '^cli' | xargs docker rm -f
docker ps -a --format '{{.Names}}' | grep '^ca' | xargs docker rm -f
TOTAL_TXS=0
sleep 5
for i in 1000 9000 90000 900000; do
  ./clolc.out -phase ce -process aud_prepare -numTXs $i | tee -a $LOG_FILE_DIR
  TOTAL_TXS=$((TOTAL_TXS + i))
  sleep 5
  for j in {1..10}; do
    echo "No: $j" >>$LOG_FILE_DIR
    ./clolc.out -phase ce -process aud_read -numTXs $TOTAL_TXS -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5

    ./clolc.out -phase ce -process aud_read_all -numTXs $TOTAL_TXS -numIter 1 | tee -a $LOG_FILE_DIR
    sleep 5
  done
  echo "Blockchain size of peer0.aud1.example.com:" >>$LOG_FILE_DIR
  docker exec -it peer0.aud1.example.com bash -c "ls -lh /var/hyperledger/production/ledgersData/chains/chains/mychannel" | tee -a $LOG_FILE_DIR
  echo "Blockchain size of peer0.aud2.example.com:" >>$LOG_FILE_DIR
  docker exec -it peer0.aud2.example.com bash -c "ls -lh /var/hyperledger/production/ledgersData/chains/chains/mychannel" | tee -a $LOG_FILE_DIR
  sleep 1
done

clean_up
rm -f $FABLO_AUD_CHAIN_CONFIG
rm clolc.out
