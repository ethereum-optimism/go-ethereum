// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package canonicaltransactionchain

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

// Lib_OVMCodecQueueElement is an auto generated low-level Go binding around an user-defined struct.
type Lib_OVMCodecQueueElement struct {
	QueueRoot   [32]byte
	Timestamp   *big.Int
	BlockNumber uint32
}

// OVMCanonicalTransactionChainABI is the input ABI used to generate the binding from.
const OVMCanonicalTransactionChainABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_libAddressManager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_forceInclusionPeriodSeconds\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_startingQueueIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_numQueueElements\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalElements\",\"type\":\"uint256\"}],\"name\":\"QueueBatchAppended\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_startingQueueIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_numQueueElements\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_totalElements\",\"type\":\"uint256\"}],\"name\":\"SequencerBatchAppended\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_l1TxOrigin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_queueIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_timestamp\",\"type\":\"uint256\"}],\"name\":\"TransactionEnqueued\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"L2_GAS_DISCOUNT_DIVISOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_ROLLUP_TX_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MIN_ROLLUP_TX_GAS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_numQueuedTransactions\",\"type\":\"uint256\"}],\"name\":\"appendQueueBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"appendSequencerBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_gasLimit\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"enqueue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getQueueElement\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"queueRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint40\",\"name\":\"timestamp\",\"type\":\"uint40\"},{\"internalType\":\"uint32\",\"name\":\"blockNumber\",\"type\":\"uint32\"}],\"internalType\":\"structLib_OVMCodec.QueueElement\",\"name\":\"_element\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalBatches\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalBatches\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTotalElements\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"_totalElements\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_name\",\"type\":\"string\"}],\"name\":\"resolve\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_contract\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_element\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"batchIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"batchRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"prevTotalElements\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"internalType\":\"structLib_OVMCodec.ChainBatchHeader\",\"name\":\"_batchHeader\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"siblings\",\"type\":\"bytes32[]\"}],\"internalType\":\"structLib_OVMCodec.ChainInclusionProof\",\"name\":\"_proof\",\"type\":\"tuple\"}],\"name\":\"verifyElement\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"_verified\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OVMCanonicalTransactionChainBin is the compiled bytecode used for deploying new contracts.
var OVMCanonicalTransactionChainBin = "0x60806040523480156200001157600080fd5b506040516200206f3803806200206f833981016040819052620000349162000217565b816200005960046002620186a060006200011860201b620008f917909392919060201c565b600780546001600160a01b0319166001600160a01b039290921691909117905560408051808201909152600d81526c27ab26afa9b2b8bab2b731b2b960991b6020820152620000a8906200015a565b600a80546001600160a01b0319166001600160a01b03929092169190911790556008819055620000ef600b606460326402540be40062000118602090811b620008f917901c565b62000110600260326000806200011860201b620008f917909392919060201c565b5050620002a7565b60028401805463ffffffff191663ffffffff9485161763ffffffff60201b19166401000000009390941692909202929092179055600482015542600390910155565b60075460405163bf40fac160e01b81526000916001600160a01b03169063bf40fac1906200018d90859060040162000251565b60206040518083038186803b158015620001a657600080fd5b505afa158015620001bb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620001e19190620001e7565b92915050565b600060208284031215620001f9578081fd5b81516001600160a01b038116811462000210578182fd5b9392505050565b600080604083850312156200022a578081fd5b82516001600160a01b038116811462000241578182fd5b6020939093015192949293505050565b6000602080835283518082850152825b818110156200027f5785810183015185820160400152820162000261565b81811115620002915783604083870101525b50601f01601f1916929092016040019392505050565b611db880620002b76000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c8063876ed5cb11610071578063876ed5cb1461012957806398fe87c814610131578063c2cf696f14610151578063d0f8934414610159578063e561dddc14610161578063facdc5da14610169576100a9565b80632a7f18be146100ae578063461a4478146100d75780636fee07e0146100f757806378f4b2f21461010c5780637aa63a8614610121575b600080fd5b6100c16100bc366004611591565b61017c565b6040516100ce9190611caa565b60405180910390f35b6100ea6100e5366004611556565b6101d9565b6040516100ce9190611617565b61010a61010536600461144d565b610260565b005b61011461038f565b6040516100ce9190611cd8565b610114610395565b6101146103ae565b61014461013f3660046114a4565b6103b4565b6040516100ce91906116b6565b61011461045e565b61010a610463565b6101146107cf565b61010a610177366004611591565b6107e6565b61018461122c565b600282026000610195600b8361093c565b905060006101a7600b6001850161093c565b6040805160608101825293845264ffffffffff8216602085015260289190911c63ffffffff1690830152509392505050565b60075460405163bf40fac160e01b81526000916001600160a01b03169063bf40fac19061020a9085906004016116f9565b60206040518083038186803b15801561022257600080fd5b505afa158015610236573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025a9190611431565b92915050565b6127108151111561028c5760405162461bcd60e51b815260040161028390611873565b60405180910390fd5b614e208210156102ae5760405162461bcd60e51b815260040161028390611bba565b600a820460005a90508181116102d65760405162461bcd60e51b815260040161028390611a98565b60005b825a830310156102eb576001016102d9565b606033878787604051602001610304949392919061162b565b60408051601f1981840301815291905280516020820120909150424360281b17610332600b838360006109f3565b600061033c610a78565b9150507f4b388aecf9fa6cc92253704e5975a6129a4f735bdbd99567df4ed0094ee4ceb5338b8b8b600186034260405161037b96959493929190611668565b60405180910390a150505050505050505050565b614e2081565b6000806103a0610a78565b5064ffffffffff1691505090565b61271081565b81516000906103c490829061093c565b6103cd84610a9f565b146103ea5760405162461bcd60e51b815260040161028390611a69565b610437836020015186868080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525050865160208801519092509050610ae5565b6104535760405162461bcd60e51b81526004016102839061183c565b506001949350505050565b600a81565b60043560d81c60093560e890811c90600c35901c61047f6107cf565b8364ffffffffff16146104a45760405162461bcd60e51b815260040161028390611977565b600a546001600160a01b031633146104ce5760405162461bcd60e51b8152600401610283906119d4565b60008162ffffff16116104f35760405162461bcd60e51b815260040161028390611736565b60008262ffffff16116105185760405162461bcd60e51b815260040161028390611b79565b60608262ffffff1667ffffffffffffffff8111801561053657600080fd5b50604051908082528060200260200182016040528015610560578160200160208202803683370190505b509050600080600f601062ffffff861602013663ffffffff82168110156105995760405162461bcd60e51b815260040161028390611c68565b60006105a3610a78565b91505060005b8762ffffff168163ffffffff161015610716576105c461124c565b6105d38263ffffffff16610b57565b90506105df8184610ba7565b60005b815163ffffffff821610156106ba57853560e81c60606042820167ffffffffffffffff8111801561061257600080fd5b506040519080825280601f01601f19166020018201604052801561063d576020820181803683370190505b50604085015160608601519192506000916020840160018153600060018201538260028201528160228201528560038d016042830137856042018120935050828e8e63ffffffff168151811061068f57fe5b60209081029190910101525050506001998a01999889019897909101600301969190910190506105e2565b5060005b81602001518163ffffffff16101561070c576106df8463ffffffff16610c6a565b898963ffffffff16815181106106f157fe5b602090810291909101015260019788019793840193016106be565b50506001016105a9565b508263ffffffff16821461073c5760405162461bcd60e51b81526004016102839061177e565b8762ffffff168563ffffffff16146107665760405162461bcd60e51b8152600401610283906117cc565b63ffffffff62ffffff89168590031661078d61078188610d0e565b8a62ffffff1683610efc565b7f602f1aeac0ca2e7a13e281a9ef0ad7838542712ce16780fa2ecffd351f05f899818363ffffffff1603826107c0610395565b60405161037b93929190611ce1565b60006107db6000610f9a565b63ffffffff16905090565b600081116108065760405162461bcd60e51b815260040161028390611a21565b6000610810610a78565b91505060608267ffffffffffffffff8111801561082c57600080fd5b50604051908082528060200260200182016040528015610856578160200160208202803683370190505b50905060005b83811015610896576108738363ffffffff16610c6a565b82828151811061087f57fe5b60209081029190910101526001928301920161085c565b506108aa6108a382610d0e565b8485610efc565b7f64d7f508348c70dea42d5302a393987e4abc20e45954ab3f9d320207751956f0838363ffffffff1603846108dd610395565b6040516108ec93929190611ce1565b60405180910390a1505050565b60028401805463ffffffff191663ffffffff9485161767ffffffff0000000019166401000000009390941692909202929092179055600482015542600390910155565b60008061094c8460010154610fa5565b63ffffffff169050808363ffffffff16106109795760405162461bcd60e51b81526004016102839061170c565b600284015463ffffffff600160401b82048116918116919091038116908416820311156109b85760405162461bcd60e51b8152600401610283906118c0565b6002840154849060009063ffffffff908116908616816109d457fe5b0663ffffffff1681526020019081526020016000205491505092915050565b6109ff84848484610fae565b6002840154600160401b900463ffffffff1615610a72576002840154600160401b900463ffffffff16600114610a4c5760028460020160089054906101000a900463ffffffff1603610a4f565b60005b8460020160086101000a81548163ffffffff021916908363ffffffff1602179055505b50505050565b6000806000610a8760006110a6565b602081901c64ffffffffff16935060481c9150509091565b60008160200151826040015183606001518460800151604051602001610ac894939291906115e3565b604051602081830303815290604052805190602001209050919050565b82516020840120600090815b8351811015610b4b576000848281518110610b0857fe5b60209081029190910101519050600186831c8116148015610b3457610b2d84836110b4565b9350610b41565b610b3e82856110b4565b93505b5050600101610af1565b50909414949350505050565b610b5f61124c565b5060408051608081018252601092909202600f81013560e890811c84526012820135901c6020840152601581013560d890811c92840192909252601a0135901c606082015290565b610bb1600b610f9a565b63ffffffff16610bc057610c66565b610bc861122c565b610bd78263ffffffff1661017c565b9050600854816020015164ffffffffff16014210610c075760405162461bcd60e51b815260040161028390611bff565b806020015164ffffffffff1683604001511115610c365760405162461bcd60e51b815260040161028390611ae3565b806040015163ffffffff1683606001511115610c645760405162461bcd60e51b815260040161028390611b2d565b505b5050565b6000610c7461122c565b610c7d8361017c565b600a549091506001600160a01b0316331480610ca8575042600854826020015164ffffffffff160111155b610cc45760405162461bcd60e51b815260040161028390611904565b610d076040518060a001604052806000151581526020018581526020016000815260200160008152602001604051806020016040528060008152508152506110e7565b9392505050565b60006001815b8351821015610d2957600191821b9101610d14565b60408051818152606081810183529160208201818036833701905050905060008060005b6002885181610d5857fe5b04811015610dc857878160020281518110610d6f57fe5b60200260200101519250878160020260010181518110610d8b57fe5b602002602001015191508260208501528160408501528380519060200120888281518110610db557fe5b6020908102919091010152600101610d4d565b5086515b600186901c811015610dfb576000801b888281518110610de857fe5b6020908102919091010152600101610dcc565b50600060028851870381610e0b57fe5b049050600186901c60025b868111610ed85760019190911c9060028304925082820360005b81811015610ea0578b8160020281518110610e4757fe5b602002602001015196508b8160020260010181518110610e6357fe5b6020026020010151955086602089015285604089015287805190602001208c8281518110610e8d57fe5b6020908102919091010152600101610e30565b50805b83811015610ece576000801b8c8281518110610ebb57fe5b6020908102919091010152600101610ea3565b5050600101610e16565b5088600081518110610ee657fe5b6020026020010151975050505050505050919050565b600080610f07610a78565b91509150610f13611274565b6040518060a00160405280610f286000610f9a565b63ffffffff1681526020018781526020018681526020018464ffffffffff1681526020016040518060200160405280600081525081525090506000610f6c82610a9f565b90506000610f8283604001518601878601611112565b9050610f906000838361111c565b5050505050505050565b600061025a82600101545b63ffffffff1690565b6000610fbd8560010154610fa5565b600286015463ffffffff918216925016600182018111611044578560040154866003015401421015611044576002860154600164010000000090910463ffffffff161161100b57600261101f565b6002860154640100000000900463ffffffff165b60028701805463ffffffff808216909301831663ffffffff1990911617908190551690505b8486600083858161105157fe5b068152602001908152602001600020819055508386600001600083856001018161107757fe5b0681526020810191909152604001600020556110966002830184611173565b8660010181905550505050505050565b6001015463ffffffff191690565b600082826040516020016110c99291906115d5565b60405160208183030381529060405280519060200120905092915050565b8051602080830151604080850151606086015160808701519251600096610ac89690959491016116c1565b60281b1760201b90565b611127838383611189565b6002830154600160401b900463ffffffff1615610c645760028301805463ffffffff600160401b8083048216600101909116026bffffffff000000000000000019909116179055505050565b63ffffffff19811663ffffffff83161792915050565b60006111988460010154610fa5565b600285015463ffffffff918216925016808214156111f15784600401548560030154014210156111f1575060028401805463ffffffff8082166401000000008304821601811663ffffffff199092169190911791829055165b838560008385816111fe57fe5b06815260208101919091526040016000205561121d6001830184611173565b85600101819055505050505050565b604080516060810182526000808252602082018190529181019190915290565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b6040518060a0016040528060008152602001600080191681526020016000815260200160008152602001606081525090565b600082601f8301126112b6578081fd5b813567ffffffffffffffff8111156112cc578182fd5b6112df601f8201601f1916602001611cf7565b91508082528360208285010111156112f657600080fd5b8060208401602084013760009082016020015292915050565b600060a08284031215611320578081fd5b61132a60a0611cf7565b905081358152602082013560208201526040820135604082015260608201356060820152608082013567ffffffffffffffff81111561136857600080fd5b611374848285016112a6565b60808301525092915050565b600060408284031215611391578081fd5b61139b6040611cf7565b90508135815260208083013567ffffffffffffffff8111156113bc57600080fd5b8301601f810185136113cd57600080fd5b80356113e06113db82611d1e565b611cf7565b81815283810190838501858402850186018910156113fd57600080fd5b600094505b83851015611420578035835260019490940193918501918501611402565b508085870152505050505092915050565b600060208284031215611442578081fd5b8151610d0781611d6a565b600080600060608486031215611461578182fd5b833561146c81611d6a565b925060208401359150604084013567ffffffffffffffff81111561148e578182fd5b61149a868287016112a6565b9150509250925092565b600080600080606085870312156114b9578081fd5b843567ffffffffffffffff808211156114d0578283fd5b818701915087601f8301126114e3578283fd5b8135818111156114f1578384fd5b886020828501011115611502578384fd5b60209283019650945090860135908082111561151c578283fd5b6115288883890161130f565b9350604087013591508082111561153d578283fd5b5061154a87828801611380565b91505092959194509250565b600060208284031215611567578081fd5b813567ffffffffffffffff81111561157d578182fd5b611589848285016112a6565b949350505050565b6000602082840312156115a2578081fd5b5035919050565b600081518084526115c1816020860160208601611d3e565b601f01601f19169290920160200192915050565b918252602082015260400190565b60008582528460208301528360408301528251611607816060850160208701611d3e565b9190910160600195945050505050565b6001600160a01b0391909116815260200190565b6001600160a01b038581168252841660208201526040810183905260806060820181905260009061165e908301846115a9565b9695505050505050565b6001600160a01b038781168252861660208201526040810185905260c06060820181905260009061169b908301866115a9565b63ffffffff9490941660808301525060a00152949350505050565b901515815260200190565b6000861515825285602083015284604083015283606083015260a060808301526116ee60a08301846115a9565b979650505050505050565b600060208252610d0760208301846115a9565b60208082526010908201526f24b73232bc103a37b7903630b933b29760811b604082015260600190565b60208082526028908201527f4d7573742070726f76696465206174206c65617374206f6e652062617463682060408201526731b7b73a32bc3a1760c11b606082015260800190565b6020808252602e908201527f4e6f7420616c6c2073657175656e636572207472616e73616374696f6e73207760408201526d32b93290383937b1b2b9b9b2b21760911b606082015260800190565b6020808252604a908201527f41637475616c207472616e73616374696f6e20696e64657820646f6573206e6f60408201527f74206d6174636820657870656374656420746f74616c20656c656d656e7473206060820152693a379030b83832b7321760b11b608082015260a00190565b60208082526018908201527f496e76616c696420696e636c7573696f6e2070726f6f662e0000000000000000604082015260600190565b6020808252602d908201527f5472616e73616374696f6e2065786365656473206d6178696d756d20726f6c6c60408201526c3ab8103230ba309039b4bd329760991b606082015260800190565b60208082526024908201527f496e64657820746f6f206f6c64202620686173206265656e206f7665727269646040820152633232b71760e11b606082015260800190565b6020808252604d908201527f5175657565207472616e73616374696f6e732063616e6e6f742062652073756260408201527f6d697474656420647572696e67207468652073657175656e63657220696e636c60608201526c3ab9b4b7b7103832b934b7b21760991b608082015260a00190565b6020808252603d908201527f41637475616c20626174636820737461727420696e64657820646f6573206e6f60408201527f74206d6174636820657870656374656420737461727420696e6465782e000000606082015260800190565b6020808252602d908201527f46756e6374696f6e2063616e206f6e6c792062652063616c6c6564206279207460408201526c34329029b2b8bab2b731b2b91760991b606082015260800190565b60208082526028908201527f4d75737420617070656e64206d6f7265207468616e207a65726f207472616e7360408201526730b1ba34b7b7399760c11b606082015260800190565b60208082526015908201527424b73b30b634b2103130ba31b4103432b0b232b91760591b604082015260600190565b6020808252602b908201527f496e73756666696369656e742067617320666f72204c322072617465206c696d60408201526a34ba34b73390313ab9371760a91b606082015260800190565b6020808252602a908201527f53657175656e636572207472616e73616374696f6e732074696d657374616d70604082015269103a37b7903434b3b41760b11b606082015260800190565b6020808252602c908201527f53657175656e636572207472616e73616374696f6e7320626c6f636b4e756d6260408201526b32b9103a37b7903434b3b41760a11b606082015260800190565b60208082526021908201527f4d75737420617070656e64206174206c65617374206f6e6520656c656d656e746040820152601760f91b606082015260800190565b60208082526025908201527f4c61796572203220676173206c696d697420746f6f206c6f7720746f20656e716040820152643ab2bab29760d91b606082015260800190565b60208082526043908201527f4f6c6465722071756575652062617463686573206d7573742062652070726f6360408201527f6573736564206265666f72652061206e65772073657175656e6365722062617460608201526231b41760e91b608082015260a00190565b60208082526022908201527f4e6f7420656e6f756768204261746368436f6e74657874732070726f76696465604082015261321760f11b606082015260800190565b8151815260208083015164ffffffffff169082015260409182015163ffffffff169181019190915260600190565b90815260200190565b9283526020830191909152604082015260600190565b60405181810167ffffffffffffffff81118282101715611d1657600080fd5b604052919050565b600067ffffffffffffffff821115611d34578081fd5b5060209081020190565b60005b83811015611d59578181015183820152602001611d41565b83811115610a725750506000910152565b6001600160a01b0381168114611d7f57600080fd5b5056fea26469706673582212206f6b5f69ca6511abb106152fb18a644046dbb3e32b41c403cfac5fcc27f25a7d64736f6c63430007000033"

