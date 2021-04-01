#!/bin/bash
./start.sh --rollup.clienthttp http://192.168.1.90:7878 --datadir ./db-verifier1 --rollup.statedumppath https://raw.githubusercontent.com/ethereum-optimism/regenesis/temp/snxPhase1GenesisExperiments/mainnet/only_EM_SimMsg.json --verifier --rpcport 8311 --eth1.l1crossdomainmessengeraddress 0xfBE93ba0a2Df92A8e8D40cE00acCF9248a6Fc812
