#!/bin/sh

# Exits if any command fails
set -e

if [ -n "$REBUILD" ]; then
  echo -e "\n\nREBUILD env var set, rebuilding...\n\n"

  make geth
  echo -e "\n\nCode built proceeding with ./entrypoint.sh...\n\n"
else
  echo -e "\n\nREBUILD env var not set, calling ./entrypoint.sh without building...\n\n"
fi

echo "Starting Geth..."
## Command to kick off geth
TARGET_GAS_LIMIT=${TARGET_GAS_LIMIT:-4294967295}
./build/bin/geth --dev --datadir $VOLUME_PATH --rpc --rpcaddr $HOSTNAME --rpcvhosts='*' --rpcport $PORT --networkid $NETWORK_ID --rpcapi 'eth,net' --gasprice '0' --targetgaslimit $TARGET_GAS_LIMIT --nousb --gcmode=archive --verbosity "6"