// DeployOVMCanonicalTransactionChain deploys a new Ethereum contract, binding an instance of OVMCanonicalTransactionChain to it.
func DeployOVMCanonicalTransactionChain(auth *bind.TransactOpts, backend bind.ContractBackend, _libAddressManager common.Address, _forceInclusionPeriodSeconds *big.Int) (common.Address, *types.Transaction, *OVMCanonicalTransactionChain, error) {
	parsed, err := abi.JSON(strings.NewReader(OVMCanonicalTransactionChainABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OVMCanonicalTransactionChainBin), backend, _libAddressManager, _forceInclusionPeriodSeconds)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OVMCanonicalTransactionChain{OVMCanonicalTransactionChainCaller: OVMCanonicalTransactionChainCaller{contract: contract}, OVMCanonicalTransactionChainTransactor: OVMCanonicalTransactionChainTransactor{contract: contract}, OVMCanonicalTransactionChainFilterer: OVMCanonicalTransactionChainFilterer{contract: contract}}, nil
}

// OVMCanonicalTransactionChain is an auto generated Go binding around an Ethereum contract.
type OVMCanonicalTransactionChain struct {
	OVMCanonicalTransactionChainCaller     // Read-only binding to the contract
	OVMCanonicalTransactionChainTransactor // Write-only binding to the contract
	OVMCanonicalTransactionChainFilterer   // Log filterer for contract events
}

// OVMCanonicalTransactionChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type OVMCanonicalTransactionChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMCanonicalTransactionChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OVMCanonicalTransactionChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMCanonicalTransactionChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OVMCanonicalTransactionChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OVMCanonicalTransactionChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OVMCanonicalTransactionChainSession struct {
	Contract     *OVMCanonicalTransactionChain // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                 // Call options to use throughout this session
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// OVMCanonicalTransactionChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OVMCanonicalTransactionChainCallerSession struct {
	Contract *OVMCanonicalTransactionChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                       // Call options to use throughout this session
}

// OVMCanonicalTransactionChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OVMCanonicalTransactionChainTransactorSession struct {
	Contract     *OVMCanonicalTransactionChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                       // Transaction auth options to use throughout this session
}

// OVMCanonicalTransactionChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type OVMCanonicalTransactionChainRaw struct {
	Contract *OVMCanonicalTransactionChain // Generic contract binding to access the raw methods on
}

// OVMCanonicalTransactionChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OVMCanonicalTransactionChainCallerRaw struct {
	Contract *OVMCanonicalTransactionChainCaller // Generic read-only contract binding to access the raw methods on
}

// OVMCanonicalTransactionChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OVMCanonicalTransactionChainTransactorRaw struct {
	Contract *OVMCanonicalTransactionChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOVMCanonicalTransactionChain creates a new instance of OVMCanonicalTransactionChain, bound to a specific deployed contract.
