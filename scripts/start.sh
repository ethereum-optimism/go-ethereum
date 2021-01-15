#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"
REPO=$DIR/..

IS_VERIFIER=
DATADIR=$HOME/.ethereum
L1_CROSS_DOMAIN_MESSENGER_ADDRESS=0x0000000000000000000000000000000000000000
ADDRESS_MANAGER_OWNER_ADDRESS=0x0000000000000000000000000000000000000000
ADDRESS_RESOLVER_ADDRESS=0x0000000000000000000000000000000000000000
ETH1_NETWORK_ID=31337
ETH1_CHAIN_ID=31337
ETH1_CTC_DEPLOYMENT_HEIGHT=8
ETH1_HTTP=http://localhost:9545
TARGET_GAS_LIMIT=9000000

USAGE="
Start the Sequencer or Verifier with most configuration pre-set.

CLI Arguments:
  -h|--help                              - help message
  -v|--verifier                          - start in verifier mode
  --eth1.http                            - eth1 http endpoint
  --eth1.networkid                       - eth1 network id
  --eth1.chainid                         - eth1 chain id
  --eth1.ctcdeploymentheight             - eth1 ctc deploy height
  --eth1.l1crossdomainmessengeraddress   - eth1 l1 xdomain messenger address
  --eth1.addressresolveraddress          - eth1 address resolver address
  --eth1.ctcdeploymentheight             - eth1 ctc deployment height
"

while (( "$#" )); do
    case "$1" in
        -h|--help)
            echo "$USAGE"
            exit 0
            ;;
        -v|--verifier)
            IS_VERIFIER=true
            shift 1
            ;;
        --datadir)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                DATADIR="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.http)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ETH1_HTTP="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.networkid)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ETH1_NETWORK_ID="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.chainid)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ETH1_CHAIN_ID="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.ctcdeploymentheight)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ETH1_CTC_DEPLOYMENT_HEIGHT="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.l1crossdomainmessengeraddress)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                L1_CROSS_DOMAIN_MESSENGER_ADDRESS="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.addressresolveraddress)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ADDRESS_RESOLVER_ADDRESS="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --eth1.ctcdeploymentheight)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                ADDRESS_MANAGER_OWNER_ADDRESS="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        --targetgaslimit)
            if [ -n "$2" ] && [ ${2:0:1} != "-" ]; then
                TARGET_GASLIMIT="$2"
                shift 2
            else
                echo "Error: Argument for $1 is missing" >&2
                exit 1
            fi
            ;;
        *)
            echo "Unknown argument $1" >&2
            shift
            ;;
    esac
done

CHAIN_ID=10
NETWORK_ID=10
ADDRESS_RESOLVER_ADDRESS=0x1De8CFD4C1A486200286073aE91DE6e8099519f1
ETH1_CHAIN_ID=1
ETH1_CONFIRMATION_DEPTH=20
ETH1_CTC_DEPLOYMENT_HEIGHT=11650235
ETH1_HTTP=https://eth-mainnet.alchemyapi.io/v2/FNlLpiCGqIUoxhz1ilNk5dOCN39qgDJF
ETH1_L1_CROSS_DOMAIN_MESSENGER_ADDRESS=0x0AEBf5161A9b57349747D078c6763a0B1d67D888
ETH1_NETWORK_ID=1
L1_NODE_WEB3_URL=https://eth-mainnet.alchemyapi.io/v2/FNlLpiCGqIUoxhz1ilNk5dOCN39qgDJF
ADDRESS_MANAGER_OWNER_ADDRESS=0xc6Dbc2DC7649c7d4292d955DA08A7C21a21e1528
ROLLUP_STATE_DUMP_PATH=https://raw.githubusercontent.com/ethereum-optimism/regenesis/master/mainnet/1.json
TARGET_GAS_LIMIT=9000000

cmd="$REPO/build/bin/geth"
cmd="$cmd --eth1.syncservice"
cmd="$cmd --datadir $HOME/.ethereum"
cmd="$cmd --eth1.http $ETH1_HTTP"
cmd="$cmd --eth1.confirmationdepth $ETH1_CONFIRMATION_DEPTH"
cmd="$cmd --eth1.networkid $ETH1_NETWORK_ID"
cmd="$cmd --eth1.chainid $ETH1_CHAIN_ID"
cmd="$cmd --eth1.l1crossdomainmessengeraddress $L1_CROSS_DOMAIN_MESSENGER_ADDRESS"
cmd="$cmd --eth1.addressresolveraddress $ADDRESS_RESOLVER_ADDRESS"
cmd="$cmd --rollup.addressmanagerowneraddress $ADDRESS_MANAGER_OWNER_ADDRESS"
cmd="$cmd --rollup.statedumppath $ROLLUP_STATE_DUMP_PATH"
cmd="$cmd --eth1.ctcdeploymentheight $ETH1_CTC_DEPLOYMENT_HEIGHT"
cmd="$cmd --rpc"
cmd="$cmd --dev"
cmd="$cmd --networkid $NETWORK_ID"
cmd="$cmd --chainid $CHAIN_ID"
cmd="$cmd --rpcaddr 0.0.0.0"
cmd="$cmd --rpccorsdomain '*'"
cmd="$cmd --wsaddr 0.0.0.0"
cmd="$cmd --wsport 8546"
cmd="$cmd --wsorigins '*'"
cmd="$cmd --networkid 420"
cmd="$cmd --rpcapi 'eth,net,rollup,web3'"
cmd="$cmd --gasprice '0'"
cmd="$cmd --nousb"
cmd="$cmd --gcmode=archive"
cmd="$cmd --ipcdisable"

if [ ! -z "$IS_VERIFIER" ]; then
    cmd="$cmd --rollup.verifier"
fi

echo -e "Running:\n$cmd"
eval env TARGET_GAS_LIMIT=$TARGET_GAS_LIMIT USING_OVM=true $cmd --verbosity=6
