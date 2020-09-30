#!/bin/sh

# Exits if any command fails
set -e

NETWORK_ID=${NETWORK_ID:-420}
PORT=${PORT:-8545}
TARGET_GAS_LIMIT=${TARGET_GAS_LIMIT:-8000000}

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
      --verbosity "6"
fi
