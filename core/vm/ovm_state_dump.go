package vm

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// AbiBytesTrue represents the ABI encoding of "true" as a byte slice
var AbiBytesTrue = common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001")

// AbiBytesFalse represents the ABI encoding of "false" as a byte slice
var AbiBytesFalse = common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000000")

var ovmStateDumpJSON = []byte(`
{
    "accounts": {
        "Proxy__OVM_L2CrossDomainMessenger": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0000",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_libAddressManager",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [],
                    "name": "messageNonce",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_target",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_sender",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_messageNonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "relayMessage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "relayedMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        }
                    ],
                    "name": "resolve",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_target",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        }
                    ],
                    "name": "sendMessage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "sentMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "successfulMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "xDomainMessageSender",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "OVM_L2CrossDomainMessenger": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001",
            "code": "0x608060405234801561001057600080fd5b50600436106100885760003560e01c806382e3702d1161005b57806382e3702d146100f3578063b1b1b20914610106578063cbd4ece914610119578063ecc704281461012c57610088565b806321d800ec1461008d5780633eae0ae0146100b6578063461a4478146100cb5780636e296e45146100eb575b600080fd5b6100a061009b3660046106d5565b610141565b6040516100ad9190610800565b60405180910390f35b6100c96100c436600461067e565b610156565b005b6100de6100d93660046106ed565b6101a5565b6040516100ad91906107af565b6100de61022c565b6100a06101013660046106d5565b61023b565b6100a06101143660046106d5565b610250565b6100c9610127366004610616565b610265565b6101346103dc565b6040516100ad91906108b0565b60006020819052908152604090205460ff1681565b60606101668433856003546103e2565b9050610172818361042f565b6003805460019081019091558151602092830120600090815260029092526040909120805460ff19169091179055505050565b60055460405163bf40fac160e01b81526000916001600160a01b03169063bf40fac1906101d690859060040161080b565b60206040518083038186803b1580156101ee57600080fd5b505afa158015610202573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061022691906105f3565b92915050565b6004546001600160a01b031681565b60026020526000908152604090205460ff1681565b60016020526000908152604090205460ff1681565b61026d610495565b15156001146102975760405162461bcd60e51b815260040161028e90610869565b60405180910390fd5b60606102a5858585856103e2565b805160208083019190912060009081526001909152604090205490915060ff16156102e25760405162461bcd60e51b815260040161028e9061081e565b600480546001600160a01b0319166001600160a01b0386811691909117909155604051600091871690610316908690610754565b6000604051808303816000865af19150503d8060008114610353576040519150601f19603f3d011682016040523d82523d6000602084013e610358565b606091505b50909150506001811515141561038e578151602080840191909120600090815260019182905260409020805460ff191690911790555b60008233436040516020016103a593929190610770565b60408051601f1981840301815291815281516020928301206000908152918290529020805460ff1916600117905550505050505050565b60035481565b6060848484846040516024016103fb94939291906107c3565b60408051601f198184030181529190526020810180516001600160e01b031663cbd4ece960e01b1790529050949350505050565b6007546040516332bea07760e21b81526001600160a01b039091169063cafa81dc9061045f90859060040161080b565b600060405180830381600087803b15801561047957600080fd5b505af115801561048d573d6000803e3d6000fd5b505050505050565b60006104d56040518060400160405280601a81526020017f4f564d5f4c3143726f7373446f6d61696e4d657373656e6765720000000000008152506101a5565b6001600160a01b0316600660009054906101000a90046001600160a01b03166001600160a01b031663d20341066040518163ffffffff1660e01b8152600401602060405180830381600087803b15801561052e57600080fd5b505af1158015610542573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061056691906105f3565b6001600160a01b031614905090565b600082601f830112610585578081fd5b813567ffffffffffffffff8082111561059c578283fd5b604051601f8301601f1916810160200182811182821017156105bc578485fd5b6040528281529250828483016020018610156105d757600080fd5b8260208601602083013760006020848301015250505092915050565b600060208284031215610604578081fd5b815161060f816108e9565b9392505050565b6000806000806080858703121561062b578283fd5b8435610636816108e9565b93506020850135610646816108e9565b9250604085013567ffffffffffffffff811115610661578283fd5b61066d87828801610575565b949793965093946060013593505050565b600080600060608486031215610692578283fd5b833561069d816108e9565b9250602084013567ffffffffffffffff8111156106b8578283fd5b6106c486828701610575565b925050604084013590509250925092565b6000602082840312156106e6578081fd5b5035919050565b6000602082840312156106fe578081fd5b813567ffffffffffffffff811115610714578182fd5b61072084828501610575565b949350505050565b600081518084526107408160208601602086016108b9565b601f01601f19169290920160200192915050565b600082516107668184602087016108b9565b9190910192915050565b600084516107828184602089016108b9565b60609490941b6bffffffffffffffffffffffff191691909301908152601481019190915260340192915050565b6001600160a01b0391909116815260200190565b6001600160a01b038581168252841660208201526080604082018190526000906107ef90830185610728565b905082606083015295945050505050565b901515815260200190565b60006020825261060f6020830184610728565b6020808252602b908201527f50726f7669646564206d6573736167652068617320616c72656164792062656560408201526a37103932b1b2b4bb32b21760a91b606082015260800190565b60208082526027908201527f50726f7669646564206d65737361676520636f756c64206e6f742062652076656040820152663934b334b2b21760c91b606082015260800190565b90815260200190565b60005b838110156108d45781810151838201526020016108bc565b838111156108e3576000848401525b50505050565b6001600160a01b03811681146108fe57600080fd5b5056fea2646970667358221220f540c4443bd3dc844d4f5ac53ee9b68863add5abd5c4f89db09816979493936964736f6c63430007000033",
            "codeHash": "0x0e9ddca7d10ea295a56994e24642ef62ad8cb0eadefea18acef747ffd806ba42",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000005": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0016"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_libAddressManager",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [],
                    "name": "messageNonce",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_target",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_sender",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_messageNonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "relayMessage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "relayedMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        }
                    ],
                    "name": "resolve",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_target",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        }
                    ],
                    "name": "sendMessage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "sentMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "successfulMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "xDomainMessageSender",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_DeployerWhitelist": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000002"
            },
            "abi": [
                {
                    "inputs": [],
                    "name": "enableArbitraryContractDeployment",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        },
                        {
                            "internalType": "bool",
                            "name": "_allowArbitraryDeployment",
                            "type": "bool"
                        }
                    ],
                    "name": "initialize",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_deployer",
                            "type": "address"
                        }
                    ],
                    "name": "isDeployerAllowed",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_allowed",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bool",
                            "name": "_allowArbitraryDeployment",
                            "type": "bool"
                        }
                    ],
                    "name": "setAllowArbitraryDeployment",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "name": "setOwner",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_deployer",
                            "type": "address"
                        },
                        {
                            "internalType": "bool",
                            "name": "_isWhitelisted",
                            "type": "bool"
                        }
                    ],
                    "name": "setWhitelistedDeployer",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "OVM_DeployerWhitelist": {
            "address": "0x4200000000000000000000000000000000000002",
            "code": "0x608060405234801561001057600080fd5b50600436106100625760003560e01c806308fd63221461006757806313af403514610097578063400ada75146100bd578063b1540a01146100eb578063bdc7b54f14610125578063d533887a1461012d575b600080fd5b6100956004803603604081101561007d57600080fd5b506001600160a01b038135169060200135151561014c565b005b610095600480360360208110156100ad57600080fd5b50356001600160a01b03166102f1565b610095600480360360408110156100d357600080fd5b506001600160a01b038135169060200135151561045d565b6101116004803603602081101561010157600080fd5b50356001600160a01b03166105ee565b604080519115158252519081900360200190f35b610095610711565b6100956004803603602081101561014357600080fd5b50351515610824565b604080516303daa95960e01b815260116004820152905133916000916101c69184916303daa9599160248082019260209290919082900301818887803b15801561019557600080fd5b505af11580156101a9573d6000803e3d6000fd5b505050506040513d60208110156101bf57600080fd5b5051610935565b9050806001600160a01b0316826001600160a01b031663735090646040518163ffffffff1660e01b815260040160206040518083038186803b15801561020b57600080fd5b505afa15801561021f573d6000803e3d6000fd5b505050506040513d602081101561023557600080fd5b50516001600160a01b03161461027c5760405162461bcd60e51b815260040180806020018281038252603a81526020018061096d603a913960400191505060405180910390fd5b33806322bd64c061028c87610938565b6102958761094d565b6040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b1580156102d257600080fd5b505af11580156102e6573d6000803e3d6000fd5b505050505050505050565b604080516303daa95960e01b8152601160048201529051339160009161033a9184916303daa9599160248082019260209290919082900301818887803b15801561019557600080fd5b9050806001600160a01b0316826001600160a01b031663735090646040518163ffffffff1660e01b815260040160206040518083038186803b15801561037f57600080fd5b505afa158015610393573d6000803e3d6000fd5b505050506040513d60208110156103a957600080fd5b50516001600160a01b0316146103f05760405162461bcd60e51b815260040180806020018281038252603a81526020018061096d603a913960400191505060405180910390fd5b33806322bd64c0601161040287610938565b6040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b15801561043f57600080fd5b505af1158015610453573d6000803e3d6000fd5b5050505050505050565b604080516303daa95960e01b815260106004820152905133916000916104d79184916303daa9599160248082019260209290919082900301818887803b1580156104a657600080fd5b505af11580156104ba573d6000803e3d6000fd5b505050506040513d60208110156104d057600080fd5b5051610967565b9050600181151514156104eb5750506105ea565b6001600160a01b0382166322bd64c06010610506600161094d565b6040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b15801561054357600080fd5b505af1158015610557573d6000803e3d6000fd5b50505050816001600160a01b03166322bd64c0601160001b61057887610938565b6040518363ffffffff1660e01b81526004018083815260200182815260200192505050600060405180830381600087803b1580156105b557600080fd5b505af11580156105c9573d6000803e3d6000fd5b50505050816001600160a01b03166322bd64c0601260001b6104028661094d565b5050565b604080516303daa95960e01b8152601060048201529051600091339183916106389184916303daa95991602480830192602092919082900301818887803b1580156104a657600080fd5b90508061064a5760019250505061070c565b6000610699836001600160a01b03166303daa959601260001b6040518263ffffffff1660e01b815260040180828152602001915050602060405180830381600087803b1580156104a657600080fd5b9050600181151514156106b2576001935050505061070c565b6000610705846001600160a01b03166303daa9596106cf89610938565b6040518263ffffffff1660e01b815260040180828152602001915050602060405180830381600087803b1580156104a657600080fd5b9450505050505b919050565b604080516303daa95960e01b8152601160048201529051339160009161075a9184916303daa9599160248082019260209290919082900301818887803b15801561019557600080fd5b9050806001600160a01b0316826001600160a01b031663735090646040518163ffffffff1660e01b815260040160206040518083038186803b15801561079f57600080fd5b505afa1580156107b3573d6000803e3d6000fd5b505050506040513d60208110156107c957600080fd5b50516001600160a01b0316146108105760405162461bcd60e51b815260040180806020018281038252603a81526020018061096d603a913960400191505060405180910390fd5b61081a6001610824565b6105ea60006102f1565b604080516303daa95960e01b8152601160048201529051339160009161086d9184916303daa9599160248082019260209290919082900301818887803b15801561019557600080fd5b9050806001600160a01b0316826001600160a01b031663735090646040518163ffffffff1660e01b815260040160206040518083038186803b1580156108b257600080fd5b505afa1580156108c6573d6000803e3d6000fd5b505050506040513d60208110156108dc57600080fd5b50516001600160a01b0316146109235760405162461bcd60e51b815260040180806020018281038252603a81526020018061096d603a913960400191505060405180910390fd5b33806322bd64c060126104028761094d565b90565b60601b6bffffffffffffffffffffffff191690565b60008161095b57600061095e565b60015b60ff1692915050565b15159056fe46756e6374696f6e2063616e206f6e6c792062652063616c6c656420627920746865206f776e6572206f66207468697320636f6e74726163742ea2646970667358221220f88b466c9bff2f68243161239393079307ce1a15f39a05357b3933f040c59f8564736f6c63430007000033",
            "codeHash": "0x91a1b614895e677b10d05e7625315510054ee4d782e6057b31252314e00c3449",
            "storage": {},
            "abi": [
                {
                    "inputs": [],
                    "name": "enableArbitraryContractDeployment",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        },
                        {
                            "internalType": "bool",
                            "name": "_allowArbitraryDeployment",
                            "type": "bool"
                        }
                    ],
                    "name": "initialize",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_deployer",
                            "type": "address"
                        }
                    ],
                    "name": "isDeployerAllowed",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_allowed",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bool",
                            "name": "_allowArbitraryDeployment",
                            "type": "bool"
                        }
                    ],
                    "name": "setAllowArbitraryDeployment",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "name": "setOwner",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_deployer",
                            "type": "address"
                        },
                        {
                            "internalType": "bool",
                            "name": "_isWhitelisted",
                            "type": "bool"
                        }
                    ],
                    "name": "setWhitelistedDeployer",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_L1MessageSender": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0004",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000001"
            },
            "abi": [
                {
                    "inputs": [],
                    "name": "getL1MessageSender",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_l1MessageSender",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "OVM_L1MessageSender": {
            "address": "0x4200000000000000000000000000000000000001",
            "code": "0x6080604052348015600f57600080fd5b506004361060285760003560e01c8063d203410614602d575b600080fd5b6033604f565b604080516001600160a01b039092168252519081900360200190f35b6000336001600160a01b0316639dc9dc936040518163ffffffff1660e01b815260040160206040518083038186803b158015608957600080fd5b505afa158015609c573d6000803e3d6000fd5b505050506040513d602081101560b157600080fd5b505190509056fea26469706673582212206075956074428a4f2a41c3b53b74d80929503a23efcb1df07cf2e8fc1714b28d64736f6c63430007000033",
            "codeHash": "0xcde5075c99d4e01c58dbf3e6d0b890e800511b53f02667a24af09942420dcefd",
            "storage": {},
            "abi": [
                {
                    "inputs": [],
                    "name": "getL1MessageSender",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_l1MessageSender",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_L2ToL1MessagePasser": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0006",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000000"
            },
            "abi": [
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": false,
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        },
                        {
                            "indexed": false,
                            "internalType": "address",
                            "name": "_sender",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "L2ToL1Message",
                    "type": "event"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        }
                    ],
                    "name": "passMessageToL1",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "sentMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "OVM_L2ToL1MessagePasser": {
            "address": "0x4200000000000000000000000000000000000000",
            "code": "0x608060405234801561001057600080fd5b50600436106100365760003560e01c806382e3702d1461003b578063cafa81dc1461006c575b600080fd5b6100586004803603602081101561005157600080fd5b5035610114565b604080519115158252519081900360200190f35b6101126004803603602081101561008257600080fd5b81019060208101813564010000000081111561009d57600080fd5b8201836020820111156100af57600080fd5b803590602001918460018302840111640100000000831117156100d157600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610129945050505050565b005b60006020819052908152604090205460ff1681565b600160008083336040516020018083805190602001908083835b602083106101625780518252601f199092019160209182019101610143565b6001836020036101000a038019825116818451168082178552505050505050905001826001600160a01b031660601b81526014019250505060405160208183030381529060405280519060200120815260200190815260200160002060006101000a81548160ff0219169083151502179055505056fea26469706673582212208a6869386aa940bc9caa7f2e3d1a2d06fc6f8f4cac9bf5eb80bbebf931d13fe264736f6c63430007000033",
            "codeHash": "0xa7da5d304c63884dbb7f11c495bf4cfb0ad6f642ebc3960cb9f072c45230347d",
            "storage": {},
            "abi": [
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": false,
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        },
                        {
                            "indexed": false,
                            "internalType": "address",
                            "name": "_sender",
                            "type": "address"
                        },
                        {
                            "indexed": false,
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "L2ToL1Message",
                    "type": "event"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_message",
                            "type": "bytes"
                        }
                    ],
                    "name": "passMessageToL1",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "",
                            "type": "bytes32"
                        }
                    ],
                    "name": "sentMessages",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_SafetyChecker": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0008",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0009"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "isBytecodeSafe",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "pure",
                    "type": "function"
                }
            ]
        },
        "OVM_SafetyChecker": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0009",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063a44eb59a14610030575b600080fd5b6100d66004803603602081101561004657600080fd5b81019060208101813564010000000081111561006157600080fd5b82018360208201111561007357600080fd5b8035906020019184600183028401116401000000008311171561009557600080fd5b91908080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152509295506100ea945050505050565b604080519115158252519081900360200190f35b60006100f4610345565b5060408051610100810182527e0101010101010101010101000000000101010101010101010101010101000081526b010101010101000000010100600160f81b016020808301919091526f0101010100000001010101010000000092820192909252630203040560e01b60608201527f0101010101010101010101010101010101010101010101010101010101010101608082015264010101010160d81b60a0820152600060c0820181905260e082015283519091741fffffffff000000000f8f000063f000013fff0ffe916a40000000000000000000026117ff60f31b039163ffffffff60601b1991870181019087015b8051600081811a880151811a82811a890151821a0182811a890151821a0182811a890151821a0182811a890151821a0182811a89015190911a01918201911a6001811b86811661032857808516610242575001605d190161032e565b808616610287575b8280600101935050825160001a915081605b141561026757610282565b6001821b851661027a57918101605e1901915b83831061024a575b610328565b8160331415610317578251602084015160e01c673350600060045af160c083901c14156102b95760088501945061030e565b817f336000905af158601d01573d60011458600c01573d6000803e3d6000fd5b60011480156102eb575080636000f35b145b156102fb5760248501945061030e565b60009a5050505050505050505050610340565b5050505061032e565b600098505050505050505050610340565b50506001015b8181106101e657600196505050505050505b919050565b604051806101000160405280600890602082028036833750919291505056fea2646970667358221220ce9ea19665ed5234a280c259228a7c5faec71cf889f6ea42656ce79736acb1f164736f6c63430007000033",
            "codeHash": "0xb2c05d3d9d991322d560d27b3aef1d1ea6d95eccb13d87aa25c475a61e92efac",
            "storage": {},
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "isBytecodeSafe",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "pure",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_ExecutionManager": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000a",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000b"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_libAddressManager",
                            "type": "address"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "minTransactionGasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "maxTransactionGasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "maxGasPerQueuePerEpoch",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "secondsPerEpoch",
                                    "type": "uint256"
                                }
                            ],
                            "internalType": "struct iOVM_ExecutionManager.GasMeterConfig",
                            "name": "_gasMeterConfig",
                            "type": "tuple"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "ovmCHAINID",
                                    "type": "uint256"
                                }
                            ],
                            "internalType": "struct iOVM_ExecutionManager.GlobalContext",
                            "name": "_globalContext",
                            "type": "tuple"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [],
                    "name": "getMaxTransactionGasLimit",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_maxTransactionGasLimit",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmADDRESS",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_ADDRESS",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmCALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmCALLER",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_CALLER",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmCHAINID",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_CHAINID",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmCREATE",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_salt",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmCREATE2",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_messageHash",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmCREATEEOA",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmDELEGATECALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_offset",
                            "type": "uint256"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_length",
                            "type": "uint256"
                        }
                    ],
                    "name": "ovmEXTCODECOPY",
                    "outputs": [
                        {
                            "internalType": "bytes",
                            "name": "_code",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "name": "ovmEXTCODEHASH",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_EXTCODEHASH",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "name": "ovmEXTCODESIZE",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_EXTCODESIZE",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmGASLIMIT",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_GASLIMIT",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmGETNONCE",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmL1QUEUEORIGIN",
                    "outputs": [
                        {
                            "internalType": "enum Lib_OVMCodec.QueueOrigin",
                            "name": "_queueOrigin",
                            "type": "uint8"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmL1TXORIGIN",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_l1TxOrigin",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmNUMBER",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_NUMBER",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmREVERT",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "ovmSETNONCE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmSLOAD",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmSSTORE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmSTATICCALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmTIMESTAMP",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_TIMESTAMP",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        }
                    ],
                    "name": "resolve",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "timestamp",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "blockNumber",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "enum Lib_OVMCodec.QueueOrigin",
                                    "name": "l1QueueOrigin",
                                    "type": "uint8"
                                },
                                {
                                    "internalType": "address",
                                    "name": "l1TxOrigin",
                                    "type": "address"
                                },
                                {
                                    "internalType": "address",
                                    "name": "entrypoint",
                                    "type": "address"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "gasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes",
                                    "name": "data",
                                    "type": "bytes"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Transaction",
                            "name": "_transaction",
                            "type": "tuple"
                        },
                        {
                            "internalType": "address",
                            "name": "_ovmStateManager",
                            "type": "address"
                        }
                    ],
                    "name": "run",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "safeCREATE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "OVM_ExecutionManager": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000b",
            "code": "0x60806040523480156200001157600080fd5b5060043610620001b45760003560e01c8063741a33eb11620000f3578063996d79a511620000a55780639dc9dc93116200007b5780639dc9dc9314620003af578063bdbf8c3614620003b9578063c1fb2ea214620003c3578063ffe7391414620003cd57620001b4565b8063996d79a5146200037757806399ccd98b14620003815780639be3ad67146200039857620001b4565b8063741a33eb14620002db578063746c32f114620002f25780638435035b14620003185780638540661f146200032f57806385979f76146200035657806390580256146200036d57620001b4565b806322bd64c0116200016b578063461a44781162000141578063461a447814620002995780634d78009214620002b05780635a98c36114620002c75780637350906414620002d157620001b4565b806322bd64c0146200025457806324749d5c146200026b5780632a2a7adb146200028257620001b4565b806303daa95914620001b95780630da449d114620001e8578063101185a4146200020157806314aa2ff7146200021a5780631c4712a7146200024057806320160f3a146200024a575b600080fd5b620001d0620001ca366004620025f9565b620003e4565b604051620001df919062002a38565b60405180910390f35b620001ff620001f9366004620025f9565b6200042e565b005b6200020b6200045d565b604051620001df919062002a74565b620002316200022b36600462002690565b62000466565b604051620001df919062002983565b620001d0620004f6565b620001d0620004fc565b620001ff620002653660046200262b565b62000502565b620001d06200027c3660046200250e565b6200056d565b620001ff6200029336600462002690565b6200058c565b62000231620002aa36600462002690565b62000599565b620001ff620002c13660046200254c565b62000620565b620001d062000801565b6200023162000807565b620001ff620002ec3660046200264d565b62000816565b6200030962000303366004620025a0565b62000993565b604051620001df919062002a5f565b620001d0620003293660046200250e565b620009cc565b6200034662000340366004620028a1565b620009e3565b604051620001df92919062002a1b565b6200034662000367366004620028a1565b62000a60565b620001d062000abe565b6200023162000ac4565b6200023162000392366004620026cf565b62000ad3565b620001ff620003a9366004620027be565b62000b5c565b6200023162000cb4565b620001d062000cc3565b620001d062000cc9565b62000346620003de366004620028a1565b62000ce4565b6000619c4060005a90506000620003fa62000ac4565b905062000408818662000d44565b93505060005a82039050808310156200042657601080548483030190555b505050919050565b6200043862000cc9565b811162000445576200045a565b6200045a6200045362000ac4565b8262000de3565b50565b60085460ff1690565b600f5460009060ff600160a01b909104161515600114156200048e576200048e600762000e5a565b8151606402619c400160005a90506000620004a862000ac4565b90506000620004c282620004bc8462000e75565b62000f08565b9050620004d0818762000fa4565b9450505060005a8203905080831015620004265760108054848303019055505050919050565b60045490565b600b5490565b600f5460ff600160a01b90910416151560011415620005275762000527600762000e5a565b61ea6060005a905060006200053b62000ac4565b90506200054a8186866200107f565b5060005a82039050808310156200056657601080548483030190555b5050505050565b6000620005846200057e83620010f1565b62001130565b90505b919050565b6200045a60028262001134565b6000805460405163bf40fac160e01b81526001600160a01b039091169063bf40fac190620005cc90859060040162002a5f565b60206040518083038186803b158015620005e557600080fd5b505afa158015620005fa573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200058491906200252d565b3330146200062e57620007fd565b62000639826200117b565b6200064a576200064a600662000e5a565b6001546040516352275acd60e11b81526001600160a01b039091169063a44eb59a906200067c90849060040162002a5f565b60206040518083038186803b1580156200069557600080fd5b505afa158015620006aa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620006d09190620025d7565b620006e157620006e1600562000e5a565b620006ec826200120e565b6000620006f9826200127b565b90506001600160a01b038116620007165762000716600862000e5a565b600060125460ff1660088111156200072a57fe5b146200074157601254620007419060ff1662000e5a565b60606200074e826200128c565b6001546040516352275acd60e11b81529192506001600160a01b03169063a44eb59a906200078190849060040162002a5f565b60206040518083038186803b1580156200079a57600080fd5b505afa158015620007af573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620007d59190620025d7565b620007e657620007e6600562000e5a565b620007fa84838380519060200120620012a6565b50505b5050565b600a5490565b600e546001600160a01b031690565b600f5460ff600160a01b909104161515600114156200083b576200083b600762000e5a565b600060018585601b0185856040516000815260200160405260405162000865949392919062002a41565b6020604051602081039080840390855afa15801562000888573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116620008c857620008c860405180606001604052806038815260200162003402603891396200058c565b620008d3816200117b565b620008df5750620007fa565b620008ea816200120e565b600f80546001600160a01b038381166001600160a01b03198316179092556040519116906000906003602160991b019062000925906200246a565b62000931919062002983565b604051809103906000f0801580156200094e573d6000803e3d6000fd5b50600f80546001600160a01b0319166001600160a01b03851617905590506200098a83826200097d816200128c565b80519060200120620012a6565b50505050505050565b6060600082600114620009a75782620009aa565b60025b9050620009c3620009bb86620010f1565b8583620012e7565b95945050505050565b600062000584620009dd83620010f1565b62001309565b600060606201388060005a9050620009fa62002478565b5060408051606081018252600f546001600160a01b0390811682528816602082015260019181018290529062000a34828a8a8a856200130d565b95509550505060005a820390508083101562000a5557601080548483030190555b505050935093915050565b60006060620186a060005a905062000a7762002478565b5060408051606081018252600f5460ff600160a01b8204161515928201929092526001600160a01b0391821681529087166020820152600062000a34828a8a8a856200130d565b60075490565b600f546001600160a01b031690565b600f5460009060ff600160a01b9091041615156001141562000afb5762000afb600762000e5a565b8251606402619c400160005a9050600062000b1562000ac4565b9050600062000b2682888862001358565b905062000b34818862000fa4565b9450505060005a820390508083101562000b5357601080548483030190555b50505092915050565b600280546001600160a01b0319166001600160a01b038381169190911791829055604051630d15d41560e41b815291169063d15d41509062000ba390339060040162002983565b60206040518083038186803b15801562000bbc57600080fd5b505afa15801562000bd1573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000bf79190620025d7565b62000c1f5760405162461bcd60e51b815260040162000c169062002afd565b60405180910390fd5b815162000c2c90620013a2565b62000c408260a001518360400151620013f0565b62000c4b57620007fd565b62000c56826200146e565b60005a905062000c7b6003600001548460a001510384608001518560c0015162000a60565b505060005a8203905062000c94818560400151620014e3565b62000c9e6200152f565b5050600280546001600160a01b03191690555050565b600d546001600160a01b031690565b60095490565b600062000cdf62000cd962000ac4565b62000e75565b905090565b60006060619c4060005a905062000cfa62002478565b5060408051606081018252600e546001600160a01b039081168252600f549081166020830152600160a01b900460ff16151591810191909152600062000a34828a8a8a856200130d565b600062000d52838362001592565b600254604051631aaf392f60e01b81526001600160a01b0390911690631aaf392f9062000d869086908690600401620029bb565b60206040518083038186803b15801562000d9f57600080fd5b505afa15801562000db4573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000dda919062002612565b90505b92915050565b62000dee82620016e4565b6002546040516374855dc360e11b81526001600160a01b039091169063e90abb869062000e229085908590600401620029bb565b600060405180830381600087803b15801562000e3d57600080fd5b505af115801562000e52573d6000803e3d6000fd5b505050505050565b6200045a816040518060200160405280600081525062001134565b600062000e8282620017ff565b60025460405163d126199f60e01b81526001600160a01b039091169063d126199f9062000eb490859060040162002983565b60206040518083038186803b15801562000ecd57600080fd5b505afa15801562000ee2573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000584919062002612565b60408051600280825260608281019093526000929190816020015b606081526020019060019003908162000f2357905050905062000f468462001956565b8160008151811062000f5457fe5b602002602001018190525062000f6a8362001982565b8160018151811062000f7857fe5b6020026020010181905250606062000f908262001999565b9050620009c38180519060200120620019c2565b600062000fcb62000fb462000ac4565b62000fc262000cd962000ac4565b60010162000de3565b62000fd562002478565b5060408051606081018252600f5460ff600160a01b8204161515928201929092526001600160a01b039182168152908416602082015260006200105b825a30888860405160240162001029929190620029f5565b60408051601f198184030181529190526020810180516001600160e01b03166326bc004960e11b1790526000620019ce565b506012805460ff1916905590508062001076576000620009c3565b50929392505050565b6200108b838362001ba9565b600254604051635c17d62960e01b81526001600160a01b0390911690635c17d62990620010c190869086908690600401620029d4565b600060405180830381600087803b158015620010dc57600080fd5b505af11580156200098a573d6000803e3d6000fd5b6000620010fe82620017ff565b600254604051637c8ee70360e01b81526001600160a01b0390911690637c8ee70390620005cc90859060040162002983565b3f90565b333b15801562001163576012805484919060ff191660018360088111156200115857fe5b021790555060016000f35b606062001171848462001d65565b9050805160208201fd5b60006200118882620017ff565b6002546040516307a1294560e01b81526001600160a01b03909116906307a1294590620011ba90859060040162002983565b60206040518083038186803b158015620011d357600080fd5b505afa158015620011e8573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620005849190620025d7565b6200121981620017ff565b600254604051637e78a4d160e11b81526001600160a01b039091169063fcf149a2906200124b90849060040162002983565b600060405180830381600087803b1580156200126657600080fd5b505af115801562000566573d6000803e3d6000fd5b60008151602083016000f092915050565b606062000584826000620012a08562001309565b620012e7565b620012b183620016e4565b6002546040516368510af960e11b81526001600160a01b039091169063d0a215f290620010c19086908690869060040162002997565b60606040519050602082018101604052818152818360208301863c9392505050565b3b90565b6000606060006064866001600160a01b03161062001336576200133086620010f1565b62001338565b855b9050620013498888838888620019ce565b92509250509550959350505050565b60008060ff60f81b858486805190602001206040516020016200137f94939291906200292c565b604051602081830303815290604052805190602001209050620009c381620019c2565b600654620013b1600062001e30565b0181106200045a57620013c660008262001e5e565b620013de6003620013d8600162001e30565b62001e5e565b6200045a6004620013d8600262001e30565b600454600090831115620014075750600062000ddd565b6003548310156200141b5750600062000ddd565b600080808460018111156200142c57fe5b141562001440575060019050600362001448565b506002905060045b60055485620014578362001e30565b620014628562001e30565b03011095945050505050565b80516009556020810151600a5560a0810151600c5560408101516008805460ff1916600183818111156200149e57fe5b02179055506060810151600d80546001600160a01b0319166001600160a01b03909216919091179055600554600b5560a0810151620014dd9062001e8b565b60115550565b600080826001811115620014f357fe5b1415620015035750600162001507565b5060025b6010546003546200152a9183918690620015218462001e30565b01010362001e5e565b505050565b600d80546001600160a01b031990811690915560006009819055600a819055600b819055600c8190556008805460ff199081169091556010829055600e8054909316909255600f80546001600160a81b0319169055601155601280549091169055565b6175305a1015620015a957620015a9600162000e5a565b600254604051630ad2267960e01b81526001600160a01b0390911690630ad2267990620015dd9085908590600401620029bb565b602060405180830381600087803b158015620015f857600080fd5b505af11580156200160d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620016339190620025d7565b620016445762001644600462000e5a565b600254604051632bcdee1960e21b81526000916001600160a01b03169063af37b86490620016799086908690600401620029bb565b602060405180830381600087803b1580156200169457600080fd5b505af1158015620016a9573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620016cf9190620025d7565b9050806200152a576200152a614e2062001ea0565b60025460405163011b1f7960e41b81526000916001600160a01b0316906311b1f790906200171790859060040162002983565b602060405180830381600087803b1580156200173257600080fd5b505af115801562001747573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200176d9190620025d7565b905080620007fd57600260009054906101000a90046001600160a01b03166001600160a01b03166333f943056040518163ffffffff1660e01b8152600401600060405180830381600087803b158015620017c657600080fd5b505af1158015620017db573d6000803e3d6000fd5b50505050620007fd6175306064620017f7620009dd86620010f1565b020162001ea0565b6175305a1015620018165762001816600162000e5a565b60025460405163c8e40fbf60e01b81526001600160a01b039091169063c8e40fbf906200184890849060040162002983565b60206040518083038186803b1580156200186157600080fd5b505afa15801562001876573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200189c9190620025d7565b620018ad57620018ad600462000e5a565b600254604051633ecdecc760e21b81526000916001600160a01b03169063fb37b31c90620018e090859060040162002983565b602060405180830381600087803b158015620018fb57600080fd5b505af115801562001910573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620019369190620025d7565b905080620007fd57620007fd6175306064620017f7620009dd86620010f1565b60408051600560a21b83186014820152603481019091526060906200197b8162001ec3565b9392505050565b606062000584620019938362001f14565b62001ec3565b606080620019a78362002028565b90506200197b620019bb825160c062002136565b8262002294565b6001600160a01b031690565b60006060620019dc62002478565b5060408051606081018252600e546001600160a01b039081168252600f549081166020830152600160a01b900460ff1615159181019190915262001a21818962002315565b601154600062001a318962001e8b565b90508060116000018190555060006060896001600160a01b03168b8a60405162001a5c919062002965565b60006040518083038160008787f1925050503d806000811462001a9c576040519150601f19603f3d011682016040523d82523d6000602084013e62001aa1565b606091505b509150915062001ab28c8662002315565b6011548262001b91576000806000606062001acd86620023cc565b92965090945092509050600484600881111562001ae657fe5b141562001af85762001af88462000e5a565b600784600881111562001b0757fe5b14801562001b1357508c155b1562001b245762001b248462000e5a565b600284600881111562001b3357fe5b148062001b4c5750600584600881111562001b4a57fe5b145b1562001b585760108290555b600284600881111562001b6757fe5b141562001b775780955062001b8a565b6040518060200160405280600081525095505b5090925050505b90920390920360115590999098509650505050505050565b6175305a101562001bc05762001bc0600162000e5a565b600254604051630ad2267960e01b81526001600160a01b0390911690630ad226799062001bf49085908590600401620029bb565b602060405180830381600087803b15801562001c0f57600080fd5b505af115801562001c24573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001c4a9190620025d7565b62001c5b5762001c5b600462000e5a565b60025460405163af3dc01160e01b81526000916001600160a01b03169063af3dc0119062001c909086908690600401620029bb565b602060405180830381600087803b15801562001cab57600080fd5b505af115801562001cc0573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062001ce69190620025d7565b9050806200152a57600260009054906101000a90046001600160a01b03166001600160a01b031663c3fd9b256040518163ffffffff1660e01b8152600401600060405180830381600087803b15801562001d3f57600080fd5b505af115801562001d54573d6000803e3d6000fd5b505050506200152a614e2062001ea0565b6060600183600881111562001d7657fe5b148062001d8f5750600883600881111562001d8d57fe5b145b1562001dab575060408051602081019091526000815262000ddd565b600483600881111562001dba57fe5b141562001dfc5760408051602080820183526000808352925162001de5938793909283920162002a89565b604051602081830303815290604052905062000ddd565b60115460105460405162001e199286929091869060200162002acb565b604051602081830303815290604052905092915050565b6000620005847306a506a506a506a506a506a506a506a506a506a583600481111562001e5857fe5b62000d44565b620007fd7306a506a506a506a506a506a506a506a506a506a583600481111562001e8457fe5b836200107f565b60005a821062001e9c575a62000584565b5090565b60115481111562001eb75762001eb7600362000e5a565b60118054919091039055565b6060808251600114801562001eed575060808360008151811062001ee357fe5b016020015160f81c105b1562001efb57508162000584565b62000dda62001f0d8451608062002136565b8462002294565b60408051602080825281830190925260609182919060208201818036833701905050905082602082015260005b602081101562001f7c5781818151811062001f5857fe5b01602001516001600160f81b0319161562001f735762001f7c565b60010162001f41565b60608160200367ffffffffffffffff8111801562001f9957600080fd5b506040519080825280601f01601f19166020018201604052801562001fc5576020820181803683370190505b50905060005b81518110156200201f57835160018401938591811062001fe757fe5b602001015160f81c60f81b82828151811062001fff57fe5b60200101906001600160f81b031916908160001a90535060010162001fcb565b50949350505050565b60608151600014156200204b575060408051600081526020810190915262000587565b6000805b835181101562002081578381815181106200206657fe5b6020026020010151518201915080806001019150506200204f565b60608267ffffffffffffffff811180156200209b57600080fd5b506040519080825280601f01601f191660200182016040528015620020c7576020820181803683370190505b50600092509050602081015b85518310156200201f576060868481518110620020ec57fe5b6020026020010151905060006020820190506200210c8382845162002424565b8785815181106200211957fe5b6020026020010151518301925050508280600101935050620020d3565b606080603884101562002193576040805160018082528183019092529060208201818036833701905050905082840160f81b816000815181106200217657fe5b60200101906001600160f81b031916908160001a90535062000dda565b600060015b808681620021a257fe5b0415620021b9576001909101906101000262002198565b8160010167ffffffffffffffff81118015620021d457600080fd5b506040519080825280601f01601f19166020018201604052801562002200576020820181803683370190505b50925084820160370160f81b836000815181106200221a57fe5b60200101906001600160f81b031916908160001a905350600190505b8181116200228b576101008183036101000a87816200225157fe5b04816200225a57fe5b0660f81b8382815181106200226b57fe5b60200101906001600160f81b031916908160001a90535060010162002236565b50509392505050565b6060806040519050835180825260208201818101602087015b81831015620022c7578051835260209283019201620022ad565b50855184518101855292509050808201602086015b81831015620022f6578051835260209283019201620022dc565b508651929092011591909101601f01601f191660405250905092915050565b805182516001600160a01b039081169116146200234e578051600e80546001600160a01b0319166001600160a01b039092169190911790555b80602001516001600160a01b031682602001516001600160a01b03161462002395576020810151600f80546001600160a01b0319166001600160a01b039092169190911790555b806040015115158260400151151514620007fd5760400151600f8054911515600160a01b0260ff60a01b1990921691909117905550565b60008060006060845160001415620023fe5750506040805160208101909152600080825260019350915081906200241d565b8480602001905181019062002414919062002716565b93509350935093505b9193509193565b8282825b6020811062002449578151835260209283019290910190601f190162002428565b905182516020929092036101000a6000190180199091169116179052505050565b6107fb8062002c0783390190565b604080516060810182526000808252602082018190529181019190915290565b803562000ddd8162002bf0565b600082601f830112620024b6578081fd5b8135620024cd620024c78262002b91565b62002b69565b9150808252836020828501011115620024e557600080fd5b8060208401602084013760009082016020015292915050565b80356002811062000ddd57600080fd5b60006020828403121562002520578081fd5b813562000dda8162002bf0565b6000602082840312156200253f578081fd5b815162000dda8162002bf0565b600080604083850312156200255f578081fd5b82356200256c8162002bf0565b9150602083013567ffffffffffffffff81111562002588578182fd5b6200259685828601620024a5565b9150509250929050565b600080600060608486031215620025b5578081fd5b8335620025c28162002bf0565b95602085013595506040909401359392505050565b600060208284031215620025e9578081fd5b8151801515811462000dda578182fd5b6000602082840312156200260b578081fd5b5035919050565b60006020828403121562002624578081fd5b5051919050565b600080604083850312156200263e578182fd5b50508035926020909101359150565b6000806000806080858703121562002663578081fd5b84359350602085013560ff811681146200267b578182fd5b93969395505050506040820135916060013590565b600060208284031215620026a2578081fd5b813567ffffffffffffffff811115620026b9578182fd5b620026c784828501620024a5565b949350505050565b60008060408385031215620026e2578182fd5b823567ffffffffffffffff811115620026f9578283fd5b6200270785828601620024a5565b95602094909401359450505050565b600080600080608085870312156200272c578182fd5b8451600981106200273b578283fd5b809450506020850151925060408501519150606085015167ffffffffffffffff81111562002767578182fd5b8501601f8101871362002778578182fd5b805162002789620024c78262002b91565b8181528860208385010111156200279e578384fd5b620027b182602083016020860162002bb6565b9598949750929550505050565b60008060408385031215620027d1578182fd5b823567ffffffffffffffff80821115620027e9578384fd5b9084019060e08287031215620027fd578384fd5b6200280960e062002b69565b8235815260208301356020820152620028268760408501620024fe565b60408201526200283a876060850162002498565b60608201526200284e876080850162002498565b608082015260a083013560a082015260c0830135828111156200286f578586fd5b6200287d88828601620024a5565b60c08301525080945050505062002898846020850162002498565b90509250929050565b600080600060608486031215620028b6578081fd5b833592506020840135620028ca8162002bf0565b9150604084013567ffffffffffffffff811115620028e6578182fd5b620028f486828701620024a5565b9150509250925092565b600081518084526200291881602086016020860162002bb6565b601f01601f19169290920160200192915050565b6001600160f81b031994909416845260609290921b6bffffffffffffffffffffffff191660018401526015830152603582015260550190565b600082516200297981846020870162002bb6565b9190910192915050565b6001600160a01b0391909116815260200190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b6001600160a01b0383168152604060208201819052600090620026c790830184620028fe565b6000831515825260406020830152620026c76040830184620028fe565b90815260200190565b93845260ff9290921660208401526040830152606082015260800190565b60006020825262000dda6020830184620028fe565b602081016002831062002a8357fe5b91905290565b600062002a968662002be5565b85825260ff8516602083015260ff841660408301526080606083015262002ac16080830184620028fe565b9695505050505050565b600062002ad88662002be5565b8582528460208301528360408301526080606083015262002ac16080830184620028fe565b60208082526046908201527f4f6e6c792061757468656e746963617465642061646472657373657320696e2060408201527f6f766d53746174654d616e616765722063616e2063616c6c20746869732066756060820152653731ba34b7b760d11b608082015260a00190565b60405181810167ffffffffffffffff8111828210171562002b8957600080fd5b604052919050565b600067ffffffffffffffff82111562002ba8578081fd5b50601f01601f191660200190565b60005b8381101562002bd357818101518382015260200162002bb9565b83811115620007fa5750506000910152565b600981106200045a57fe5b6001600160a01b03811681146200045a57600080fdfe608060405234801561001057600080fd5b506040516107fb3803806107fb8339818101604052602081101561003357600080fd5b505161003e81610044565b506101a8565b610069336000801b8360601b6001600160601b03191661006c60201b6103781760201c565b50565b604080516024810184905260448082018490528251808303909101815260649091019091526020810180516001600160e01b03908116628af59360e61b179091526100b99185916100bf16565b50505050565b60606100cc835a846100d3565b9392505050565b606060006060856001600160a01b031685856040518082805190602001908083835b602083106101145780518252601f1990920191602091820191016100f5565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038160008787f1925050503d8060008114610177576040519150601f19603f3d011682016040523d82523d6000602084013e61017c565b606091505b5090925090508161018f57805160208201fd5b80516001141561019f5760016000f35b91506100cc9050565b610644806101b76000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80630900f01014610099575b60006060610079335a61003c6100c1565b6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506100d592505050565b91509150811561008b57805160208201f35b6100953382610279565b5050005b6100bf600480360360208110156100af57600080fd5b50356001600160a01b0316610325565b005b60006100cd33826103c7565b60601c905090565b600060608061019d8787878760405160240180848152602001836001600160a01b0316815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561013657818101518382015260200161011e565b50505050905090810190601f1680156101635780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b03166001620631bb60e21b031917905294506104319350505050565b90508080602001905160408110156101b457600080fd5b8151602083018051604051929492938301929190846401000000008211156101db57600080fd5b9083019060208201858111156101f057600080fd5b825164010000000081118282018810171561020a57600080fd5b82525081516020918201929091019080838360005b8381101561023757818101518382015260200161021f565b50505050905090810190601f1680156102645780820380516001836020036101000a031916815260200191505b50604052505050925092505094509492505050565b61032082826040516024018080602001828103825283818151815260200191508051906020019080838360005b838110156102be5781810151838201526020016102a6565b50505050905090810190601f1680156102eb5780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b0316632a2a7adb60e01b1790529250610431915050565b505050565b61036c3361033233610445565b6001600160a01b03166103443361049f565b6001600160a01b0316146040518060600160405280603281526020016105dd603291396104d9565b610375816104e8565b50565b604080516024810184905260448082018490528251808303909101815260649091019091526020810180516001600160e01b0316628af59360e61b1790526103c1908490610431565b50505050565b6040805160248082018490528251808303909101815260449091019091526020810180516001600160e01b03166303daa95960e01b179052600090606090610410908590610431565b905080806020019051602081101561042757600080fd5b5051949350505050565b606061043e835a84610507565b9392505050565b6040805160048152602481019091526020810180516001600160e01b0316631cd4241960e21b17905260009060609061047f908490610431565b905080806020019051602081101561049657600080fd5b50519392505050565b6040805160048152602481019091526020810180516001600160e01b031663996d79a560e01b17905260009060609061047f908490610431565b81610320576103208382610279565b6103753360006bffffffffffffffffffffffff19606085901b16610378565b606060006060856001600160a01b031685856040518082805190602001908083835b602083106105485780518252601f199092019160209182019101610529565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038160008787f1925050503d80600081146105ab576040519150601f19603f3d011682016040523d82523d6000602084013e6105b0565b606091505b509092509050816105c357805160208201fd5b8051600114156105d35760016000f35b915061043e905056fe454f41732063616e206f6e6c792075706772616465207468656972206f776e20454f4120696d706c656d656e746174696f6ea2646970667358221220c680cc9bdbb40315bbbfe1e6943345abaa87dba44eb60d7eaaeb57236ce2403264736f6c634300070000335369676e61747572652070726f766964656420666f7220454f4120636f6e7472616374206372656174696f6e20697320696e76616c69642ea2646970667358221220f3cf5303978067c4fac9bfaae96ba11b70b6615b4e2fc1c983871ae5a756524d64736f6c63430007000033",
            "codeHash": "0x3fffb6e9dd895ccef90d54c2da5ea567636a23e4fe9c36196a6f87aceb82746a",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0016",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0008",
                "0x0000000000000000000000000000000000000000000000000000000000000004": "0x3b9aca00",
                "0x0000000000000000000000000000000000000000000000000000000000000005": "0xe8d4a51000",
                "0x0000000000000000000000000000000000000000000000000000000000000006": "0x0258",
                "0x0000000000000000000000000000000000000000000000000000000000000007": "0x01a4"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_libAddressManager",
                            "type": "address"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "minTransactionGasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "maxTransactionGasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "maxGasPerQueuePerEpoch",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "secondsPerEpoch",
                                    "type": "uint256"
                                }
                            ],
                            "internalType": "struct iOVM_ExecutionManager.GasMeterConfig",
                            "name": "_gasMeterConfig",
                            "type": "tuple"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "ovmCHAINID",
                                    "type": "uint256"
                                }
                            ],
                            "internalType": "struct iOVM_ExecutionManager.GlobalContext",
                            "name": "_globalContext",
                            "type": "tuple"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [],
                    "name": "getMaxTransactionGasLimit",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_maxTransactionGasLimit",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmADDRESS",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_ADDRESS",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmCALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmCALLER",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_CALLER",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmCHAINID",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_CHAINID",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmCREATE",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_salt",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmCREATE2",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_messageHash",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmCREATEEOA",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmDELEGATECALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_offset",
                            "type": "uint256"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_length",
                            "type": "uint256"
                        }
                    ],
                    "name": "ovmEXTCODECOPY",
                    "outputs": [
                        {
                            "internalType": "bytes",
                            "name": "_code",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "name": "ovmEXTCODEHASH",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_EXTCODEHASH",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "name": "ovmEXTCODESIZE",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_EXTCODESIZE",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmGASLIMIT",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_GASLIMIT",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmGETNONCE",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmL1QUEUEORIGIN",
                    "outputs": [
                        {
                            "internalType": "enum Lib_OVMCodec.QueueOrigin",
                            "name": "_queueOrigin",
                            "type": "uint8"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmL1TXORIGIN",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_l1TxOrigin",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmNUMBER",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_NUMBER",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmREVERT",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "ovmSETNONCE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmSLOAD",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "name": "ovmSSTORE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_calldata",
                            "type": "bytes"
                        }
                    ],
                    "name": "ovmSTATICCALL",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmTIMESTAMP",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_TIMESTAMP",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        }
                    ],
                    "name": "resolve",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "timestamp",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "blockNumber",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "enum Lib_OVMCodec.QueueOrigin",
                                    "name": "l1QueueOrigin",
                                    "type": "uint8"
                                },
                                {
                                    "internalType": "address",
                                    "name": "l1TxOrigin",
                                    "type": "address"
                                },
                                {
                                    "internalType": "address",
                                    "name": "entrypoint",
                                    "type": "address"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "gasLimit",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes",
                                    "name": "data",
                                    "type": "bytes"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Transaction",
                            "name": "_transaction",
                            "type": "tuple"
                        },
                        {
                            "internalType": "address",
                            "name": "_ovmStateManager",
                            "type": "address"
                        }
                    ],
                    "name": "run",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_bytecode",
                            "type": "bytes"
                        }
                    ],
                    "name": "safeCREATE",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_StateManager": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000c",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000d"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "commitAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountCommitted",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "commitContractStorage",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageCommitted",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_ethAddress",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_codeHash",
                            "type": "bytes32"
                        }
                    ],
                    "name": "commitPendingAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccount",
                    "outputs": [
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "nonce",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "balance",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "storageRoot",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "codeHash",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "address",
                                    "name": "ethAddress",
                                    "type": "address"
                                },
                                {
                                    "internalType": "bool",
                                    "name": "isFresh",
                                    "type": "bool"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Account",
                            "name": "_account",
                            "type": "tuple"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountEthAddress",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_ethAddress",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountNonce",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountStorageRoot",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_storageRoot",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "getContractStorage",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "getTotalUncommittedAccounts",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_total",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "getTotalUncommittedContractStorage",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_total",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "hasAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "hasContractStorage",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "hasEmptyAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "incrementTotalUncommittedAccounts",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "incrementTotalUncommittedContractStorage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "initPendingAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "isAuthenticated",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmExecutionManager",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "owner",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "nonce",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "balance",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "storageRoot",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "codeHash",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "address",
                                    "name": "ethAddress",
                                    "type": "address"
                                },
                                {
                                    "internalType": "bool",
                                    "name": "isFresh",
                                    "type": "bool"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Account",
                            "name": "_account",
                            "type": "tuple"
                        }
                    ],
                    "name": "putAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "name": "putContractStorage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "putEmptyAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "setAccountNonce",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_ovmExecutionManager",
                            "type": "address"
                        }
                    ],
                    "name": "setExecutionManager",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "testAndSetAccountChanged",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountAlreadyChanged",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "testAndSetAccountLoaded",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountAlreadyLoaded",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "testAndSetContractStorageChanged",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageAlreadyChanged",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "testAndSetContractStorageLoaded",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageAlreadyLoaded",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "OVM_StateManager": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000d",
            "code": "0x608060405234801561001057600080fd5b50600436106101c45760003560e01c806399056ba9116100f9578063d126199f11610097578063e90abb8611610071578063e90abb8614610381578063fb37b31c14610394578063fbcbc0f1146103a7578063fcf149a2146103c7576101c4565b8063d126199f14610353578063d15d415014610366578063d7bd4a2a14610379576101c4565b8063c3fd9b25116100d3578063c3fd9b2514610312578063c7650bf21461031a578063c8e40fbf1461032d578063d0a215f214610340576101c4565b806399056ba9146102e4578063af37b864146102ec578063af3dc011146102ff576101c4565b806333f94305116101665780636c87ad20116101405780636c87ad20146102a15780637c8ee703146102b65780638da5cb5b146102c95780638f3b9647146102d1576101c4565b806333f94305146102735780635c17d6291461027b5780636b18e4e81461028e576101c4565b80631381ba4d116101a25780631381ba4d14610218578063167020d21461022d5780631aaf392f1461024057806326dc5b1214610260576101c4565b806307a12945146101c95780630ad22679146101f257806311b1f79014610205575b600080fd5b6101dc6101d7366004610eb8565b6103da565b6040516101e99190611075565b60405180910390f35b6101dc610200366004610f10565b61041c565b6101dc610213366004610eb8565b610479565b61022b610226366004610eb8565b6104f3565b005b6101dc61023b366004610eb8565b610554565b61025361024e366004610f10565b610619565b6040516101e99190611080565b61025361026e366004610eb8565b6106c7565b61022b6106e6565b61022b610289366004610f3a565b610730565b61022b61029c366004610eb8565b610803565b6102a9610880565b6040516101e99190611061565b6102a96102c4366004610eb8565b61088f565b6102a96108b0565b61022b6102df366004610f79565b6108bf565b610253610975565b6101dc6102fa366004610f10565b61097b565b6101dc61030d366004610f10565b6109ee565b61022b610a44565b6101dc610328366004610f10565b610a8e565b6101dc61033b366004610eb8565b610b56565b61022b61034e366004610ed3565b610b76565b610253610361366004610eb8565b610bef565b6101dc610374366004610eb8565b610c0a565b610253610c37565b61022b61038f366004610f10565b610c3d565b6101dc6103a2366004610eb8565b610c98565b6103ba6103b5366004610eb8565b610cec565b6040516101e991906110df565b61022b6103d5366004610eb8565b610d60565b6001600160a01b0381166000908152600260205260409020600301547d4b1dc0de000000004b1dc0de000000004b1dc0de000000004b1dc0de0000145b919050565b6001600160a01b038216600090815260046020908152604080832084845290915281205460ff168061047057506001600160a01b038316600090815260026020526040902060040154600160a01b900460ff165b90505b92915050565b600080546001600160a01b031633148061049d57506001546001600160a01b031633145b6104c25760405162461bcd60e51b81526004016104b990611089565b60405180910390fd5b610473826040516020016104d69190611022565b604051602081830303815290604052805190602001206002610df6565b6000546001600160a01b031633148061051657506001546001600160a01b031633145b6105325760405162461bcd60e51b81526004016104b990611089565b600180546001600160a01b0319166001600160a01b0392909216919091179055565b600080546001600160a01b031633148061057857506001546001600160a01b031633145b6105945760405162461bcd60e51b81526004016104b990611089565b6000826040516020016105a79190611022565b60408051601f1981840301815291905280516020909101209050600260008281526005602052604090205460ff1660038111156105e057fe5b146105ef576000915050610417565b6000908152600560205260409020805460ff19166003179055505060068054600019019055600190565b6001600160a01b038216600090815260046020908152604080832084845290915281205460ff1615801561066f57506001600160a01b038316600090815260026020526040902060040154600160a01b900460ff165b1561067c57506000610473565b506001600160a01b0391909116600090815260036020908152604080832093835292905220547ffeedfacecafebeeffeedfacecafebeeffeedfacecafebeeffeedfacecafebeef1890565b6001600160a01b03166000908152600260208190526040909120015490565b6000546001600160a01b031633148061070957506001546001600160a01b031633145b6107255760405162461bcd60e51b81526004016104b990611089565b600680546001019055565b6000546001600160a01b031633148061075357506001546001600160a01b031633145b61076f5760405162461bcd60e51b81526004016104b990611089565b6001600160a01b038316600081815260036020908152604080832086845282528083207ffeedfacecafebeeffeedfacecafebeeffeedfacecafebeeffeedfacecafebeef86189055928252600481528282208583529052205460ff166107fe576001600160a01b03831660009081526004602090815260408083208584529091529020805460ff191660011790555b505050565b6000546001600160a01b031633148061082657506001546001600160a01b031633145b6108425760405162461bcd60e51b81526004016104b990611089565b6001600160a01b031660009081526002602052604090207d4b1dc0de000000004b1dc0de000000004b1dc0de000000004b1dc0de0000600390910155565b6001546001600160a01b031681565b6001600160a01b039081166000908152600260205260409020600401541690565b6000546001600160a01b031681565b6000546001600160a01b03163314806108e257506001546001600160a01b031633145b6108fe5760405162461bcd60e51b81526004016104b990611089565b6001600160a01b039182166000908152600260208181526040928390208451815590840151600182015591830151908201556060820151600382015560808201516004909101805460a0909301516001600160a01b0319909316919093161760ff60a01b1916600160a01b91151591909102179055565b60075490565b600080546001600160a01b031633148061099f57506001546001600160a01b031633145b6109bb5760405162461bcd60e51b81526004016104b990611089565b61047083836040516020016109d192919061103f565b604051602081830303815290604052805190602001206001610df6565b600080546001600160a01b0316331480610a1257506001546001600160a01b031633145b610a2e5760405162461bcd60e51b81526004016104b990611089565b61047083836040516020016104d692919061103f565b6000546001600160a01b0316331480610a6757506001546001600160a01b031633145b610a835760405162461bcd60e51b81526004016104b990611089565b600780546001019055565b600080546001600160a01b0316331480610ab257506001546001600160a01b031633145b610ace5760405162461bcd60e51b81526004016104b990611089565b60008383604051602001610ae392919061103f565b60408051601f1981840301815291905280516020909101209050600260008281526005602052604090205460ff166003811115610b1c57fe5b14610b2b576000915050610473565b6000908152600560205260409020805460ff1916600317905550506007805460001901905550600190565b6001600160a01b0316600090815260026020526040902060030154151590565b6000546001600160a01b0316331480610b9957506001546001600160a01b031633145b610bb55760405162461bcd60e51b81526004016104b990611089565b6001600160a01b0392831660009081526002602052604090206004810180546001600160a01b031916939094169290921790925560030155565b6001600160a01b031660009081526002602052604090205490565b600080546001600160a01b03838116911614806104735750506001546001600160a01b0390811691161490565b60065490565b6000546001600160a01b0316331480610c6057506001546001600160a01b031633145b610c7c5760405162461bcd60e51b81526004016104b990611089565b6001600160a01b03909116600090815260026020526040902055565b600080546001600160a01b0316331480610cbc57506001546001600160a01b031633145b610cd85760405162461bcd60e51b81526004016104b990611089565b610473826040516020016109d19190611022565b610cf4610e5c565b506001600160a01b03908116600090815260026020818152604092839020835160c08101855281548152600182015492810192909252918201549281019290925260038101546060830152600401549182166080820152600160a01b90910460ff16151560a082015290565b6000546001600160a01b0316331480610d8357506001546001600160a01b031633145b610d9f5760405162461bcd60e51b81526004016104b990611089565b6001600160a01b03166000908152600260205260409020600181557fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a4706003820155600401805460ff60a01b1916600160a01b179055565b600080826003811115610e0557fe5b60008581526005602052604090205460ff166003811115610e2257fe5b1015905080610470576000848152600560205260409020805484919060ff19166001836003811115610e5057fe5b02179055509392505050565b6040805160c081018252600080825260208201819052918101829052606081018290526080810182905260a081019190915290565b80356001600160a01b038116811461047357600080fd5b8035801515811461047357600080fd5b600060208284031215610ec9578081fd5b6104708383610e91565b600080600060608486031215610ee7578182fd5b610ef18585610e91565b9250610f008560208601610e91565b9150604084013590509250925092565b60008060408385031215610f22578182fd5b610f2c8484610e91565b946020939093013593505050565b600080600060608486031215610f4e578283fd5b83356001600160a01b0381168114610f64578384fd5b95602085013595506040909401359392505050565b60008082840360e0811215610f8c578283fd5b610f968585610e91565b925060c0601f1982011215610fa9578182fd5b5060405160c0810181811067ffffffffffffffff82111715610fc9578283fd5b8060405250602084013581526040840135602082015260608401356040820152608084013560608201526110008560a08601610e91565b60808201526110128560c08601610ea8565b60a0820152809150509250929050565b60609190911b6bffffffffffffffffffffffff1916815260140190565b60609290921b6bffffffffffffffffffffffff19168252601482015260340190565b6001600160a01b0391909116815260200190565b901515815260200190565b90815260200190565b60208082526036908201527f46756e6374696f6e2063616e206f6e6c792062652063616c6c65642062792061604082015275757468656e746963617465642061646472657373657360501b606082015260800190565b815181526020808301519082015260408083015190820152606080830151908201526080808301516001600160a01b03169082015260a09182015115159181019190915260c0019056fea2646970667358221220f23d04471ab92169b68e21b5fb1ac47d850ecdc096c103995de61e3a5b00081b64736f6c63430007000033",
            "codeHash": "0xd5c31f5f067037a667c5a398a1c053a2722d66e82d5729104e84b60809a57fb8",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000b"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "constructor"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "commitAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountCommitted",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "commitContractStorage",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageCommitted",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_ethAddress",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_codeHash",
                            "type": "bytes32"
                        }
                    ],
                    "name": "commitPendingAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccount",
                    "outputs": [
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "nonce",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "balance",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "storageRoot",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "codeHash",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "address",
                                    "name": "ethAddress",
                                    "type": "address"
                                },
                                {
                                    "internalType": "bool",
                                    "name": "isFresh",
                                    "type": "bool"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Account",
                            "name": "_account",
                            "type": "tuple"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountEthAddress",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "_ethAddress",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountNonce",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "getAccountStorageRoot",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_storageRoot",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "getContractStorage",
                    "outputs": [
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "getTotalUncommittedAccounts",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_total",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "getTotalUncommittedContractStorage",
                    "outputs": [
                        {
                            "internalType": "uint256",
                            "name": "_total",
                            "type": "uint256"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "hasAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "hasContractStorage",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "hasEmptyAccount",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_exists",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "incrementTotalUncommittedAccounts",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "incrementTotalUncommittedContractStorage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "initPendingAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "isAuthenticated",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "ovmExecutionManager",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "owner",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "components": [
                                {
                                    "internalType": "uint256",
                                    "name": "nonce",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "uint256",
                                    "name": "balance",
                                    "type": "uint256"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "storageRoot",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "bytes32",
                                    "name": "codeHash",
                                    "type": "bytes32"
                                },
                                {
                                    "internalType": "address",
                                    "name": "ethAddress",
                                    "type": "address"
                                },
                                {
                                    "internalType": "bool",
                                    "name": "isFresh",
                                    "type": "bool"
                                }
                            ],
                            "internalType": "struct Lib_OVMCodec.Account",
                            "name": "_account",
                            "type": "tuple"
                        }
                    ],
                    "name": "putAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_value",
                            "type": "bytes32"
                        }
                    ],
                    "name": "putContractStorage",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "putEmptyAccount",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        },
                        {
                            "internalType": "uint256",
                            "name": "_nonce",
                            "type": "uint256"
                        }
                    ],
                    "name": "setAccountNonce",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_ovmExecutionManager",
                            "type": "address"
                        }
                    ],
                    "name": "setExecutionManager",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "testAndSetAccountChanged",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountAlreadyChanged",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "testAndSetAccountLoaded",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasAccountAlreadyLoaded",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "testAndSetContractStorageChanged",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageAlreadyChanged",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_contract",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_key",
                            "type": "bytes32"
                        }
                    ],
                    "name": "testAndSetContractStorageLoaded",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_wasContractStorageAlreadyLoaded",
                            "type": "bool"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_ECDSAContractAccount": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000e",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000003"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_transaction",
                            "type": "bytes"
                        },
                        {
                            "internalType": "enum Lib_OVMCodec.EOASignatureType",
                            "name": "_signatureType",
                            "type": "uint8"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "execute",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "OVM_ECDSAContractAccount": {
            "address": "0x4200000000000000000000000000000000000003",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063d1be05c214610030575b600080fd5b61004361003e366004610cf7565b61005a565b604051610051929190610eac565b60405180910390f35b600060603382600188600181111561006e57fe5b1490506100bc3361007e84610194565b6001600160a01b03166100948c858c8c8c6101ed565b6001600160a01b0316146040518060600160405280603c8152602001611267603c913961025a565b6100c4610bb6565b6100ce8a8361026e565b9050610101336100dd856103c7565b8360000151146040518060600160405280603481526020016112a36034913961025a565b60608101516001600160a01b03166101595760006101288483604001518460a00151610417565b905060018160405160200161013d9190610e98565b604051602081830303815290604052955095505050505061018a565b61016a838260000151600101610481565b61018283826040015183606001518460a001516104c5565b945094505050505b9550959350505050565b6040805160048152602481019091526020810180516001600160e01b031663996d79a560e01b1790526000906060906101ce908490610536565b9050808060200190518101906101e49190610c88565b9150505b919050565b6000806101fa8787610543565b905060018186601b018686604051600081526020016040526040516102229493929190610ecf565b6020604051602081039080840390855afa158015610244573d6000803e3d6000fd5b5050604051601f19015198975050505050505050565b81610269576102698382610564565b505050565b610276610bb6565b81156102e557600080600080600060608880602001905181019061029a9190610dba565b6040805160e0810182529687526020870194909452928501939093526001600160a01b0390921660608401526000608084015260a083015260c082015296506103c195505050505050565b60606102f0846105a8565b90506040518060e0016040528061031a8360008151811061030d57fe5b60200260200101516105bb565b815260200161032f8360018151811061030d57fe5b81526020016103448360028151811061030d57fe5b81526020016103668360038151811061035957fe5b60200260200101516105c6565b6001600160a01b031681526020016103848360048151811061030d57fe5b81526020016103a68360058151811061039957fe5b602002602001015161060e565b81526020016103bb8360068151811061030d57fe5b90529150505b92915050565b6040805160048152602481019091526020810180516001600160e01b03166360fd975160e11b179052600090606090610401908490610536565b9050808060200190518101906101e49190610da2565b600060606104608585856040516024016104319190610eed565b60408051601f198184030181529190526020810180516001600160e01b03166314aa2ff760e01b179052610668565b9050808060200190518101906104769190610c88565b9150505b9392505050565b610269828260405160240161049691906111a0565b60408051601f198184030181529190526020810180516001600160e01b0316630da449d160e01b179052610536565b6000606080610512878787876040516024016104e3939291906111a9565b60408051601f198184030181529190526020810180516001600160e01b03166342cbcfbb60e11b179052610536565b9050808060200190518101906105289190610ca4565b925092505094509492505050565b606061047a835a84610668565b6000811561055b57610554836106f6565b90506103c1565b61047a83610766565b61026982826040516024016105799190610eed565b60408051601f198184030181529190526020810180516001600160e01b0316632a2a7adb60e01b179052610536565b60606103c16105b683610771565b610796565b60006103c1826108b8565b8051600090600114156105db575060006101e8565b81516015146106055760405162461bcd60e51b81526004016105fc9061108d565b60405180910390fd5b6103c1826105bb565b6060600080600061061e85610948565b91945092509050600081600181111561063357fe5b146106505760405162461bcd60e51b81526004016105fc90611132565b61065f85602001518484610b05565b95945050505050565b606060006060856001600160a01b031685856040516106879190610e5a565b60006040518083038160008787f1925050503d80600081146106c5576040519150601f19603f3d011682016040523d82523d6000602084013e6106ca565b606091505b509092509050816106dd57805160208201fd5b8051600114156106ed5760016000f35b915061047a9050565b604080518082018252601c81527f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080830191909152835184820120925160009391610747918491849101610e76565b6040516020818303038152906040528051906020012092505050919050565b805160209091012090565b610779610bfc565b506040805180820190915281518152602082810190820152919050565b60606000806107a484610948565b919350909150600190508160018111156107ba57fe5b146107d75760405162461bcd60e51b81526004016105fc9061101f565b6040805160208082526104208201909252606091816020015b6107f8610bfc565b8152602001906001900390816107f05790505090506000835b86518110156108ad576020821061083a5760405162461bcd60e51b81526004016105fc90610f9e565b6000806108666040518060400160405280858c60000151038152602001858c6020015101815250610948565b509150915060405180604001604052808383018152602001848b602001510181525085858151811061089457fe5b6020908102919091010152600193909301920101610811565b508152949350505050565b60006021826000015111156108df5760405162461bcd60e51b81526004016105fc90610fe8565b60008060006108ed85610948565b91945092509050600081600181111561090257fe5b1461091f5760405162461bcd60e51b81526004016105fc90610fe8565b60208086015184018051909184101561093e5760208490036101000a90045b9695505050505050565b6000806000808460000151116109705760405162461bcd60e51b81526004016105fc906110c4565b6020840151805160001a607f8111610995576000600160009450945094505050610afe565b60b781116109d5578551607f1982019081106109c35760405162461bcd60e51b81526004016105fc90611056565b60019550935060009250610afe915050565b60bf8111610a4f57855160b6198201908110610a035760405162461bcd60e51b81526004016105fc90610f67565b6000816020036101000a6001850151049050808201886000015111610a3a5760405162461bcd60e51b81526004016105fc906110fb565b60019091019550935060009250610afe915050565b60f78111610a8e57855160bf198201908110610a7d5760405162461bcd60e51b81526004016105fc90611169565b600195509350849250610afe915050565b855160f6198201908110610ab45760405162461bcd60e51b81526004016105fc90610f00565b6000816020036101000a6001850151049050808201886000015111610aeb5760405162461bcd60e51b81526004016105fc90610f37565b6001918201965094509250610afe915050565b9193909250565b6060808267ffffffffffffffff81118015610b1f57600080fd5b506040519080825280601f01601f191660200182016040528015610b4a576020820181803683370190505b509050805160001415610b5e57905061047a565b8484016020820160005b60208604811015610b89578251825260209283019290910190600101610b68565b5060006001602087066020036101000a039050808251168119845116178252839450505050509392505050565b6040518060e0016040528060008152602001600081526020016000815260200160006001600160a01b031681526020016000815260200160608152602001600081525090565b604051806040016040528060008152602001600081525090565b600082601f830112610c26578081fd5b8151610c39610c34826111fa565b6111d3565b9150808252836020828501011115610c5057600080fd5b610c6181602084016020860161121e565b5092915050565b8035600281106103c157600080fd5b803560ff811681146103c157600080fd5b600060208284031215610c99578081fd5b815161047a8161124e565b60008060408385031215610cb6578081fd5b82518015158114610cc5578182fd5b602084015190925067ffffffffffffffff811115610ce1578182fd5b610ced85828601610c16565b9150509250929050565b600080600080600060a08688031215610d0e578081fd5b853567ffffffffffffffff811115610d24578182fd5b8601601f81018813610d34578182fd5b8035610d42610c34826111fa565b818152896020838501011115610d56578384fd5b816020840160208301378360208383010152809750505050610d7b8760208801610c68565b9350610d8a8760408801610c77565b94979396509394606081013594506080013592915050565b600060208284031215610db3578081fd5b5051919050565b60008060008060008060c08789031215610dd2578081fd5b865195506020870151945060408701519350606087015192506080870151610df98161124e565b60a088015190925067ffffffffffffffff811115610e15578182fd5b610e2189828a01610c16565b9150509295509295509295565b60008151808452610e4681602086016020860161121e565b601f01601f19169290920160200192915050565b60008251610e6c81846020870161121e565b9190910192915050565b60008351610e8881846020880161121e565b9190910191825250602001919050565b6001600160a01b0391909116815260200190565b6000831515825260406020830152610ec76040830184610e2e565b949350505050565b93845260ff9290921660208401526040830152606082015260800190565b60006020825261047a6020830184610e2e565b6020808252601d908201527f496e76616c696420524c50206c6f6e67206c697374206c656e6774682e000000604082015260600190565b60208082526016908201527524b73b30b634b210292628103637b733903634b9ba1760511b604082015260600190565b6020808252601f908201527f496e76616c696420524c50206c6f6e6720737472696e67206c656e6774682e00604082015260600190565b6020808252602a908201527f50726f766964656420524c50206c6973742065786365656473206d6178206c6960408201526939ba103632b733ba341760b11b606082015260800190565b6020808252601a908201527f496e76616c696420524c5020627974657333322076616c75652e000000000000604082015260600190565b60208082526017908201527f496e76616c696420524c50206c6973742076616c75652e000000000000000000604082015260600190565b60208082526019908201527f496e76616c696420524c502073686f727420737472696e672e00000000000000604082015260600190565b6020808252601a908201527f496e76616c696420524c5020616464726573732076616c75652e000000000000604082015260600190565b60208082526018908201527f524c50206974656d2063616e6e6f74206265206e756c6c2e0000000000000000604082015260600190565b60208082526018908201527f496e76616c696420524c50206c6f6e6720737472696e672e0000000000000000604082015260600190565b60208082526018908201527f496e76616c696420524c502062797465732076616c75652e0000000000000000604082015260600190565b60208082526017908201527f496e76616c696420524c502073686f7274206c6973742e000000000000000000604082015260600190565b90815260200190565b8381526001600160a01b038316602082015260606040820181905260009061065f90830184610e2e565b60405181810167ffffffffffffffff811182821017156111f257600080fd5b604052919050565b600067ffffffffffffffff821115611210578081fd5b50601f01601f191660200190565b60005b83811015611239578181015183820152602001611221565b83811115611248576000848401525b50505050565b6001600160a01b038116811461126357600080fd5b5056fe5369676e61747572652070726f766964656420666f7220454f41207472616e73616374696f6e20657865637574696f6e20697320696e76616c69642e5472616e73616374696f6e206e6f6e636520646f6573206e6f74206d6174636820746865206578706563746564206e6f6e63652ea264697066735822122067c01f3f2f90427142d2f87a21e288bf686b1f9997b0815a37c3d9f66628dd1764736f6c63430007000033",
            "codeHash": "0x72b5bf696dbeca2ecd7c2f5457fdbcd10c61901fb46c17a9d779987014a3cddd",
            "storage": {},
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_transaction",
                            "type": "bytes"
                        },
                        {
                            "internalType": "enum Lib_OVMCodec.EOASignatureType",
                            "name": "_signatureType",
                            "type": "uint8"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "execute",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Proxy__OVM_SequencerEntrypoint": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0010",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000005"
            },
            "abi": [
                {
                    "stateMutability": "nonpayable",
                    "type": "fallback"
                }
            ]
        },
        "OVM_SequencerEntrypoint": {
            "address": "0x4200000000000000000000000000000000000005",
            "code": "0x608060405234801561001057600080fd5b50600061005b6100566000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525092506102cf915050565b61037a565b905060006100ad6100a86000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060019250602091506103e49050565b6104a3565b905060006100fa6100a86000368080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525060219250602091506103e49050565b905060006101416000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250604192506102cf915050565b905060606101886000368080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250604292506104aa915050565b90506000600186600181111561019a57fe5b14905060606101b16101ab846104db565b8361056e565b905060006101c28284878a8a61082f565b90506101ce33826108b1565b6101ef5760006101de8385610926565b90506101ed3382888b8b610947565b505b60608284878a8a604051602401808060200186151581526020018560ff168152602001848152602001838152602001828103825287818151815260200191508051906020019080838360005b8381101561025357818101518382015260200161023b565b50505050905090810190601f1680156102805780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b03166368df02e160e11b17905297506102c2965033955050505050505a84846109a9565b5050505050505050505050005b60008182600101101561031c576040805162461bcd60e51b815260206004820152601060248201526f746f55696e74385f6f766572666c6f7760801b604482015290519081900360640190fd5b816001018351101561036b576040805162461bcd60e51b8152602060048201526013602482015272746f55696e74385f6f75744f66426f756e647360681b604482015290519081900360640190fd5b50818101600101515b92915050565b600060ff821661038c575060006103df565b8160ff16600214156103a0575060016103df565b6103df336040518060400160405280601f81526020017f5472616e73616374696f6e2074797065206d7573742062652030206f72203200815250610b4a565b919050565b606081830184511015610433576040805162461bcd60e51b815260206004820152601260248201527152656164206f7574206f6620626f756e647360701b604482015290519081900360640190fd5b60608215801561044e57604051915060208201604052610498565b6040519150601f8416801560200281840101858101878315602002848b0101015b8183101561048757805183526020928301920161046f565b5050858452601f01601f1916604052505b5090505b9392505050565b6020015190565b606081835103600014156104cd5750604080516020810190915260008152610374565b61049c8383848651036103e4565b6104e36113a3565b6040518060e001604052806104f9846006610bf6565b62ffffff16815260200161050e846003610bf6565b62ffffff16620f4240028152602001610528846000610bf6565b62ffffff16815260200161053d846009610c9d565b6001600160a01b031681526020016000815260200161055d84601d6104aa565b81526101a460209091015292915050565b60608115610648578260000151836040015184602001518560c0015186606001518760a0015160405160200180878152602001868152602001858152602001848152602001836001600160a01b0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156105fb5781810151838201526020016105e3565b50505050905090810190601f1680156106285780820380516001836020036101000a031916815260200191505b509750505050505050506040516020818303038152906040529050610374565b6040805160098082526101408201909252606091816020015b6060815260200190600190039081610661575050845190915061068390610d4d565b8160008151811061069057fe5b60200260200101819052506106a88460200151610d4d565b816001815181106106b557fe5b60200260200101819052506106cd8460400151610d4d565b816002815181106106da57fe5b602090810291909101015260608401516001600160a01b031661072c5761070f60405180602001604052806000815250610d5b565b8160038151811061071c57fe5b6020026020010181905250610752565b6107398460600151610da4565b8160038151811061074657fe5b60200260200101819052505b61075c6000610d4d565b8160048151811061076957fe5b60200260200101819052506107818460a00151610d5b565b8160058151811061078e57fe5b60200260200101819052506107a68460c00151610d4d565b816006815181106107b357fe5b60200260200101819052506107d660405180602001604052806000815250610d5b565b816007815181106107e357fe5b602002602001018190525061080660405180602001604052806000815250610d5b565b8160088151811061081357fe5b602002602001018190525061082781610dc7565b915050610374565b60008061083c8787610926565b905060018186601b01868660405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa15801561089b573d6000803e3d6000fd5b5050604051601f19015198975050505050505050565b604080516001600160a01b0383166024808301919091528251808303909101815260449091019091526020810180516001600160e01b0316638435035b60e01b179052600090606090610905908590610dea565b905080806020019051602081101561091c57600080fd5b5051949350505050565b6000811561093e5761093783610df7565b9050610374565b61049c83610ea8565b604080516024810186905260ff851660448201526064810184905260848082018490528251808303909101815260a49091019091526020810180516001600160e01b031663741a33eb60e01b1790526109a1908690610dea565b505050505050565b6000606080610a6e8787878760405160240180848152602001836001600160a01b0316815260200180602001828103825283818151815260200191508051906020019080838360005b83811015610a0a5781810151838201526020016109f2565b50505050905090810190601f168015610a375780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b03166342cbcfbb60e11b1790529450610dea9350505050565b9050808060200190516040811015610a8557600080fd5b815160208301805160405192949293830192919084640100000000821115610aac57600080fd5b908301906020820185811115610ac157600080fd5b8251640100000000811182820188101715610adb57600080fd5b82525081516020918201929091019080838360005b83811015610b08578181015183820152602001610af0565b50505050905090810190601f168015610b355780820380516001836020036101000a031916815260200191505b50604052505050925092505094509492505050565b610bf182826040516024018080602001828103825283818151815260200191508051906020019080838360005b83811015610b8f578181015183820152602001610b77565b50505050905090810190601f168015610bbc5780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b0316632a2a7adb60e01b1790529250610dea915050565b505050565b600081826003011015610c44576040805162461bcd60e51b8152602060048201526011602482015270746f55696e7432345f6f766572666c6f7760781b604482015290519081900360640190fd5b8160030183511015610c94576040805162461bcd60e51b8152602060048201526014602482015273746f55696e7432345f6f75744f66426f756e647360601b604482015290519081900360640190fd5b50016003015190565b600081826014011015610cec576040805162461bcd60e51b8152602060048201526012602482015271746f416464726573735f6f766572666c6f7760701b604482015290519081900360640190fd5b8160140183511015610d3d576040805162461bcd60e51b8152602060048201526015602482015274746f416464726573735f6f75744f66426f756e647360581b604482015290519081900360640190fd5b500160200151600160601b900490565b6060610374610d5b83610eb3565b60608082516001148015610d835750608083600081518110610d7957fe5b016020015160f81c105b15610d8f575081610374565b61049c610d9e84516080610fbc565b8461110c565b60408051600560a21b831860148201526034810190915260609061049c81610d5b565b606080610dd383611189565b905061049c610de4825160c0610fbc565b8261110c565b606061049c835a8461128a565b604080518082018252601c8082527f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080840191825285518682012094516000959385938593929092019182918083835b60208310610e685780518252601f199092019160209182019101610e49565b51815160209384036101000a6000190180199092169116179052920193845250604080518085038152938201905282519201919091209695505050505050565b805160209091012090565b60408051602080825281830190925260609182919060208201818036833701905050905082602082015260005b6020811015610f1657818181518110610ef557fe5b01602001516001600160f81b03191615610f0e57610f16565b600101610ee0565b60608160200367ffffffffffffffff81118015610f3257600080fd5b506040519080825280601f01601f191660200182016040528015610f5d576020820181803683370190505b50905060005b8151811015610fb3578351600184019385918110610f7d57fe5b602001015160f81c60f81b828281518110610f9457fe5b60200101906001600160f81b031916908160001a905350600101610f63565b50949350505050565b6060806038841015611016576040805160018082528183019092529060208201818036833701905050905082840160f81b81600081518110610ffa57fe5b60200101906001600160f81b031916908160001a90535061049c565b600060015b80868161102457fe5b0415611039576001909101906101000261101b565b8160010167ffffffffffffffff8111801561105357600080fd5b506040519080825280601f01601f19166020018201604052801561107e576020820181803683370190505b50925084820160370160f81b8360008151811061109757fe5b60200101906001600160f81b031916908160001a905350600190505b818111611103576101008183036101000a87816110cc57fe5b04816110d457fe5b0660f81b8382815181106110e457fe5b60200101906001600160f81b031916908160001a9053506001016110b3565b50509392505050565b6060806040519050835180825260208201818101602087015b8183101561113d578051835260209283019201611125565b50855184518101855292509050808201602086015b8183101561116a578051835260209283019201611152565b508651929092011591909101601f01601f191660405250905092915050565b60608151600014156111aa57506040805160008152602081019091526103df565b6000805b83518110156111dd578381815181106111c357fe5b6020026020010151518201915080806001019150506111ae565b60608267ffffffffffffffff811180156111f657600080fd5b506040519080825280601f01601f191660200182016040528015611221576020820181803683370190505b50600092509050602081015b8551831015610fb357606086848151811061124457fe5b6020026020010151905060006020820190506112628382845161135f565b87858151811061126e57fe5b602002602001015151830192505050828060010193505061122d565b606060006060856001600160a01b031685856040518082805190602001908083835b602083106112cb5780518252601f1990920191602091820191016112ac565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038160008787f1925050503d806000811461132e576040519150601f19603f3d011682016040523d82523d6000602084013e611333565b606091505b5090925090508161134657805160208201fd5b8051600114156113565760016000f35b915061049c9050565b8282825b60208110611382578151835260209283019290910190601f1901611363565b905182516020929092036101000a6000190180199091169116179052505050565b6040518060e0016040528060008152602001600081526020016000815260200160006001600160a01b03168152602001600081526020016060815260200160008152509056fea2646970667358221220b6e25b83df4f7020239c68addc0e4e4c2b0c9888971bd1de4180bb74a08a20ea64736f6c63430007000033",
            "codeHash": "0x8ac92938dfae95fcd8c7a9b8d74dd9dc4bff4f49a7d3b2ea4206cdc537cf7254",
            "storage": {},
            "abi": [
                {
                    "stateMutability": "nonpayable",
                    "type": "fallback"
                }
            ]
        },
        "Proxy__OVM_ProxySequencerEntrypoint": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0012",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0x4200000000000000000000000000000000000004"
            },
            "abi": [
                {
                    "stateMutability": "nonpayable",
                    "type": "fallback"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_implementation",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "name": "init",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_implementation",
                            "type": "address"
                        }
                    ],
                    "name": "upgrade",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "OVM_ProxySequencerEntrypoint": {
            "address": "0x4200000000000000000000000000000000000004",
            "code": "0x608060405234801561001057600080fd5b50600436106100365760003560e01c80630900f01014610084578063f09a4016146100ac575b610080335a6100436100da565b6000368080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152506100ee92505050565b5050005b6100aa6004803603602081101561009a57600080fd5b50356001600160a01b0316610292565b005b6100aa600480360360408110156100c257600080fd5b506001600160a01b03813581169160200135166102e4565b60006100e63382610330565b60601c905090565b60006060806101b68787878760405160240180848152602001836001600160a01b0316815260200180602001828103825283818151815260200191508051906020019080838360005b8381101561014f578181015183820152602001610137565b50505050905090810190601f16801561017c5780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b03166001620631bb60e21b0319179052945061039a9350505050565b90508080602001905160408110156101cd57600080fd5b8151602083018051604051929492938301929190846401000000008211156101f457600080fd5b90830190602082018581111561020957600080fd5b825164010000000081118282018810171561022357600080fd5b82525081516020918201929091019080838360005b83811015610250578181015183820152602001610238565b50505050905090810190601f16801561027d5780820380516001836020036101000a031916815260200191505b50604052505050925092505094509492505050565b6102d83361029f336103ae565b6001600160a01b03166102b0610408565b6001600160a01b03161460405180606001604052806025815260200161063360259139610415565b6102e181610429565b50565b61031a3360006102f2610408565b6001600160a01b03161460405180606001604052806027815260200161065860279139610415565b61032381610448565b61032c82610429565b5050565b6040805160248082018490528251808303909101815260449091019091526020810180516001600160e01b03166303daa95960e01b17905260009060609061037990859061039a565b905080806020019051602081101561039057600080fd5b5051949350505050565b60606103a7835a84610467565b9392505050565b6040805160048152602481019091526020810180516001600160e01b0316631cd4241960e21b1790526000906060906103e890849061039a565b90508080602001905160208110156103ff57600080fd5b50519392505050565b60006100e6336001610330565b8161042457610424838261053c565b505050565b6102e13360006bffffffffffffffffffffffff19606085901b166105e3565b6102e13360016bffffffffffffffffffffffff19606085901b166105e3565b606060006060856001600160a01b031685856040518082805190602001908083835b602083106104a85780518252601f199092019160209182019101610489565b6001836020036101000a03801982511681845116808217855250505050505090500191505060006040518083038160008787f1925050503d806000811461050b576040519150601f19603f3d011682016040523d82523d6000602084013e610510565b606091505b5090925090508161052357805160208201fd5b8051600114156105335760016000f35b91506103a79050565b61042482826040516024018080602001828103825283818151815260200191508051906020019080838360005b83811015610581578181015183820152602001610569565b50505050905090810190601f1680156105ae5780820380516001836020036101000a031916815260200191505b5060408051601f198184030181529190526020810180516001600160e01b0316632a2a7adb60e01b179052925061039a915050565b604080516024810184905260448082018490528251808303909101815260649091019091526020810180516001600160e01b0316628af59360e61b17905261062c90849061039a565b5050505056fe4f6e6c79206f776e65722063616e20757067726164652074686520456e747279706f696e7450726f7879456e747279706f696e742068617320616c7265616479206265656e20696e69746564a2646970667358221220d461d950907c6ecd1cd1bacb1366fd46e7610cb9ca15a13e3bd1f7edd8197e8b64736f6c63430007000033",
            "codeHash": "0x0b61b00fce08a0d19fbdad141253e9e8a30b9af8ec5f8fe6c9e66fe64ed7c07b",
            "storage": {},
            "abi": [
                {
                    "stateMutability": "nonpayable",
                    "type": "fallback"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_implementation",
                            "type": "address"
                        },
                        {
                            "internalType": "address",
                            "name": "_owner",
                            "type": "address"
                        }
                    ],
                    "name": "init",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_implementation",
                            "type": "address"
                        }
                    ],
                    "name": "upgrade",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Proxy__mockOVM_ECDSAContractAccount": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0014",
            "code": "0x608060405234801561001057600080fd5b506004361061002b5760003560e01c8063776d1a0114610077575b60015460408051602036601f8101829004820283018201909352828252610075936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b005b6100756004803603602081101561008d57600080fd5b50356001600160a01b031661015d565b60006060836001600160a01b0316836040518082805190602001908083835b602083106100db5780518252601f1990920191602091820191016100bc565b6001836020036101000a0380198251168184511680821785525050505050509050019150506000604051808303816000865af19150503d806000811461013d576040519150601f19603f3d011682016040523d82523d6000602084013e610142565b606091505b5091509150811561015557805160208201f35b805160208201fd5b6000546001600160a01b031633141561019057600180546001600160a01b0319166001600160a01b0383161790556101da565b60015460408051602036601f81018290048202830182019093528282526101da936001600160a01b0316926000918190840183828082843760009201919091525061009d92505050565b5056fea2646970667358221220293887d48c4c1c34de868edf3e9a6be82327946c76d71f7c2023e67f556c6ecb64736f6c63430007000033",
            "codeHash": "0x0033b946bc1a66d1a2a7bd76e67701e9245080b0eb8e940316e638252c6551d7",
            "storage": {
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0000000000000000000000000000000000000000000000000000000000000001": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0015"
            },
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_transaction",
                            "type": "bytes"
                        },
                        {
                            "internalType": "enum Lib_OVMCodec.EOASignatureType",
                            "name": "_signatureType",
                            "type": "uint8"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "execute",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "qall",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "mockOVM_ECDSAContractAccount": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0015",
            "code": "0x608060405234801561001057600080fd5b50600436106100365760003560e01c8063ac4340511461003b578063d1be05c214610065575b600080fd5b61004e610049366004610c49565b610078565b60405161005c929190610d70565b60405180910390f35b61004e610073366004610bb7565b610094565b6000606061008833868686610183565b91509150935093915050565b60006060338260018860018111156100a857fe5b1490506100b3610a48565b6100bd8a836101f4565b90506100f0336100cc8561034d565b83600001511460405180606001604052806034815260200161110d603491396103a6565b60608101516001600160a01b03166101485760006101178483604001518460a001516103ba565b905060018160405160200161012c9190610d5c565b6040516020818303038152906040529550955050505050610179565b610159838260000151600101610424565b61017183826040015183606001518460a00151610183565b945094505050505b9550959350505050565b60006060806101d0878787876040516024016101a19392919061104f565b60408051601f198184030181529190526020810180516001600160e01b03166342cbcfbb60e11b179052610464565b9050808060200190518101906101e69190610b64565b925092505094509492505050565b6101fc610a48565b811561026b5760008060008060006060888060200190518101906102209190610ca0565b6040805160e0810182529687526020870194909452928501939093526001600160a01b0390921660608401526000608084015260a083015260c0820152965061034795505050505050565b606061027684610471565b90506040518060e001604052806102a08360008151811061029357fe5b6020026020010151610484565b81526020016102b58360018151811061029357fe5b81526020016102ca8360028151811061029357fe5b81526020016102ec836003815181106102df57fe5b602002602001015161048f565b6001600160a01b0316815260200161030a8360048151811061029357fe5b815260200161032c8360058151811061031f57fe5b60200260200101516104d7565b81526020016103418360068151811061029357fe5b90529150505b92915050565b6040805160048152602481019091526020810180516001600160e01b03166360fd975160e11b179052600090606090610387908490610464565b90508080602001905181019061039d9190610c31565b9150505b919050565b816103b5576103b58382610531565b505050565b600060606104038585856040516024016103d49190610d93565b60408051601f198184030181529190526020810180516001600160e01b03166314aa2ff760e01b179052610575565b9050808060200190518101906104199190610b48565b9150505b9392505050565b6103b582826040516024016104399190611046565b60408051601f198184030181529190526020810180516001600160e01b0316630da449d160e01b1790525b606061041d835a84610575565b606061034761047f83610603565b610628565b60006103478261074a565b8051600090600114156104a4575060006103a1565b81516015146104ce5760405162461bcd60e51b81526004016104c590610f33565b60405180910390fd5b61034782610484565b606060008060006104e7856107da565b9194509250905060008160018111156104fc57fe5b146105195760405162461bcd60e51b81526004016104c590610fd8565b61052885602001518484610997565b95945050505050565b6103b582826040516024016105469190610d93565b60408051601f198184030181529190526020810180516001600160e01b0316632a2a7adb60e01b179052610464565b606060006060856001600160a01b031685856040516105949190610d40565b60006040518083038160008787f1925050503d80600081146105d2576040519150601f19603f3d011682016040523d82523d6000602084013e6105d7565b606091505b509092509050816105ea57805160208201fd5b8051600114156105fa5760016000f35b915061041d9050565b61060b610a8e565b506040805180820190915281518152602082810190820152919050565b6060600080610636846107da565b9193509091506001905081600181111561064c57fe5b146106695760405162461bcd60e51b81526004016104c590610ec5565b6040805160208082526104208201909252606091816020015b61068a610a8e565b8152602001906001900390816106825790505090506000835b865181101561073f57602082106106cc5760405162461bcd60e51b81526004016104c590610e44565b6000806106f86040518060400160405280858c60000151038152602001858c60200151018152506107da565b509150915060405180604001604052808383018152602001848b602001510181525085858151811061072657fe5b60209081029190910101526001939093019201016106a3565b508152949350505050565b60006021826000015111156107715760405162461bcd60e51b81526004016104c590610e8e565b600080600061077f856107da565b91945092509050600081600181111561079457fe5b146107b15760405162461bcd60e51b81526004016104c590610e8e565b6020808601518401805190918410156107d05760208490036101000a90045b9695505050505050565b6000806000808460000151116108025760405162461bcd60e51b81526004016104c590610f6a565b6020840151805160001a607f8111610827576000600160009450945094505050610990565b60b78111610867578551607f1982019081106108555760405162461bcd60e51b81526004016104c590610efc565b60019550935060009250610990915050565b60bf81116108e157855160b61982019081106108955760405162461bcd60e51b81526004016104c590610e0d565b6000816020036101000a60018501510490508082018860000151116108cc5760405162461bcd60e51b81526004016104c590610fa1565b60019091019550935060009250610990915050565b60f7811161092057855160bf19820190811061090f5760405162461bcd60e51b81526004016104c59061100f565b600195509350849250610990915050565b855160f61982019081106109465760405162461bcd60e51b81526004016104c590610da6565b6000816020036101000a600185015104905080820188600001511161097d5760405162461bcd60e51b81526004016104c590610ddd565b6001918201965094509250610990915050565b9193909250565b6060808267ffffffffffffffff811180156109b157600080fd5b506040519080825280601f01601f1916602001820160405280156109dc576020820181803683370190505b5090508051600014156109f057905061041d565b8484016020820160005b60208604811015610a1b5782518252602092830192909101906001016109fa565b5060006001602087066020036101000a039050808251168119845116178252839450505050509392505050565b6040518060e0016040528060008152602001600081526020016000815260200160006001600160a01b031681526020016000815260200160608152602001600081525090565b604051806040016040528060008152602001600081525090565b600082601f830112610ab8578081fd5b8135610acb610ac6826110a0565b611079565b9150808252836020828501011115610ae257600080fd5b8060208401602084013760009082016020015292915050565b600082601f830112610b0b578081fd5b8151610b19610ac6826110a0565b9150808252836020828501011115610b3057600080fd5b610b418160208401602086016110c4565b5092915050565b600060208284031215610b59578081fd5b815161041d816110f4565b60008060408385031215610b76578081fd5b82518015158114610b85578182fd5b602084015190925067ffffffffffffffff811115610ba1578182fd5b610bad85828601610afb565b9150509250929050565b600080600080600060a08688031215610bce578081fd5b853567ffffffffffffffff811115610be4578182fd5b610bf088828901610aa8565b955050602086013560028110610c04578182fd5b9350604086013560ff81168114610c19578182fd5b94979396509394606081013594506080013592915050565b600060208284031215610c42578081fd5b5051919050565b600080600060608486031215610c5d578283fd5b833592506020840135610c6f816110f4565b9150604084013567ffffffffffffffff811115610c8a578182fd5b610c9686828701610aa8565b9150509250925092565b60008060008060008060c08789031215610cb8578081fd5b865195506020870151945060408701519350606087015192506080870151610cdf816110f4565b60a088015190925067ffffffffffffffff811115610cfb578182fd5b610d0789828a01610afb565b9150509295509295509295565b60008151808452610d2c8160208601602086016110c4565b601f01601f19169290920160200192915050565b60008251610d528184602087016110c4565b9190910192915050565b6001600160a01b0391909116815260200190565b6000831515825260406020830152610d8b6040830184610d14565b949350505050565b60006020825261041d6020830184610d14565b6020808252601d908201527f496e76616c696420524c50206c6f6e67206c697374206c656e6774682e000000604082015260600190565b60208082526016908201527524b73b30b634b210292628103637b733903634b9ba1760511b604082015260600190565b6020808252601f908201527f496e76616c696420524c50206c6f6e6720737472696e67206c656e6774682e00604082015260600190565b6020808252602a908201527f50726f766964656420524c50206c6973742065786365656473206d6178206c6960408201526939ba103632b733ba341760b11b606082015260800190565b6020808252601a908201527f496e76616c696420524c5020627974657333322076616c75652e000000000000604082015260600190565b60208082526017908201527f496e76616c696420524c50206c6973742076616c75652e000000000000000000604082015260600190565b60208082526019908201527f496e76616c696420524c502073686f727420737472696e672e00000000000000604082015260600190565b6020808252601a908201527f496e76616c696420524c5020616464726573732076616c75652e000000000000604082015260600190565b60208082526018908201527f524c50206974656d2063616e6e6f74206265206e756c6c2e0000000000000000604082015260600190565b60208082526018908201527f496e76616c696420524c50206c6f6e6720737472696e672e0000000000000000604082015260600190565b60208082526018908201527f496e76616c696420524c502062797465732076616c75652e0000000000000000604082015260600190565b60208082526017908201527f496e76616c696420524c502073686f7274206c6973742e000000000000000000604082015260600190565b90815260200190565b8381526001600160a01b038316602082015260606040820181905260009061052890830184610d14565b60405181810167ffffffffffffffff8111828210171561109857600080fd5b604052919050565b600067ffffffffffffffff8211156110b6578081fd5b50601f01601f191660200190565b60005b838110156110df5781810151838201526020016110c7565b838111156110ee576000848401525b50505050565b6001600160a01b038116811461110957600080fd5b5056fe5472616e73616374696f6e206e6f6e636520646f6573206e6f74206d6174636820746865206578706563746564206e6f6e63652ea2646970667358221220cfaacc112a94005c62c12da059c39e22a2394b3344e0b1e638c42674d502340964736f6c63430007000033",
            "codeHash": "0xbe6d65fd3396999ce8dac663418fd477a14c84d3744b36636861f407a7969a59",
            "storage": {},
            "abi": [
                {
                    "inputs": [
                        {
                            "internalType": "bytes",
                            "name": "_transaction",
                            "type": "bytes"
                        },
                        {
                            "internalType": "enum Lib_OVMCodec.EOASignatureType",
                            "name": "_signatureType",
                            "type": "uint8"
                        },
                        {
                            "internalType": "uint8",
                            "name": "_v",
                            "type": "uint8"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_r",
                            "type": "bytes32"
                        },
                        {
                            "internalType": "bytes32",
                            "name": "_s",
                            "type": "bytes32"
                        }
                    ],
                    "name": "execute",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "uint256",
                            "name": "_gasLimit",
                            "type": "uint256"
                        },
                        {
                            "internalType": "address",
                            "name": "_to",
                            "type": "address"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_data",
                            "type": "bytes"
                        }
                    ],
                    "name": "qall",
                    "outputs": [
                        {
                            "internalType": "bool",
                            "name": "_success",
                            "type": "bool"
                        },
                        {
                            "internalType": "bytes",
                            "name": "_returndata",
                            "type": "bytes"
                        }
                    ],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        },
        "Lib_AddressManager": {
            "address": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0016",
            "code": "0x608060405234801561001057600080fd5b50600436106100575760003560e01c8063715018a61461005c5780638da5cb5b146100665780639b2ea4bd1461008a578063bf40fac11461013b578063f2fde38b146101e1575b600080fd5b610064610207565b005b61006e6102b0565b604080516001600160a01b039092168252519081900360200190f35b610064600480360360408110156100a057600080fd5b8101906020810181356401000000008111156100bb57600080fd5b8201836020820111156100cd57600080fd5b803590602001918460018302840111640100000000831117156100ef57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550505090356001600160a01b031691506102bf9050565b61006e6004803603602081101561015157600080fd5b81019060208101813564010000000081111561016c57600080fd5b82018360208201111561017e57600080fd5b803590602001918460018302840111640100000000831117156101a057600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250929550610362945050505050565b610064600480360360208110156101f757600080fd5b50356001600160a01b0316610391565b6000546001600160a01b03163314610266576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000546001600160a01b031681565b6000546001600160a01b0316331461031e576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b806001600061032c85610490565b815260200190815260200160002060006101000a8154816001600160a01b0302191690836001600160a01b031602179055505050565b60006001600061037184610490565b81526020810191909152604001600020546001600160a01b031692915050565b6000546001600160a01b031633146103f0576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6001600160a01b0381166104355760405162461bcd60e51b815260040180806020018281038252602d815260200180610508602d913960400191505060405180910390fd5b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6000816040516020018082805190602001908083835b602083106104c55780518252601f1990920191602091820191016104a6565b6001836020036101000a03801982511681845116808217855250505050505090500191505060405160208183030381529060405280519060200120905091905056fe4f776e61626c653a206e6577206f776e65722063616e6e6f7420626520746865207a65726f2061646472657373a26469706673582212204367ffc2e6671623708150e2d0cff4c12cf566722a26b4748555d789953e2d2264736f6c63430007000033",
            "codeHash": "0x47fa60e704defda58d5b162cbc036d760788fbd5ce6730de562af406e0db37a8",
            "storage": {
                "0x24e095abd8bf5f81f3350e6cb0d49574e94e998bfb6341a6ed085c6e3ef4d7fe": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0004",
                "0x4a268d14639fa54a62da41e53d5cfed7d8ef15ff1108a54747e0fd38d7741a68": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000a",
                "0x5c2e827bedec24adf1d781771ca0503c801b1637965c73d197cb2ea8857f2921": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000c",
                "0x9dc316a765d11a12b06619d367ef78fecac216d290033f772936da756c0d28fe": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead000e",
                "0xb73b2537b0fac790040c3ef6c5d622006013c6e62c05ff3c8275f38003cd72a1": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0014",
                "0xde24ca96c4b0b6ed2c73bb46c1053b6edd9470cda80c625493502cc81a3ccfa7": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0006",
                "0x0000000000000000000000000000000000000000000000000000000000000000": "0x17ec8597ff92c3f44523bdc65bf0f1be632917ff",
                "0x0248c104bff13515d06afb602d097ac0d52680c2d14e6c66219633a4b949f2ef": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0000",
                "0x0b198951118b45b895fd138b1229db341527c87de0bd478d658ea055cd73802f": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0012",
                "0x0cc4bd6bd0492462730f0bcc5303174d0a2af52b1ae68b25e2c7daada2292362": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0002",
                "0xf0b64a30864ef1e4b0c96bb2c6ba336fd423add8e4f685027042faf4a65c6112": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0010",
                "0xf56747885613486d091c4459f3b37706019a79fb2cf73bde37750a936fe58e30": "0xdeaddeaddeaddeaddeaddeaddeaddeaddead0008"
            },
            "abi": [
                {
                    "anonymous": false,
                    "inputs": [
                        {
                            "indexed": true,
                            "internalType": "address",
                            "name": "previousOwner",
                            "type": "address"
                        },
                        {
                            "indexed": true,
                            "internalType": "address",
                            "name": "newOwner",
                            "type": "address"
                        }
                    ],
                    "name": "OwnershipTransferred",
                    "type": "event"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        }
                    ],
                    "name": "getAddress",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "owner",
                    "outputs": [
                        {
                            "internalType": "address",
                            "name": "",
                            "type": "address"
                        }
                    ],
                    "stateMutability": "view",
                    "type": "function"
                },
                {
                    "inputs": [],
                    "name": "renounceOwnership",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "string",
                            "name": "_name",
                            "type": "string"
                        },
                        {
                            "internalType": "address",
                            "name": "_address",
                            "type": "address"
                        }
                    ],
                    "name": "setAddress",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                },
                {
                    "inputs": [
                        {
                            "internalType": "address",
                            "name": "_newOwner",
                            "type": "address"
                        }
                    ],
                    "name": "transferOwnership",
                    "outputs": [],
                    "stateMutability": "nonpayable",
                    "type": "function"
                }
            ]
        }
    }
}
`)

type ovmDumpAccount struct {
	Address  common.Address         `json:"address"`
	Code     string                 `json:"code"`
	CodeHash string                 `json:"codeHash"`
	Storage  map[common.Hash]string `json:"storage"`
	ABI      abi.ABI                `json:"abi"`
}

type ovmDump struct {
	Accounts map[string]ovmDumpAccount `json:"accounts"`
}

// OvmStateDump is the full (parsed) OVM state dump object.
var OvmStateDump ovmDump

// OvmStateManager is the account corresponding to the OVM_StateManager.
var OvmStateManager ovmDumpAccount

// OvmExecutionManager is the account corresponding to the OVM_ExecutionManager.
var OvmExecutionManager ovmDumpAccount

// UsingOVM is used to enable or disable functionality necessary for the OVM.
var UsingOVM bool

func init() {
	var err error

	err = json.Unmarshal(ovmStateDumpJSON, &OvmStateDump)
	if err != nil {
		panic(fmt.Errorf("could not decode OVM state dump: %v", err))
	}

	OvmStateManager = OvmStateDump.Accounts["OVM_StateManager"]
	OvmExecutionManager = OvmStateDump.Accounts["OVM_ExecutionManager"]
	UsingOVM = os.Getenv("USING_OVM") == "true"
}
