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

// OVMStateCommitmentChainABI is the input ABI used to generate the binding from.
const OVMStateCommitmentChainABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_libAddressManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"FRAUD_PROOF_WINDOW\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_batch\",\"type\":\"bytes32[]\"}],\"name\":\"appendStateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"}],\"name\":\"deleteStateBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBatches\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalBatches\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalElements\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalElements\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"}],\"name\":\"insideFraudProofWindow\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_inside\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_element\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"verifyElement\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_verified\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OVMStateCommitmentChainBin is the compiled bytecode used for deploying new contracts.
var OVMStateCommitmentChainBin = "0x60806040523480156200001157600080fd5b50604051620017f4380380620017f4833981016040819052620000349162000201565b806200005960046002620186a060006200013260201b620004f517909392919060201c565b600780546001600160a01b0319166001600160a01b039290921691909117905560408051808201909152601d81527f4f564d5f43616e6f6e6963616c5472616e73616374696f6e436861696e0000006020820152620000b89062000174565b600880546001600160a01b0319166001600160a01b039290921691909117905560408051808201909152601181527027ab26afa33930bab22b32b934b334b2b960791b60208201526200010b9062000174565b600980546001600160a01b0319166001600160a01b03929092169190911790555062000287565b60028401805463ffffffff191663ffffffff9485161763ffffffff60201b19166401000000009390941692909202929092179055600482015542600390910155565b60075460405163bf40fac160e01b81526000916001600160a01b03169063bf40fac190620001a790859060040162000231565b60206040518083038186803b158015620001c057600080fd5b505afa158015620001d5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001fb919062000201565b92915050565b60006020828403121562000213578081fd5b81516001600160a01b03811681146200022a578182fd5b9392505050565b6000602080835283518082850152825b818110156200025f5785810183015185820160400152820162000241565b81811115620002715783604083870101525b50601f01601f1916929092016040019392505050565b61155d80620002976000396000f3fe608060405234801561001057600080fd5b50600436106100885760003560e01c806398fe87c81161005b57806398fe87c814610100578063b8e189ac14610113578063c17b291b14610126578063e561dddc1461012e57610088565b8063461a44781461008d5780634bb10367146100b65780637aa63a86146100cb5780639418bddd146100e0575b600080fd5b6100a061009b366004611079565b610136565b6040516100ad9190611140565b60405180910390f35b6100c96100c4366004610f8e565b6101bf565b005b6100d3610374565b6040516100ad91906110f5565b6100f36100ee3660046110ab565b610388565b6040516100ad9190611154565b6100f361010e366004610fc8565b6103d2565b6100c96101213660046110ab565b61047c565b6100d36104d7565b6100d36104de565b60075460405163bf40fac160e01b81526000916001600160a01b03169063bf40fac19061016790859060040161115f565b60206040518083038186803b15801561017f57600080fd5b505afa158015610193573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906101b79190610f67565b90505b919050565b60008151116101e95760405162461bcd60e51b81526004016101e090611460565b60405180910390fd5b600860009054906101000a90046001600160a01b03166001600160a01b0316637aa63a866040518163ffffffff1660e01b815260040160206040518083038186803b15801561023757600080fd5b505afa15801561024b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061026f91906110dd565b8151610279610374565b0111156102985760405162461bcd60e51b81526004016101e090611200565b606081516001600160401b03811180156102b157600080fd5b506040519080825280602002602001820160405280156102e557816020015b60608152602001906001900390816102d05790505b50905060005b82518110156103465782818151811061030057fe5b602002602001015160405160200161031891906110f5565b60405160208183030381529060405282828151811061033357fe5b60209081029190910101526001016102eb565b50610370814260405160200161035c91906110f5565b604051602081830303815290604052610538565b5050565b60006103806000610590565b60201c905090565b60008082608001518060200190518101906103a391906110dd565b9050806103c25760405162461bcd60e51b81526004016101e09061141b565b4262093a80820111915050919050565b81516000906103e290829061059e565b6103eb84610655565b146104085760405162461bcd60e51b81526004016101e090611331565b610455836020015186868080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152505086516020880151909250905061069b565b6104715760405162461bcd60e51b81526004016101e09061126f565b506001949350505050565b6009546001600160a01b031633146104a65760405162461bcd60e51b81526004016101e0906113be565b6104af81610388565b6104cb5760405162461bcd60e51b81526004016101e090611360565b6104d48161070d565b50565b62093a8081565b60006104ea6000610793565b63ffffffff16905090565b60028401805463ffffffff191663ffffffff9485161767ffffffff0000000019166401000000009390941692909202929092179055600482015542600390910155565b610540610d8c565b6040518060a001604052806105556000610793565b63ffffffff168152602001610569856107a2565b8152845160208201526006546040820152606001839052905061058b81610841565b505050565b6001015463ffffffff191690565b6000806105ae846001015461086c565b63ffffffff169050808363ffffffff16106105db5760405162461bcd60e51b81526004016101e090611192565b600284015463ffffffff600160401b820481169181169190910381169084168203111561061a5760405162461bcd60e51b81526004016101e0906112a6565b6002840154849060009063ffffffff9081169086168161063657fe5b0663ffffffff1681526020019081526020016000205491505092915050565b6000816020015182604001518360600151846080015160405160200161067e949392919061110c565b604051602081830303815290604052805190602001209050919050565b82516020840120600090815b83518110156107015760008482815181106106be57fe5b60209081029190910101519050600186831c81161480156106ea576106e38483610875565b93506106f7565b6106f48285610875565b93505b50506001016106a7565b50909414949350505050565b6107176000610793565b63ffffffff1681600001511061073f5760405162461bcd60e51b81526004016101e0906114a3565b805161074d9060009061059e565b61075682610655565b146107735760405162461bcd60e51b81526004016101e090611331565b6060810151600681905581516104d491600091600019019060201b6108a8565b60006101b7826001015461086c565b6000606082516001600160401b03811180156107bd57600080fd5b506040519080825280602002602001820160405280156107e7578160200160208202803683370190505b50905060005b83518110156108305783818151811061080257fe5b60200260200101518051906020012082828151811061081d57fe5b60209081029190910101526001016107ed565b5061083a81610977565b9392505050565b600061084c82610655565b905061037081836040015161085f610374565b600092910160201b610b8c565b63ffffffff1690565b6000828260405160200161088a9291906110fe565b60405160208183030381529060405280519060200120905092915050565b60006108b7846001015461086c565b90506000836001019050600081830390506000818760020160089054906101000a900463ffffffff160190508363ffffffff168363ffffffff161061090e5760405162461bcd60e51b81526004016101e090611192565b600287015463ffffffff908116908216111561093c5760405162461bcd60e51b81526004016101e0906112ea565b60028701805463ffffffff60401b1916600160401b63ffffffff8416021790556109668386610bde565b876001018190555050505050505050565b6000808251116109995760405162461bcd60e51b81526004016101e0906111bc565b8151600114156109bf57816000815181106109b057fe5b602002602001015190506101ba565b60606109cb8351610bf4565b835190915083906002900660011415610a665783516001016001600160401b03811180156109f857600080fd5b50604051908082528060200260200182016040528015610a22578160200160208202803683370190505b50905060005b8451811015610a6457848181518110610a3d57fe5b6020026020010151828281518110610a5157fe5b6020908102919091010152600101610a28565b505b83516000906002810660011415610aa657838281518110610a8357fe5b6020026020010151838281518110610a9757fe5b60209081029190910101526001015b6001811115610b6c5760018201915060005b60028204811015610b1b57610afc848260020281518110610ad557fe5b6020026020010151858360020260010181518110610aef57fe5b6020026020010151610875565b848281518110610b0857fe5b6020908102919091010152600101610ab8565b50600290046001808216148015610b33575080600114155b15610b6757838281518110610b4457fe5b6020026020010151838281518110610b5857fe5b60209081029190910101526001015b610aa6565b82600081518110610b7957fe5b6020026020010151945050505050919050565b610b97838383610ce9565b6002830154600160401b900463ffffffff161561058b5760028301805463ffffffff600160401b80830482166001019091160263ffffffff60401b19909116179055505050565b63ffffffff19811663ffffffff83161792915050565b606080826001600160401b0381118015610c0d57600080fd5b50604051908082528060200260200182016040528015610c37578160200160208202803683370190505b5090506000604051602001610c4c91906110f5565b6040516020818303038152906040528051906020012081600081518110610c6f57fe5b602090810291909101015260015b8151811015610ce257816001820381518110610c9557fe5b6020026020010151604051602001610cad91906110f5565b60405160208183030381529060405280519060200120828281518110610ccf57fe5b6020908102919091010152600101610c7d565b5092915050565b6000610cf8846001015461086c565b600285015463ffffffff91821692501680821415610d51578460040154856003015401421015610d51575060028401805463ffffffff8082166401000000008304821601811663ffffffff199092169190911791829055165b83856000838581610d5e57fe5b068152602081019190915260400160002055610d7d6001830184610bde565b85600101819055505050505050565b6040518060a0016040528060008152602001600080191681526020016000815260200160008152602001606081525090565b600082601f830112610dce578081fd5b81356001600160401b03811115610de3578182fd5b6020808202610df38282016114d1565b83815293508184018583018287018401881015610e0f57600080fd5b600092505b84831015610e32578035825260019290920191908301908301610e14565b505050505092915050565b600082601f830112610e4d578081fd5b81356001600160401b03811115610e62578182fd5b610e75601f8201601f19166020016114d1565b9150808252836020828501011115610e8c57600080fd5b8060208401602084013760009082016020015292915050565b600060a08284031215610eb6578081fd5b610ec060a06114d1565b90508135815260208201356020820152604082013560408201526060820135606082015260808201356001600160401b03811115610efd57600080fd5b610f0984828501610e3d565b60808301525092915050565b600060408284031215610f26578081fd5b610f3060406114d1565b90508135815260208201356001600160401b03811115610f4f57600080fd5b610f5b84828501610dbe565b60208301525092915050565b600060208284031215610f78578081fd5b81516001600160a01b038116811461083a578182fd5b600060208284031215610f9f578081fd5b81356001600160401b03811115610fb4578182fd5b610fc084828501610dbe565b949350505050565b60008060008060608587031215610fdd578283fd5b84356001600160401b0380821115610ff3578485fd5b818701915087601f830112611006578485fd5b813581811115611014578586fd5b886020828501011115611025578586fd5b60209283019650945090860135908082111561103f578384fd5b61104b88838901610ea5565b93506040870135915080821115611060578283fd5b5061106d87828801610f15565b91505092959194509250565b60006020828403121561108a578081fd5b81356001600160401b0381111561109f578182fd5b610fc084828501610e3d565b6000602082840312156110bc578081fd5b81356001600160401b038111156110d1578182fd5b610fc084828501610ea5565b6000602082840312156110ee578081fd5b5051919050565b90815260200190565b918252602082015260400190565b600085825284602083015283604083015282516111308160608501602087016114f7565b9190910160600195945050505050565b6001600160a01b0391909116815260200190565b901515815260200190565b600060208252825180602084015261117e8160408501602087016114f7565b601f01601f19169190910160400192915050565b60208082526010908201526f24b73232bc103a37b7903630b933b29760811b604082015260600190565b60208082526024908201527f4d7573742070726f76696465206174206c65617374206f6e65206c656166206860408201526330b9b41760e11b606082015260800190565b60208082526049908201527f4e756d626572206f6620737461746520726f6f74732063616e6e6f742065786360408201527f65656420746865206e756d626572206f662063616e6f6e6963616c207472616e60608201526839b0b1ba34b7b7399760b91b608082015260a00190565b60208082526018908201527f496e76616c696420696e636c7573696f6e2070726f6f662e0000000000000000604082015260600190565b60208082526024908201527f496e64657820746f6f206f6c64202620686173206265656e206f7665727269646040820152633232b71760e11b606082015260800190565b60208082526027908201527f417474656d7074696e6720746f2064656c65746520746f6f206d616e7920656c60408201526632b6b2b73a399760c91b606082015260800190565b60208082526015908201527424b73b30b634b2103130ba31b4103432b0b232b91760591b604082015260600190565b602080825260409082018190527f537461746520626174636865732063616e206f6e6c792062652064656c657465908201527f642077697468696e207468652066726175642070726f6f662077696e646f772e606082015260800190565b6020808252603b908201527f537461746520626174636865732063616e206f6e6c792062652064656c65746560408201527f6420627920746865204f564d5f467261756456657269666965722e0000000000606082015260800190565b60208082526025908201527f4261746368206865616465722074696d657374616d702063616e6e6f74206265604082015264207a65726f60d81b606082015260800190565b60208082526023908201527f43616e6e6f74207375626d697420616e20656d7074792073746174652062617460408201526231b41760e91b606082015260800190565b60208082526014908201527324b73b30b634b2103130ba31b41034b73232bc1760611b604082015260600190565b6040518181016001600160401b03811182821017156114ef57600080fd5b604052919050565b60005b838110156115125781810151838201526020016114fa565b83811115611521576000848401525b5050505056fea2646970667358221220413e0f41d29705dc421437240cc99460980bba41f5dea2ab8961521770bc1e4c64736f6c63430007000033"

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

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCaller) VerifyElement(opts *bind.CallOpts, _element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _OVMStateCommitmentChain.contract.Call(opts, out, "verifyElement", _element, _batchHeader, _proof)
	return *ret0, err
}

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) VerifyElement(_element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMStateCommitmentChain.Contract.VerifyElement(&_OVMStateCommitmentChain.CallOpts, _element, _batchHeader, _proof)
}

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMStateCommitmentChain *OVMStateCommitmentChainCallerSession) VerifyElement(_element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMStateCommitmentChain.Contract.VerifyElement(&_OVMStateCommitmentChain.CallOpts, _element, _batchHeader, _proof)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x4bb10367.
//
// Solidity: function appendStateBatch(bytes32[] _batch) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactor) AppendStateBatch(opts *bind.TransactOpts, _batch [][32]byte) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.contract.Transact(opts, "appendStateBatch", _batch)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x4bb10367.
//
// Solidity: function appendStateBatch(bytes32[] _batch) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainSession) AppendStateBatch(_batch [][32]byte) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.AppendStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batch)
}

// AppendStateBatch is a paid mutator transaction binding the contract method 0x4bb10367.
//
// Solidity: function appendStateBatch(bytes32[] _batch) returns()
func (_OVMStateCommitmentChain *OVMStateCommitmentChainTransactorSession) AppendStateBatch(_batch [][32]byte) (*types.Transaction, error) {
	return _OVMStateCommitmentChain.Contract.AppendStateBatch(&_OVMStateCommitmentChain.TransactOpts, _batch)
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

