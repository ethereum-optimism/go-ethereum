// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package statecommitmentchain

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Lib_OVMCodecChainBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type Lib_OVMCodecChainBatchHeader struct {
	BatchIndex        *big.Int
	BatchRoot         [32]byte
	BatchSize         *big.Int
	PrevTotalElements *big.Int
	ExtraData         []byte
}

// Lib_OVMCodecChainInclusionProof is an auto generated low-level Go binding around an user-defined struct.
type Lib_OVMCodecChainInclusionProof struct {
	Index    *big.Int
	Siblings [][32]byte
}

// Lib_OVMCodecTransaction is an auto generated low-level Go binding around an user-defined struct.
type Lib_OVMCodecTransaction struct {
	Timestamp     *big.Int
	BlockNumber   *big.Int
	L1QueueOrigin uint8
	L1TxOrigin    common.Address
	Entrypoint    common.Address
	GasLimit      *big.Int
	Data          []byte
}

// Lib_OVMCodecTransactionChainElement is an auto generated low-level Go binding around an user-defined struct.
type Lib_OVMCodecTransactionChainElement struct {
	IsSequenced bool
	QueueIndex  *big.Int
	Timestamp   *big.Int
	BlockNumber *big.Int
	TxData      []byte
}

// OVMStateCommitmentChainABI is the input ABI used to generate the binding from.
const OVMStateCommitmentChainABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_libAddressManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"FRAUD_PROOF_WINDOW\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SEQUENCER_PUBLISH_WINDOW\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_batch\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"_shouldStartAtElement\",\"type\":\"uint256\"}],\"name\":\"appendStateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_id\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"canOverwrite\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"}],\"name\":\"deleteStateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastSequencerTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_lastSequencerTimestamp\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBatches\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalBatches\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalElements\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalElements\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"}],\"name\":\"insideFraudProofWindow\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_inside\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_stateBatchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"enumLib_OVMCodec.QueueOrigin\",\"name\":\"l1QueueOrigin\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"l1TxOrigin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"entrypoint\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.Transaction\",\"name\":\"_transaction\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"isSequenced\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"queueIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.TransactionChainElement\",\"name\":\"_txChainElement\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_txBatchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"_txInclusionProof\",\"type\":\"tuple\"}],\"name\":\"setLastOverwritableIndex\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_element\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"verifyStateCommitment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OVMStateCommitmentChainBin is the compiled bytecode used for deploying new contracts.
var OVMStateCommitmentChainBin = "0x60806040523480156200001157600080fd5b506040516200247738038062002477833981016040819052620000349162000257565b600080546001600160a01b0319166001600160a01b03831617905560408051808201909152601d81527f4f564d5f43616e6f6e6963616c5472616e73616374696f6e436861696e00000060208201526200008e9062000195565b600b80546001600160a01b0319166001600160a01b039290921691909117905560408051808201909152601181527027ab26afa33930bab22b32b934b334b2b960791b6020820152620000e19062000195565b600c80546001600160a01b0319166001600160a01b039290921691909117905560408051808201909152600f81526e27ab26afa137b73226b0b730b3b2b960891b6020820152620001329062000195565b600d80546001600160a01b0319166001600160a01b03929092169190911790556200018e600360107f96df3abc26f419f0cc8d819984a2b87820d08c41bf1b84a59ce36f5d7336d1913062000222602090811b6200076117901c565b50620002dd565b6000805460405163bf40fac160e01b81526001600160a01b039091169063bf40fac190620001c890859060040162000287565b60206040518083038186803b158015620001e157600080fd5b505afa158015620001f6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200021c919062000257565b92915050565b6004840183905560068401929092558255600190910180546001600160a01b0319166001600160a01b03909216919091179055565b60006020828403121562000269578081fd5b81516001600160a01b038116811462000280578182fd5b9392505050565b6000602080835283518082850152825b81811015620002b55785810183015185820160400152820162000297565b81811115620002c75783604083870101525b50601f01601f1916929092016040019392505050565b61218a80620002ed6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c806381eb62ef1161007157806381eb62ef146101475780638ca5cbb91461014f5780639418bddd14610162578063b8e189ac14610175578063c17b291b14610188578063e561dddc14610190576100b4565b80632979761b146100b9578063461a4478146100ce5780634d69ee57146100f7578063677f5aff146101175780637aa63a861461012a5780637ad168a01461013f575b600080fd5b6100cc6100c736600461185b565b610198565b005b6100e16100dc3660046117ef565b6102be565b6040516100ee9190611aeb565b60405180910390f35b61010a610105366004611765565b610347565b6040516100ee9190611aff565b61010a6101253660046117ce565b6103ac565b61013261047e565b6040516100ee9190611ad4565b610132610497565b6101326104b0565b6100cc61015d366004611707565b6104b6565b61010a610170366004611829565b610677565b6100cc610183366004611829565b6106c2565b610132610742565b610132610749565b6101a185610796565b6101c65760405162461bcd60e51b81526004016101bd90611d42565b60405180910390fd5b6101cf85610677565b156101ec5760405162461bcd60e51b81526004016101bd90611c79565b60015485511161020e5760405162461bcd60e51b81526004016101bd90611fa3565b600b546040516326f2b4e760e11b81526001600160a01b0390911690634de569ce9061024490879087908790879060040161202e565b60206040518083038186803b15801561025c57600080fd5b505afa158015610270573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102949190611749565b6102b05760405162461bcd60e51b81526004016101bd90611f03565b505091516001555160025550565b6000805460405163bf40fac160e01b81526001600160a01b039091169063bf40fac1906102ef908590600401611b4e565b60206040518083038186803b15801561030757600080fd5b505afa15801561031b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061033f91906116eb565b90505b919050565b600061035283610796565b61036e5760405162461bcd60e51b81526004016101bd90611d42565b610386836020015185846000015185602001516107be565b6103a25760405162461bcd60e51b81526004016101bd90611c14565b5060019392505050565b60007f30a907da349b6916f6bc60eb25a37176d5705fba414069bdb6f625ebf8bb6c558314156104715760028054600b5460405163153f8c5f60e11b815291926001600160a01b0390911691632a7f18be9161040e9190870490600401611ad4565b60606040518083038186803b15801561042657600080fd5b505afa15801561043a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061045e9190611926565b6020015164ffffffffff16109050610478565b5060015481105b92915050565b600080610489610829565b5064ffffffffff1691505090565b6000806104a2610829565b64ffffffffff169250505090565b61070881565b6104be61047e565b81146104dc5760405162461bcd60e51b81526004016101bd90611ce5565b600d54604051630156a69560e11b81526001600160a01b03909116906302ad4d2a9061050c903390600401611aeb565b60206040518083038186803b15801561052457600080fd5b505afa158015610538573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061055c9190611749565b6105785760405162461bcd60e51b81526004016101bd90611eb4565b60008251116105995760405162461bcd60e51b81526004016101bd90611e71565b600b60009054906101000a90046001600160a01b03166001600160a01b0316637aa63a866040518163ffffffff1660e01b815260040160206040518083038186803b1580156105e757600080fd5b505afa1580156105fb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061061f9190611970565b825161062961047e565b0111156106485760405162461bcd60e51b81526004016101bd90611ba5565b61067382423360405160200161065f9291906120ea565b604051602081830303815290604052610850565b5050565b60008082608001518060200190518101906106929190611988565b509050806106b25760405162461bcd60e51b81526004016101bd90611e2c565b4262093a80820111915050919050565b600c546001600160a01b031633146106ec5760405162461bcd60e51b81526004016101bd90611dcf565b6106f581610796565b6107115760405162461bcd60e51b81526004016101bd90611d42565b61071a81610677565b6107365760405162461bcd60e51b81526004016101bd90611d71565b61073f816109fc565b50565b62093a8081565b60006107556003610a74565b64ffffffffff16905090565b6004840183905560068401929092558255600190910180546001600160a01b0319166001600160a01b03909216919091179055565b80516000906107ae9060039064ffffffffff16610a8f565b6107b783610bed565b1492915050565b600083815b835181101561081d5760008482815181106107da57fe5b60209081029190910101519050600186831c8116148015610806576107ff8483610c33565b9350610813565b6108108285610c33565b93505b50506001016107c3565b50909414949350505050565b60008060006108386003610c66565b64ffffffffff602882901c16935060501c9150509091565b60006108806040518060400160405280600d81526020016c27ab26afa9b2b8bab2b731b2b960991b8152506102be565b905060008061088d610829565b9092509050336001600160a01b03841614156108aa5750426108d4565b426107088264ffffffffff1601106108d45760405162461bcd60e51b81526004016101bd90611f3a565b606085516001600160401b03811180156108ed57600080fd5b5060405190808252806020026020018201604052801561092157816020015b606081526020019060019003908161090c5790505b50905060005b86518110156109825786818151811061093c57fe5b60200260200101516040516020016109549190611ad4565b60405160208183030381529060405282828151811061096f57fe5b6020908102919091010152600101610927565b5061098b6113b4565b6040518060a0016040528061099e610749565b81526020016109ac84610c84565b81528351602082015264ffffffffff8616604082015260600187905290506109f36109d682610bed565b6109ea836040015184606001510186610d23565b60039190610d31565b50505050505050565b610a066003610a74565b64ffffffffff16816000015110610a2f5760405162461bcd60e51b81526004016101bd90612000565b610a3881610796565b610a545760405162461bcd60e51b81526004016101bd90611d42565b61073f8160000151610a6b83606001516000610d23565b60039190610eb6565b6000610a7e6113e6565b610a8783610fc6565b519392505050565b6000610a996113e6565b610aa284610fc6565b805190915064ffffffffff168310610acc5760405162461bcd60e51b81526004016101bd90611c4b565b6000610aee82604001516001600160401b03168661102890919063ffffffff16565b90506000610b1583604001516001016001600160401b03168761102890919063ffffffff16565b9050826080015164ffffffffff168510610b76576080830151825464ffffffffff9091168603908110610b5a5760405162461bcd60e51b81526004016101bd90611c4b565b6000908152600190920160205250604090205491506104789050565b6080830151606084015164ffffffffff9182168781039290911610610bad5760405162461bcd60e51b81526004016101bd90611c4b565b8154811115610bce5760405162461bcd60e51b81526004016101bd90611c4b565b8154036000908152600190910160205260409020549250610478915050565b60008160200151826040015183606001518460800151604051602001610c169493929190611b0a565b604051602081830303815290604052805190602001209050919050565b60008282604051602001610c48929190611add565b60405160208183030381529060405280519060200120905092915050565b6000610c706113e6565b610c7983610fc6565b602001519392505050565b6000606082516001600160401b0381118015610c9f57600080fd5b50604051908082528060200260200182016040528015610cc9578160200160208202803683370190505b50905060005b8351811015610d1257838181518110610ce457fe5b602002602001015180519060200120828281518110610cff57fe5b6020908102919091010152600101610ccf565b50610d1c81611044565b9392505050565b602890811b91909117901b90565b610d396113e6565b610d4284610fc6565b90506000610d6682604001516001600160401b03168661102890919063ffffffff16565b8054909150610d7457601081555b8054608083015183510364ffffffffff1610610e6c5760018501548554608084015160405163677f5aff60e01b81526000936001600160a01b03169263677f5aff92610dc292600401611b39565b602060405180830381600087803b158015610ddc57600080fd5b505af1925050508015610e0c575060408051601f3d908101601f19168201909252610e0991810190611749565b60015b610e1857506000610e1b565b90505b8015610e62576040830180516001016001600160401b03169081905260808401805164ffffffffff90811660608701528551169052610e5b908790611028565b9150610e6a565b815460020282555b505b608082015182510364ffffffffff9081166000818152600184810160209081526040909220889055855101909216845264ffffffffff198516918401919091526109f38684611259565b610ebe6113e6565b610ec784610fc6565b9050806000015164ffffffffff168364ffffffffff16108015610efc5750806060015164ffffffffff168364ffffffffff1610155b610f185760405162461bcd60e51b81526004016101bd90611c4b565b6000610f3a82604001516001600160401b03168661102890919063ffffffff16565b90506000610f6183604001516001016001600160401b03168761102890919063ffffffff16565b9050826080015164ffffffffff168564ffffffffff161015610fa457604083018051600019016001600160401b03169052606083015164ffffffffff1660808401525b64ffffffffff8516835264ffffffffff19841660208401526109f38684611259565b610fce6113e6565b5060028101546003909101546040805160a08101825264ffffffffff808516825264ffffffffff1990941660208201526001600160401b038316818301529082901c8316606082015260689190911c909116608082015290565b6000600282061561103c5782600601610d1c565b505060040190565b6000808251116110665760405162461bcd60e51b81526004016101bd90611b61565b81516001141561108c578160008151811061107d57fe5b60200260200101519050610342565b606061109883516112bf565b8351909150839060029006600114156111335783516001016001600160401b03811180156110c557600080fd5b506040519080825280602002602001820160405280156110ef578160200160208202803683370190505b50905060005b84518110156111315784818151811061110a57fe5b602002602001015182828151811061111e57fe5b60209081029190910101526001016110f5565b505b835160009060028106600114156111735783828151811061115057fe5b602002602001015183828151811061116457fe5b60209081029190910101526001015b60018111156112395760018201915060005b600282048110156111e8576111c98482600202815181106111a257fe5b60200260200101518583600202600101815181106111bc57fe5b6020026020010151610c33565b8482815181106111d557fe5b6020908102919091010152600101611185565b50600290046001808216148015611200575080600114155b156112345783828151811061121157fe5b602002602001015183828151811061122557fe5b60209081029190910101526001015b611173565b8260008151811061124657fe5b6020026020010151945050505050919050565b80516020820151604080840151606085015160808601516002880154600096868117969584901b8517606884901b17959094909390929091871461129f5760028a018790555b858a60030154146112b25760038a018690555b5050505050505092915050565b606080826001600160401b03811180156112d857600080fd5b50604051908082528060200260200182016040528015611302578160200160208202803683370190505b50905060006040516020016113179190611ad4565b604051602081830303815290604052805190602001208160008151811061133a57fe5b602090810291909101015260015b81518110156113ad5781600182038151811061136057fe5b60200260200101516040516020016113789190611ad4565b6040516020818303038152906040528051906020012082828151811061139a57fe5b6020908102919091010152600101611348565b5092915050565b6040518060a0016040528060008152602001600080191681526020016000815260200160008152602001606081525090565b6040805160a08101825260008082526020820181905291810182905260608101829052608081019190915290565b803561047881612131565b600082601f83011261142f578081fd5b81356001600160401b03811115611444578182fd5b6020808202611454828201612101565b8381529350818401858301828701840188101561147057600080fd5b600092505b84831015611493578035825260019290920191908301908301611475565b505050505092915050565b600082601f8301126114ae578081fd5b81356001600160401b038111156114c3578182fd5b6114d6601f8201601f1916602001612101565b91508082528360208285010111156114ed57600080fd5b8060208401602084013760009082016020015292915050565b80356002811061047857600080fd5b600060a08284031215611526578081fd5b61153060a0612101565b90508135815260208201356020820152604082013560408201526060820135606082015260808201356001600160401b0381111561156d57600080fd5b6115798482850161149e565b60808301525092915050565b600060408284031215611596578081fd5b6115a06040612101565b90508135815260208201356001600160401b038111156115bf57600080fd5b6115cb8482850161141f565b60208301525092915050565b600060a082840312156115e8578081fd5b6115f260a0612101565b905081356115ff81612146565b8082525060208201356020820152604082013560408201526060820135606082015260808201356001600160401b0381111561156d57600080fd5b600060e0828403121561164b578081fd5b61165560e0612101565b905081358152602082013560208201526116728360408401611506565b60408201526116848360608401611414565b60608201526116968360808401611414565b608082015260a082013560a082015260c08201356001600160401b038111156116be57600080fd5b6116ca8482850161149e565b60c08301525092915050565b805164ffffffffff8116811461047857600080fd5b6000602082840312156116fc578081fd5b8151610d1c81612131565b60008060408385031215611719578081fd5b82356001600160401b0381111561172e578182fd5b61173a8582860161141f565b95602094909401359450505050565b60006020828403121561175a578081fd5b8151610d1c81612146565b600080600060608486031215611779578081fd5b8335925060208401356001600160401b0380821115611796578283fd5b6117a287838801611515565b935060408601359150808211156117b7578283fd5b506117c486828701611585565b9150509250925092565b600080604083850312156117e0578182fd5b50508035926020909101359150565b600060208284031215611800578081fd5b81356001600160401b03811115611815578182fd5b6118218482850161149e565b949350505050565b60006020828403121561183a578081fd5b81356001600160401b0381111561184f578182fd5b61182184828501611515565b600080600080600060a08688031215611872578283fd5b85356001600160401b0380821115611888578485fd5b61189489838a01611515565b965060208801359150808211156118a9578485fd5b6118b589838a0161163a565b955060408801359150808211156118ca578485fd5b6118d689838a016115d7565b945060608801359150808211156118eb578283fd5b6118f789838a01611515565b9350608088013591508082111561190c578283fd5b5061191988828901611585565b9150509295509295909350565b600060608284031215611937578081fd5b6119416060612101565b8251815261195284602085016116d6565b602082015261196484604085016116d6565b60408201529392505050565b600060208284031215611981578081fd5b5051919050565b6000806040838503121561199a578182fd5b8251915060208301516119ac81612131565b809150509250929050565b60008151808452815b818110156119dc576020818501810151868301820152016119c0565b818111156119ed5782602083870101525b50601f01601f19169290920160200192915050565b600081518352602082015160208401526040820151604084015260608201516060840152608082015160a0608085015261182160a08501826119b7565b6000604083018251845260208084015160408287015282815180855260608801915083830194508592505b80831015611a8a5784518252938301936001929092019190830190611a6a565b509695505050505050565b6000815115158352602082015160208401526040820151604084015260608201516060840152608082015160a0608085015261182160a08501826119b7565b90815260200190565b918252602082015260400190565b6001600160a01b0391909116815260200190565b901515815260200190565b600085825284602083015283604083015260806060830152611b2f60808301846119b7565b9695505050505050565b91825264ffffffffff16602082015260400190565b600060208252610d1c60208301846119b7565b60208082526024908201527f4d7573742070726f76696465206174206c65617374206f6e65206c656166206860408201526330b9b41760e11b606082015260800190565b60208082526049908201527f4e756d626572206f6620737461746520726f6f74732063616e6e6f742065786360408201527f65656420746865206e756d626572206f662063616e6f6e6963616c207472616e60608201526839b0b1ba34b7b7399760b91b608082015260a00190565b60208082526018908201527f496e76616c696420696e636c7573696f6e2070726f6f662e0000000000000000604082015260600190565b60208082526014908201527324b73232bc1037baba1037b3103137bab732399760611b604082015260600190565b60208082526046908201527f426174636820686561646572206d757374206265206f757473696465206f662060408201527f66726175642070726f6f662077696e646f7720746f206265206f7665727772696060820152653a30b136329760d11b608082015260a00190565b6020808252603d908201527f41637475616c20626174636820737461727420696e64657820646f6573206e6f60408201527f74206d6174636820657870656374656420737461727420696e6465782e000000606082015260800190565b60208082526015908201527424b73b30b634b2103130ba31b4103432b0b232b91760591b604082015260600190565b602080825260409082018190527f537461746520626174636865732063616e206f6e6c792062652064656c657465908201527f642077697468696e207468652066726175642070726f6f662077696e646f772e606082015260800190565b6020808252603b908201527f537461746520626174636865732063616e206f6e6c792062652064656c65746560408201527f6420627920746865204f564d5f467261756456657269666965722e0000000000606082015260800190565b60208082526025908201527f4261746368206865616465722074696d657374616d702063616e6e6f74206265604082015264207a65726f60d81b606082015260800190565b60208082526023908201527f43616e6e6f74207375626d697420616e20656d7074792073746174652062617460408201526231b41760e91b606082015260800190565b6020808252602f908201527f50726f706f73657220646f6573206e6f74206861766520656e6f75676820636f60408201526e1b1b185d195c985b081c1bdcdd1959608a1b606082015260800190565b6020808252601a908201527f496e76616c6964207472616e73616374696f6e2070726f6f662e000000000000604082015260600190565b60208082526043908201527f43616e6e6f74207075626c69736820737461746520726f6f747320776974686960408201527f6e207468652073657175656e636572207075626c69636174696f6e2077696e6460608201526237bb9760e91b608082015260a00190565b60208082526039908201527f426174636820696e646578206d7573742062652067726561746572207468616e60408201527f206c617374206f7665727772697461626c6520696e6465782e00000000000000606082015260800190565b60208082526014908201527324b73b30b634b2103130ba31b41034b73232bc1760611b604082015260600190565b60006080825285516080830152602086015160a0830152604086015161205381612127565b8060c084015250606086015160018060a01b0380821660e085015280608089015116610100850152505060a086015161012083015260c086015160e06101408401526120a36101608401826119b7565b905082810360208401526120b78187611a95565b905082810360408401526120cb8186611a02565b905082810360608401526120df8185611a3f565b979650505050505050565b9182526001600160a01b0316602082015260400190565b6040518181016001600160401b038111828210171561211f57600080fd5b604052919050565b6002811061073f57fe5b6001600160a01b038116811461073f57600080fd5b801515811461073f57600080fdfea26469706673582212200cb074be623b3d17d845200459ebe861f1e590cbb64feb9f0299aeb43a02bffb64736f6c63430007000033"