func NewOVMCanonicalTransactionChain(address common.Address, backend bind.ContractBackend) (*OVMCanonicalTransactionChain, error) {
	contract, err := bindOVMCanonicalTransactionChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChain{OVMCanonicalTransactionChainCaller: OVMCanonicalTransactionChainCaller{contract: contract}, OVMCanonicalTransactionChainTransactor: OVMCanonicalTransactionChainTransactor{contract: contract}, OVMCanonicalTransactionChainFilterer: OVMCanonicalTransactionChainFilterer{contract: contract}}, nil
}

// NewOVMCanonicalTransactionChainCaller creates a new read-only instance of OVMCanonicalTransactionChain, bound to a specific deployed contract.
func NewOVMCanonicalTransactionChainCaller(address common.Address, caller bind.ContractCaller) (*OVMCanonicalTransactionChainCaller, error) {
	contract, err := bindOVMCanonicalTransactionChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainCaller{contract: contract}, nil
}

// NewOVMCanonicalTransactionChainTransactor creates a new write-only instance of OVMCanonicalTransactionChain, bound to a specific deployed contract.
func NewOVMCanonicalTransactionChainTransactor(address common.Address, transactor bind.ContractTransactor) (*OVMCanonicalTransactionChainTransactor, error) {
	contract, err := bindOVMCanonicalTransactionChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainTransactor{contract: contract}, nil
}

