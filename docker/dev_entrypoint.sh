#!/bin/sh

# Exits if any command fails
set -e

NETWORK_ID=${NETWORK_ID:-420}
PORT=${PORT:-8545}
TARGET_GAS_LIMIT=${TARGET_GAS_LIMIT:-8000000}
TX_INGESTION=${TX_INGESTION:-false}
TX_INGESTION_DB_HOST=${TX_INGESTION_DB_HOST:-localhost}
TX_INGESTION_POLL_INTERVAL=${TX_INGESTION_POLL_INTERVAL:-3s}

if [ -n "$REBUILD" ]; then
    echo -e "\nRebuilding geth\n"
    make geth
else
    echo "Starting Geth..."
    ./build/bin/geth --dev \
      --datadir $VOLUME_PATH \
      --rpc \
      --rpcaddr $HOSTNAME \
      --rpcvhosts='*' \
      --rpccorsdomain='*' \
      --rpcport $PORT \
      --networkid $NETWORK_ID \
      --rpcapi 'eth,net' \
      --gasprice '0' \
      --targetgaslimit $TARGET_GAS_LIMIT \
      --nousb \
      --gcmode=archive \
      --verbosity "6" \
      --txingestion.enable="$TX_INGESTION" \
      --txingestion.dbhost "$TX_INGESTION_DB_HOST" \
      --txingestion.pollinterval "$TX_INGESTION_POLL_INTERVAL"
fi