// DeployOVMStateCommitmentChain deploys a new Ethereum contract, binding an instance of OVMStateCommitmentChain to it.
func DeployOVMStateCommitmentChain(auth *bind.TransactOpts, backend bind.ContractBackend, _libAddressManager common.Address) (common.Address, *types.Transaction, *OVMStateCommitmentChain, error) {
	parsed, err := abi.JSON(strings.NewReader(OVMStateCommitmentChainABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OVMStateCommitmentChainBin), backend, _libAddressManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OVMStateCommitmentChain{OVMStateCommitmentChainCaller: OVMStateCommitmentChainCaller{contract: contract}, OVMStateCommitmentChainTransactor: OVMStateCommitmentChainTransactor{contract: contract}, OVMStateCommitmentChainFilterer: OVMStateCommitmentChainFilterer{contract: contract}}, nil
}

// OVMStateCommitmentChain is an auto generated Go binding around an Ethereum contract.
type OVMStateCommitmentChain struct {
	OVMStateCommitmentChainCaller     // Read-only binding to the contract
	OVMStateCommitmentChainTransactor // Write-only binding to the contract
	OVMStateCommitmentChainFilterer   // Log filterer for contract events
}

// OVMStateCommitmentChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type OVMStateCommitmentChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMStateCommitmentChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OVMStateCommitmentChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMStateCommitmentChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OVMStateCommitmentChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMStateCommitmentChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OVMStateCommitmentChainSession struct {
	Contract     *OVMStateCommitmentChain // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// OVMStateCommitmentChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OVMStateCommitmentChainCallerSession struct {
	Contract *OVMStateCommitmentChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// OVMStateCommitmentChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OVMStateCommitmentChainTransactorSession struct {
	Contract     *OVMStateCommitmentChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// OVMStateCommitmentChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type OVMStateCommitmentChainRaw struct {
	Contract *OVMStateCommitmentChain // Generic contract binding to access the raw methods on
}

// OVMStateCommitmentChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OVMStateCommitmentChainCallerRaw struct {
	Contract *OVMStateCommitmentChainCaller // Generic read-only contract binding to access the raw methods on
}

// OVMStateCommitmentChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OVMStateCommitmentChainTransactorRaw struct {
	Contract *OVMStateCommitmentChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOVMStateCommitmentChain creates a new instance of OVMStateCommitmentChain, bound to a specific deployed contract.
func NewOVMStateCommitmentChain(address common.Address, backend bind.ContractBackend) (*OVMStateCommitmentChain, error) {
	contract, err := bindOVMStateCommitmentChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OVMStateCommitmentChain{OVMStateCommitmentChainCaller: OVMStateCommitmentChainCaller{contract: contract}, OVMStateCommitmentChainTransactor: OVMStateCommitmentChainTransactor{contract: contract}, OVMStateCommitmentChainFilterer: OVMStateCommitmentChainFilterer{contract: contract}}, nil
}

// NewOVMStateCommitmentChainCaller creates a new read-only instance of OVMStateCommitmentChain, bound to a specific deployed contract.
func NewOVMStateCommitmentChainCaller(address common.Address, caller bind.ContractCaller) (*OVMStateCommitmentChainCaller, error) {
	contract, err := bindOVMStateCommitmentChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OVMStateCommitmentChainCaller{contract: contract}, nil
}