// NewOVMCanonicalTransactionChainFilterer creates a new log filterer instance of OVMCanonicalTransactionChain, bound to a specific deployed contract.
func NewOVMCanonicalTransactionChainFilterer(address common.Address, filterer bind.ContractFilterer) (*OVMCanonicalTransactionChainFilterer, error) {
	contract, err := bindOVMCanonicalTransactionChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainFilterer{contract: contract}, nil
}

// bindOVMCanonicalTransactionChain binds a generic wrapper to an already deployed contract.
func bindOVMCanonicalTransactionChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OVMCanonicalTransactionChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OVMCanonicalTransactionChain.Contract.OVMCanonicalTransactionChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.OVMCanonicalTransactionChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.OVMCanonicalTransactionChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OVMCanonicalTransactionChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.contract.Transact(opts, method, params...)
}

// L2GASDISCOUNTDIVISOR is a free data retrieval call binding the contract method 0xc2cf696f.
//
// Solidity: function L2_GAS_DISCOUNT_DIVISOR() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) L2GASDISCOUNTDIVISOR(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "L2_GAS_DISCOUNT_DIVISOR")
	return *ret0, err
}

// L2GASDISCOUNTDIVISOR is a free data retrieval call binding the contract method 0xc2cf696f.
//
// Solidity: function L2_GAS_DISCOUNT_DIVISOR() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) L2GASDISCOUNTDIVISOR() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.L2GASDISCOUNTDIVISOR(&_OVMCanonicalTransactionChain.CallOpts)
}

// L2GASDISCOUNTDIVISOR is a free data retrieval call binding the contract method 0xc2cf696f.
//
// Solidity: function L2_GAS_DISCOUNT_DIVISOR() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) L2GASDISCOUNTDIVISOR() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.L2GASDISCOUNTDIVISOR(&_OVMCanonicalTransactionChain.CallOpts)
}

// MAXROLLUPTXSIZE is a free data retrieval call binding the contract method 0x876ed5cb.
//
// Solidity: function MAX_ROLLUP_TX_SIZE() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) MAXROLLUPTXSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "MAX_ROLLUP_TX_SIZE")
	return *ret0, err
}

// MAXROLLUPTXSIZE is a free data retrieval call binding the contract method 0x876ed5cb.
//
// Solidity: function MAX_ROLLUP_TX_SIZE() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) MAXROLLUPTXSIZE() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.MAXROLLUPTXSIZE(&_OVMCanonicalTransactionChain.CallOpts)
}

// MAXROLLUPTXSIZE is a free data retrieval call binding the contract method 0x876ed5cb.
//
// Solidity: function MAX_ROLLUP_TX_SIZE() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) MAXROLLUPTXSIZE() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.MAXROLLUPTXSIZE(&_OVMCanonicalTransactionChain.CallOpts)
}

// MINROLLUPTXGAS is a free data retrieval call binding the contract method 0x78f4b2f2.
//
// Solidity: function MIN_ROLLUP_TX_GAS() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) MINROLLUPTXGAS(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "MIN_ROLLUP_TX_GAS")
	return *ret0, err
}

