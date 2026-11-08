// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = errors.New
	_ = big.NewInt
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

// AttestedSettlementDeposit is an auto generated low-level Go binding around an user-defined struct.
type AttestedSettlementDeposit struct {
	Buyer             common.Address
	Seller            common.Address
	Amount            *big.Int
	Hashlock          [32]byte
	AssetIdHash       [32]byte
	RecipientIdHash   [32]byte
	Timelock          uint64
	MinFabricDeadline uint64
	State             uint8
}

// FabricLockAttestation is an auto generated low-level Go binding around an user-defined struct.
type FabricLockAttestation struct {
	SwapId          [32]byte
	Hashlock        [32]byte
	AssetIdHash     [32]byte
	RecipientIdHash [32]byte
	Deadline        uint64
	ChannelHash     [32]byte
	NamespaceHash   [32]byte
	FabricTxId      [32]byte
	BlockNumber     uint64
}

// FabricProofVerifierAttestorKey is an auto generated low-level Go binding around an user-defined struct.
type FabricProofVerifierAttestorKey struct {
	Qx       [32]byte
	Qy       [32]byte
	CertHash [32]byte
	Active   bool
}

// FabricProofVerifierOrg is an auto generated low-level Go binding around an user-defined struct.
type FabricProofVerifierOrg struct {
	Registered bool
	Active     bool
	KeyCount   uint32
	MspId      string
}

// HTLCPaymentLock is an auto generated low-level Go binding around an user-defined struct.
type HTLCPaymentLock struct {
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  uint64
	State     uint8
}

// OrgSignature is an auto generated low-level Go binding around an user-defined struct.
type OrgSignature struct {
	MspIdHash [32]byte
	KeyIndex  uint32
	R         [32]byte
	S         [32]byte
}

// AttestedSettlementMetaData contains all meta data concerning the AttestedSettlement contract.
var AttestedSettlementMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token_\",\"type\":\"address\"},{\"internalType\":\"contractFabricProofVerifier\",\"name\":\"verifier_\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelHash_\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"namespaceHash_\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"channelHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minFabricDeadline\",\"type\":\"uint64\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"getDeposit\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"minFabricDeadline\",\"type\":\"uint64\"},{\"internalType\":\"enumAttestedSettlement.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structAttestedSettlement.Deposit\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"}],\"name\":\"hashlockUsed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"namespaceHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"revealedPreimage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"preimage\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"channelHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"namespaceHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fabricTxId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"}],\"internalType\":\"structFabricLockAttestation\",\"name\":\"att\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structOrgSignature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalDeposited\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractFabricProofVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"minFabricDeadline\",\"type\":\"uint64\"}],\"name\":\"PaymentDeposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentRefunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"preimage\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"fabricTxId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentSettled\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"field\",\"type\":\"string\"}],\"name\":\"AttestationMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"Expired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"}],\"name\":\"HashlockReused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParams\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"enumAttestedSettlement.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"NotDeposited\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"NotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"SwapExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"TimelockInPast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"UnexpectedBalanceDelta\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"WrongPreimage\",\"type\":\"error\"}]",
	ID:  "7b5da917fa8fdf4bdfb95e0aef5cc02298",
	Bin: "0x610100604052348015610010575f5ffd5b50604051611b7f380380611b7f83398101604081905261002f916100c8565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556001600160a01b038416158061007057506001600160a01b038316155b1561008e57604051635435b28960e11b815260040160405180910390fd5b6001600160a01b039384166080529190921660a05260c09190915260e05261010d565b6001600160a01b03811681146100c5575f5ffd5b50565b5f5f5f5f608085870312156100db575f5ffd5b84516100e6816100b1565b60208601519094506100f7816100b1565b6040860151606090960151949790965092505050565b60805160a05160c05160e051611a086101775f395f81816101ba01526106fc01525f8181610193015261066f01525f818160ff01526107c401525f81816101f4015281816108b301528181610ada015281816110d90152818161117301526111e20152611a085ff3fe608060405234801561000f575f5ffd5b50600436106100c4575f3560e01c8063b4d46db61161007d578063fc0c546a11610058578063fc0c546a146101ef578063ff50abdc14610216578063ffdfa6de1461021f575f5ffd5b8063b4d46db61461018e578063bff579c4146101b5578063cdc04541146101dc575f5ffd5b8063579bfbcd116100ad578063579bfbcd146101465780637249fbb61461015b5780637a86983f1461016e575f5ffd5b80631335f9bb146100c85780632b7ac3f3146100fa575b5f5ffd5b6100e76100d63660046115a5565b60026020525f908152604090205481565b6040519081526020015b60405180910390f35b6101217f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f1565b6101596101543660046115bc565b610251565b005b6101596101693660046115a5565b610972565b61018161017c3660046115a5565b610b81565b6040516100f191906116df565b6100e77f000000000000000000000000000000000000000000000000000000000000000081565b6100e77f000000000000000000000000000000000000000000000000000000000000000081565b6101596101ea3660046117b8565b610cac565b6101217f000000000000000000000000000000000000000000000000000000000000000081565b6100e760035481565b61024161022d3660046115a5565b60016020525f908152604090205460ff1681565b60405190151581526020016100f1565b61025961133d565b5f85815260208190526040902060016006820154700100000000000000000000000000000000900460ff16600381111561029557610295611679565b146102ea57858160060160109054906101000a900460ff166040517fbdf562960000000000000000000000000000000000000000000000000000000081526004016102e1929190611841565b60405180910390fd5b600681015467ffffffffffffffff1642106103495760068101546040517f91cc37bc00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201524260248201526044016102e1565b806003015460028660405160200161036391815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261039b9161185c565b602060405180830381855afa1580156103b6573d5f5f3e3d5ffd5b5050506040513d601f19601f820116820180604052508101906103d99190611872565b14610413576040517f1fa2ea6e000000000000000000000000000000000000000000000000000000008152600481018790526024016102e1565b8335861461047d576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600660248201527f737761704964000000000000000000000000000000000000000000000000000060448201526064016102e1565b80600301548460200135146104ee576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600860248201527f686173686c6f636b00000000000000000000000000000000000000000000000060448201526064016102e1565b806004015484604001351461055f576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600760248201527f617373657449640000000000000000000000000000000000000000000000000060448201526064016102e1565b80600501548460600135146105d0576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f726563697069656e74000000000000000000000000000000000000000000000060448201526064016102e1565b600681015468010000000000000000900467ffffffffffffffff166105fb60a0860160808701611889565b67ffffffffffffffff16101561066d576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600860248201527f646561646c696e6500000000000000000000000000000000000000000000000060448201526064016102e1565b7f00000000000000000000000000000000000000000000000000000000000000008460a00135146106fa576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600760248201527f6368616e6e656c0000000000000000000000000000000000000000000000000060448201526064016102e1565b7f00000000000000000000000000000000000000000000000000000000000000008460c0013514610787576040517f8f0ad4b900000000000000000000000000000000000000000000000000000000815260206004820152600960248201527f6e616d657370616365000000000000000000000000000000000000000000000060448201526064016102e1565b6040517fca5a8b3400000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063ca5a8b34906107fd90879087908790600401611908565b5f6040518083038186803b158015610813575f5ffd5b505afa158015610825573d5f5f3e3d5ffd5b505050506006810180547fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff167002000000000000000000000000000000001790555f8681526002602081905260408220879055820154600380549192839261088e9084906119cf565b909155505060018201546108dc9073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691168361136b565b60018201546040805188815260e0880135602082015290810183905273ffffffffffffffffffffffffffffffffffffffff9091169088907ff7c445e59d8c2176674d5923afc02b47a9daa940d51b4c961f9fd171b428a20c9060600160405180910390a3505061096b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050505050565b61097a61133d565b5f81815260208190526040902060016006820154700100000000000000000000000000000000900460ff1660038111156109b6576109b6611679565b14610a0257818160060160109054906101000a900460ff166040517fbdf562960000000000000000000000000000000000000000000000000000000081526004016102e1929190611841565b600681015467ffffffffffffffff16421015610a625760068101546040517f8788802000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff90911660048201524260248201526044016102e1565b6006810180547fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff167003000000000000000000000000000000001790556002810154600380548291905f90610ab89084906119cf565b90915550508154610b039073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000811691168361136b565b815460405182815273ffffffffffffffffffffffffffffffffffffffff9091169084907f3b28817a65441bed1ac3a5496c6b71346aa302e80c7418c12e1c49eaa17211769060200160405180910390a35050610b7e60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b50565b610bcd60408051610120810182525f80825260208201819052918101829052606081018290526080810182905260a0810182905260c0810182905260e081018290529061010082015290565b5f8281526020818152604091829020825161012081018452815473ffffffffffffffffffffffffffffffffffffffff908116825260018301541692810192909252600281015492820192909252600380830154606083015260048301546080830152600583015460a0830152600683015467ffffffffffffffff80821660c08501526801000000000000000082041660e084015291929161010084019170010000000000000000000000000000000090910460ff1690811115610c9257610c92611679565b6003811115610ca357610ca3611679565b90525092915050565b610cb461133d565b871580610cbf575084155b80610cde575073ffffffffffffffffffffffffffffffffffffffff8716155b80610cfe575073ffffffffffffffffffffffffffffffffffffffff871633145b80610d07575085155b80610d10575083155b80610d19575082155b15610d50576040517fa86b651200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8167ffffffffffffffff168167ffffffffffffffff1611610d9d576040517fa86b651200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5f89815260208190526040902060060154700100000000000000000000000000000000900460ff166003811115610dd757610dd7611679565b14610e11576040517fafa955b9000000000000000000000000000000000000000000000000000000008152600481018990526024016102e1565b5f8581526001602052604090205460ff1615610e5c576040517ff41a15ea000000000000000000000000000000000000000000000000000000008152600481018690526024016102e1565b428267ffffffffffffffff1611610eb1576040517f874ac8d000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff831660048201524260248201526044016102e1565b604080516101208101825233815273ffffffffffffffffffffffffffffffffffffffff89166020820152908101879052606081018690526080810185905260a0810184905267ffffffffffffffff80841660c0830152821660e08201526101008101600190525f89815260208181526040918290208351815473ffffffffffffffffffffffffffffffffffffffff9182167fffffffffffffffffffffffff0000000000000000000000000000000000000000918216178355928501516001830180549190921693169290921790915590820151600282015560608201516003808301919091556080830151600483015560a0830151600583015560c083015160068301805460e086015167ffffffffffffffff90811668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000000000000000000090921693169290921791909117808255610100850151927fffffffffffffffffffffffffffffff00ffffffffffffffffffffffffffffffff9091169070010000000000000000000000000000000090849081111561105357611053611679565b021790555050505f858152600160208190526040822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00169091179055600380548892906110a49084906119e8565b90915550506040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa158015611133573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906111579190611872565b905061119b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633308a6113cb565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f90829073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015611227573d5f5f3e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061124b9190611872565b61125591906119cf565b905087811461129a576040517f55e5b61d00000000000000000000000000000000000000000000000000000000815260048101899052602481018290526044016102e1565b604080518981526020810189905267ffffffffffffffff8681168284015285166060820152905173ffffffffffffffffffffffffffffffffffffffff8b169133918d917f9b203b05f08c3d118093c139610248577823b7bec08bc0303c19f8efa3e0556b919081900360800190a4505061133360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050505050505050565b61134561142d565b60027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b611378838383600161148a565b6113c6576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016102e1565b505050565b6113d9848484846001611512565b611427576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024016102e1565b50505050565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0054600203611488576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6040517fa9059cbb000000000000000000000000000000000000000000000000000000005f81815273ffffffffffffffffffffffffffffffffffffffff8616600452602485905291602083604481808b5af1925060015f511483166115065783831516156114fa573d5f823e3d81fd5b5f873b113d1516831692505b60405250949350505050565b6040517f23b872dd000000000000000000000000000000000000000000000000000000005f81815273ffffffffffffffffffffffffffffffffffffffff8781166004528616602452604485905291602083606481808c5af1925060015f51148316611594578383151615611588573d5f823e3d81fd5b5f883b113d1516831692505b604052505f60605295945050505050565b5f602082840312156115b5575f5ffd5b5035919050565b5f5f5f5f5f8587036101808112156115d2575f5ffd5b86359550602087013594506101207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc08201121561160d575f5ffd5b5060408601925061016086013567ffffffffffffffff81111561162e575f5ffd5b8601601f8101881361163e575f5ffd5b803567ffffffffffffffff811115611654575f5ffd5b8860208260071b8401011115611668575f5ffd5b959894975092955050506020019190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b600481106116db577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b9052565b5f6101208201905073ffffffffffffffffffffffffffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff602084015116602083015260408301516040830152606083015160608301526080830151608083015260a083015160a083015260c083015161176460c084018267ffffffffffffffff169052565b5060e083015161178060e084018267ffffffffffffffff169052565b506101008301516117956101008401826116a6565b5092915050565b803567ffffffffffffffff811681146117b3575f5ffd5b919050565b5f5f5f5f5f5f5f5f610100898b0312156117d0575f5ffd5b88359750602089013573ffffffffffffffffffffffffffffffffffffffff811681146117fa575f5ffd5b965060408901359550606089013594506080890135935060a0890135925061182460c08a0161179c565b915061183260e08a0161179c565b90509295985092959890939650565b8281526040810161185560208301846116a6565b9392505050565b5f82518060208501845e5f920191825250919050565b5f60208284031215611882575f5ffd5b5051919050565b5f60208284031215611899575f5ffd5b6118558261179c565b8183526020830192505f815f5b848110156118fe5781358652602082013563ffffffff81168082146118d2575f5ffd5b6020880152506040828101359087015260608083013590870152608095860195909101906001016118af565b5093949350505050565b833581526020808501359082015260408085013590820152606080850135908201525f67ffffffffffffffff6119406080870161179c565b16608083015260a0858101359083015260c0808601359083015260e0808601359083015267ffffffffffffffff61197a610100870161179c565b16610100830152610140610120830152611999610140830184866118a2565b95945050505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b818103818111156119e2576119e26119a2565b92915050565b808201808211156119e2576119e26119a256fea164736f6c634300081e000a",
}

// AttestedSettlement is an auto generated Go binding around an Ethereum contract.
type AttestedSettlement struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *AttestedSettlement) GetABI() abi.ABI {
	return c.abi
}

// NewAttestedSettlement creates a new instance of AttestedSettlement.
func NewAttestedSettlement() *AttestedSettlement {
	parsed, err := AttestedSettlementMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &AttestedSettlement{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *AttestedSettlement) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address token_, address verifier_, bytes32 channelHash_, bytes32 namespaceHash_) returns()