// NewOVMStateCommitmentChainTransactor creates a new write-only instance of OVMStateCommitmentChain, bound to a specific deployed contract.
func NewOVMStateCommitmentChainTransactor(address common.Address, transactor bind.ContractTransactor) (*OVMStateCommitmentChainTransactor, error) {
	contract, err := bindOVMStateCommitmentChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OVMStateCommitmentChainTransactor{contract: contract}, nil
}

// NewOVMStateCommitmentChainFilterer creates a new log filterer instance of OVMStateCommitmentChain, bound to a specific deployed contract.
func NewOVMStateCommitmentChainFilterer(address common.Address, filterer bind.ContractFilterer) (*OVMStateCommitmentChainFilterer, error) {
	contract, err := bindOVMStateCommitmentChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OVMStateCommitmentChainFilterer{contract: contract}, nil
}

// bindOVMStateCommitmentChain binds a generic wrapper to an already deployed contract.
func bindOVMStateCommitmentChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OVMStateCommitmentChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OVMStateCommitmentChain.Contract.OVMStateCommitmentChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.OVMStateCommitmentChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.OVMStateCommitmentChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OVMStateCommitmentChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.contract.Transact(opts, method, params...)
}

// FRAUDPROOFWINDOW is a free data retrieval call binding the contract method 0xc17b291b.
//
// Solidity: function FRAUD_PROOF_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) FRAUDPROOFWINDOW(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "FRAUD_PROOF_WINDOW")
	return *ret0, err
}