// MINROLLUPTXGAS is a free data retrieval call binding the contract method 0x78f4b2f2.
//
// Solidity: function MIN_ROLLUP_TX_GAS() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) MINROLLUPTXGAS() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.MINROLLUPTXGAS(&_OVMCanonicalTransactionChain.CallOpts)
}

// MINROLLUPTXGAS is a free data retrieval call binding the contract method 0x78f4b2f2.
//
// Solidity: function MIN_ROLLUP_TX_GAS() constant returns(uint256)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) MINROLLUPTXGAS() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.MINROLLUPTXGAS(&_OVMCanonicalTransactionChain.CallOpts)
}

// GetQueueElement is a free data retrieval call binding the contract method 0x2a7f18be.
//
// Solidity: function getQueueElement(uint256 _index) constant returns(Lib_OVMCodecQueueElement _element)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) GetQueueElement(opts *bind.CallOpts, _index *big.Int) (Lib_OVMCodecQueueElement, error) {
	var (
		ret0 = new(Lib_OVMCodecQueueElement)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "getQueueElement", _index)
	return *ret0, err
}

// GetQueueElement is a free data retrieval call binding the contract method 0x2a7f18be.
//
// Solidity: function getQueueElement(uint256 _index) constant returns(Lib_OVMCodecQueueElement _element)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) GetQueueElement(_index *big.Int) (Lib_OVMCodecQueueElement, error) {
	return _OVMCanonicalTransactionChain.Contract.GetQueueElement(&_OVMCanonicalTransactionChain.CallOpts, _index)
}

// GetQueueElement is a free data retrieval call binding the contract method 0x2a7f18be.
//
// Solidity: function getQueueElement(uint256 _index) constant returns(Lib_OVMCodecQueueElement _element)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) GetQueueElement(_index *big.Int) (Lib_OVMCodecQueueElement, error) {
	return _OVMCanonicalTransactionChain.Contract.GetQueueElement(&_OVMCanonicalTransactionChain.CallOpts, _index)
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) GetTotalBatches(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "getTotalBatches")
	return *ret0, err
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) GetTotalBatches() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.GetTotalBatches(&_OVMCanonicalTransactionChain.CallOpts)
}

// GetTotalBatches is a free data retrieval call binding the contract method 0xe561dddc.
//
// Solidity: function getTotalBatches() constant returns(uint256 _totalBatches)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) GetTotalBatches() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.GetTotalBatches(&_OVMCanonicalTransactionChain.CallOpts)
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) GetTotalElements(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "getTotalElements")
	return *ret0, err
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) GetTotalElements() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.GetTotalElements(&_OVMCanonicalTransactionChain.CallOpts)
}

// GetTotalElements is a free data retrieval call binding the contract method 0x7aa63a86.
//
// Solidity: function getTotalElements() constant returns(uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) GetTotalElements() (*big.Int, error) {
	return _OVMCanonicalTransactionChain.Contract.GetTotalElements(&_OVMCanonicalTransactionChain.CallOpts)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) Resolve(opts *bind.CallOpts, _name string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "resolve", _name)
	return *ret0, err
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) Resolve(_name string) (common.Address, error) {
	return _OVMCanonicalTransactionChain.Contract.Resolve(&_OVMCanonicalTransactionChain.CallOpts, _name)
}

// Resolve is a free data retrieval call binding the contract method 0x461a4478.
//
// Solidity: function resolve(string _name) constant returns(address _contract)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) Resolve(_name string) (common.Address, error) {
	return _OVMCanonicalTransactionChain.Contract.Resolve(&_OVMCanonicalTransactionChain.CallOpts, _name)
}

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCaller) VerifyElement(opts *bind.CallOpts, _element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _OVMCanonicalTransactionChain.contract.Call(opts, out, "verifyElement", _element, _batchHeader, _proof)
	return *ret0, err
}

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) VerifyElement(_element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMCanonicalTransactionChain.Contract.VerifyElement(&_OVMCanonicalTransactionChain.CallOpts, _element, _batchHeader, _proof)
}

// VerifyElement is a free data retrieval call binding the contract method 0x98fe87c8.
//
// Solidity: function verifyElement(bytes _element, Lib_OVMCodecChainBatchHeader _batchHeader, Lib_OVMCodecChainInclusionProof _proof) constant returns(bool _verified)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainCallerSession) VerifyElement(_element []byte, _batchHeader Lib_OVMCodecChainBatchHeader, _proof Lib_OVMCodecChainInclusionProof) (bool, error) {
	return _OVMCanonicalTransactionChain.Contract.VerifyElement(&_OVMCanonicalTransactionChain.CallOpts, _element, _batchHeader, _proof)
}

// AppendQueueBatch is a paid mutator transaction binding the contract method 0xfacdc5da.
//
// Solidity: function appendQueueBatch(uint256 _numQueuedTransactions) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactor) AppendQueueBatch(opts *bind.TransactOpts, _numQueuedTransactions *big.Int) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.contract.Transact(opts, "appendQueueBatch", _numQueuedTransactions)
}

// AppendQueueBatch is a paid mutator transaction binding the contract method 0xfacdc5da.
//
// Solidity: function appendQueueBatch(uint256 _numQueuedTransactions) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) AppendQueueBatch(_numQueuedTransactions *big.Int) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.AppendQueueBatch(&_OVMCanonicalTransactionChain.TransactOpts, _numQueuedTransactions)
}

