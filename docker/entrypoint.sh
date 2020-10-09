#!/bin/sh

## Passed in from environment variables:
HOSTNAME=${HOSTNAME:-0.0.0.0}
PORT=${PORT:-8545}
NETWORK_ID=${NETWORK_ID:-420}
VOLUME_PATH=${VOLUME_PATH:-/mnt/l2geth}

CLEAR_DATA_FILE_PATH="${VOLUME_PATH}/.clear_data_key_${CLEAR_DATA_KEY}"
TARGET_GAS_LIMIT=${TARGET_GAS_LIMIT:-8000000}
TX_INGESTION=${TX_INGESTION:-false}
TX_INGESTION_DB_HOST=${TX_INGESTION_DB_HOST:-localhost}
TX_INGESTION_DB_PASSWORD=${TX_INGESTION_DB_PASSWORD:-test}
TX_INGESTION_DB_USER=${TX_INGESTION_DB_USER:-test}
TX_INGESTION_POLL_INTERVAL=${TX_INGESTION_POLL_INTERVAL:-3s}

if [[ -n "$CLEAR_DATA_KEY" && ! -f "$CLEAR_DATA_FILE_PATH" ]]; then
  echo "Detected change in CLEAR_DATA_KEY. Purging data."
  rm -rf ${VOLUME_PATH}/*
  rm -rf ${VOLUME_PATH}/.clear_data_key_*
  echo "Local data cleared from '${VOLUME_PATH}/*'"
  echo "Contents of volume dir: $(ls -alh $VOLUME_PATH)"
  touch $CLEAR_DATA_FILE_PATH
fi

echo "Starting Geth..."
## Command to kick off geth
geth --dev \
  --datadir $VOLUME_PATH \
  --rpc \
  --rpcaddr $HOSTNAME \
  --rpcvhosts='*' \
  --rpccorsdomain='*' \
  --rpcport $PORT \
  --networkid $NETWORK_ID \
  --ipcdisable \
  --rpcapi 'eth,net' \
  --gasprice '0' \
  --targetgaslimit $TARGET_GAS_LIMIT \
  --nousb \
  --gcmode=archive \
  --verbosity "6" \
  --txingestion.enable=$TX_INGESTION \
  --txingestion.dbhost=$TX_INGESTION_DB_HOST \
  --txingestion.pollinterval=$TX_INGESTION_POLL_INTERVAL \
  --txingestion.dbuser=$TX_INGESTION_DB_USER \
  --txingestion.dbpassword=$TX_INGESTION_DB_PASSWORD