// FRAUDPROOFWINDOW is a free data retrieval call binding the contract method 0xc17b291b.
//
// Solidity: function FRAUD_PROOF_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) FRAUDPROOFWINDOW() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.FRAUDPROOFWINDOW(&_OVMStateCommitmentChain.CallOpts)
}

// FRAUDPROOFWINDOW is a free data retrieval call binding the contract method 0xc17b291b.
//
// Solidity: function FRAUD_PROOF_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) FRAUDPROOFWINDOW() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.FRAUDPROOFWINDOW(&_OVMStateCommitmentChain.CallOpts)
}

// SEQUENCERPUBLISHWINDOW is a free data retrieval call binding the contract method 0x81eb62ef.
//
// Solidity: function SEQUENCER_PUBLISH_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) SEQUENCERPUBLISHWINDOW(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "SEQUENCER_PUBLISH_WINDOW")
	return *ret0, err
}

// SEQUENCERPUBLISHWINDOW is a free data retrieval call binding the contract method 0x81eb62ef.
//
// Solidity: function SEQUENCER_PUBLISH_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) SEQUENCERPUBLISHWINDOW() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.SEQUENCERPUBLISHWINDOW(&_OVMStateCommitmentChain.CallOpts)
}

// SEQUENCERPUBLISHWINDOW is a free data retrieval call binding the contract method 0x81eb62ef.
//
// Solidity: function SEQUENCER_PUBLISH_WINDOW() constant returns(uint256)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) SEQUENCERPUBLISHWINDOW() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.SEQUENCERPUBLISHWINDOW(&_OVMStateCommitmentChain.CallOpts)
}