// AppendQueueBatch is a paid mutator transaction binding the contract method 0xfacdc5da.
//
// Solidity: function appendQueueBatch(uint256 _numQueuedTransactions) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactorSession) AppendQueueBatch(_numQueuedTransactions *big.Int) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.AppendQueueBatch(&_OVMCanonicalTransactionChain.TransactOpts, _numQueuedTransactions)
}

// AppendSequencerBatch is a paid mutator transaction binding the contract method 0xd0f89344.
//
// Solidity: function appendSequencerBatch() returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactor) AppendSequencerBatch(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.contract.Transact(opts, "appendSequencerBatch")
}

// AppendSequencerBatch is a paid mutator transaction binding the contract method 0xd0f89344.
//
// Solidity: function appendSequencerBatch() returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) AppendSequencerBatch() (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.AppendSequencerBatch(&_OVMCanonicalTransactionChain.TransactOpts)
}

// AppendSequencerBatch is a paid mutator transaction binding the contract method 0xd0f89344.
//
// Solidity: function appendSequencerBatch() returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactorSession) AppendSequencerBatch() (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.AppendSequencerBatch(&_OVMCanonicalTransactionChain.TransactOpts)
}

// Enqueue is a paid mutator transaction binding the contract method 0x6fee07e0.
//
// Solidity: function enqueue(address _target, uint256 _gasLimit, bytes _data) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactor) Enqueue(opts *bind.TransactOpts, _target common.Address, _gasLimit *big.Int, _data []byte) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.contract.Transact(opts, "enqueue", _target, _gasLimit, _data)
}

// Enqueue is a paid mutator transaction binding the contract method 0x6fee07e0.
//
// Solidity: function enqueue(address _target, uint256 _gasLimit, bytes _data) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainSession) Enqueue(_target common.Address, _gasLimit *big.Int, _data []byte) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.Enqueue(&_OVMCanonicalTransactionChain.TransactOpts, _target, _gasLimit, _data)
}

// Enqueue is a paid mutator transaction binding the contract method 0x6fee07e0.
//
// Solidity: function enqueue(address _target, uint256 _gasLimit, bytes _data) returns()
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainTransactorSession) Enqueue(_target common.Address, _gasLimit *big.Int, _data []byte) (*types.Transaction, error) {
	return _OVMCanonicalTransactionChain.Contract.Enqueue(&_OVMCanonicalTransactionChain.TransactOpts, _target, _gasLimit, _data)
}

// OVMCanonicalTransactionChainQueueBatchAppendedIterator is returned from FilterQueueBatchAppended and is used to iterate over the raw logs and unpacked data for QueueBatchAppended events raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainQueueBatchAppendedIterator struct {
	Event *OVMCanonicalTransactionChainQueueBatchAppended // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OVMCanonicalTransactionChainQueueBatchAppendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OVMCanonicalTransactionChainQueueBatchAppended)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OVMCanonicalTransactionChainQueueBatchAppended)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OVMCanonicalTransactionChainQueueBatchAppendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OVMCanonicalTransactionChainQueueBatchAppendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OVMCanonicalTransactionChainQueueBatchAppended represents a QueueBatchAppended event raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainQueueBatchAppended struct {
	StartingQueueIndex *big.Int
	NumQueueElements   *big.Int
	TotalElements      *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterQueueBatchAppended is a free log retrieval operation binding the contract event 0x64d7f508348c70dea42d5302a393987e4abc20e45954ab3f9d320207751956f0.
//
// Solidity: event QueueBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) FilterQueueBatchAppended(opts *bind.FilterOpts) (*OVMCanonicalTransactionChainQueueBatchAppendedIterator, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.FilterLogs(opts, "QueueBatchAppended")
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainQueueBatchAppendedIterator{contract: _OVMCanonicalTransactionChain.contract, event: "QueueBatchAppended", logs: logs, sub: sub}, nil
}

// WatchQueueBatchAppended is a free log subscription operation binding the contract event 0x64d7f508348c70dea42d5302a393987e4abc20e45954ab3f9d320207751956f0.
//
// Solidity: event QueueBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) WatchQueueBatchAppended(opts *bind.WatchOpts, sink chan<- *OVMCanonicalTransactionChainQueueBatchAppended) (event.Subscription, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.WatchLogs(opts, "QueueBatchAppended")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OVMCanonicalTransactionChainQueueBatchAppended)
				if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "QueueBatchAppended", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseQueueBatchAppended is a log parse operation binding the contract event 0x64d7f508348c70dea42d5302a393987e4abc20e45954ab3f9d320207751956f0.