func (attestedSettlement *AttestedSettlement) PackConstructor(token_ common.Address, verifier_ common.Address, channelHash_ [32]byte, namespaceHash_ [32]byte) []byte {
	enc, err := attestedSettlement.abi.Pack("", token_, verifier_, channelHash_, namespaceHash_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackChannelHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb4d46db6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function channelHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) PackChannelHash() []byte {
	enc, err := attestedSettlement.abi.Pack("channelHash")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackChannelHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xb4d46db6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function channelHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) TryPackChannelHash() ([]byte, error) {
	return attestedSettlement.abi.Pack("channelHash")
}

// UnpackChannelHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xb4d46db6.
//
// Solidity: function channelHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) UnpackChannelHash(data []byte) ([32]byte, error) {
	out, err := attestedSettlement.abi.Unpack("channelHash", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcdc04541.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function deposit(bytes32 swapId, address seller, uint256 amount, bytes32 hashlock, bytes32 assetIdHash, bytes32 recipientIdHash, uint64 timelock, uint64 minFabricDeadline) returns()
func (attestedSettlement *AttestedSettlement) PackDeposit(swapId [32]byte, seller common.Address, amount *big.Int, hashlock [32]byte, assetIdHash [32]byte, recipientIdHash [32]byte, timelock uint64, minFabricDeadline uint64) []byte {
	enc, err := attestedSettlement.abi.Pack("deposit", swapId, seller, amount, hashlock, assetIdHash, recipientIdHash, timelock, minFabricDeadline)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcdc04541.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function deposit(bytes32 swapId, address seller, uint256 amount, bytes32 hashlock, bytes32 assetIdHash, bytes32 recipientIdHash, uint64 timelock, uint64 minFabricDeadline) returns()
func (attestedSettlement *AttestedSettlement) TryPackDeposit(swapId [32]byte, seller common.Address, amount *big.Int, hashlock [32]byte, assetIdHash [32]byte, recipientIdHash [32]byte, timelock uint64, minFabricDeadline uint64) ([]byte, error) {
	return attestedSettlement.abi.Pack("deposit", swapId, seller, amount, hashlock, assetIdHash, recipientIdHash, timelock, minFabricDeadline)
}

// PackGetDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a86983f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getDeposit(bytes32 swapId) view returns((address,address,uint256,bytes32,bytes32,bytes32,uint64,uint64,uint8))
func (attestedSettlement *AttestedSettlement) PackGetDeposit(swapId [32]byte) []byte {
	enc, err := attestedSettlement.abi.Pack("getDeposit", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetDeposit is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7a86983f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getDeposit(bytes32 swapId) view returns((address,address,uint256,bytes32,bytes32,bytes32,uint64,uint64,uint8))
func (attestedSettlement *AttestedSettlement) TryPackGetDeposit(swapId [32]byte) ([]byte, error) {
	return attestedSettlement.abi.Pack("getDeposit", swapId)
}

// UnpackGetDeposit is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7a86983f.
//
// Solidity: function getDeposit(bytes32 swapId) view returns((address,address,uint256,bytes32,bytes32,bytes32,uint64,uint64,uint8))
func (attestedSettlement *AttestedSettlement) UnpackGetDeposit(data []byte) (AttestedSettlementDeposit, error) {
	out, err := attestedSettlement.abi.Unpack("getDeposit", data)
	if err != nil {
		return *new(AttestedSettlementDeposit), err
	}
	out0 := *abi.ConvertType(out[0], new(AttestedSettlementDeposit)).(*AttestedSettlementDeposit)
	return out0, nil
}

// PackHashlockUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffdfa6de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (attestedSettlement *AttestedSettlement) PackHashlockUsed(hashlock [32]byte) []byte {
	enc, err := attestedSettlement.abi.Pack("hashlockUsed", hashlock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashlockUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffdfa6de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (attestedSettlement *AttestedSettlement) TryPackHashlockUsed(hashlock [32]byte) ([]byte, error) {
	return attestedSettlement.abi.Pack("hashlockUsed", hashlock)
}

// UnpackHashlockUsed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffdfa6de.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (attestedSettlement *AttestedSettlement) UnpackHashlockUsed(data []byte) (bool, error) {
	out, err := attestedSettlement.abi.Unpack("hashlockUsed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackNamespaceHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbff579c4.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function namespaceHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) PackNamespaceHash() []byte {
	enc, err := attestedSettlement.abi.Pack("namespaceHash")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackNamespaceHash is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbff579c4.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function namespaceHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) TryPackNamespaceHash() ([]byte, error) {
	return attestedSettlement.abi.Pack("namespaceHash")
}

// UnpackNamespaceHash is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbff579c4.
//
// Solidity: function namespaceHash() view returns(bytes32)
func (attestedSettlement *AttestedSettlement) UnpackNamespaceHash(data []byte) ([32]byte, error) {
	out, err := attestedSettlement.abi.Unpack("namespaceHash", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7249fbb6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refund(bytes32 swapId) returns()
func (attestedSettlement *AttestedSettlement) PackRefund(swapId [32]byte) []byte {
	enc, err := attestedSettlement.abi.Pack("refund", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7249fbb6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refund(bytes32 swapId) returns()
func (attestedSettlement *AttestedSettlement) TryPackRefund(swapId [32]byte) ([]byte, error) {
	return attestedSettlement.abi.Pack("refund", swapId)
}

// PackRevealedPreimage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1335f9bb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (attestedSettlement *AttestedSettlement) PackRevealedPreimage(swapId [32]byte) []byte {
	enc, err := attestedSettlement.abi.Pack("revealedPreimage", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevealedPreimage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1335f9bb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (attestedSettlement *AttestedSettlement) TryPackRevealedPreimage(swapId [32]byte) ([]byte, error) {
	return attestedSettlement.abi.Pack("revealedPreimage", swapId)
}

// UnpackRevealedPreimage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1335f9bb.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (attestedSettlement *AttestedSettlement) UnpackRevealedPreimage(data []byte) ([32]byte, error) {
	out, err := attestedSettlement.abi.Unpack("revealedPreimage", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x579bfbcd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function settle(bytes32 swapId, bytes32 preimage, (bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) att, (bytes32,uint32,bytes32,bytes32)[] sigs) returns()
func (attestedSettlement *AttestedSettlement) PackSettle(swapId [32]byte, preimage [32]byte, att FabricLockAttestation, sigs []OrgSignature) []byte {
	enc, err := attestedSettlement.abi.Pack("settle", swapId, preimage, att, sigs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSettle is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x579bfbcd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function settle(bytes32 swapId, bytes32 preimage, (bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) att, (bytes32,uint32,bytes32,bytes32)[] sigs) returns()
func (attestedSettlement *AttestedSettlement) TryPackSettle(swapId [32]byte, preimage [32]byte, att FabricLockAttestation, sigs []OrgSignature) ([]byte, error) {
	return attestedSettlement.abi.Pack("settle", swapId, preimage, att, sigs)
}

// PackToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc0c546a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function token() view returns(address)
func (attestedSettlement *AttestedSettlement) PackToken() []byte {
	enc, err := attestedSettlement.abi.Pack("token")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc0c546a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function token() view returns(address)
func (attestedSettlement *AttestedSettlement) TryPackToken() ([]byte, error) {
	return attestedSettlement.abi.Pack("token")
}

// UnpackToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (attestedSettlement *AttestedSettlement) UnpackToken(data []byte) (common.Address, error) {
	out, err := attestedSettlement.abi.Unpack("token", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTotalDeposited is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xff50abdc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalDeposited() view returns(uint256)
func (attestedSettlement *AttestedSettlement) PackTotalDeposited() []byte {
	enc, err := attestedSettlement.abi.Pack("totalDeposited")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalDeposited is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xff50abdc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalDeposited() view returns(uint256)
func (attestedSettlement *AttestedSettlement) TryPackTotalDeposited() ([]byte, error) {
	return attestedSettlement.abi.Pack("totalDeposited")
}

// UnpackTotalDeposited is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xff50abdc.
//
// Solidity: function totalDeposited() view returns(uint256)
func (attestedSettlement *AttestedSettlement) UnpackTotalDeposited(data []byte) (*big.Int, error) {
	out, err := attestedSettlement.abi.Unpack("totalDeposited", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b7ac3f3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function verifier() view returns(address)
func (attestedSettlement *AttestedSettlement) PackVerifier() []byte {
	enc, err := attestedSettlement.abi.Pack("verifier")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackVerifier is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2b7ac3f3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function verifier() view returns(address)
func (attestedSettlement *AttestedSettlement) TryPackVerifier() ([]byte, error) {
	return attestedSettlement.abi.Pack("verifier")
}

// UnpackVerifier is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (attestedSettlement *AttestedSettlement) UnpackVerifier(data []byte) (common.Address, error) {
	out, err := attestedSettlement.abi.Unpack("verifier", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// AttestedSettlementPaymentDeposited represents a PaymentDeposited event raised by the AttestedSettlement contract.
type AttestedSettlementPaymentDeposited struct {
	SwapId            [32]byte
	Buyer             common.Address
	Seller            common.Address
	Amount            *big.Int
	Hashlock          [32]byte
	Timelock          uint64
	MinFabricDeadline uint64
	Raw               *types.Log // Blockchain specific contextual infos
}

const AttestedSettlementPaymentDepositedEventName = "PaymentDeposited"

// ContractEventName returns the user-defined event name.
func (AttestedSettlementPaymentDeposited) ContractEventName() string {
	return AttestedSettlementPaymentDepositedEventName
}

// UnpackPaymentDepositedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentDeposited(bytes32 indexed swapId, address indexed buyer, address indexed seller, uint256 amount, bytes32 hashlock, uint64 timelock, uint64 minFabricDeadline)
func (attestedSettlement *AttestedSettlement) UnpackPaymentDepositedEvent(log *types.Log) (*AttestedSettlementPaymentDeposited, error) {
	event := "PaymentDeposited"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != attestedSettlement.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(AttestedSettlementPaymentDeposited)
	if len(log.Data) > 0 {
		if err := attestedSettlement.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range attestedSettlement.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AttestedSettlementPaymentRefunded represents a PaymentRefunded event raised by the AttestedSettlement contract.
type AttestedSettlementPaymentRefunded struct {
	SwapId [32]byte
	Buyer  common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const AttestedSettlementPaymentRefundedEventName = "PaymentRefunded"

// ContractEventName returns the user-defined event name.
func (AttestedSettlementPaymentRefunded) ContractEventName() string {
	return AttestedSettlementPaymentRefundedEventName
}

// UnpackPaymentRefundedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentRefunded(bytes32 indexed swapId, address indexed buyer, uint256 amount)
func (attestedSettlement *AttestedSettlement) UnpackPaymentRefundedEvent(log *types.Log) (*AttestedSettlementPaymentRefunded, error) {
	event := "PaymentRefunded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != attestedSettlement.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(AttestedSettlementPaymentRefunded)
	if len(log.Data) > 0 {
		if err := attestedSettlement.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range attestedSettlement.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// AttestedSettlementPaymentSettled represents a PaymentSettled event raised by the AttestedSettlement contract.
type AttestedSettlementPaymentSettled struct {
	SwapId     [32]byte
	Preimage   [32]byte
	FabricTxId [32]byte
	Seller     common.Address
	Amount     *big.Int
	Raw        *types.Log // Blockchain specific contextual infos
}

const AttestedSettlementPaymentSettledEventName = "PaymentSettled"

// ContractEventName returns the user-defined event name.
func (AttestedSettlementPaymentSettled) ContractEventName() string {
	return AttestedSettlementPaymentSettledEventName
}

// UnpackPaymentSettledEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentSettled(bytes32 indexed swapId, bytes32 preimage, bytes32 fabricTxId, address indexed seller, uint256 amount)
func (attestedSettlement *AttestedSettlement) UnpackPaymentSettledEvent(log *types.Log) (*AttestedSettlementPaymentSettled, error) {
	event := "PaymentSettled"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != attestedSettlement.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(AttestedSettlementPaymentSettled)
	if len(log.Data) > 0 {
		if err := attestedSettlement.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range attestedSettlement.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (attestedSettlement *AttestedSettlement) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["AttestationMismatch"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackAttestationMismatchError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["Expired"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["HashlockReused"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackHashlockReusedError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["InvalidParams"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackInvalidParamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["NotDeposited"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackNotDepositedError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["NotExpired"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackNotExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["SwapExists"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackSwapExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["TimelockInPast"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackTimelockInPastError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["UnexpectedBalanceDelta"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackUnexpectedBalanceDeltaError(raw[4:])
	}
	if bytes.Equal(raw[:4], attestedSettlement.abi.Errors["WrongPreimage"].ID.Bytes()[:4]) {
		return attestedSettlement.UnpackWrongPreimageError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// AttestedSettlementAttestationMismatch represents a AttestationMismatch error raised by the AttestedSettlement contract.
type AttestedSettlementAttestationMismatch struct {
	Field string
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AttestationMismatch(string field)
func AttestedSettlementAttestationMismatchErrorID() common.Hash {
	return common.HexToHash("0x8f0ad4b9814f94e848934cab6a6810929743f7830b355a4759d71918fe8ab7a6")
}

// UnpackAttestationMismatchError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AttestationMismatch(string field)
func (attestedSettlement *AttestedSettlement) UnpackAttestationMismatchError(raw []byte) (*AttestedSettlementAttestationMismatch, error) {
	out := new(AttestedSettlementAttestationMismatch)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "AttestationMismatch", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementExpired represents a Expired error raised by the AttestedSettlement contract.
type AttestedSettlementExpired struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Expired(uint64 timelock, uint256 nowTs)
func AttestedSettlementExpiredErrorID() common.Hash {
	return common.HexToHash("0x91cc37bcaeec164e5bc377adb7e09e2a00c91458581814d0294b11a1ab47958a")
}

// UnpackExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Expired(uint64 timelock, uint256 nowTs)
func (attestedSettlement *AttestedSettlement) UnpackExpiredError(raw []byte) (*AttestedSettlementExpired, error) {
	out := new(AttestedSettlementExpired)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "Expired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementHashlockReused represents a HashlockReused error raised by the AttestedSettlement contract.
type AttestedSettlementHashlockReused struct {
	Hashlock [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HashlockReused(bytes32 hashlock)
func AttestedSettlementHashlockReusedErrorID() common.Hash {
	return common.HexToHash("0xf41a15eade40d555d3009b515211d8eac2ddf768b0eaad4389bb9a2534e66a3c")
}

// UnpackHashlockReusedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HashlockReused(bytes32 hashlock)
func (attestedSettlement *AttestedSettlement) UnpackHashlockReusedError(raw []byte) (*AttestedSettlementHashlockReused, error) {
	out := new(AttestedSettlementHashlockReused)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "HashlockReused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementInvalidParams represents a InvalidParams error raised by the AttestedSettlement contract.
type AttestedSettlementInvalidParams struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidParams()
func AttestedSettlementInvalidParamsErrorID() common.Hash {
	return common.HexToHash("0xa86b65120b98ec2b50fcad62782dbb454a5c0604e0cf2db01fd94e5048d895b9")
}

// UnpackInvalidParamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidParams()
func (attestedSettlement *AttestedSettlement) UnpackInvalidParamsError(raw []byte) (*AttestedSettlementInvalidParams, error) {
	out := new(AttestedSettlementInvalidParams)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "InvalidParams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementNotDeposited represents a NotDeposited error raised by the AttestedSettlement contract.
type AttestedSettlementNotDeposited struct {
	SwapId [32]byte
	State  uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotDeposited(bytes32 swapId, uint8 state)
func AttestedSettlementNotDepositedErrorID() common.Hash {
	return common.HexToHash("0xbdf562960a887b648a830143996a61884b77fc4cc481c20880d38e9582380e5b")
}

// UnpackNotDepositedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotDeposited(bytes32 swapId, uint8 state)
func (attestedSettlement *AttestedSettlement) UnpackNotDepositedError(raw []byte) (*AttestedSettlementNotDeposited, error) {
	out := new(AttestedSettlementNotDeposited)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "NotDeposited", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementNotExpired represents a NotExpired error raised by the AttestedSettlement contract.
type AttestedSettlementNotExpired struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotExpired(uint64 timelock, uint256 nowTs)
func AttestedSettlementNotExpiredErrorID() common.Hash {
	return common.HexToHash("0x878880200d6e2e3220bbbcf837eb40e26654b6a531b808079cf8848829ff512f")
}

// UnpackNotExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotExpired(uint64 timelock, uint256 nowTs)
func (attestedSettlement *AttestedSettlement) UnpackNotExpiredError(raw []byte) (*AttestedSettlementNotExpired, error) {
	out := new(AttestedSettlementNotExpired)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "NotExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the AttestedSettlement contract.
type AttestedSettlementReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func AttestedSettlementReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (attestedSettlement *AttestedSettlement) UnpackReentrancyGuardReentrantCallError(raw []byte) (*AttestedSettlementReentrancyGuardReentrantCall, error) {
	out := new(AttestedSettlementReentrancyGuardReentrantCall)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the AttestedSettlement contract.
type AttestedSettlementSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func AttestedSettlementSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (attestedSettlement *AttestedSettlement) UnpackSafeERC20FailedOperationError(raw []byte) (*AttestedSettlementSafeERC20FailedOperation, error) {
	out := new(AttestedSettlementSafeERC20FailedOperation)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementSwapExists represents a SwapExists error raised by the AttestedSettlement contract.
type AttestedSettlementSwapExists struct {
	SwapId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapExists(bytes32 swapId)
func AttestedSettlementSwapExistsErrorID() common.Hash {
	return common.HexToHash("0xafa955b93ee380941bdede9cdd884e4748f2da7f8e02336dda9c3bc5efeb9523")
}

// UnpackSwapExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapExists(bytes32 swapId)
func (attestedSettlement *AttestedSettlement) UnpackSwapExistsError(raw []byte) (*AttestedSettlementSwapExists, error) {
	out := new(AttestedSettlementSwapExists)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "SwapExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementTimelockInPast represents a TimelockInPast error raised by the AttestedSettlement contract.
type AttestedSettlementTimelockInPast struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimelockInPast(uint64 timelock, uint256 nowTs)
func AttestedSettlementTimelockInPastErrorID() common.Hash {
	return common.HexToHash("0x874ac8d085b75da5651ef95027f3391b3208fd3cc8a677f09d139a5760db7bc6")
}

// UnpackTimelockInPastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimelockInPast(uint64 timelock, uint256 nowTs)
func (attestedSettlement *AttestedSettlement) UnpackTimelockInPastError(raw []byte) (*AttestedSettlementTimelockInPast, error) {
	out := new(AttestedSettlementTimelockInPast)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "TimelockInPast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementUnexpectedBalanceDelta represents a UnexpectedBalanceDelta error raised by the AttestedSettlement contract.
type AttestedSettlementUnexpectedBalanceDelta struct {
	Expected *big.Int
	Received *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnexpectedBalanceDelta(uint256 expected, uint256 received)
func AttestedSettlementUnexpectedBalanceDeltaErrorID() common.Hash {
	return common.HexToHash("0x55e5b61d2c182dd91a8abf6fae0f353c2e31558869e50172aae83e146d5db184")
}

// UnpackUnexpectedBalanceDeltaError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnexpectedBalanceDelta(uint256 expected, uint256 received)
func (attestedSettlement *AttestedSettlement) UnpackUnexpectedBalanceDeltaError(raw []byte) (*AttestedSettlementUnexpectedBalanceDelta, error) {
	out := new(AttestedSettlementUnexpectedBalanceDelta)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "UnexpectedBalanceDelta", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// AttestedSettlementWrongPreimage represents a WrongPreimage error raised by the AttestedSettlement contract.
type AttestedSettlementWrongPreimage struct {
	SwapId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongPreimage(bytes32 swapId)
func AttestedSettlementWrongPreimageErrorID() common.Hash {
	return common.HexToHash("0x1fa2ea6e5b8169ee4d614d5d9fbb8e5441ae35f50318e1f1a35549322cb46d98")
}

// UnpackWrongPreimageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongPreimage(bytes32 swapId)
func (attestedSettlement *AttestedSettlement) UnpackWrongPreimageError(raw []byte) (*AttestedSettlementWrongPreimage, error) {
	out := new(AttestedSettlementWrongPreimage)
	if err := attestedSettlement.abi.UnpackIntoInterface(out, "WrongPreimage", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenMetaData contains all meta data concerning the DrexToken contract.
var DrexTokenMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"COMPLIANCE_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ISSUER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"AllowlistUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"allowance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientAllowance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"needed\",\"type\":\"uint256\"}],\"name\":\"ERC20InsufficientBalance\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC20InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"ERC20InvalidSpender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"NotAllowlisted\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ZeroAddress\",\"type\":\"error\"}]",
	ID:  "cfdb2bf1d92b35672e8eaece47a0bf6299",
	Bin: "0x608060405234801561000f575f5ffd5b5060405161160138038061160183398101604081905261002e916101d0565b604051806040016040528060128152602001714472657820284a616e7573206c6f63616c2960701b81525060405180604001604052806004815260200163088a48ab60e31b81525081600390816100859190610295565b5060046100928282610295565b5050506001600160a01b0381166100bc5760405163d92e233d60e01b815260040160405180910390fd5b6100c65f82610123565b506100f17f114e74f6ea3bd819998f78687bfcb11b140da08e9b7d222fa9c1f1ba1f2aa12282610123565b5061011c7f442a94f1a1fac79af32856af2a64f63648cfa2ef3b98610a5bb7cbec4cee698582610123565b505061034f565b5f8281526005602090815260408083206001600160a01b038516845290915281205460ff166101c7575f8381526005602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561017f3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016101ca565b505f5b92915050565b5f602082840312156101e0575f5ffd5b81516001600160a01b03811681146101f6575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061022557607f821691505b60208210810361024357634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561029057805f5260205f20601f840160051c8101602085101561026e5750805b601f840160051c820191505b8181101561028d575f815560010161027a565b50505b505050565b81516001600160401b038111156102ae576102ae6101fd565b6102c2816102bc8454610211565b84610249565b6020601f8211600181146102f4575f83156102dd5750848201515b5f19600385901b1c1916600184901b17845561028d565b5f84815260208120601f198516915b828110156103235787850151825560209485019460019092019101610303565b508482101561034057868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b6112a58061035c5f395ff3fe608060405234801561000f575f5ffd5b5060043610610179575f3560e01c80634697f05d116100d25780639dc29fac11610088578063babcc53911610063578063babcc53914610372578063d547741f14610394578063dd62ed3e146103a7575f5ffd5b80639dc29fac14610345578063a217fddf14610358578063a9059cbb1461035f575f5ffd5b806382aefa24116100b857806382aefa24146102d157806391d14854146102f857806395d89b411461033d575f5ffd5b80634697f05d1461028957806370a082311461029c575f5ffd5b806323b872dd11610132578063313ce5671161010d578063313ce5671461025457806336568abe1461026357806340c10f1914610276575f5ffd5b806323b872dd1461020a578063248a9ca31461021d5780632f2ff15d1461023f575f5ffd5b806306fdde031161016257806306fdde03146101da578063095ea7b3146101ef57806318160ddd14610202575f5ffd5b806301ffc9a71461017d578063062d3bd7146101a5575b5f5ffd5b61019061018b366004611031565b6103ec565b60405190151581526020015b60405180910390f35b6101cc7f442a94f1a1fac79af32856af2a64f63648cfa2ef3b98610a5bb7cbec4cee698581565b60405190815260200161019c565b6101e2610484565b60405161019c9190611077565b6101906101fd3660046110f2565b610514565b6002546101cc565b61019061021836600461111a565b61052b565b6101cc61022b366004611154565b5f9081526005602052604090206001015490565b61025261024d36600461116b565b61054e565b005b6040516002815260200161019c565b61025261027136600461116b565b610578565b6102526102843660046110f2565b6105d6565b610252610297366004611195565b61060a565b6101cc6102aa3660046111ce565b73ffffffffffffffffffffffffffffffffffffffff165f9081526020819052604090205490565b6101cc7f114e74f6ea3bd819998f78687bfcb11b140da08e9b7d222fa9c1f1ba1f2aa12281565b61019061030636600461116b565b5f91825260056020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b6101e261070b565b6102526103533660046110f2565b61071a565b6101cc5f81565b61019061036d3660046110f2565b61074e565b6101906103803660046111ce565b60066020525f908152604090205460ff1681565b6102526103a236600461116b565b61075b565b6101cc6103b53660046111e7565b73ffffffffffffffffffffffffffffffffffffffff9182165f90815260016020908152604080832093909416825291909152205490565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061047e57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b6060600380546104939061120f565b80601f01602080910402602001604051908101604052809291908181526020018280546104bf9061120f565b801561050a5780601f106104e15761010080835404028352916020019161050a565b820191905f5260205f20905b8154815290600101906020018083116104ed57829003601f168201915b5050505050905090565b5f3361052181858561077f565b5060019392505050565b5f3361053885828561078c565b610543858585610859565b506001949350505050565b5f8281526005602052604090206001015461056881610902565b610572838361090f565b50505050565b73ffffffffffffffffffffffffffffffffffffffff811633146105c7576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6105d18282610a0c565b505050565b7f114e74f6ea3bd819998f78687bfcb11b140da08e9b7d222fa9c1f1ba1f2aa12261060081610902565b6105d18383610ac9565b7f442a94f1a1fac79af32856af2a64f63648cfa2ef3b98610a5bb7cbec4cee698561063481610902565b73ffffffffffffffffffffffffffffffffffffffff8316610681576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff83165f8181526006602090815260409182902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001686151590811790915591519182527f13518841ff4d3053cb7703afaa39b145c6331829b982d42f4d4fd7568b2e8e24910160405180910390a2505050565b6060600480546104939061120f565b7f114e74f6ea3bd819998f78687bfcb11b140da08e9b7d222fa9c1f1ba1f2aa12261074481610902565b6105d18383610b27565b5f33610521818585610859565b5f8281526005602052604090206001015461077581610902565b6105728383610a0c565b6105d18383836001610b81565b73ffffffffffffffffffffffffffffffffffffffff8381165f908152600160209081526040808320938616835292905220547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff811015610572578181101561084b576040517ffb8f41b200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8416600482015260248101829052604481018390526064015b60405180910390fd5b61057284848484035f610b81565b73ffffffffffffffffffffffffffffffffffffffff83166108a8576040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b73ffffffffffffffffffffffffffffffffffffffff82166108f7576040517fec442f050000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b6105d1838383610cc6565b61090c8133610e03565b50565b5f82815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff16610a05575f83815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556109a33390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161047e565b505f61047e565b5f82815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff1615610a05575f83815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161047e565b73ffffffffffffffffffffffffffffffffffffffff8216610b18576040517fec442f050000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b610b235f8383610cc6565b5050565b73ffffffffffffffffffffffffffffffffffffffff8216610b76576040517f96c6fd1e0000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b610b23825f83610cc6565b73ffffffffffffffffffffffffffffffffffffffff8416610bd0576040517fe602df050000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b73ffffffffffffffffffffffffffffffffffffffff8316610c1f576040517f94280d620000000000000000000000000000000000000000000000000000000081525f6004820152602401610842565b73ffffffffffffffffffffffffffffffffffffffff8085165f9081526001602090815260408083209387168352929052208290558015610572578273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92584604051610cb891815260200190565b60405180910390a350505050565b73ffffffffffffffffffffffffffffffffffffffff831615801590610d10575073ffffffffffffffffffffffffffffffffffffffff83165f9081526006602052604090205460ff16155b15610d5f576040517f6f7a527500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff84166004820152602401610842565b73ffffffffffffffffffffffffffffffffffffffff821615801590610da9575073ffffffffffffffffffffffffffffffffffffffff82165f9081526006602052604090205460ff16155b15610df8576040517f6f7a527500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff83166004820152602401610842565b6105d1838383610e8a565b5f82815260056020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16610b23576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff8216600482015260248101839052604401610842565b73ffffffffffffffffffffffffffffffffffffffff8316610ec1578060025f828254610eb69190611260565b90915550610f719050565b73ffffffffffffffffffffffffffffffffffffffff83165f9081526020819052604090205481811015610f46576040517fe450d38c00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024810182905260448101839052606401610842565b73ffffffffffffffffffffffffffffffffffffffff84165f9081526020819052604090209082900390555b73ffffffffffffffffffffffffffffffffffffffff8216610f9a57600280548290039055610fc5565b73ffffffffffffffffffffffffffffffffffffffff82165f9081526020819052604090208054820190555b8173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8360405161102491815260200190565b60405180910390a3505050565b5f60208284031215611041575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114611070575f5ffd5b9392505050565b602081525f82518060208401528060208501604085015e5f6040828501015260407fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505092915050565b803573ffffffffffffffffffffffffffffffffffffffff811681146110ed575f5ffd5b919050565b5f5f60408385031215611103575f5ffd5b61110c836110ca565b946020939093013593505050565b5f5f5f6060848603121561112c575f5ffd5b611135846110ca565b9250611143602085016110ca565b929592945050506040919091013590565b5f60208284031215611164575f5ffd5b5035919050565b5f5f6040838503121561117c575f5ffd5b8235915061118c602084016110ca565b90509250929050565b5f5f604083850312156111a6575f5ffd5b6111af836110ca565b9150602083013580151581146111c3575f5ffd5b809150509250929050565b5f602082840312156111de575f5ffd5b611070826110ca565b5f5f604083850312156111f8575f5ffd5b611201836110ca565b915061118c602084016110ca565b600181811c9082168061122357607f821691505b60208210810361125a577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b8082018082111561047e577f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffdfea164736f6c634300081e000a",
}

// DrexToken is an auto generated Go binding around an Ethereum contract.
type DrexToken struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *DrexToken) GetABI() abi.ABI {
	return c.abi
}

// NewDrexToken creates a new instance of DrexToken.
func NewDrexToken() *DrexToken {
	parsed, err := DrexTokenMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &DrexToken{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *DrexToken) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin) returns()
func (drexToken *DrexToken) PackConstructor(admin common.Address) []byte {
	enc, err := drexToken.abi.Pack("", admin)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackCOMPLIANCEROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x062d3bd7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (drexToken *DrexToken) PackCOMPLIANCEROLE() []byte {
	enc, err := drexToken.abi.Pack("COMPLIANCE_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackCOMPLIANCEROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x062d3bd7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (drexToken *DrexToken) TryPackCOMPLIANCEROLE() ([]byte, error) {
	return drexToken.abi.Pack("COMPLIANCE_ROLE")
}

// UnpackCOMPLIANCEROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x062d3bd7.
//
// Solidity: function COMPLIANCE_ROLE() view returns(bytes32)
func (drexToken *DrexToken) UnpackCOMPLIANCEROLE(data []byte) ([32]byte, error) {
	out, err := drexToken.abi.Unpack("COMPLIANCE_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (drexToken *DrexToken) PackDEFAULTADMINROLE() []byte {
	enc, err := drexToken.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (drexToken *DrexToken) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return drexToken.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (drexToken *DrexToken) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := drexToken.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackISSUERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82aefa24.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function ISSUER_ROLE() view returns(bytes32)
func (drexToken *DrexToken) PackISSUERROLE() []byte {
	enc, err := drexToken.abi.Pack("ISSUER_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackISSUERROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x82aefa24.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function ISSUER_ROLE() view returns(bytes32)
func (drexToken *DrexToken) TryPackISSUERROLE() ([]byte, error) {
	return drexToken.abi.Pack("ISSUER_ROLE")
}

// UnpackISSUERROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x82aefa24.
//
// Solidity: function ISSUER_ROLE() view returns(bytes32)
func (drexToken *DrexToken) UnpackISSUERROLE(data []byte) ([32]byte, error) {
	out, err := drexToken.abi.Unpack("ISSUER_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (drexToken *DrexToken) PackAllowance(owner common.Address, spender common.Address) []byte {
	enc, err := drexToken.abi.Pack("allowance", owner, spender)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAllowance is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xdd62ed3e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (drexToken *DrexToken) TryPackAllowance(owner common.Address, spender common.Address) ([]byte, error) {
	return drexToken.abi.Pack("allowance", owner, spender)
}

// UnpackAllowance is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (drexToken *DrexToken) UnpackAllowance(data []byte) (*big.Int, error) {
	out, err := drexToken.abi.Unpack("allowance", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (drexToken *DrexToken) PackApprove(spender common.Address, value *big.Int) []byte {
	enc, err := drexToken.abi.Pack("approve", spender, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackApprove is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x095ea7b3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (drexToken *DrexToken) TryPackApprove(spender common.Address, value *big.Int) ([]byte, error) {
	return drexToken.abi.Pack("approve", spender, value)
}

// UnpackApprove is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (drexToken *DrexToken) UnpackApprove(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("approve", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (drexToken *DrexToken) PackBalanceOf(account common.Address) []byte {
	enc, err := drexToken.abi.Pack("balanceOf", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBalanceOf is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x70a08231.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (drexToken *DrexToken) TryPackBalanceOf(account common.Address) ([]byte, error) {
	return drexToken.abi.Pack("balanceOf", account)
}

// UnpackBalanceOf is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (drexToken *DrexToken) UnpackBalanceOf(data []byte) (*big.Int, error) {
	out, err := drexToken.abi.Unpack("balanceOf", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9dc29fac.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (drexToken *DrexToken) PackBurn(from common.Address, amount *big.Int) []byte {
	enc, err := drexToken.abi.Pack("burn", from, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackBurn is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9dc29fac.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function burn(address from, uint256 amount) returns()
func (drexToken *DrexToken) TryPackBurn(from common.Address, amount *big.Int) ([]byte, error) {
	return drexToken.abi.Pack("burn", from, amount)
}

// PackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function decimals() pure returns(uint8)
func (drexToken *DrexToken) PackDecimals() []byte {
	enc, err := drexToken.abi.Pack("decimals")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDecimals is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x313ce567.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function decimals() pure returns(uint8)
func (drexToken *DrexToken) TryPackDecimals() ([]byte, error) {
	return drexToken.abi.Pack("decimals")
}

// UnpackDecimals is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x313ce567.
//
// Solidity: function decimals() pure returns(uint8)
func (drexToken *DrexToken) UnpackDecimals(data []byte) (uint8, error) {
	out, err := drexToken.abi.Unpack("decimals", data)
	if err != nil {
		return *new(uint8), err
	}
	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (drexToken *DrexToken) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := drexToken.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (drexToken *DrexToken) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return drexToken.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (drexToken *DrexToken) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := drexToken.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (drexToken *DrexToken) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := drexToken.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (drexToken *DrexToken) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return drexToken.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (drexToken *DrexToken) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := drexToken.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (drexToken *DrexToken) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return drexToken.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (drexToken *DrexToken) UnpackHasRole(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbabcc539.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isAllowed(address account) view returns(bool)
func (drexToken *DrexToken) PackIsAllowed(account common.Address) []byte {
	enc, err := drexToken.abi.Pack("isAllowed", account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xbabcc539.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isAllowed(address account) view returns(bool)
func (drexToken *DrexToken) TryPackIsAllowed(account common.Address) ([]byte, error) {
	return drexToken.abi.Pack("isAllowed", account)
}

// UnpackIsAllowed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xbabcc539.
//
// Solidity: function isAllowed(address account) view returns(bool)
func (drexToken *DrexToken) UnpackIsAllowed(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("isAllowed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (drexToken *DrexToken) PackMint(to common.Address, amount *big.Int) []byte {
	enc, err := drexToken.abi.Pack("mint", to, amount)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackMint is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x40c10f19.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function mint(address to, uint256 amount) returns()
func (drexToken *DrexToken) TryPackMint(to common.Address, amount *big.Int) ([]byte, error) {
	return drexToken.abi.Pack("mint", to, amount)
}

// PackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function name() view returns(string)
func (drexToken *DrexToken) PackName() []byte {
	enc, err := drexToken.abi.Pack("name")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackName is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x06fdde03.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function name() view returns(string)
func (drexToken *DrexToken) TryPackName() ([]byte, error) {
	return drexToken.abi.Pack("name")
}

// UnpackName is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (drexToken *DrexToken) UnpackName(data []byte) (string, error) {
	out, err := drexToken.abi.Unpack("name", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (drexToken *DrexToken) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := drexToken.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (drexToken *DrexToken) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return drexToken.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (drexToken *DrexToken) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := drexToken.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (drexToken *DrexToken) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return drexToken.abi.Pack("revokeRole", role, account)
}

// PackSetAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4697f05d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setAllowed(address account, bool allowed) returns()
func (drexToken *DrexToken) PackSetAllowed(account common.Address, allowed bool) []byte {
	enc, err := drexToken.abi.Pack("setAllowed", account, allowed)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetAllowed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4697f05d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setAllowed(address account, bool allowed) returns()
func (drexToken *DrexToken) TryPackSetAllowed(account common.Address, allowed bool) ([]byte, error) {
	return drexToken.abi.Pack("setAllowed", account, allowed)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (drexToken *DrexToken) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := drexToken.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (drexToken *DrexToken) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return drexToken.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (drexToken *DrexToken) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function symbol() view returns(string)
func (drexToken *DrexToken) PackSymbol() []byte {
	enc, err := drexToken.abi.Pack("symbol")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSymbol is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x95d89b41.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function symbol() view returns(string)
func (drexToken *DrexToken) TryPackSymbol() ([]byte, error) {
	return drexToken.abi.Pack("symbol")
}

// UnpackSymbol is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (drexToken *DrexToken) UnpackSymbol(data []byte) (string, error) {
	out, err := drexToken.abi.Unpack("symbol", data)
	if err != nil {
		return *new(string), err
	}
	out0 := *abi.ConvertType(out[0], new(string)).(*string)
	return out0, nil
}

// PackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalSupply() view returns(uint256)
func (drexToken *DrexToken) PackTotalSupply() []byte {
	enc, err := drexToken.abi.Pack("totalSupply")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalSupply is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x18160ddd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalSupply() view returns(uint256)
func (drexToken *DrexToken) TryPackTotalSupply() ([]byte, error) {
	return drexToken.abi.Pack("totalSupply")
}

// UnpackTotalSupply is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (drexToken *DrexToken) UnpackTotalSupply(data []byte) (*big.Int, error) {
	out, err := drexToken.abi.Unpack("totalSupply", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// PackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (drexToken *DrexToken) PackTransfer(to common.Address, value *big.Int) []byte {
	enc, err := drexToken.abi.Pack("transfer", to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransfer is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa9059cbb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (drexToken *DrexToken) TryPackTransfer(to common.Address, value *big.Int) ([]byte, error) {
	return drexToken.abi.Pack("transfer", to, value)
}

// UnpackTransfer is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (drexToken *DrexToken) UnpackTransfer(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("transfer", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (drexToken *DrexToken) PackTransferFrom(from common.Address, to common.Address, value *big.Int) []byte {
	enc, err := drexToken.abi.Pack("transferFrom", from, to, value)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTransferFrom is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x23b872dd.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (drexToken *DrexToken) TryPackTransferFrom(from common.Address, to common.Address, value *big.Int) ([]byte, error) {
	return drexToken.abi.Pack("transferFrom", from, to, value)
}

// UnpackTransferFrom is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (drexToken *DrexToken) UnpackTransferFrom(data []byte) (bool, error) {
	out, err := drexToken.abi.Unpack("transferFrom", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// DrexTokenAllowlistUpdated represents a AllowlistUpdated event raised by the DrexToken contract.
type DrexTokenAllowlistUpdated struct {
	Account common.Address
	Allowed bool
	Raw     *types.Log // Blockchain specific contextual infos
}

const DrexTokenAllowlistUpdatedEventName = "AllowlistUpdated"

// ContractEventName returns the user-defined event name.
func (DrexTokenAllowlistUpdated) ContractEventName() string {
	return DrexTokenAllowlistUpdatedEventName
}

// UnpackAllowlistUpdatedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AllowlistUpdated(address indexed account, bool allowed)
func (drexToken *DrexToken) UnpackAllowlistUpdatedEvent(log *types.Log) (*DrexTokenAllowlistUpdated, error) {
	event := "AllowlistUpdated"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenAllowlistUpdated)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DrexTokenApproval represents a Approval event raised by the DrexToken contract.
type DrexTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     *types.Log // Blockchain specific contextual infos
}

const DrexTokenApprovalEventName = "Approval"

// ContractEventName returns the user-defined event name.
func (DrexTokenApproval) ContractEventName() string {
	return DrexTokenApprovalEventName
}

// UnpackApprovalEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (drexToken *DrexToken) UnpackApprovalEvent(log *types.Log) (*DrexTokenApproval, error) {
	event := "Approval"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenApproval)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DrexTokenRoleAdminChanged represents a RoleAdminChanged event raised by the DrexToken contract.
type DrexTokenRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const DrexTokenRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (DrexTokenRoleAdminChanged) ContractEventName() string {
	return DrexTokenRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (drexToken *DrexToken) UnpackRoleAdminChangedEvent(log *types.Log) (*DrexTokenRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DrexTokenRoleGranted represents a RoleGranted event raised by the DrexToken contract.
type DrexTokenRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const DrexTokenRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (DrexTokenRoleGranted) ContractEventName() string {
	return DrexTokenRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (drexToken *DrexToken) UnpackRoleGrantedEvent(log *types.Log) (*DrexTokenRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenRoleGranted)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DrexTokenRoleRevoked represents a RoleRevoked event raised by the DrexToken contract.
type DrexTokenRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const DrexTokenRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (DrexTokenRoleRevoked) ContractEventName() string {
	return DrexTokenRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (drexToken *DrexToken) UnpackRoleRevokedEvent(log *types.Log) (*DrexTokenRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenRoleRevoked)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// DrexTokenTransfer represents a Transfer event raised by the DrexToken contract.
type DrexTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   *types.Log // Blockchain specific contextual infos
}

const DrexTokenTransferEventName = "Transfer"

// ContractEventName returns the user-defined event name.
func (DrexTokenTransfer) ContractEventName() string {
	return DrexTokenTransferEventName
}

// UnpackTransferEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (drexToken *DrexToken) UnpackTransferEvent(log *types.Log) (*DrexTokenTransfer, error) {
	event := "Transfer"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != drexToken.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(DrexTokenTransfer)
	if len(log.Data) > 0 {
		if err := drexToken.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range drexToken.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (drexToken *DrexToken) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], drexToken.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return drexToken.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return drexToken.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InsufficientAllowance"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InsufficientAllowanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InsufficientBalance"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InsufficientBalanceError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InvalidApprover"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InvalidApproverError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InvalidReceiver"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InvalidReceiverError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InvalidSender"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InvalidSenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ERC20InvalidSpender"].ID.Bytes()[:4]) {
		return drexToken.UnpackERC20InvalidSpenderError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["NotAllowlisted"].ID.Bytes()[:4]) {
		return drexToken.UnpackNotAllowlistedError(raw[4:])
	}
	if bytes.Equal(raw[:4], drexToken.abi.Errors["ZeroAddress"].ID.Bytes()[:4]) {
		return drexToken.UnpackZeroAddressError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// DrexTokenAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the DrexToken contract.
type DrexTokenAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func DrexTokenAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (drexToken *DrexToken) UnpackAccessControlBadConfirmationError(raw []byte) (*DrexTokenAccessControlBadConfirmation, error) {
	out := new(DrexTokenAccessControlBadConfirmation)
	if err := drexToken.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the DrexToken contract.
type DrexTokenAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func DrexTokenAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (drexToken *DrexToken) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*DrexTokenAccessControlUnauthorizedAccount, error) {
	out := new(DrexTokenAccessControlUnauthorizedAccount)
	if err := drexToken.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InsufficientAllowance represents a ERC20InsufficientAllowance error raised by the DrexToken contract.
type DrexTokenERC20InsufficientAllowance struct {
	Spender   common.Address
	Allowance *big.Int
	Needed    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func DrexTokenERC20InsufficientAllowanceErrorID() common.Hash {
	return common.HexToHash("0xfb8f41b23e99d2101d86da76cdfa87dd51c82ed07d3cb62cbc473e469dbc75c3")
}

// UnpackERC20InsufficientAllowanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientAllowance(address spender, uint256 allowance, uint256 needed)
func (drexToken *DrexToken) UnpackERC20InsufficientAllowanceError(raw []byte) (*DrexTokenERC20InsufficientAllowance, error) {
	out := new(DrexTokenERC20InsufficientAllowance)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InsufficientAllowance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InsufficientBalance represents a ERC20InsufficientBalance error raised by the DrexToken contract.
type DrexTokenERC20InsufficientBalance struct {
	Sender  common.Address
	Balance *big.Int
	Needed  *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func DrexTokenERC20InsufficientBalanceErrorID() common.Hash {
	return common.HexToHash("0xe450d38cd8d9f7d95077d567d60ed49c7254716e6ad08fc9872816c97e0ffec6")
}

// UnpackERC20InsufficientBalanceError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InsufficientBalance(address sender, uint256 balance, uint256 needed)
func (drexToken *DrexToken) UnpackERC20InsufficientBalanceError(raw []byte) (*DrexTokenERC20InsufficientBalance, error) {
	out := new(DrexTokenERC20InsufficientBalance)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InsufficientBalance", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InvalidApprover represents a ERC20InvalidApprover error raised by the DrexToken contract.
type DrexTokenERC20InvalidApprover struct {
	Approver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidApprover(address approver)
func DrexTokenERC20InvalidApproverErrorID() common.Hash {
	return common.HexToHash("0xe602df05cc75712490294c6c104ab7c17f4030363910a7a2626411c6d3118847")
}

// UnpackERC20InvalidApproverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidApprover(address approver)
func (drexToken *DrexToken) UnpackERC20InvalidApproverError(raw []byte) (*DrexTokenERC20InvalidApprover, error) {
	out := new(DrexTokenERC20InvalidApprover)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InvalidApprover", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InvalidReceiver represents a ERC20InvalidReceiver error raised by the DrexToken contract.
type DrexTokenERC20InvalidReceiver struct {
	Receiver common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func DrexTokenERC20InvalidReceiverErrorID() common.Hash {
	return common.HexToHash("0xec442f055133b72f3b2f9f0bb351c406b178527de2040a7d1feb4e058771f613")
}

// UnpackERC20InvalidReceiverError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidReceiver(address receiver)
func (drexToken *DrexToken) UnpackERC20InvalidReceiverError(raw []byte) (*DrexTokenERC20InvalidReceiver, error) {
	out := new(DrexTokenERC20InvalidReceiver)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InvalidReceiver", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InvalidSender represents a ERC20InvalidSender error raised by the DrexToken contract.
type DrexTokenERC20InvalidSender struct {
	Sender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSender(address sender)
func DrexTokenERC20InvalidSenderErrorID() common.Hash {
	return common.HexToHash("0x96c6fd1edd0cd6ef7ff0ecc0facdf53148dc0048b57fe58af65755250a7a96bd")
}

// UnpackERC20InvalidSenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSender(address sender)
func (drexToken *DrexToken) UnpackERC20InvalidSenderError(raw []byte) (*DrexTokenERC20InvalidSender, error) {
	out := new(DrexTokenERC20InvalidSender)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InvalidSender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenERC20InvalidSpender represents a ERC20InvalidSpender error raised by the DrexToken contract.
type DrexTokenERC20InvalidSpender struct {
	Spender common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ERC20InvalidSpender(address spender)
func DrexTokenERC20InvalidSpenderErrorID() common.Hash {
	return common.HexToHash("0x94280d62c347d8d9f4d59a76ea321452406db88df38e0c9da304f58b57b373a2")
}

// UnpackERC20InvalidSpenderError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ERC20InvalidSpender(address spender)
func (drexToken *DrexToken) UnpackERC20InvalidSpenderError(raw []byte) (*DrexTokenERC20InvalidSpender, error) {
	out := new(DrexTokenERC20InvalidSpender)
	if err := drexToken.abi.UnpackIntoInterface(out, "ERC20InvalidSpender", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenNotAllowlisted represents a NotAllowlisted error raised by the DrexToken contract.
type DrexTokenNotAllowlisted struct {
	Account common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotAllowlisted(address account)
func DrexTokenNotAllowlistedErrorID() common.Hash {
	return common.HexToHash("0x6f7a52759a0e05af1260e0db2b4d083b603c97568e62767ad42503ce98076109")
}

// UnpackNotAllowlistedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotAllowlisted(address account)
func (drexToken *DrexToken) UnpackNotAllowlistedError(raw []byte) (*DrexTokenNotAllowlisted, error) {
	out := new(DrexTokenNotAllowlisted)
	if err := drexToken.abi.UnpackIntoInterface(out, "NotAllowlisted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// DrexTokenZeroAddress represents a ZeroAddress error raised by the DrexToken contract.
type DrexTokenZeroAddress struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ZeroAddress()
func DrexTokenZeroAddressErrorID() common.Hash {
	return common.HexToHash("0xd92e233df2717d4a40030e20904abd27b68fcbeede117eaaccbbdac9618c8c73")
}

// UnpackZeroAddressError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ZeroAddress()
func (drexToken *DrexToken) UnpackZeroAddressError(raw []byte) (*DrexTokenZeroAddress, error) {
	out := new(DrexTokenZeroAddress)
	if err := drexToken.abi.UnpackIntoInterface(out, "ZeroAddress", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierMetaData contains all meta data concerning the FabricProofVerifier contract.
var FabricProofVerifierMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_TAG\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GOVERNOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activeOrgCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"qx\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"qy\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"}],\"name\":\"addAttestorKey\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"channelHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"namespaceHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fabricTxId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"}],\"internalType\":\"structFabricLockAttestation\",\"name\":\"a\",\"type\":\"tuple\"}],\"name\":\"digest\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"getAttestorKey\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"qx\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"qy\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"internalType\":\"structFabricProofVerifier.AttestorKey\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"}],\"name\":\"getOrg\",\"outputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"registered\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"keyCount\",\"type\":\"uint32\"},{\"internalType\":\"string\",\"name\":\"mspId\",\"type\":\"string\"}],\"internalType\":\"structFabricProofVerifier.Org\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"channelHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"namespaceHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fabricTxId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"}],\"internalType\":\"structFabricLockAttestation\",\"name\":\"a\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structOrgSignature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"isValid\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"orgs\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"mspId\",\"type\":\"string\"}],\"name\":\"registerOrg\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"assetIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"recipientIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"deadline\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"channelHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"namespaceHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"fabricTxId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"blockNumber\",\"type\":\"uint64\"}],\"internalType\":\"structFabricLockAttestation\",\"name\":\"a\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structOrgSignature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"requireValid\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"revokeAttestorKey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"setOrgActive\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"newThreshold\",\"type\":\"uint32\"}],\"name\":\"setThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"threshold\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"qx\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"qy\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"certHash\",\"type\":\"bytes32\"}],\"name\":\"AttestorKeyAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"AttestorKeyRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"mspId\",\"type\":\"string\"}],\"name\":\"OrgRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"active\",\"type\":\"bool\"}],\"name\":\"OrgStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"threshold\",\"type\":\"uint32\"}],\"name\":\"ThresholdChanged\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"InactiveSigner\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidPublicKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"InvalidSignature\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"threshold\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"activeOrgs\",\"type\":\"uint32\"}],\"name\":\"InvalidThreshold\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"}],\"name\":\"OrgAlreadyRegistered\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SignaturesNotSorted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"valid\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"threshold\",\"type\":\"uint32\"}],\"name\":\"ThresholdNotMet\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint32\",\"name\":\"keyIndex\",\"type\":\"uint32\"}],\"name\":\"UnknownKey\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"mspIdHash\",\"type\":\"bytes32\"}],\"name\":\"UnknownOrg\",\"type\":\"error\"}]",
	ID:  "07ae657b114866e70da145d49cb27579e9",
	Bin: "0x608060405234801561000f575f5ffd5b5060405161299e38038061299e83398101604081905261002e91610113565b6100385f8261006a565b506100637f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f558261006a565b5050610140565b5f828152602081815260408083206001600160a01b038516845290915281205460ff1661010a575f838152602081815260408083206001600160a01b03861684529091529020805460ff191660011790556100c23390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161010d565b505f5b92915050565b5f60208284031215610123575f5ffd5b81516001600160a01b0381168114610139575f5ffd5b9392505050565b6128518061014d5f395ff3fe608060405234801561000f575f5ffd5b5060043610610179575f3560e01c806391d14854116100d2578063cab19cdc11610088578063d73b018e11610063578063d73b018e146103b0578063e32e13eb146103c3578063f31cc6bf146103cb575f5ffd5b8063cab19cdc14610361578063ccc5749014610376578063d547741f1461039d575f5ffd5b8063a217fddf116100b8578063a217fddf14610334578063c2312b841461033b578063ca5a8b341461034e575f5ffd5b806391d14854146102de5780639499ffee14610321575f5ffd5b80634f971ba31161013257806371b312da1161010d57806371b312da146102985780637500c13d146102ab57806381e84589146102be575f5ffd5b80634f971ba31461022a5780636693f2531461023d57806366bf3b2314610250575f5ffd5b80632f2ff15d116101625780632f2ff15d146101d557806336568abe146101ea57806342cde4e8146101fd575f5ffd5b806301ffc9a71461017d578063248a9ca3146101a5575b5f5ffd5b61019061018b36600461209c565b6103db565b60405190151581526020015b60405180910390f35b6101c76101b33660046120db565b5f9081526020819052604090206001015490565b60405190815260200161019c565b6101e86101e33660046120f2565b610473565b005b6101e86101f83660046120f2565b61049d565b60045461021590640100000000900463ffffffff1681565b60405163ffffffff909116815260200161019c565b6101c761023836600461214f565b6104fb565b6101e861024b366004612182565b6106ad565b61026361025e36600461219b565b6107bb565b60405161019c919081518152602080830151908201526040808301519082015260609182015115159181019190915260800190565b6101e86102a63660046121c5565b61083b565b6101906102b93660046121ec565b610a04565b6102d16102cc3660046120db565b610a83565b60405161019c9190612270565b6101906102ec3660046120f2565b5f9182526020828152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b6101e861032f36600461219b565b610b8a565b6101c75f81565b6101c76103493660046122f4565b610c7b565b6101e861035c3660046121ec565b610edc565b610369611155565b60405161019c9190612362565b6101c77f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f5581565b6101e86103ab3660046120f2565b6111ab565b6102156103be3660046123a4565b6111cf565b6101c76113ab565b6004546102159063ffffffff1681565b5f7fffffffff0000000000000000000000000000000000000000000000000000000082167f7965db0b00000000000000000000000000000000000000000000000000000000148061046d57507f01ffc9a7000000000000000000000000000000000000000000000000000000007fffffffff000000000000000000000000000000000000000000000000000000008316145b92915050565b5f8281526020819052604090206001015461048d81611418565b6104978383611425565b50505050565b73ffffffffffffffffffffffffffffffffffffffff811633146104ec576040517f6697b23200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6104f6828261151e565b505050565b5f60028060405161052f907f6a616e75732e6661627269632d6c6f636b2d6174746573746174696f6e2e7631815260200190565b602060405180830381855afa15801561054a573d5f5f3e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061056d91906123d3565b4630853560208701356040880135606089013561059060a08b0160808c01612401565b8a60a001358b60c001358c60e001358d6101000160208101906105b39190612401565b60408051602081019d909d528c019a909a5273ffffffffffffffffffffffffffffffffffffffff90981660608b015260808a019690965260a089019490945260c088019290925260e087015267ffffffffffffffff908116610100870152610120860191909152610140850191909152610160840191909152166101808201526101a001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261066f9161241a565b602060405180830381855afa15801561068a573d5f5f3e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061046d91906123d3565b7f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f556106d781611418565b63ffffffff821615806106f5575060045463ffffffff908116908316115b1561074757600480546040517f07fb268a00000000000000000000000000000000000000000000000000000000815263ffffffff80861693820193909352911660248201526044015b60405180910390fd5b600480547fffffffffffffffffffffffffffffffffffffffffffffffff00000000ffffffff1664010000000063ffffffff8516908102919091179091556040519081527f844eb4d915c4a4493f4e43c43daab74aa3069d5732c9b77658d96be54b9a51459060200160405180910390a15050565b6107e460405180608001604052805f81526020015f81526020015f81526020015f151581525090565b505f82815260026020818152604080842063ffffffff8616855282529283902083516080810185528154815260018201549281019290925291820154928101929092526003015460ff161515606082015292915050565b7f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f5561086581611418565b5f838152600160205260409020805460ff166108b0576040517f88f42f0f0000000000000000000000000000000000000000000000000000000081526004810185905260240161073e565b805483151561010090910460ff161515036108cb5750505050565b80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101008415150217815582610918576004546109139060019063ffffffff1661245d565b61092c565b60045461092c9063ffffffff166001612479565b600480547fffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000001663ffffffff9283169081179182905564010000000090910490911611156109c257600480546040517f07fb268a00000000000000000000000000000000000000000000000000000000815263ffffffff64010000000083048116938201939093529116602482015260440161073e565b837feb37fa5e3e0d1aa836eb26e1f27b2b96b8aef29d37c09082e6d4b28a11b66172846040516109f6911515815260200190565b60405180910390a250505050565b6040517fca5a8b340000000000000000000000000000000000000000000000000000000081525f90309063ca5a8b3490610a46908790879087906004016124f5565b5f6040518083038186803b158015610a5c575f5ffd5b505afa925050508015610a6d575060015b610a7857505f610a7c565b5060015b9392505050565b604080516080810182525f80825260208201819052918101919091526060808201525f828152600160208181526040928390208351608081018552815460ff808216151583526101008204161515938201939093526201000090920463ffffffff16938201939093529082018054919291606084019190610b0390612586565b80601f0160208091040260200160405190810160405280929190818152602001828054610b2f90612586565b8015610b7a5780601f10610b5157610100808354040283529160200191610b7a565b820191905f5260205f20905b815481529060010190602001808311610b5d57829003601f168201915b5050505050815250509050919050565b7f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f55610bb481611418565b5f83815260026020908152604080832063ffffffff8616845290915290208054610c19576040517f758035110000000000000000000000000000000000000000000000000000000081526004810185905263ffffffff8416602482015260440161073e565b6003810180547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905560405163ffffffff84169085907f5d70e1d81e25d4e614dfcd3d2bb25d6cd967470fb50bb9cb15d89a3061b0ea0c905f90a350505050565b5f7f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f55610ca681611418565b8383604051610cb69291906125d1565b60408051918290039091205f8181526001602052919091205490925060ff1615610d0f576040517f16ab79840000000000000000000000000000000000000000000000000000000081526004810183905260240161073e565b60405180608001604052806001151581526020016001151581526020015f63ffffffff16815260200185858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201829052509390945250508481526001602081815260409283902085518154928701519487015163ffffffff1662010000027fffffffffffffffffffffffffffffffffffffffffffffffffffff00000000ffff951515610100027fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff921515929092167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00009094169390931717939093161782556060840151919250820190610e2b9082612658565b505060038054600181810183555f9283527fc2575a0e9e593c00f959f8c92f12db2869c3395a3b0502d05e2516446f71f85b909101859055600480549193509190610e7d90849063ffffffff16612479565b92506101000a81548163ffffffff021916908363ffffffff160217905550817f14727d186cc4010886531e53e6d4020256e23086e738cc52c3ac2a64043e477a8585604051610ecd92919061276f565b60405180910390a25092915050565b600454640100000000900463ffffffff165f819003610f3a57600480546040517f07fb268a0000000000000000000000000000000000000000000000000000000081525f9281019290925263ffffffff16602482015260440161073e565b5f610f44856104fb565b90505f80805b858110156110f55736878783818110610f6557610f656127bb565b90506080020190505f82118015610f7d575080358410155b15610fb4576040517f01eba55100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80355f8181526001602090815260408083206002835281842094985093918391610fe391908701908701612182565b63ffffffff16815260208101919091526040015f208254909150610100900460ff1615806110165750600381015460ff16155b1561106e57823561102d6040850160208601612182565b6040517f5661712c000000000000000000000000000000000000000000000000000000008152600481019290925263ffffffff16602482015260440161073e565b61108a8784604001358560600135845f015485600101546115d7565b6110e15782356110a06040850160208601612182565b6040517f3b8b21a2000000000000000000000000000000000000000000000000000000008152600481019290925263ffffffff16602482015260440161073e565b846001019450505050806001019050610f4a565b508363ffffffff168163ffffffff16101561114c576040517f0a5563a100000000000000000000000000000000000000000000000000000000815263ffffffff80831660048301528516602482015260440161073e565b50505050505050565b606060038054806020026020016040519081016040528092919081815260200182805480156111a157602002820191905f5260205f20905b81548152602001906001019080831161118d575b5050505050905090565b5f828152602081905260409020600101546111c581611418565b610497838361151e565b5f7f7935bd0ae54bc31f548c14dba4d37c5c64b3f8ca900cb468fb8abd54d5894f556111fa81611418565b5f868152600160205260409020805460ff16611245576040517f88f42f0f0000000000000000000000000000000000000000000000000000000081526004810188905260240161073e565b61124f8686611612565b611285576040517fa2d0fee800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805462010000900463ffffffff1681600261129f836127e8565b91906101000a81548163ffffffff021916908363ffffffff160217905550925060405180608001604052808781526020018681526020018581526020016001151581525060025f8981526020019081526020015f205f8563ffffffff1663ffffffff1681526020019081526020015f205f820151815f015560208201518160010155604082015181600201556060820151816003015f6101000a81548160ff0219169083151502179055509050508263ffffffff16877fd592f8b6157285258913023b5145c77ceaaa14c320930ac94332bfac7fe25017888888604051611399939291909283526020830191909152604082015260600190565b60405180910390a35050949350505050565b6040517f6a616e75732e6661627269632d6c6f636b2d6174746573746174696f6e2e76318152600290602001602060405180830381855afa1580156113f2573d5f5f3e3d5ffd5b5050506040513d601f19601f8201168201806040525081019061141591906123d3565b81565b6114228133611697565b50565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff16611517575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff86168452909152902080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556114b53390565b73ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161046d565b505f61046d565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915281205460ff1615611517575f8381526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8616808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a450600161046d565b5f5f5f6115e78888888888611720565b9150915080611602576115fd88888888886117fb565b611604565b815b925050505b95945050505050565b5f7fffffffff00000001000000000000000000000000ffffffffffffffffffffffff80838409817f5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b8387857fffffffff00000001000000000000000000000000fffffffffffffffffffffffc878b8c0908090894821191909310169190921416919050565b5f8281526020818152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1661171c576040517fe2517d3f00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024810183905260440161073e565b5050565b5f5f61172c86866118f7565b158061173f575061173d8484611612565b155b1561174f57505f905060016117f1565b61175c8787878787611961565b1561176c575060019050806117f1565b6117db7fbb5a52f42f9c9261ed4361f59422a1e30036e7c32b270c8807a419feca605023600560017fa71af64de5126a4a4e02b7922d66ce9415ce88a4c9d25514d91082c8725ac9577f5d47723c8fbe580bb369fec9c2665d8e30a435b9932645482e7c9f11e872296b611961565b156117eb57505f905060016117f1565b505f9050805b9550959350505050565b5f61180685856118f7565b158061181957506118178383611612565b155b1561182557505f611609565b5f61183084846119a1565b90505f61185d867fffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551611b94565b90505f7fffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551828a0990505f7fffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551838a0990505f6118b9858484611ba3565b509050896118e77fffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc6325518361280c565b149b9a5050505050505050505050565b5f821580159061192657507fffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc63255183105b801561193157508115155b8015610a7c5750507f7fffffff800000007fffffffffffffffde737d56d38bcf4279dce5617e3192a81015919050565b5f6040518681528560208201528460408201528360608201528260808201525f5f5260205f60a0836101005afa61199457fe5b50505f5195945050505050565b6119a9612055565b60405180606001604052805f81526020015f81526020015f815250815f601081106119d6576119d66127bb565b60200201819052506040518060600160405280848152602001838152602001600181525081600160108110611a0d57611a0d6127bb565b602002018190525060405180606001604052807f6b17d1f2e12c4247f8bce6e563a440f277037d812deb33a0f4a13945d898c29681526020017f4fe342e2fe1a7f9b8ee7eb4a7c0f9e162bce33576b315ececbb6406837bf51f58152602001600181525081600460108110611a8457611a846127bb565b6020020152611a9a8160015b6020020151611cc7565b6040820152611aaa816004611a90565b6101008201526020810151611ac7908260025b6020020151611d23565b6060820152611ae08160015b6020020151826004611abd565b60a0820152611af0816002611ad3565b60c0820152611b00816003611ad3565b60e0820152611b198160015b6020020151826008611abd565b610120820152611b2a816002611b0c565b610140820152611b3b816003611b0c565b610160820152611b4c816004611b0c565b610180820152611b668160015b602002015182600c611abd565b6101a0820152611b77816002611b59565b6101c0820152611b88816003611b59565b6101e082015292915050565b5f610a7c836002840384611d81565b5f80808080805b6080811015611cab578115611bde57611bc4848484611da2565b91955093509150611bd6848484611da2565b919550935091505b600c60fc89901c1660fe88901c17898160108110611bfe57611bfe6127bb565b60200201516040015115611c9757825f03611c6d57898160108110611c2557611c256127bb565b6020020151518a8260108110611c3d57611c3d6127bb565b6020020151602001518b8360108110611c5857611c586127bb565b60200201516040015191965094509250611c97565b611c8f8a8260108110611c8257611c826127bb565b6020020151868686611e46565b919650945092505b50600297881b979690961b95600101611baa565b50611cb7838383611f93565b945094505050505b935093915050565b611ce860405180606001604052805f81526020015f81526020015f81525090565b5f5f5f611d01855f015186602001518760400151611da2565b6040805160608101825293845260208401929092529082015295945050505050565b611d4460405180606001604052805f81526020015f81526020015f81525090565b5f5f5f611d5e86865f015187602001518860400151611e46565b604080516060810182529384526020840192909252908201529695505050505050565b5f5f5f611d8f868686611ff2565b9150915081611609576116096012612044565b5f5f5f7fffffffff00000001000000000000000000000000ffffffffffffffffffffffff80868709818687098283848384097fffffffff00000001000000000000000000000000fffffffffffffffffffffffc0984858c8d096003090890508283838b09600409838482600209850385848509089650838485858609600809850385868a880385088509089550505050808186880960020991505093509350939050565b5f5f5f7fffffffff00000001000000000000000000000000ffffffffffffffffffffffff604088015181818209828388858a8b090960208c0151098381850385868686098c090884858a8b098d51098581870387868f0908935081158415168015611eb85760018114611efe57611f82565b868586098788898386096002098903898a848a098b038b88890908089a5087888983890987098903898a8e8c038c8689090887090899505086878c880986099750611f82565b8c8c8c898283098a8283098b8c8d8384097fffffffff00000001000000000000000000000000fffffffffffffffffffffffc098d8e8889096003090890508b8c83870960040994508b8c866002098d038d838409089e508b8c8384096008098c0391508b8f8d03860894508b828d878409089d505050898a8284096002099a505050505b505050505050509450945094915050565b5f5f825f03611fa657505f905080611cbf565b7fffffffff00000001000000000000000000000000ffffffffffffffffffffffff5f611fd28583611b94565b905081818209828189099450828383830988099350505050935093915050565b5f5f825f0361200557505f905080611cbf565b60405160208152602080820152602060408201528560608201528460808201528360a082015260205f60c08360055afa9250505f519050935093915050565b634e487b715f52806020526024601cfd5b6040518061020001604052806010905b61208660405180606001604052805f81526020015f81526020015f81525090565b8152602001906001900390816120655790505090565b5f602082840312156120ac575f5ffd5b81357fffffffff0000000000000000000000000000000000000000000000000000000081168114610a7c575f5ffd5b5f602082840312156120eb575f5ffd5b5035919050565b5f5f60408385031215612103575f5ffd5b82359150602083013573ffffffffffffffffffffffffffffffffffffffff8116811461212d575f5ffd5b809150509250929050565b5f6101208284031215612149575f5ffd5b50919050565b5f6101208284031215612160575f5ffd5b610a7c8383612138565b803563ffffffff8116811461217d575f5ffd5b919050565b5f60208284031215612192575f5ffd5b610a7c8261216a565b5f5f604083850312156121ac575f5ffd5b823591506121bc6020840161216a565b90509250929050565b5f5f604083850312156121d6575f5ffd5b823591506020830135801515811461212d575f5ffd5b5f5f5f61014084860312156121ff575f5ffd5b6122098585612138565b925061012084013567ffffffffffffffff811115612225575f5ffd5b8401601f81018613612235575f5ffd5b803567ffffffffffffffff81111561224b575f5ffd5b8660208260071b840101111561225f575f5ffd5b939660209190910195509293505050565b6020815281511515602082015260208201511515604082015263ffffffff60408301511660608201525f606083015160808084015280518060a0850152806020830160c086015e5f60c0828601015260c07fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8301168501019250505092915050565b5f5f60208385031215612305575f5ffd5b823567ffffffffffffffff81111561231b575f5ffd5b8301601f8101851361232b575f5ffd5b803567ffffffffffffffff811115612341575f5ffd5b856020828401011115612352575f5ffd5b6020919091019590945092505050565b602080825282518282018190525f918401906040840190835b8181101561239957835183526020938401939092019160010161237b565b509095945050505050565b5f5f5f5f608085870312156123b7575f5ffd5b5050823594602084013594506040840135936060013592509050565b5f602082840312156123e3575f5ffd5b5051919050565b803567ffffffffffffffff8116811461217d575f5ffd5b5f60208284031215612411575f5ffd5b610a7c826123ea565b5f82518060208501845e5f920191825250919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b63ffffffff828116828216039081111561046d5761046d612430565b63ffffffff818116838216019081111561046d5761046d612430565b8183526020830192505f815f5b848110156124eb578135865263ffffffff6124bf6020840161216a565b1660208701526040828101359087015260608083013590870152608095860195909101906001016124a2565b5093949350505050565b833581526020808501359082015260408085013590820152606080850135908201525f67ffffffffffffffff61252d608087016123ea565b16608083015260a0858101359083015260c0808601359083015260e0808601359083015267ffffffffffffffff61256761010087016123ea565b1661010083015261014061012083015261160961014083018486612495565b600181811c9082168061259a57607f821691505b602082108103612149577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b818382375f9101908152919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b601f8211156104f657805f5260205f20601f840160051c810160208510156126325750805b601f840160051c820191505b81811015612651575f815560010161263e565b5050505050565b815167ffffffffffffffff811115612672576126726125e0565b612686816126808454612586565b8461260d565b6020601f8211600181146126d7575f83156126a15750848201515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600385901b1c1916600184901b178455612651565b5f848152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08516915b828110156127245787850151825560209485019460019092019101612704565b508482101561276057868401517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600387901b60f8161c191681555b50505050600190811b01905550565b60208152816020820152818360408301375f818301604090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0160101919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f63ffffffff821663ffffffff810361280357612803612430565b60010192915050565b5f8261283f577f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b50069056fea164736f6c634300081e000a",
}

// FabricProofVerifier is an auto generated Go binding around an Ethereum contract.
type FabricProofVerifier struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *FabricProofVerifier) GetABI() abi.ABI {
	return c.abi
}

// NewFabricProofVerifier creates a new instance of FabricProofVerifier.
func NewFabricProofVerifier() *FabricProofVerifier {
	parsed, err := FabricProofVerifierMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &FabricProofVerifier{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *FabricProofVerifier) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address admin) returns()
func (fabricProofVerifier *FabricProofVerifier) PackConstructor(admin common.Address) []byte {
	enc, err := fabricProofVerifier.abi.Pack("", admin)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) PackDEFAULTADMINROLE() []byte {
	enc, err := fabricProofVerifier.abi.Pack("DEFAULT_ADMIN_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDEFAULTADMINROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xa217fddf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) TryPackDEFAULTADMINROLE() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("DEFAULT_ADMIN_ROLE")
}

// UnpackDEFAULTADMINROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) UnpackDEFAULTADMINROLE(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("DEFAULT_ADMIN_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackDOMAINTAG is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe32e13eb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function DOMAIN_TAG() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) PackDOMAINTAG() []byte {
	enc, err := fabricProofVerifier.abi.Pack("DOMAIN_TAG")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDOMAINTAG is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xe32e13eb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function DOMAIN_TAG() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) TryPackDOMAINTAG() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("DOMAIN_TAG")
}

// UnpackDOMAINTAG is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xe32e13eb.
//
// Solidity: function DOMAIN_TAG() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) UnpackDOMAINTAG(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("DOMAIN_TAG", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGOVERNORROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccc57490.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function GOVERNOR_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) PackGOVERNORROLE() []byte {
	enc, err := fabricProofVerifier.abi.Pack("GOVERNOR_ROLE")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGOVERNORROLE is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xccc57490.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function GOVERNOR_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) TryPackGOVERNORROLE() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("GOVERNOR_ROLE")
}

// UnpackGOVERNORROLE is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xccc57490.
//
// Solidity: function GOVERNOR_ROLE() view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) UnpackGOVERNORROLE(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("GOVERNOR_ROLE", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackActiveOrgCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf31cc6bf.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function activeOrgCount() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) PackActiveOrgCount() []byte {
	enc, err := fabricProofVerifier.abi.Pack("activeOrgCount")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackActiveOrgCount is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xf31cc6bf.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function activeOrgCount() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) TryPackActiveOrgCount() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("activeOrgCount")
}

// UnpackActiveOrgCount is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xf31cc6bf.
//
// Solidity: function activeOrgCount() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) UnpackActiveOrgCount(data []byte) (uint32, error) {
	out, err := fabricProofVerifier.abi.Unpack("activeOrgCount", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackAddAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd73b018e.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function addAttestorKey(bytes32 mspIdHash, bytes32 qx, bytes32 qy, bytes32 certHash) returns(uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) PackAddAttestorKey(mspIdHash [32]byte, qx [32]byte, qy [32]byte, certHash [32]byte) []byte {
	enc, err := fabricProofVerifier.abi.Pack("addAttestorKey", mspIdHash, qx, qy, certHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackAddAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd73b018e.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function addAttestorKey(bytes32 mspIdHash, bytes32 qx, bytes32 qy, bytes32 certHash) returns(uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) TryPackAddAttestorKey(mspIdHash [32]byte, qx [32]byte, qy [32]byte, certHash [32]byte) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("addAttestorKey", mspIdHash, qx, qy, certHash)
}

// UnpackAddAttestorKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd73b018e.
//
// Solidity: function addAttestorKey(bytes32 mspIdHash, bytes32 qx, bytes32 qy, bytes32 certHash) returns(uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) UnpackAddAttestorKey(data []byte) (uint32, error) {
	out, err := fabricProofVerifier.abi.Unpack("addAttestorKey", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// PackDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f971ba3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function digest((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) PackDigest(a FabricLockAttestation) []byte {
	enc, err := fabricProofVerifier.abi.Pack("digest", a)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackDigest is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x4f971ba3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function digest((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) TryPackDigest(a FabricLockAttestation) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("digest", a)
}

// UnpackDigest is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x4f971ba3.
//
// Solidity: function digest((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) UnpackDigest(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("digest", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGetAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66bf3b23.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getAttestorKey(bytes32 mspIdHash, uint32 keyIndex) view returns((bytes32,bytes32,bytes32,bool))
func (fabricProofVerifier *FabricProofVerifier) PackGetAttestorKey(mspIdHash [32]byte, keyIndex uint32) []byte {
	enc, err := fabricProofVerifier.abi.Pack("getAttestorKey", mspIdHash, keyIndex)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x66bf3b23.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getAttestorKey(bytes32 mspIdHash, uint32 keyIndex) view returns((bytes32,bytes32,bytes32,bool))
func (fabricProofVerifier *FabricProofVerifier) TryPackGetAttestorKey(mspIdHash [32]byte, keyIndex uint32) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("getAttestorKey", mspIdHash, keyIndex)
}

// UnpackGetAttestorKey is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x66bf3b23.
//
// Solidity: function getAttestorKey(bytes32 mspIdHash, uint32 keyIndex) view returns((bytes32,bytes32,bytes32,bool))
func (fabricProofVerifier *FabricProofVerifier) UnpackGetAttestorKey(data []byte) (FabricProofVerifierAttestorKey, error) {
	out, err := fabricProofVerifier.abi.Unpack("getAttestorKey", data)
	if err != nil {
		return *new(FabricProofVerifierAttestorKey), err
	}
	out0 := *abi.ConvertType(out[0], new(FabricProofVerifierAttestorKey)).(*FabricProofVerifierAttestorKey)
	return out0, nil
}

// PackGetOrg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81e84589.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getOrg(bytes32 mspIdHash) view returns((bool,bool,uint32,string))
func (fabricProofVerifier *FabricProofVerifier) PackGetOrg(mspIdHash [32]byte) []byte {
	enc, err := fabricProofVerifier.abi.Pack("getOrg", mspIdHash)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetOrg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x81e84589.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getOrg(bytes32 mspIdHash) view returns((bool,bool,uint32,string))
func (fabricProofVerifier *FabricProofVerifier) TryPackGetOrg(mspIdHash [32]byte) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("getOrg", mspIdHash)
}

// UnpackGetOrg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x81e84589.
//
// Solidity: function getOrg(bytes32 mspIdHash) view returns((bool,bool,uint32,string))
func (fabricProofVerifier *FabricProofVerifier) UnpackGetOrg(data []byte) (FabricProofVerifierOrg, error) {
	out, err := fabricProofVerifier.abi.Unpack("getOrg", data)
	if err != nil {
		return *new(FabricProofVerifierOrg), err
	}
	out0 := *abi.ConvertType(out[0], new(FabricProofVerifierOrg)).(*FabricProofVerifierOrg)
	return out0, nil
}

// PackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) PackGetRoleAdmin(role [32]byte) []byte {
	enc, err := fabricProofVerifier.abi.Pack("getRoleAdmin", role)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetRoleAdmin is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x248a9ca3.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) TryPackGetRoleAdmin(role [32]byte) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("getRoleAdmin", role)
}

// UnpackGetRoleAdmin is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (fabricProofVerifier *FabricProofVerifier) UnpackGetRoleAdmin(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("getRoleAdmin", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (fabricProofVerifier *FabricProofVerifier) PackGrantRole(role [32]byte, account common.Address) []byte {
	enc, err := fabricProofVerifier.abi.Pack("grantRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGrantRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x2f2ff15d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackGrantRole(role [32]byte, account common.Address) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("grantRole", role, account)
}

// PackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) PackHasRole(role [32]byte, account common.Address) []byte {
	enc, err := fabricProofVerifier.abi.Pack("hasRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHasRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x91d14854.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) TryPackHasRole(role [32]byte, account common.Address) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("hasRole", role, account)
}

// UnpackHasRole is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) UnpackHasRole(data []byte) (bool, error) {
	out, err := fabricProofVerifier.abi.Unpack("hasRole", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackIsValid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7500c13d.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function isValid((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a, (bytes32,uint32,bytes32,bytes32)[] sigs) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) PackIsValid(a FabricLockAttestation, sigs []OrgSignature) []byte {
	enc, err := fabricProofVerifier.abi.Pack("isValid", a, sigs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackIsValid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7500c13d.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function isValid((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a, (bytes32,uint32,bytes32,bytes32)[] sigs) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) TryPackIsValid(a FabricLockAttestation, sigs []OrgSignature) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("isValid", a, sigs)
}

// UnpackIsValid is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x7500c13d.
//
// Solidity: function isValid((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a, (bytes32,uint32,bytes32,bytes32)[] sigs) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) UnpackIsValid(data []byte) (bool, error) {
	out, err := fabricProofVerifier.abi.Unpack("isValid", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackOrgs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcab19cdc.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function orgs() view returns(bytes32[])
func (fabricProofVerifier *FabricProofVerifier) PackOrgs() []byte {
	enc, err := fabricProofVerifier.abi.Pack("orgs")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackOrgs is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xcab19cdc.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function orgs() view returns(bytes32[])
func (fabricProofVerifier *FabricProofVerifier) TryPackOrgs() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("orgs")
}

// UnpackOrgs is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xcab19cdc.
//
// Solidity: function orgs() view returns(bytes32[])
func (fabricProofVerifier *FabricProofVerifier) UnpackOrgs(data []byte) ([][32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("orgs", data)
	if err != nil {
		return *new([][32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)
	return out0, nil
}

// PackRegisterOrg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc2312b84.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function registerOrg(string mspId) returns(bytes32 mspIdHash)
func (fabricProofVerifier *FabricProofVerifier) PackRegisterOrg(mspId string) []byte {
	enc, err := fabricProofVerifier.abi.Pack("registerOrg", mspId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRegisterOrg is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xc2312b84.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function registerOrg(string mspId) returns(bytes32 mspIdHash)
func (fabricProofVerifier *FabricProofVerifier) TryPackRegisterOrg(mspId string) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("registerOrg", mspId)
}

// UnpackRegisterOrg is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xc2312b84.
//
// Solidity: function registerOrg(string mspId) returns(bytes32 mspIdHash)
func (fabricProofVerifier *FabricProofVerifier) UnpackRegisterOrg(data []byte) ([32]byte, error) {
	out, err := fabricProofVerifier.abi.Unpack("registerOrg", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (fabricProofVerifier *FabricProofVerifier) PackRenounceRole(role [32]byte, callerConfirmation common.Address) []byte {
	enc, err := fabricProofVerifier.abi.Pack("renounceRole", role, callerConfirmation)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRenounceRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x36568abe.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackRenounceRole(role [32]byte, callerConfirmation common.Address) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("renounceRole", role, callerConfirmation)
}

// PackRequireValid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xca5a8b34.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function requireValid((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a, (bytes32,uint32,bytes32,bytes32)[] sigs) view returns()
func (fabricProofVerifier *FabricProofVerifier) PackRequireValid(a FabricLockAttestation, sigs []OrgSignature) []byte {
	enc, err := fabricProofVerifier.abi.Pack("requireValid", a, sigs)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRequireValid is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xca5a8b34.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function requireValid((bytes32,bytes32,bytes32,bytes32,uint64,bytes32,bytes32,bytes32,uint64) a, (bytes32,uint32,bytes32,bytes32)[] sigs) view returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackRequireValid(a FabricLockAttestation, sigs []OrgSignature) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("requireValid", a, sigs)
}

// PackRevokeAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9499ffee.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeAttestorKey(bytes32 mspIdHash, uint32 keyIndex) returns()
func (fabricProofVerifier *FabricProofVerifier) PackRevokeAttestorKey(mspIdHash [32]byte, keyIndex uint32) []byte {
	enc, err := fabricProofVerifier.abi.Pack("revokeAttestorKey", mspIdHash, keyIndex)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeAttestorKey is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x9499ffee.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeAttestorKey(bytes32 mspIdHash, uint32 keyIndex) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackRevokeAttestorKey(mspIdHash [32]byte, keyIndex uint32) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("revokeAttestorKey", mspIdHash, keyIndex)
}

// PackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (fabricProofVerifier *FabricProofVerifier) PackRevokeRole(role [32]byte, account common.Address) []byte {
	enc, err := fabricProofVerifier.abi.Pack("revokeRole", role, account)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevokeRole is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd547741f.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackRevokeRole(role [32]byte, account common.Address) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("revokeRole", role, account)
}

// PackSetOrgActive is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71b312da.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setOrgActive(bytes32 mspIdHash, bool active) returns()
func (fabricProofVerifier *FabricProofVerifier) PackSetOrgActive(mspIdHash [32]byte, active bool) []byte {
	enc, err := fabricProofVerifier.abi.Pack("setOrgActive", mspIdHash, active)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetOrgActive is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x71b312da.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setOrgActive(bytes32 mspIdHash, bool active) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackSetOrgActive(mspIdHash [32]byte, active bool) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("setOrgActive", mspIdHash, active)
}

// PackSetThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6693f253.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function setThreshold(uint32 newThreshold) returns()
func (fabricProofVerifier *FabricProofVerifier) PackSetThreshold(newThreshold uint32) []byte {
	enc, err := fabricProofVerifier.abi.Pack("setThreshold", newThreshold)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSetThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x6693f253.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function setThreshold(uint32 newThreshold) returns()
func (fabricProofVerifier *FabricProofVerifier) TryPackSetThreshold(newThreshold uint32) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("setThreshold", newThreshold)
}

// PackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) PackSupportsInterface(interfaceId [4]byte) []byte {
	enc, err := fabricProofVerifier.abi.Pack("supportsInterface", interfaceId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackSupportsInterface is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x01ffc9a7.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) TryPackSupportsInterface(interfaceId [4]byte) ([]byte, error) {
	return fabricProofVerifier.abi.Pack("supportsInterface", interfaceId)
}

// UnpackSupportsInterface is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (fabricProofVerifier *FabricProofVerifier) UnpackSupportsInterface(data []byte) (bool, error) {
	out, err := fabricProofVerifier.abi.Unpack("supportsInterface", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42cde4e8.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function threshold() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) PackThreshold() []byte {
	enc, err := fabricProofVerifier.abi.Pack("threshold")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackThreshold is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x42cde4e8.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function threshold() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) TryPackThreshold() ([]byte, error) {
	return fabricProofVerifier.abi.Pack("threshold")
}

// UnpackThreshold is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x42cde4e8.
//
// Solidity: function threshold() view returns(uint32)
func (fabricProofVerifier *FabricProofVerifier) UnpackThreshold(data []byte) (uint32, error) {
	out, err := fabricProofVerifier.abi.Unpack("threshold", data)
	if err != nil {
		return *new(uint32), err
	}
	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)
	return out0, nil
}

// FabricProofVerifierAttestorKeyAdded represents a AttestorKeyAdded event raised by the FabricProofVerifier contract.
type FabricProofVerifierAttestorKeyAdded struct {
	MspIdHash [32]byte
	KeyIndex  uint32
	Qx        [32]byte
	Qy        [32]byte
	CertHash  [32]byte
	Raw       *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierAttestorKeyAddedEventName = "AttestorKeyAdded"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierAttestorKeyAdded) ContractEventName() string {
	return FabricProofVerifierAttestorKeyAddedEventName
}

// UnpackAttestorKeyAddedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AttestorKeyAdded(bytes32 indexed mspIdHash, uint32 indexed keyIndex, bytes32 qx, bytes32 qy, bytes32 certHash)
func (fabricProofVerifier *FabricProofVerifier) UnpackAttestorKeyAddedEvent(log *types.Log) (*FabricProofVerifierAttestorKeyAdded, error) {
	event := "AttestorKeyAdded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierAttestorKeyAdded)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierAttestorKeyRevoked represents a AttestorKeyRevoked event raised by the FabricProofVerifier contract.
type FabricProofVerifierAttestorKeyRevoked struct {
	MspIdHash [32]byte
	KeyIndex  uint32
	Raw       *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierAttestorKeyRevokedEventName = "AttestorKeyRevoked"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierAttestorKeyRevoked) ContractEventName() string {
	return FabricProofVerifierAttestorKeyRevokedEventName
}

// UnpackAttestorKeyRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event AttestorKeyRevoked(bytes32 indexed mspIdHash, uint32 indexed keyIndex)
func (fabricProofVerifier *FabricProofVerifier) UnpackAttestorKeyRevokedEvent(log *types.Log) (*FabricProofVerifierAttestorKeyRevoked, error) {
	event := "AttestorKeyRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierAttestorKeyRevoked)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierOrgRegistered represents a OrgRegistered event raised by the FabricProofVerifier contract.
type FabricProofVerifierOrgRegistered struct {
	MspIdHash [32]byte
	MspId     string
	Raw       *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierOrgRegisteredEventName = "OrgRegistered"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierOrgRegistered) ContractEventName() string {
	return FabricProofVerifierOrgRegisteredEventName
}

// UnpackOrgRegisteredEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OrgRegistered(bytes32 indexed mspIdHash, string mspId)
func (fabricProofVerifier *FabricProofVerifier) UnpackOrgRegisteredEvent(log *types.Log) (*FabricProofVerifierOrgRegistered, error) {
	event := "OrgRegistered"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierOrgRegistered)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierOrgStatusChanged represents a OrgStatusChanged event raised by the FabricProofVerifier contract.
type FabricProofVerifierOrgStatusChanged struct {
	MspIdHash [32]byte
	Active    bool
	Raw       *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierOrgStatusChangedEventName = "OrgStatusChanged"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierOrgStatusChanged) ContractEventName() string {
	return FabricProofVerifierOrgStatusChangedEventName
}

// UnpackOrgStatusChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event OrgStatusChanged(bytes32 indexed mspIdHash, bool active)
func (fabricProofVerifier *FabricProofVerifier) UnpackOrgStatusChangedEvent(log *types.Log) (*FabricProofVerifierOrgStatusChanged, error) {
	event := "OrgStatusChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierOrgStatusChanged)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierRoleAdminChanged represents a RoleAdminChanged event raised by the FabricProofVerifier contract.
type FabricProofVerifierRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierRoleAdminChangedEventName = "RoleAdminChanged"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierRoleAdminChanged) ContractEventName() string {
	return FabricProofVerifierRoleAdminChangedEventName
}

// UnpackRoleAdminChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (fabricProofVerifier *FabricProofVerifier) UnpackRoleAdminChangedEvent(log *types.Log) (*FabricProofVerifierRoleAdminChanged, error) {
	event := "RoleAdminChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierRoleAdminChanged)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierRoleGranted represents a RoleGranted event raised by the FabricProofVerifier contract.
type FabricProofVerifierRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierRoleGrantedEventName = "RoleGranted"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierRoleGranted) ContractEventName() string {
	return FabricProofVerifierRoleGrantedEventName
}

// UnpackRoleGrantedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (fabricProofVerifier *FabricProofVerifier) UnpackRoleGrantedEvent(log *types.Log) (*FabricProofVerifierRoleGranted, error) {
	event := "RoleGranted"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierRoleGranted)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierRoleRevoked represents a RoleRevoked event raised by the FabricProofVerifier contract.
type FabricProofVerifierRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierRoleRevokedEventName = "RoleRevoked"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierRoleRevoked) ContractEventName() string {
	return FabricProofVerifierRoleRevokedEventName
}

// UnpackRoleRevokedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (fabricProofVerifier *FabricProofVerifier) UnpackRoleRevokedEvent(log *types.Log) (*FabricProofVerifierRoleRevoked, error) {
	event := "RoleRevoked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierRoleRevoked)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// FabricProofVerifierThresholdChanged represents a ThresholdChanged event raised by the FabricProofVerifier contract.
type FabricProofVerifierThresholdChanged struct {
	Threshold uint32
	Raw       *types.Log // Blockchain specific contextual infos
}

const FabricProofVerifierThresholdChangedEventName = "ThresholdChanged"

// ContractEventName returns the user-defined event name.
func (FabricProofVerifierThresholdChanged) ContractEventName() string {
	return FabricProofVerifierThresholdChangedEventName
}

// UnpackThresholdChangedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event ThresholdChanged(uint32 threshold)
func (fabricProofVerifier *FabricProofVerifier) UnpackThresholdChangedEvent(log *types.Log) (*FabricProofVerifierThresholdChanged, error) {
	event := "ThresholdChanged"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != fabricProofVerifier.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(FabricProofVerifierThresholdChanged)
	if len(log.Data) > 0 {
		if err := fabricProofVerifier.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range fabricProofVerifier.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (fabricProofVerifier *FabricProofVerifier) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["AccessControlBadConfirmation"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackAccessControlBadConfirmationError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["AccessControlUnauthorizedAccount"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackAccessControlUnauthorizedAccountError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["InactiveSigner"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackInactiveSignerError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["InvalidPublicKey"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackInvalidPublicKeyError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["InvalidSignature"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackInvalidSignatureError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["InvalidThreshold"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackInvalidThresholdError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["OrgAlreadyRegistered"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackOrgAlreadyRegisteredError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["SignaturesNotSorted"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackSignaturesNotSortedError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["ThresholdNotMet"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackThresholdNotMetError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["UnknownKey"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackUnknownKeyError(raw[4:])
	}
	if bytes.Equal(raw[:4], fabricProofVerifier.abi.Errors["UnknownOrg"].ID.Bytes()[:4]) {
		return fabricProofVerifier.UnpackUnknownOrgError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// FabricProofVerifierAccessControlBadConfirmation represents a AccessControlBadConfirmation error raised by the FabricProofVerifier contract.
type FabricProofVerifierAccessControlBadConfirmation struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlBadConfirmation()
func FabricProofVerifierAccessControlBadConfirmationErrorID() common.Hash {
	return common.HexToHash("0x6697b23232a647058342c0724fe7c415cab25915b54e5dbc03f233173d37b41c")
}

// UnpackAccessControlBadConfirmationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlBadConfirmation()
func (fabricProofVerifier *FabricProofVerifier) UnpackAccessControlBadConfirmationError(raw []byte) (*FabricProofVerifierAccessControlBadConfirmation, error) {
	out := new(FabricProofVerifierAccessControlBadConfirmation)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "AccessControlBadConfirmation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierAccessControlUnauthorizedAccount represents a AccessControlUnauthorizedAccount error raised by the FabricProofVerifier contract.
type FabricProofVerifierAccessControlUnauthorizedAccount struct {
	Account    common.Address
	NeededRole [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func FabricProofVerifierAccessControlUnauthorizedAccountErrorID() common.Hash {
	return common.HexToHash("0xe2517d3fbfae6f8515ef5ff1ccedc3933ab0cbbda0b492c06eb54ad10ef03b3e")
}

// UnpackAccessControlUnauthorizedAccountError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error AccessControlUnauthorizedAccount(address account, bytes32 neededRole)
func (fabricProofVerifier *FabricProofVerifier) UnpackAccessControlUnauthorizedAccountError(raw []byte) (*FabricProofVerifierAccessControlUnauthorizedAccount, error) {
	out := new(FabricProofVerifierAccessControlUnauthorizedAccount)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "AccessControlUnauthorizedAccount", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierInactiveSigner represents a InactiveSigner error raised by the FabricProofVerifier contract.
type FabricProofVerifierInactiveSigner struct {
	MspIdHash [32]byte
	KeyIndex  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InactiveSigner(bytes32 mspIdHash, uint32 keyIndex)
func FabricProofVerifierInactiveSignerErrorID() common.Hash {
	return common.HexToHash("0x5661712cdf4f4fdf8a48e6281c78db3b3f928455039e55f8408da3be943a6fc7")
}

// UnpackInactiveSignerError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InactiveSigner(bytes32 mspIdHash, uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) UnpackInactiveSignerError(raw []byte) (*FabricProofVerifierInactiveSigner, error) {
	out := new(FabricProofVerifierInactiveSigner)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "InactiveSigner", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierInvalidPublicKey represents a InvalidPublicKey error raised by the FabricProofVerifier contract.
type FabricProofVerifierInvalidPublicKey struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidPublicKey()
func FabricProofVerifierInvalidPublicKeyErrorID() common.Hash {
	return common.HexToHash("0xa2d0fee854985cb3d3a394146b5df1c41d5b4d6fcede7fc57a3c05a49b2f3e73")
}

// UnpackInvalidPublicKeyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidPublicKey()
func (fabricProofVerifier *FabricProofVerifier) UnpackInvalidPublicKeyError(raw []byte) (*FabricProofVerifierInvalidPublicKey, error) {
	out := new(FabricProofVerifierInvalidPublicKey)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "InvalidPublicKey", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierInvalidSignature represents a InvalidSignature error raised by the FabricProofVerifier contract.
type FabricProofVerifierInvalidSignature struct {
	MspIdHash [32]byte
	KeyIndex  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidSignature(bytes32 mspIdHash, uint32 keyIndex)
func FabricProofVerifierInvalidSignatureErrorID() common.Hash {
	return common.HexToHash("0x3b8b21a221f897d38931f0080c1e20d0316a9acbb45dc9898e085748f3b4e9c5")
}

// UnpackInvalidSignatureError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidSignature(bytes32 mspIdHash, uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) UnpackInvalidSignatureError(raw []byte) (*FabricProofVerifierInvalidSignature, error) {
	out := new(FabricProofVerifierInvalidSignature)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "InvalidSignature", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierInvalidThreshold represents a InvalidThreshold error raised by the FabricProofVerifier contract.
type FabricProofVerifierInvalidThreshold struct {
	Threshold  uint32
	ActiveOrgs uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidThreshold(uint32 threshold, uint32 activeOrgs)
func FabricProofVerifierInvalidThresholdErrorID() common.Hash {
	return common.HexToHash("0x07fb268af061aaa27b25fffaa659e706ad1ff2f0607fba0d0a8a0e9d7930ce98")
}

// UnpackInvalidThresholdError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidThreshold(uint32 threshold, uint32 activeOrgs)
func (fabricProofVerifier *FabricProofVerifier) UnpackInvalidThresholdError(raw []byte) (*FabricProofVerifierInvalidThreshold, error) {
	out := new(FabricProofVerifierInvalidThreshold)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "InvalidThreshold", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierOrgAlreadyRegistered represents a OrgAlreadyRegistered error raised by the FabricProofVerifier contract.
type FabricProofVerifierOrgAlreadyRegistered struct {
	MspIdHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error OrgAlreadyRegistered(bytes32 mspIdHash)
func FabricProofVerifierOrgAlreadyRegisteredErrorID() common.Hash {
	return common.HexToHash("0x16ab79843b69b67624976d2982db790296b482239d36085376a235236c202d2b")
}

// UnpackOrgAlreadyRegisteredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error OrgAlreadyRegistered(bytes32 mspIdHash)
func (fabricProofVerifier *FabricProofVerifier) UnpackOrgAlreadyRegisteredError(raw []byte) (*FabricProofVerifierOrgAlreadyRegistered, error) {
	out := new(FabricProofVerifierOrgAlreadyRegistered)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "OrgAlreadyRegistered", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierSignaturesNotSorted represents a SignaturesNotSorted error raised by the FabricProofVerifier contract.
type FabricProofVerifierSignaturesNotSorted struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SignaturesNotSorted()
func FabricProofVerifierSignaturesNotSortedErrorID() common.Hash {
	return common.HexToHash("0x01eba551ab5ccd3c63b3ca1f0cfd2b3551d237fc0b5bfaf4ddd918473b76ef10")
}

// UnpackSignaturesNotSortedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SignaturesNotSorted()
func (fabricProofVerifier *FabricProofVerifier) UnpackSignaturesNotSortedError(raw []byte) (*FabricProofVerifierSignaturesNotSorted, error) {
	out := new(FabricProofVerifierSignaturesNotSorted)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "SignaturesNotSorted", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierThresholdNotMet represents a ThresholdNotMet error raised by the FabricProofVerifier contract.
type FabricProofVerifierThresholdNotMet struct {
	Valid     uint32
	Threshold uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ThresholdNotMet(uint32 valid, uint32 threshold)
func FabricProofVerifierThresholdNotMetErrorID() common.Hash {
	return common.HexToHash("0x0a5563a16a78ae3409b36f4d673672e3ed4a3fef7d09e94461986dba50f52db9")
}

// UnpackThresholdNotMetError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ThresholdNotMet(uint32 valid, uint32 threshold)
func (fabricProofVerifier *FabricProofVerifier) UnpackThresholdNotMetError(raw []byte) (*FabricProofVerifierThresholdNotMet, error) {
	out := new(FabricProofVerifierThresholdNotMet)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "ThresholdNotMet", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierUnknownKey represents a UnknownKey error raised by the FabricProofVerifier contract.
type FabricProofVerifierUnknownKey struct {
	MspIdHash [32]byte
	KeyIndex  uint32
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownKey(bytes32 mspIdHash, uint32 keyIndex)
func FabricProofVerifierUnknownKeyErrorID() common.Hash {
	return common.HexToHash("0x7580351102732209f3817ca01b56459063c7b6b7d61901c63c97af602608c7e0")
}

// UnpackUnknownKeyError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownKey(bytes32 mspIdHash, uint32 keyIndex)
func (fabricProofVerifier *FabricProofVerifier) UnpackUnknownKeyError(raw []byte) (*FabricProofVerifierUnknownKey, error) {
	out := new(FabricProofVerifierUnknownKey)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "UnknownKey", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// FabricProofVerifierUnknownOrg represents a UnknownOrg error raised by the FabricProofVerifier contract.
type FabricProofVerifierUnknownOrg struct {
	MspIdHash [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnknownOrg(bytes32 mspIdHash)
func FabricProofVerifierUnknownOrgErrorID() common.Hash {
	return common.HexToHash("0x88f42f0f5cf9b7b7651e707375f16bb4b277f40afc305b01b2e9f7791f68dd7b")
}

// UnpackUnknownOrgError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnknownOrg(bytes32 mspIdHash)
func (fabricProofVerifier *FabricProofVerifier) UnpackUnknownOrgError(raw []byte) (*FabricProofVerifierUnknownOrg, error) {
	out := new(FabricProofVerifierUnknownOrg)
	if err := fabricProofVerifier.abi.UnpackIntoInterface(out, "UnknownOrg", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentMetaData contains all meta data concerning the HTLCPayment contract.
var HTLCPaymentMetaData = bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"token_\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"preimage\",\"type\":\"bytes32\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"getLock\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"enumHTLCPayment.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"internalType\":\"structHTLCPayment.Lock\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"}],\"name\":\"hashlockUsed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"revealedPreimage\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalLocked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"preimage\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"}],\"name\":\"PaymentLocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"PaymentRefunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"Expired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashlock\",\"type\":\"bytes32\"}],\"name\":\"HashlockReused\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidParams\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"NotExpired\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"},{\"internalType\":\"enumHTLCPayment.State\",\"name\":\"state\",\"type\":\"uint8\"}],\"name\":\"NotLocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"SwapExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"timelock\",\"type\":\"uint64\"},{\"internalType\":\"uint256\",\"name\":\"nowTs\",\"type\":\"uint256\"}],\"name\":\"TimelockInPast\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"expected\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"received\",\"type\":\"uint256\"}],\"name\":\"UnexpectedBalanceDelta\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"swapId\",\"type\":\"bytes32\"}],\"name\":\"WrongPreimage\",\"type\":\"error\"}]",
	ID:  "74bc3769e9ca42993ab454fea0d19287ec",
	Bin: "0x60a060405234801561000f575f5ffd5b5060405161130c38038061130c83398101604081905261002e9161008a565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f00556001600160a01b03811661007957604051635435b28960e11b815260040160405180910390fd5b6001600160a01b03166080526100b7565b5f6020828403121561009a575f5ffd5b81516001600160a01b03811681146100b0575f5ffd5b9392505050565b60805161121a6100f25f395f818161012401528181610548015281816105e201528181610651015281816108fd0152610bdb015261121a5ff3fe608060405234801561000f575f5ffd5b5060043610610085575f3560e01c806384cc9dfb1161005857806384cc9dfb146100ec578063d6f27b58146100ff578063fc0c546a1461011f578063ffdfa6de1461016b575f5ffd5b80631335f9bb1461008957806356891412146100bb5780635a1d305c146100c45780637249fbb6146100d9575b5f5ffd5b6100a8610097366004610fd7565b60026020525f908152604090205481565b6040519081526020015b60405180910390f35b6100a860035481565b6100d76100d2366004610fee565b61019d565b005b6100d76100e7366004610fd7565b6107a1565b6100d76100fa366004611064565b6109a4565b61011261010d366004610fd7565b610c8c565b6040516100b291906110ea565b6101467f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100b2565b61018d610179366004610fd7565b60016020525f908152604090205460ff1681565b60405190151581526020016100b2565b6101a5610d6f565b8415806101b0575081155b806101cf575073ffffffffffffffffffffffffffffffffffffffff8416155b806101d8575082155b1561020f576040517fa86b651200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff85160361025e576040517fa86b651200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f5f8681526020819052604090206004015468010000000000000000900460ff16600381111561029057610290611084565b146102cf576040517fafa955b9000000000000000000000000000000000000000000000000000000008152600481018690526024015b60405180910390fd5b5f8281526001602052604090205460ff161561031a576040517ff41a15ea000000000000000000000000000000000000000000000000000000008152600481018390526024016102c6565b428167ffffffffffffffff161161036f576040517f874ac8d000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff821660048201524260248201526044016102c6565b6040805160c08101825233815273ffffffffffffffffffffffffffffffffffffffff861660208201529081018490526060810183905267ffffffffffffffff8216608082015260a08101600190525f8681526020818152604091829020835181547fffffffffffffffffffffffff000000000000000000000000000000000000000090811673ffffffffffffffffffffffffffffffffffffffff92831617835592850151600183018054909416911617909155908201516002820155606082015160038083019190915560808301516004830180547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000811667ffffffffffffffff909316928317825560a08601519391927fffffffffffffffffffffffffffffffffffffffffffffff0000000000000000009091161790680100000000000000009084908111156104c2576104c2611084565b021790555050505f828152600160208190526040822080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016909117905560038054859290610513908490611199565b90915550506040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906370a0823190602401602060405180830381865afa1580156105a2573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105c691906111b2565b905061060a73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016333087610d9d565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f90829073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015610696573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906106ba91906111b2565b6106c491906111c9565b9050848114610709576040517f55e5b61d00000000000000000000000000000000000000000000000000000000815260048101869052602481018290526044016102c6565b604080518681526020810186905267ffffffffffffffff851681830152905173ffffffffffffffffffffffffffffffffffffffff88169133918a917f2b4be6c27d101b53534e41f243bea0a316551d2aca70bfb8d2c6d4a4ec3a7e13919081900360600190a4505061079a60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050505050565b6107a9610d6f565b5f8181526020819052604090206001600482015468010000000000000000900460ff1660038111156107dd576107dd611084565b1461082a576004818101546040517f172b3e0d0000000000000000000000000000000000000000000000000000000081526102c692859268010000000000000000900460ff1691016111dc565b600481015467ffffffffffffffff1642101561088d576004818101546040517f8788802000000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116918101919091524260248201526044016102c6565b6004810180547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680300000000000000001790556002810154600380548291905f906108db9084906111c9565b909155505081546109269073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911683610dff565b815460405182815273ffffffffffffffffffffffffffffffffffffffff9091169084907f3b28817a65441bed1ac3a5496c6b71346aa302e80c7418c12e1c49eaa17211769060200160405180910390a350506109a160017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b50565b6109ac610d6f565b5f8281526020819052604090206001600482015468010000000000000000900460ff1660038111156109e0576109e0611084565b14610a2d576004818101546040517f172b3e0d0000000000000000000000000000000000000000000000000000000081526102c692869268010000000000000000900460ff1691016111dc565b600481015467ffffffffffffffff164210610a8f576004818101546040517f91cc37bc00000000000000000000000000000000000000000000000000000000815267ffffffffffffffff909116918101919091524260248201526044016102c6565b8060030154600283604051602001610aa991815260200190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815290829052610ae1916111f7565b602060405180830381855afa158015610afc573d5f5f3e3d5ffd5b5050506040513d601f19601f82011682018060405250810190610b1f91906111b2565b14610b59576040517f1fa2ea6e000000000000000000000000000000000000000000000000000000008152600481018490526024016102c6565b6004810180547fffffffffffffffffffffffffffffffffffffffffffffff00ffffffffffffffff16680200000000000000001790555f83815260026020819052604082208490558201546003805491928392610bb69084906111c9565b90915550506001820154610c049073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000008116911683610dff565b6001820154604080518581526020810184905273ffffffffffffffffffffffffffffffffffffffff9092169186917f92c652097517b25aee57ddbd37e720bc85a39abd88c21e62855fc8c347700e66910160405180910390a35050610c8860017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050565b610cc16040805160c0810182525f8082526020820181905291810182905260608101829052608081018290529060a082015290565b5f8281526020818152604091829020825160c081018452815473ffffffffffffffffffffffffffffffffffffffff9081168252600183015416928101929092526002810154928201929092526003808301546060830152600483015467ffffffffffffffff8116608084015291929160a08401916801000000000000000090910460ff1690811115610d5557610d55611084565b6003811115610d6657610d66611084565b90525092915050565b610d77610e5f565b60027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b610dab848484846001610ebc565b610df9576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff851660048201526024016102c6565b50505050565b610e0c8383836001610f4f565b610e5a576040517f5274afe700000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff841660048201526024016102c6565b505050565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0054600203610eba576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6040517f23b872dd000000000000000000000000000000000000000000000000000000005f81815273ffffffffffffffffffffffffffffffffffffffff8781166004528616602452604485905291602083606481808c5af1925060015f51148316610f3e578383151615610f32573d5f823e3d81fd5b5f883b113d1516831692505b604052505f60605295945050505050565b6040517fa9059cbb000000000000000000000000000000000000000000000000000000005f81815273ffffffffffffffffffffffffffffffffffffffff8616600452602485905291602083604481808b5af1925060015f51148316610fcb578383151615610fbf573d5f823e3d81fd5b5f873b113d1516831692505b60405250949350505050565b5f60208284031215610fe7575f5ffd5b5035919050565b5f5f5f5f5f60a08688031215611002575f5ffd5b85359450602086013573ffffffffffffffffffffffffffffffffffffffff8116811461102c575f5ffd5b93506040860135925060608601359150608086013567ffffffffffffffff81168114611056575f5ffd5b809150509295509295909350565b5f5f60408385031215611075575f5ffd5b50508035926020909101359150565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b600481106110e6577f4e487b71000000000000000000000000000000000000000000000000000000005f52602160045260245ffd5b9052565b5f60c08201905073ffffffffffffffffffffffffffffffffffffffff835116825273ffffffffffffffffffffffffffffffffffffffff6020840151166020830152604083015160408301526060830151606083015267ffffffffffffffff608084015116608083015260a083015161116560a08401826110b1565b5092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b808201808211156111ac576111ac61116c565b92915050565b5f602082840312156111c2575f5ffd5b5051919050565b818103818111156111ac576111ac61116c565b828152604081016111f060208301846110b1565b9392505050565b5f82518060208501845e5f92019182525091905056fea164736f6c634300081e000a",
}

// HTLCPayment is an auto generated Go binding around an Ethereum contract.
type HTLCPayment struct {
	abi abi.ABI
}

// GetABI returns the ABI associated with this contract binding.
func (c *HTLCPayment) GetABI() abi.ABI {
	return c.abi
}

// NewHTLCPayment creates a new instance of HTLCPayment.
func NewHTLCPayment() *HTLCPayment {
	parsed, err := HTLCPaymentMetaData.ParseABI()
	if err != nil {
		panic(errors.New("invalid ABI: " + err.Error()))
	}
	return &HTLCPayment{abi: *parsed}
}

// Instance creates a wrapper for a deployed contract instance at the given address.
// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
func (c *HTLCPayment) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
	return bind.NewBoundContract(addr, c.abi, backend, backend, backend)
}

// PackConstructor is the Go binding used to pack the parameters required for
// contract deployment.
//
// Solidity: constructor(address token_) returns()
func (hTLCPayment *HTLCPayment) PackConstructor(token_ common.Address) []byte {
	enc, err := hTLCPayment.abi.Pack("", token_)
	if err != nil {
		panic(err)
	}
	return enc
}

// PackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84cc9dfb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function claim(bytes32 swapId, bytes32 preimage) returns()
func (hTLCPayment *HTLCPayment) PackClaim(swapId [32]byte, preimage [32]byte) []byte {
	enc, err := hTLCPayment.abi.Pack("claim", swapId, preimage)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackClaim is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x84cc9dfb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function claim(bytes32 swapId, bytes32 preimage) returns()
func (hTLCPayment *HTLCPayment) TryPackClaim(swapId [32]byte, preimage [32]byte) ([]byte, error) {
	return hTLCPayment.abi.Pack("claim", swapId, preimage)
}

// PackGetLock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6f27b58.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function getLock(bytes32 swapId) view returns((address,address,uint256,bytes32,uint64,uint8))
func (hTLCPayment *HTLCPayment) PackGetLock(swapId [32]byte) []byte {
	enc, err := hTLCPayment.abi.Pack("getLock", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackGetLock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xd6f27b58.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function getLock(bytes32 swapId) view returns((address,address,uint256,bytes32,uint64,uint8))
func (hTLCPayment *HTLCPayment) TryPackGetLock(swapId [32]byte) ([]byte, error) {
	return hTLCPayment.abi.Pack("getLock", swapId)
}

// UnpackGetLock is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xd6f27b58.
//
// Solidity: function getLock(bytes32 swapId) view returns((address,address,uint256,bytes32,uint64,uint8))
func (hTLCPayment *HTLCPayment) UnpackGetLock(data []byte) (HTLCPaymentLock, error) {
	out, err := hTLCPayment.abi.Unpack("getLock", data)
	if err != nil {
		return *new(HTLCPaymentLock), err
	}
	out0 := *abi.ConvertType(out[0], new(HTLCPaymentLock)).(*HTLCPaymentLock)
	return out0, nil
}

// PackHashlockUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffdfa6de.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (hTLCPayment *HTLCPayment) PackHashlockUsed(hashlock [32]byte) []byte {
	enc, err := hTLCPayment.abi.Pack("hashlockUsed", hashlock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackHashlockUsed is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xffdfa6de.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (hTLCPayment *HTLCPayment) TryPackHashlockUsed(hashlock [32]byte) ([]byte, error) {
	return hTLCPayment.abi.Pack("hashlockUsed", hashlock)
}

// UnpackHashlockUsed is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xffdfa6de.
//
// Solidity: function hashlockUsed(bytes32 hashlock) view returns(bool)
func (hTLCPayment *HTLCPayment) UnpackHashlockUsed(data []byte) (bool, error) {
	out, err := hTLCPayment.abi.Unpack("hashlockUsed", data)
	if err != nil {
		return *new(bool), err
	}
	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	return out0, nil
}

// PackLock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a1d305c.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function lock(bytes32 swapId, address recipient, uint256 amount, bytes32 hashlock, uint64 timelock) returns()
func (hTLCPayment *HTLCPayment) PackLock(swapId [32]byte, recipient common.Address, amount *big.Int, hashlock [32]byte, timelock uint64) []byte {
	enc, err := hTLCPayment.abi.Pack("lock", swapId, recipient, amount, hashlock, timelock)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackLock is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x5a1d305c.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function lock(bytes32 swapId, address recipient, uint256 amount, bytes32 hashlock, uint64 timelock) returns()
func (hTLCPayment *HTLCPayment) TryPackLock(swapId [32]byte, recipient common.Address, amount *big.Int, hashlock [32]byte, timelock uint64) ([]byte, error) {
	return hTLCPayment.abi.Pack("lock", swapId, recipient, amount, hashlock, timelock)
}

// PackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7249fbb6.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function refund(bytes32 swapId) returns()
func (hTLCPayment *HTLCPayment) PackRefund(swapId [32]byte) []byte {
	enc, err := hTLCPayment.abi.Pack("refund", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRefund is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x7249fbb6.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function refund(bytes32 swapId) returns()
func (hTLCPayment *HTLCPayment) TryPackRefund(swapId [32]byte) ([]byte, error) {
	return hTLCPayment.abi.Pack("refund", swapId)
}

// PackRevealedPreimage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1335f9bb.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (hTLCPayment *HTLCPayment) PackRevealedPreimage(swapId [32]byte) []byte {
	enc, err := hTLCPayment.abi.Pack("revealedPreimage", swapId)
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackRevealedPreimage is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x1335f9bb.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (hTLCPayment *HTLCPayment) TryPackRevealedPreimage(swapId [32]byte) ([]byte, error) {
	return hTLCPayment.abi.Pack("revealedPreimage", swapId)
}

// UnpackRevealedPreimage is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x1335f9bb.
//
// Solidity: function revealedPreimage(bytes32 swapId) view returns(bytes32)
func (hTLCPayment *HTLCPayment) UnpackRevealedPreimage(data []byte) ([32]byte, error) {
	out, err := hTLCPayment.abi.Unpack("revealedPreimage", data)
	if err != nil {
		return *new([32]byte), err
	}
	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	return out0, nil
}

// PackToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc0c546a.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function token() view returns(address)
func (hTLCPayment *HTLCPayment) PackToken() []byte {
	enc, err := hTLCPayment.abi.Pack("token")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackToken is the Go binding used to pack the parameters required for calling
// the contract method with ID 0xfc0c546a.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function token() view returns(address)
func (hTLCPayment *HTLCPayment) TryPackToken() ([]byte, error) {
	return hTLCPayment.abi.Pack("token")
}

// UnpackToken is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (hTLCPayment *HTLCPayment) UnpackToken(data []byte) (common.Address, error) {
	out, err := hTLCPayment.abi.Unpack("token", data)
	if err != nil {
		return *new(common.Address), err
	}
	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	return out0, nil
}

// PackTotalLocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56891412.  This method will panic if any
// invalid/nil inputs are passed.
//
// Solidity: function totalLocked() view returns(uint256)
func (hTLCPayment *HTLCPayment) PackTotalLocked() []byte {
	enc, err := hTLCPayment.abi.Pack("totalLocked")
	if err != nil {
		panic(err)
	}
	return enc
}

// TryPackTotalLocked is the Go binding used to pack the parameters required for calling
// the contract method with ID 0x56891412.  This method will return an error
// if any inputs are invalid/nil.
//
// Solidity: function totalLocked() view returns(uint256)
func (hTLCPayment *HTLCPayment) TryPackTotalLocked() ([]byte, error) {
	return hTLCPayment.abi.Pack("totalLocked")
}

// UnpackTotalLocked is the Go binding that unpacks the parameters returned
// from invoking the contract method with ID 0x56891412.
//
// Solidity: function totalLocked() view returns(uint256)
func (hTLCPayment *HTLCPayment) UnpackTotalLocked(data []byte) (*big.Int, error) {
	out, err := hTLCPayment.abi.Unpack("totalLocked", data)
	if err != nil {
		return new(big.Int), err
	}
	out0 := abi.ConvertType(out[0], new(big.Int)).(*big.Int)
	return out0, nil
}

// HTLCPaymentPaymentClaimed represents a PaymentClaimed event raised by the HTLCPayment contract.
type HTLCPaymentPaymentClaimed struct {
	SwapId    [32]byte
	Preimage  [32]byte
	Recipient common.Address
	Amount    *big.Int
	Raw       *types.Log // Blockchain specific contextual infos
}

const HTLCPaymentPaymentClaimedEventName = "PaymentClaimed"

// ContractEventName returns the user-defined event name.
func (HTLCPaymentPaymentClaimed) ContractEventName() string {
	return HTLCPaymentPaymentClaimedEventName
}

// UnpackPaymentClaimedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentClaimed(bytes32 indexed swapId, bytes32 preimage, address indexed recipient, uint256 amount)
func (hTLCPayment *HTLCPayment) UnpackPaymentClaimedEvent(log *types.Log) (*HTLCPaymentPaymentClaimed, error) {
	event := "PaymentClaimed"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != hTLCPayment.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(HTLCPaymentPaymentClaimed)
	if len(log.Data) > 0 {
		if err := hTLCPayment.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range hTLCPayment.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// HTLCPaymentPaymentLocked represents a PaymentLocked event raised by the HTLCPayment contract.
type HTLCPaymentPaymentLocked struct {
	SwapId    [32]byte
	Sender    common.Address
	Recipient common.Address
	Amount    *big.Int
	Hashlock  [32]byte
	Timelock  uint64
	Raw       *types.Log // Blockchain specific contextual infos
}

const HTLCPaymentPaymentLockedEventName = "PaymentLocked"

// ContractEventName returns the user-defined event name.
func (HTLCPaymentPaymentLocked) ContractEventName() string {
	return HTLCPaymentPaymentLockedEventName
}

// UnpackPaymentLockedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentLocked(bytes32 indexed swapId, address indexed sender, address indexed recipient, uint256 amount, bytes32 hashlock, uint64 timelock)
func (hTLCPayment *HTLCPayment) UnpackPaymentLockedEvent(log *types.Log) (*HTLCPaymentPaymentLocked, error) {
	event := "PaymentLocked"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != hTLCPayment.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(HTLCPaymentPaymentLocked)
	if len(log.Data) > 0 {
		if err := hTLCPayment.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range hTLCPayment.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// HTLCPaymentPaymentRefunded represents a PaymentRefunded event raised by the HTLCPayment contract.
type HTLCPaymentPaymentRefunded struct {
	SwapId [32]byte
	Sender common.Address
	Amount *big.Int
	Raw    *types.Log // Blockchain specific contextual infos
}

const HTLCPaymentPaymentRefundedEventName = "PaymentRefunded"

// ContractEventName returns the user-defined event name.
func (HTLCPaymentPaymentRefunded) ContractEventName() string {
	return HTLCPaymentPaymentRefundedEventName
}

// UnpackPaymentRefundedEvent is the Go binding that unpacks the event data emitted
// by contract.
//
// Solidity: event PaymentRefunded(bytes32 indexed swapId, address indexed sender, uint256 amount)
func (hTLCPayment *HTLCPayment) UnpackPaymentRefundedEvent(log *types.Log) (*HTLCPaymentPaymentRefunded, error) {
	event := "PaymentRefunded"
	if len(log.Topics) == 0 {
		return nil, bind.ErrNoEventSignature
	}
	if log.Topics[0] != hTLCPayment.abi.Events[event].ID {
		return nil, bind.ErrEventSignatureMismatch
	}
	out := new(HTLCPaymentPaymentRefunded)
	if len(log.Data) > 0 {
		if err := hTLCPayment.abi.UnpackIntoInterface(out, event, log.Data); err != nil {
			return nil, err
		}
	}
	var indexed abi.Arguments
	for _, arg := range hTLCPayment.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
		return nil, err
	}
	out.Raw = log
	return out, nil
}

// UnpackError attempts to decode the provided error data using user-defined
// error definitions.
func (hTLCPayment *HTLCPayment) UnpackError(raw []byte) (any, error) {
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["Expired"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["HashlockReused"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackHashlockReusedError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["InvalidParams"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackInvalidParamsError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["NotExpired"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackNotExpiredError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["NotLocked"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackNotLockedError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["ReentrancyGuardReentrantCall"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackReentrancyGuardReentrantCallError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["SafeERC20FailedOperation"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackSafeERC20FailedOperationError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["SwapExists"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackSwapExistsError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["TimelockInPast"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackTimelockInPastError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["UnexpectedBalanceDelta"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackUnexpectedBalanceDeltaError(raw[4:])
	}
	if bytes.Equal(raw[:4], hTLCPayment.abi.Errors["WrongPreimage"].ID.Bytes()[:4]) {
		return hTLCPayment.UnpackWrongPreimageError(raw[4:])
	}
	return nil, errors.New("Unknown error")
}

// HTLCPaymentExpired represents a Expired error raised by the HTLCPayment contract.
type HTLCPaymentExpired struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error Expired(uint64 timelock, uint256 nowTs)
func HTLCPaymentExpiredErrorID() common.Hash {
	return common.HexToHash("0x91cc37bcaeec164e5bc377adb7e09e2a00c91458581814d0294b11a1ab47958a")
}

// UnpackExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error Expired(uint64 timelock, uint256 nowTs)
func (hTLCPayment *HTLCPayment) UnpackExpiredError(raw []byte) (*HTLCPaymentExpired, error) {
	out := new(HTLCPaymentExpired)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "Expired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentHashlockReused represents a HashlockReused error raised by the HTLCPayment contract.
type HTLCPaymentHashlockReused struct {
	Hashlock [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error HashlockReused(bytes32 hashlock)
func HTLCPaymentHashlockReusedErrorID() common.Hash {
	return common.HexToHash("0xf41a15eade40d555d3009b515211d8eac2ddf768b0eaad4389bb9a2534e66a3c")
}

// UnpackHashlockReusedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error HashlockReused(bytes32 hashlock)
func (hTLCPayment *HTLCPayment) UnpackHashlockReusedError(raw []byte) (*HTLCPaymentHashlockReused, error) {
	out := new(HTLCPaymentHashlockReused)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "HashlockReused", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentInvalidParams represents a InvalidParams error raised by the HTLCPayment contract.
type HTLCPaymentInvalidParams struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error InvalidParams()
func HTLCPaymentInvalidParamsErrorID() common.Hash {
	return common.HexToHash("0xa86b65120b98ec2b50fcad62782dbb454a5c0604e0cf2db01fd94e5048d895b9")
}

// UnpackInvalidParamsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error InvalidParams()
func (hTLCPayment *HTLCPayment) UnpackInvalidParamsError(raw []byte) (*HTLCPaymentInvalidParams, error) {
	out := new(HTLCPaymentInvalidParams)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "InvalidParams", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentNotExpired represents a NotExpired error raised by the HTLCPayment contract.
type HTLCPaymentNotExpired struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotExpired(uint64 timelock, uint256 nowTs)
func HTLCPaymentNotExpiredErrorID() common.Hash {
	return common.HexToHash("0x878880200d6e2e3220bbbcf837eb40e26654b6a531b808079cf8848829ff512f")
}

// UnpackNotExpiredError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotExpired(uint64 timelock, uint256 nowTs)
func (hTLCPayment *HTLCPayment) UnpackNotExpiredError(raw []byte) (*HTLCPaymentNotExpired, error) {
	out := new(HTLCPaymentNotExpired)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "NotExpired", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentNotLocked represents a NotLocked error raised by the HTLCPayment contract.
type HTLCPaymentNotLocked struct {
	SwapId [32]byte
	State  uint8
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error NotLocked(bytes32 swapId, uint8 state)
func HTLCPaymentNotLockedErrorID() common.Hash {
	return common.HexToHash("0x172b3e0dd6c6c8c37927b521f9a4f246b89021494018735e50e066a362eb354a")
}

// UnpackNotLockedError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error NotLocked(bytes32 swapId, uint8 state)
func (hTLCPayment *HTLCPayment) UnpackNotLockedError(raw []byte) (*HTLCPaymentNotLocked, error) {
	out := new(HTLCPaymentNotLocked)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "NotLocked", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentReentrancyGuardReentrantCall represents a ReentrancyGuardReentrantCall error raised by the HTLCPayment contract.
type HTLCPaymentReentrancyGuardReentrantCall struct {
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error ReentrancyGuardReentrantCall()
func HTLCPaymentReentrancyGuardReentrantCallErrorID() common.Hash {
	return common.HexToHash("0x3ee5aeb571de7fc460830b4d0017439a1ca56fb0bc39062227ade4fe4a24c1ca")
}

// UnpackReentrancyGuardReentrantCallError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error ReentrancyGuardReentrantCall()
func (hTLCPayment *HTLCPayment) UnpackReentrancyGuardReentrantCallError(raw []byte) (*HTLCPaymentReentrancyGuardReentrantCall, error) {
	out := new(HTLCPaymentReentrancyGuardReentrantCall)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "ReentrancyGuardReentrantCall", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentSafeERC20FailedOperation represents a SafeERC20FailedOperation error raised by the HTLCPayment contract.
type HTLCPaymentSafeERC20FailedOperation struct {
	Token common.Address
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SafeERC20FailedOperation(address token)
func HTLCPaymentSafeERC20FailedOperationErrorID() common.Hash {
	return common.HexToHash("0x5274afe73c98b4749fc91ffae6b7b574e7842cb2144a159e9377a5f20b32edf9")
}

// UnpackSafeERC20FailedOperationError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SafeERC20FailedOperation(address token)
func (hTLCPayment *HTLCPayment) UnpackSafeERC20FailedOperationError(raw []byte) (*HTLCPaymentSafeERC20FailedOperation, error) {
	out := new(HTLCPaymentSafeERC20FailedOperation)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "SafeERC20FailedOperation", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentSwapExists represents a SwapExists error raised by the HTLCPayment contract.
type HTLCPaymentSwapExists struct {
	SwapId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error SwapExists(bytes32 swapId)
func HTLCPaymentSwapExistsErrorID() common.Hash {
	return common.HexToHash("0xafa955b93ee380941bdede9cdd884e4748f2da7f8e02336dda9c3bc5efeb9523")
}

// UnpackSwapExistsError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error SwapExists(bytes32 swapId)
func (hTLCPayment *HTLCPayment) UnpackSwapExistsError(raw []byte) (*HTLCPaymentSwapExists, error) {
	out := new(HTLCPaymentSwapExists)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "SwapExists", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentTimelockInPast represents a TimelockInPast error raised by the HTLCPayment contract.
type HTLCPaymentTimelockInPast struct {
	Timelock uint64
	NowTs    *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error TimelockInPast(uint64 timelock, uint256 nowTs)
func HTLCPaymentTimelockInPastErrorID() common.Hash {
	return common.HexToHash("0x874ac8d085b75da5651ef95027f3391b3208fd3cc8a677f09d139a5760db7bc6")
}

// UnpackTimelockInPastError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error TimelockInPast(uint64 timelock, uint256 nowTs)
func (hTLCPayment *HTLCPayment) UnpackTimelockInPastError(raw []byte) (*HTLCPaymentTimelockInPast, error) {
	out := new(HTLCPaymentTimelockInPast)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "TimelockInPast", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentUnexpectedBalanceDelta represents a UnexpectedBalanceDelta error raised by the HTLCPayment contract.
type HTLCPaymentUnexpectedBalanceDelta struct {
	Expected *big.Int
	Received *big.Int
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error UnexpectedBalanceDelta(uint256 expected, uint256 received)
func HTLCPaymentUnexpectedBalanceDeltaErrorID() common.Hash {
	return common.HexToHash("0x55e5b61d2c182dd91a8abf6fae0f353c2e31558869e50172aae83e146d5db184")
}

// UnpackUnexpectedBalanceDeltaError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error UnexpectedBalanceDelta(uint256 expected, uint256 received)
func (hTLCPayment *HTLCPayment) UnpackUnexpectedBalanceDeltaError(raw []byte) (*HTLCPaymentUnexpectedBalanceDelta, error) {
	out := new(HTLCPaymentUnexpectedBalanceDelta)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "UnexpectedBalanceDelta", raw); err != nil {
		return nil, err
	}
	return out, nil
}

// HTLCPaymentWrongPreimage represents a WrongPreimage error raised by the HTLCPayment contract.
type HTLCPaymentWrongPreimage struct {
	SwapId [32]byte
}

// ErrorID returns the hash of canonical representation of the error's signature.
//
// Solidity: error WrongPreimage(bytes32 swapId)
func HTLCPaymentWrongPreimageErrorID() common.Hash {
	return common.HexToHash("0x1fa2ea6e5b8169ee4d614d5d9fbb8e5441ae35f50318e1f1a35549322cb46d98")
}

// UnpackWrongPreimageError is the Go binding used to decode the provided
// error data into the corresponding Go error struct.
//
// Solidity: error WrongPreimage(bytes32 swapId)
func (hTLCPayment *HTLCPayment) UnpackWrongPreimageError(raw []byte) (*HTLCPaymentWrongPreimage, error) {
	out := new(HTLCPaymentWrongPreimage)
	if err := hTLCPayment.abi.UnpackIntoInterface(out, "WrongPreimage", raw); err != nil {
		return nil, err
	}
	return out, nil
}