// CanOverwrite is a free data retrieval call binding the contract method 0x677f5aff.
//
// Solidity: function canOverwrite(bytes32 _id, uint256 _index) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) CanOverwrite(opts *bind.CallOpts, _id [32]byte, _index *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "canOverwrite", _id, _index)
	return *ret0, err
}

// CanOverwrite is a free data retrieval call binding the contract method 0x677f5aff.
//
// Solidity: function canOverwrite(bytes32 _id, uint256 _index) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) CanOverwrite(_id [32]byte, _index *big.Int) (bool, error) {
	return _OVMStateCommitmentChain.Contract.CanOverwrite(&_OVMStateCommitmentChain.CallOpts, _id, _index)
}

// CanOverwrite is a free data retrieval call binding the contract method 0x677f5aff.
//
// Solidity: function canOverwrite(bytes32 _id, uint256 _index) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) CanOverwrite(_id [32]byte, _index *big.Int) (bool, error) {
	return _OVMStateCommitmentChain.Contract.CanOverwrite(&_OVMStateCommitmentChain.CallOpts, _id, _index)
}

// GetLastSequencerTimestamp is a free data retrieval call binding the contract method 0x7ad168a0.
//
// Solidity: function getLastSequencerTimestamp() constant returns(uint256 _lastSequencerTimestamp)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) GetLastSequencerTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "getLastSequencerTimestamp")
	return *ret0, err
}

