#!/usr/bin/env bash

function clean_up() {
  echo "Clean up"
  docker ps -aq | xargs docker stop | xargs docker rm
  ./fablo prune
  rm -rf fablo-target
  docker volume prune -f
  docker network prune -f
  docker container prune -f
  rm -rf wallet
  rm -rf audWallet
  rm -rf orgWallet
  rm -rf keystore
  rm lc_tx_id.log
  rm oc_tx_id.log
  rm ac_tx_id.log
}

clean_up
