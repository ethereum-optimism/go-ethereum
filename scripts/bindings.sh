#!/bin/bash

# Generate bindings for contracts from the latest master branch
# of the contracts repository.
# Copyright 2020 Optimism PBC
# https://github.com/ethereum-optimism

SCRIPTS_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
BASE_DIR="$SCRIPTS_DIR/.."
TEMP_DIR=/tmp/geth-bindings
mkdir -p $TEMP_DIR

IMAGE=ethereumoptimism/go-ethereum-devtools

REPO=contracts-v2
REPO_URL=https://github.com/ethereum-optimism/$REPO
CONTRACT_PATH="artifacts"

GIT_HASH=$(git rev-parse)
HAS_IMAGE=$(docker images "$IMAGE" --format='{{.ID}}')

# Builds bindings for each of the targets
TARGETS="OVM_CanonicalTransactionChain
Lib_AddressResolver"

if [ ! command -v docker &>/dev/null ]; then
    echo "Please install docker"
    exit 1
fi

if [ ! command -v yarn &>/dev/null ]; then
    echo "Please install yarn"
    exit 1
fi

if [ -z "$HAS_IMAGE" ]; then
    docker build \
        -t $IMAGE:latest \
        -f $BASE_DIR/Dockerfile.alltools \
        --label="io.optimism.repo.git.hash=$GIT_HASH" \
        $BASE_DIR
fi

if [ ! -d $TEMP_DIR/$REPO ]; then
    git clone --depth 1 $REPO_URL $TEMP_DIR/$REPO
fi

(
    cd $TEMP_DIR/$REPO
    if [ ! -d $TEMP_DIR/$REPO/$CONTRACT_PATH ]; then
        yarn --frozen-lockfile --ignore-engines
        yarn build
    fi
)

while read FILE; do
    FILE_PATH="$TEMP_DIR/$REPO/$CONTRACT_PATH/$FILE.json"
    if [ -f "$FILE_PATH" ]; then
        cat "$FILE_PATH" \
            | docker run -i --rm stedolan/jq '.bytecode' \
            | tr -d '"' > $TEMP_DIR/$FILE-bytecode.bin

        PACKAGE=$(echo $FILE \
            | cut -d '_' -f2 \
            | tr '[:upper:]' '[:lower:]')

        mkdir -p $BASE_DIR/contracts/$PACKAGE

        cat "$FILE_PATH" \
            | docker run -i --rm stedolan/jq '.abi' \
            | docker run -i --rm \
                -v $TEMP_DIR/$FILE-bytecode.bin:/mnt/$FILE-bytecode.bin \
                --entrypoint abigen $IMAGE \
                --pkg $PACKAGE \
                --abi - --type $FILE \
                --bin /mnt/$FILE-bytecode.bin > $BASE_DIR/contracts/$PACKAGE/$FILE.go

        rm $TEMP_DIR/$FILE-bytecode.bin
    fi
done <<< "$TARGETS"