// GetLastSequencerTimestamp is a free data retrieval call binding the contract method 0x7ad168a0.
//
// Solidity: function getLastSequencerTimestamp() constant returns(uint256 _lastSequencerTimestamp)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) GetLastSequencerTimestamp() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetLastSequencerTimestamp(&_OVMStateCommitmentChain.CallOpts)
}

// GetLastSequencerTimestamp is a free data retrieval call binding the contract method 0x7ad168a0.
//
// Solidity: function getLastSequencerTimestamp() constant returns(uint256 _lastSequencerTimestamp)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) GetLastSequencerTimestamp() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetLastSequencerTimestamp(&_OVMStateCommitmentChain.CallOpts)
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) GetTotalBatches(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "getTotalBatches")
	return *ret0, err
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) GetTotalBatches() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetTotalBatches(&_OVMStateCommitmentChain.CallOpts)
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) GetTotalBatches() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetTotalBatches(&_OVMStateCommitmentChain.CallOpts)
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) GetTotalElements(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "getTotalElements")
	return *ret0, err
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) GetTotalElements() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetTotalElements(&_OVMStateCommitmentChain.CallOpts)
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) GetTotalElements() (*big.Int, error) {
	return _OVMStateCommitmentChain.Contract.GetTotalElements(&_OVMStateCommitmentChain.CallOpts)
}

// InsideFraudProofWindow is a free data retrieval call binding the contract method 0x9418bddd.
//
// Solidity: function insideFraudProofWindow(Lib_OVMCodecChainBatchHeader _batchHeader) constant returns(bool _inside)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) InsideFraudProofWindow(opts *bind.CallOpts, _batchHeader Lib_OVMCodecChainBatchHeader) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "insideFraudProofWindow", _batchHeader)
	return *ret0, err
}

// InsideFraudProofWindow is a free data retrieval call binding the contract method 0x9418bddd.
//
// Solidity: function insideFraudProofWindow(Lib_OVMCodecChainBatchHeader _batchHeader) constant returns(bool _inside)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) InsideFraudProofWindow(_batchHeader Lib_OVMCodecChainBatchHeader) (bool, error) {
	return _OVMStateCommitmentChain.Contract.InsideFraudProofWindow(&_OVMStateCommitmentChain.CallOpts, _batchHeader)
}

// InsideFraudProofWindow is a free data retrieval call binding the contract method 0x9418bddd.
//
// Solidity: function insideFraudProofWindow(Lib_OVMCodecChainBatchHeader _batchHeader) constant returns(bool _inside)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) InsideFraudProofWindow(_batchHeader Lib_OVMCodecChainBatchHeader) (bool, error) {
	return _OVMStateCommitmentChain.Contract.InsideFraudProofWindow(&_OVMStateCommitmentChain.CallOpts, _batchHeader)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) Resolve(opts *bind.CallOpts, _name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "resolve", _name)
	return *ret0, err
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) Resolve(_name string) (common.Address, error) {
	return _OVMStateCommitmentChain.Contract.Resolve(&_OVMStateCommitmentChain.CallOpts, _name)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) Resolve(_name string) (common.Address, error) {
	return _OVMStateCommitmentChain.Contract.Resolve(&_OVMStateCommitmentChain.CallOpts, _name)
}

// VerifyStateCommitment is a free data retrieval call binding the contract method 0x4d69ee57.
//
// Solidity: function verifyStateCommitment(bytes32 _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) VerifyStateCommitment(opts *bind.CallOpts, _element [32]byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "verifyStateCommitment", _element, _batchHeader, _proof)
	return *ret0, err
}

// VerifyStateCommitment is a free data retrieval call binding the contract method 0x4d69ee57.
//
// Solidity: function verifyStateCommitment(bytes32 _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) VerifyStateCommitment(_element [32]byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMStateCommitmentChain.Contract.VerifyStateCommitment(&_OVMStateCommitmentChain.CallOpts, _element, _batchHeader, _proof)
}