//
// Solidity: event QueueBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) ParseQueueBatchAppended(log types.Log) (*OVMCanonicalTransactionChainQueueBatchAppended, error) {
	event := new(OVMCanonicalTransactionChainQueueBatchAppended)
	if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "QueueBatchAppended", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OVMCanonicalTransactionChainSequencerBatchAppendedIterator is returned from FilterSequencerBatchAppended and is used to iterate over the raw logs and unpacked data for SequencerBatchAppended events raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainSequencerBatchAppendedIterator struct {
	Event *OVMCanonicalTransactionChainSequencerBatchAppended // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OVMCanonicalTransactionChainSequencerBatchAppendedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OVMCanonicalTransactionChainSequencerBatchAppended)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OVMCanonicalTransactionChainSequencerBatchAppended)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OVMCanonicalTransactionChainSequencerBatchAppendedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OVMCanonicalTransactionChainSequencerBatchAppendedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OVMCanonicalTransactionChainSequencerBatchAppended represents a SequencerBatchAppended event raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainSequencerBatchAppended struct {
	StartingQueueIndex *big.Int
	NumQueueElements   *big.Int
	TotalElements      *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterSequencerBatchAppended is a free log retrieval operation binding the contract event 0x602f1aeac0ca2e7a13e281a9ef0ad7838542712ce16780fa2ecffd351f05f899.
//
// Solidity: event SequencerBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) FilterSequencerBatchAppended(opts *bind.FilterOpts) (*OVMCanonicalTransactionChainSequencerBatchAppendedIterator, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.FilterLogs(opts, "SequencerBatchAppended")
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainSequencerBatchAppendedIterator{contract: _OVMCanonicalTransactionChain.contract, event: "SequencerBatchAppended", logs: logs, sub: sub}, nil
}

// WatchSequencerBatchAppended is a free log subscription operation binding the contract event 0x602f1aeac0ca2e7a13e281a9ef0ad7838542712ce16780fa2ecffd351f05f899.
//
// Solidity: event SequencerBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) WatchSequencerBatchAppended(opts *bind.WatchOpts, sink chan<- *OVMCanonicalTransactionChainSequencerBatchAppended) (event.Subscription, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.WatchLogs(opts, "SequencerBatchAppended")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OVMCanonicalTransactionChainSequencerBatchAppended)
				if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "SequencerBatchAppended", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSequencerBatchAppended is a log parse operation binding the contract event 0x602f1aeac0ca2e7a13e281a9ef0ad7838542712ce16780fa2ecffd351f05f899.
//
// Solidity: event SequencerBatchAppended(uint256 _startingQueueIndex, uint256 _numQueueElements, uint256 _totalElements)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) ParseSequencerBatchAppended(log types.Log) (*OVMCanonicalTransactionChainSequencerBatchAppended, error) {
	event := new(OVMCanonicalTransactionChainSequencerBatchAppended)
	if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "SequencerBatchAppended", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OVMCanonicalTransactionChainTransactionEnqueuedIterator is returned from FilterTransactionEnqueued and is used to iterate over the raw logs and unpacked data for TransactionEnqueued events raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainTransactionEnqueuedIterator struct {
	Event *OVMCanonicalTransactionChainTransactionEnqueued // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OVMCanonicalTransactionChainTransactionEnqueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OVMCanonicalTransactionChainTransactionEnqueued)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OVMCanonicalTransactionChainTransactionEnqueued)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OVMCanonicalTransactionChainTransactionEnqueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OVMCanonicalTransactionChainTransactionEnqueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OVMCanonicalTransactionChainTransactionEnqueued represents a TransactionEnqueued event raised by the OVMCanonicalTransactionChain contract.
type OVMCanonicalTransactionChainTransactionEnqueued struct {
	L1TxOrigin common.Address
	Target     common.Address
	GasLimit   *big.Int
	Data       []byte
	QueueIndex *big.Int
	Timestamp  *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTransactionEnqueued is a free log retrieval operation binding the contract event 0x4b388aecf9fa6cc92253704e5975a6129a4f735bdbd99567df4ed0094ee4ceb5.
//
// Solidity: event TransactionEnqueued(address _l1TxOrigin, address _target, uint256 _gasLimit, bytes _data, uint256 _queueIndex, uint256 _timestamp)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) FilterTransactionEnqueued(opts *bind.FilterOpts) (*OVMCanonicalTransactionChainTransactionEnqueuedIterator, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.FilterLogs(opts, "TransactionEnqueued")
	if err != nil {
		return nil, err
	}
	return &OVMCanonicalTransactionChainTransactionEnqueuedIterator{contract: _OVMCanonicalTransactionChain.contract, event: "TransactionEnqueued", logs: logs, sub: sub}, nil
}

// WatchTransactionEnqueued is a free log subscription operation binding the contract event 0x4b388aecf9fa6cc92253704e5975a6129a4f735bdbd99567df4ed0094ee4ceb5.
//
// Solidity: event TransactionEnqueued(address _l1TxOrigin, address _target, uint256 _gasLimit, bytes _data, uint256 _queueIndex, uint256 _timestamp)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) WatchTransactionEnqueued(opts *bind.WatchOpts, sink chan<- *OVMCanonicalTransactionChainTransactionEnqueued) (event.Subscription, error) {

	logs, sub, err := _OVMCanonicalTransactionChain.contract.WatchLogs(opts, "TransactionEnqueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OVMCanonicalTransactionChainTransactionEnqueued)
				if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "TransactionEnqueued", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransactionEnqueued is a log parse operation binding the contract event 0x4b388aecf9fa6cc92253704e5975a6129a4f735bdbd99567df4ed0094ee4ceb5.
//
// Solidity: event TransactionEnqueued(address _l1TxOrigin, address _target, uint256 _gasLimit, bytes _data, uint256 _queueIndex, uint256 _timestamp)
func (_OVMCanonicalTransactionChain *OVMCanonicalTransactionChainFilterer) ParseTransactionEnqueued(log types.Log) (*OVMCanonicalTransactionChainTransactionEnqueued, error) {
	event := new(OVMCanonicalTransactionChainTransactionEnqueued)
	if err := _OVMCanonicalTransactionChain.contract.UnpackLog(event, "TransactionEnqueued", log); err != nil {
		return nil, err
	}
	return event, nil
}

