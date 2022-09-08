package drivechain

/*
#cgo LDFLAGS: ./drivechain/target/debug/libdrivechain_eth.a -ldl -lm
#include "./bindings.h"
*/
import "C"
import (
	"github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/log"
	"math/big"
	"strings"
	"unsafe"
)

const THIS_SIDECHAIN = 1

func Init(dbPath, rpcUser, rpcPassword string) {
	cDbPath := C.CString(dbPath)
	cRpcUser := C.CString(rpcUser)
	cRpcPassword := C.CString(rpcPassword)
	C.init(cDbPath, C.ulong(THIS_SIDECHAIN), cRpcUser, cRpcPassword)
	C.free(unsafe.Pointer(cDbPath))
	C.free(unsafe.Pointer(cRpcUser))
	C.free(unsafe.Pointer(cRpcPassword))
}

func GetPrevMainBlockHash(mainBlockHash string) string {
	var cMainBlockHash = C.CString(mainBlockHash)
	var cPrevMainBlockHash = C.get_prev_main_block_hash(cMainBlockHash)
	var prevMainBlockHash = C.GoString(cPrevMainBlockHash)
	C.free_string(cPrevMainBlockHash)
	C.free(unsafe.Pointer(cMainBlockHash))
	return prevMainBlockHash
}

type RawDeposit struct {
	address string
	amount  uint64
}

func getDepositOutputs() []RawDeposit {
	ptrDeposits := C.get_deposit_outputs()
	cDeposits := unsafe.Slice(ptrDeposits.ptr, ptrDeposits.len)
	deposits := make([]RawDeposit, 0, ptrDeposits.len)
	for _, cDeposit := range cDeposits {
		deposit := RawDeposit{
			address: C.GoString(cDeposit.address),
			amount:  uint64(cDeposit.amount),
		}
		deposits = append(deposits, deposit)
	}
	C.free_deposits(ptrDeposits)
	return deposits
}

type Deposit struct {
	Address common.Address
	Amount  *big.Int
}

func GetDepositOutputs() []Deposit {
	rawDeposits := getDepositOutputs()
	deposits := make([]Deposit, 0, len(rawDeposits))
	for _, rawDeposit := range rawDeposits {
		deposits = append(deposits, Deposit{
			Address: common.HexToAddress(rawDeposit.address),
			Amount:  big.NewInt(int64(rawDeposit.amount)),
		})
	}
	return deposits
}

func ConnectBlock(deposits []Deposit, just_checking bool) bool {
	depositsArray := C.malloc(C.size_t(len(deposits)) * C.size_t(unsafe.Sizeof(C.Deposit{})))
	a := (*[1<<30 - 1]C.Deposit)(depositsArray)
	for i, deposit := range deposits {
		cDeposit := C.Deposit{
			address: C.CString(strings.ToLower(deposit.Address.String())),
			amount:  C.ulong(deposit.Amount.Uint64()),
		}
		a[i] = cDeposit
	}
	cDeposits := C.Deposits{
		ptr: &a[0],
		len: C.ulong(len(deposits)),
	}
	return bool(C.connect_block(cDeposits, C.bool(just_checking)))
}

func FormatDepositAddress(address string) string {
	cAddress := C.CString(address)
	cDepositAddress := C.format_deposit_address(cAddress)
	depositAddress := C.GoString(cDepositAddress)
	C.free(unsafe.Pointer(cAddress))
	C.free_string(cDepositAddress)
	return depositAddress
}

func attemptBmm(criticalHash string, amount uint64) {
	cCriticalHash := C.CString(criticalHash)
	C.attempt_bmm(cCriticalHash, C.ulong(amount))
	C.free(unsafe.Pointer(cCriticalHash))
}

func AttemptBmm(criticalHash common.Hash, amount uint64) {
	attemptBmm(criticalHash.Hex()[2:], amount)
}

type BmmState uint

const (
	Succeded BmmState = iota
	Failed
	Pending
)

func ConfirmBmm() BmmState {
	return BmmState(C.confirm_bmm())
}

func verifyBmm(mainBlockHash string, criticalHash string) bool {
	cMainBlockHash := C.CString(mainBlockHash)
	cCriticalHash := C.CString(criticalHash)
	result := bool(C.verify_bmm(cMainBlockHash, cCriticalHash))
	C.free(unsafe.Pointer(cMainBlockHash))
	C.free(unsafe.Pointer(cCriticalHash))
	return result
}

func VerifyBmm(mainBlockHash common.Hash, criticalHash common.Hash) bool {
	return verifyBmm(mainBlockHash.Hex()[2:], criticalHash.Hex()[2:])
}


// NOTE: Treasure account idea makes a lot of sense for ethereum!