// VerifyStateCommitment is a free data retrieval call binding the contract method 0x4d69ee57.
//
// Solidity: function verifyStateCommitment(bytes32 _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) VerifyStateCommitment(_element [32]byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMStateCommitmentChain.Contract.VerifyStateCommitment(&_OVMStateCommitmentChain.CallOpts, _element, _batchHeader, _proof)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x8ca5cbb9.
//
// Solidity: function appendStateBatch(bytes32[] _batch, uint256 _shouldStartAtElement) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactor) AppendStateBatch(opts *bind.TransactOpts, _batch [][32]byte, _shouldStartAtElement *big.Int) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.contract.Transact(opts, "appendStateBatch", _batch, _shouldStartAtElement)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x8ca5cbb9.
//
// Solidity: function appendStateBatch(bytes32[] _batch, uint256 _shouldStartAtElement) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) AppendStateBatch(_batch [][32]byte, _shouldStartAtElement *big.Int) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.AppendStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batch, _shouldStartAtElement)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x8ca5cbb9.
//
// Solidity: function appendStateBatch(bytes32[] _batch, uint256 _shouldStartAtElement) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorSession) AppendStateBatch(_batch [][32]byte, _shouldStartAtElement *big.Int) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.AppendStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batch, _shouldStartAtElement)
}

// DeleteStateBatch is a paid mutator transaction binding the contract method 0xb8e189ac.
//
// Solidity: function deleteStateBatch(Lib_OVMCodecChainBatchHeader _batchHeader) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactor) DeleteStateBatch(opts *bind.TransactOpts, _batchHeader Lib_OVMCodecChainBatchHeader) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.contract.Transact(opts, "deleteStateBatch", _batchHeader)
}

// DeleteStateBatch is a paid mutator transaction binding the contract method 0xb8e189ac.
//
// Solidity: function deleteStateBatch(Lib_OVMCodecChainBatchHeader _batchHeader) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) DeleteStateBatch(_batchHeader Lib_OVMCodecChainBatchHeader) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.DeleteStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batchHeader)
}

// DeleteStateBatch is a paid mutator transaction binding the contract method 0xb8e189ac.
//
// Solidity: function deleteStateBatch(Lib_OVMCodecChainBatchHeader _batchHeader) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorSession) DeleteStateBatch(_batchHeader Lib_OVMCodecChainBatchHeader) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.DeleteStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batchHeader)
}

// SetLastOverwritableIndex is a paid mutator transaction binding the contract method 0x2979761b.
//
// Solidity: function setLastOverwritableIndex(Lib_OVMCodecChainBatchHeader _stateBatchHeader, Lib_OVMCodecTransaction _transaction, Lib_OVMCodecTransactionChainElement _txChainElement, Lib_OVMCodecChainBatchHeader _txBatchHeader, Lib_OVMCodecChainInclusionProof _txInclusionProof) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactor) SetLastOverwritableIndex(opts *bind.TransactOpts, _stateBatchHeader Lib_OVMCodecChainBatchHeader, _transaction Lib_OVMCodecTransaction, _txChainElement Lib_OVMCodecTransactionChainElement, _txBatchHeader Lib_OVMCodecChainBatchHeader, _txInclusionProof Lib_OVMCodecChainInclusionProof) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.contract.Transact(opts, "setLastOverwritableIndex", _stateBatchHeader, _transaction, _txChainElement, _txBatchHeader, _txInclusionProof)
}

// SetLastOverwritableIndex is a paid mutator transaction binding the contract method 0x2979761b.
//
// Solidity: function setLastOverwritableIndex(Lib_OVMCodecChainBatchHeader _stateBatchHeader, Lib_OVMCodecTransaction _transaction, Lib_OVMCodecTransactionChainElement _txChainElement, Lib_OVMCodecChainBatchHeader _txBatchHeader, Lib_OVMCodecChainInclusionProof _txInclusionProof) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) SetLastOverwritableIndex(_stateBatchHeader Lib_OVMCodecChainBatchHeader, _transaction Lib_OVMCodecTransaction, _txChainElement Lib_OVMCodecTransactionChainElement, _txBatchHeader Lib_OVMCodecChainBatchHeader, _txInclusionProof Lib_OVMCodecChainInclusionProof) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.SetLastOverwritableIndex(&_OVMStateCommitmentChain.TransactOpts, _stateBatchHeader, _transaction, _txChainElement, _txBatchHeader, _txInclusionProof)
}

// SetLastOverwritableIndex is a paid mutator transaction binding the contract method 0x2979761b.
//
// Solidity: function setLastOverwritableIndex(Lib_OVMCodecChainBatchHeader _stateBatchHeader, Lib_OVMCodecTransaction _transaction, Lib_OVMCodecTransactionChainElement _txChainElement, Lib_OVMCodecChainBatchHeader _txBatchHeader, Lib_OVMCodecChainInclusionProof _txInclusionProof) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorSession) SetLastOverwritableIndex(_stateBatchHeader Lib_OVMCodecChainBatchHeader, _transaction Lib_OVMCodecTransaction, _txChainElement Lib_OVMCodecTransactionChainElement, _txBatchHeader Lib_OVMCodecChainBatchHeader, _txInclusionProof Lib_OVMCodecChainInclusionProof) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.SetLastOverwritableIndex(&_OVMStateCommitmentChain.TransactOpts, _stateBatchHeader, _transaction, _txChainElement, _txBatchHeader, _txInclusionProof)
}

