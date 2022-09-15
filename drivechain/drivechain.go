package drivechain

/*
#cgo LDFLAGS: ./drivechain/target/debug/libdrivechain_eth.a -ldl -lm
#include "./bindings.h"
*/
import "C"
import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"
	"unsafe"
)

const THIS_SIDECHAIN = 1

// A publicly known "private key" to the treasury account, that holds 21M BTC.
// There are special consensus rules for this account.
//
// The only transfers from this account that are valid correspond to deposits on
// mainchain or to refunds of earlier withdrawal.
//
// Withdrawals are transfers to this account with special data.
//
// Transfering funds to this account without the special withdrawal data will
// burn the coins. They will never show up on mainchain and there will be no way
// to refund them.
const TREASURY_PRIVATE_KEY = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
const TREASURY_ACCOUNT = "0xc96aaa54e2d44c299564da76e1cd3184a2386b8d"

// There are 10,000,000,000 Wei in one Satoshi
var Satoshi = big.NewInt(10_000_000_000)

// There are 10^8 satoshi in one BTC
// There are 10^18 Wei in one Ether.
//
// So let 1 BTC = 1 "Ether" and 1 satoshi = 10^10 Wei.
//
// So there should be 21 * 10 ^ 6 * 10 ^ 18 = 21 * 10^24 "Wei" in the treasury account.

func Init(dbPath, rpcUser, rpcPassword string) {
	privKey, err := crypto.HexToECDSA(TREASURY_PRIVATE_KEY)
	if err != nil {
		panic(fmt.Sprintf("can't get treasury private key: %s", err))
	}
	address := crypto.PubkeyToAddress(*privKey.Public().(*ecdsa.PublicKey))
	actualTreasuryAccount := strings.ToLower(address.Hex())
	if TREASURY_ACCOUNT != actualTreasuryAccount {
		panic(fmt.Sprintf("treasury account: %s != actual treasury account: %s", TREASURY_ACCOUNT))
	}
	cDbPath := C.CString(dbPath)
	cRpcUser := C.CString(rpcUser)
	cRpcPassword := C.CString(rpcPassword)
	C.init(cDbPath, C.ulong(THIS_SIDECHAIN), cRpcUser, cRpcPassword)
	C.free(unsafe.Pointer(cDbPath))
	C.free(unsafe.Pointer(cRpcUser))
	C.free(unsafe.Pointer(cRpcPassword))
}

func GetMainchainTip() common.Hash {
	var cMainchainTip = C.get_mainchain_tip()
	var mainchainTip = C.GoString(cMainchainTip)
	C.free_string(cMainchainTip)
	return common.HexToHash(mainchainTip)
}

func GetPrevMainBlockHash(mainBlockHash common.Hash) common.Hash {
	var cMainBlockHash = C.CString(mainBlockHash.Hex()[2:])
	var cPrevMainBlockHash = C.get_prev_main_block_hash(cMainBlockHash)
	var prevMainBlockHash = C.GoString(cPrevMainBlockHash)
	C.free_string(cPrevMainBlockHash)
	C.free(unsafe.Pointer(cMainBlockHash))
	return common.HexToHash(prevMainBlockHash)
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

type Withdrawal struct {
	Address [20]C.uchar
	Amount  *big.Int
	Fee     *big.Int
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

// common.Hash here is for transaction hashes.
func ConnectBlock(deposits []Deposit, withdrawals map[common.Hash]Withdrawal, refunds []common.Hash, just_checking bool) bool {
	depositsMemory := C.malloc(C.size_t(len(deposits)) * C.size_t(unsafe.Sizeof(C.Deposit{})))
	depositsSlice := (*[1<<30 - 1]C.Deposit)(depositsMemory)
	for i, deposit := range deposits {
		cDeposit := C.Deposit{
			address: C.CString(strings.ToLower(deposit.Address.String())),
			amount:  C.ulong(deposit.Amount.Uint64()),
		}
		depositsSlice[i] = cDeposit
	}
	cDeposits := C.Deposits{
		ptr: &depositsSlice[0],
		len: C.ulong(len(deposits)),
	}
	withdrawalsMemory := C.malloc(C.size_t(len(withdrawals)) * C.size_t(unsafe.Sizeof(C.Withdrawal{})))
	withdrawalsSlice := (*[1<<30 - 1]C.Withdrawal)(withdrawalsMemory)
	{
		i := 0
		for id, w := range withdrawals {
			cWithdrawal := C.Withdrawal{
				id:      C.CString(id.Hex()),
				address: w.Address,
				amount:  C.ulong(w.Amount.Uint64()),
				fee:     C.ulong(w.Fee.Uint64()),
			}
			withdrawalsSlice[i] = cWithdrawal
			i += 1
		}
	}
	cWithdrawals := C.Withdrawals{
		ptr: &withdrawalsSlice[0],
		len: C.ulong(len(withdrawals)),
	}
	// this is an array of C strings
	refundsMemory := C.malloc(C.size_t(len(withdrawals)) * C.size_t(unsafe.Sizeof(C.Refund{})))
	refundsSlice := (*[1<<30 - 1]C.Refund)(refundsMemory)
	for i, r := range refunds {
		cRefund := C.Refund{
			id: C.CString(r.Hex()),
		}
		refundsSlice[i] = cRefund
	}
	cRefunds := C.Refunds{
		ptr: &refundsSlice[0],
		len: C.ulong(len(refunds)),
	}
	return bool(C.connect_block(cDeposits, cWithdrawals, cRefunds, C.bool(just_checking)))
}

func FormatDepositAddress(address string) string {
	cAddress := C.CString(address)
	cDepositAddress := C.format_deposit_address(cAddress)
	depositAddress := C.GoString(cDepositAddress)
	C.free(unsafe.Pointer(cAddress))
	C.free_string(cDepositAddress)
	return depositAddress
}

func CreateDeposit(address common.Address, amount uint64, fee uint64) bool {
	// var weiInSatoshi big.Int
	// weiInSatoshi.Exp(big.NewInt(10), big.NewInt(10), nil) // There are 10^10 Wei in 1 Satoshi.
	// var satAmount big.Int
	// var satFee big.Int
	// satAmount.Div(amount, &weiInSatoshi)
	// satFee.Div(fee, &weiInSatoshi)
	cAddress := C.CString(strings.ToLower(address.Hex()))

	cAmount := C.ulong(amount)
	cFee := C.ulong(fee)
	result := C.create_deposit(cAddress, cAmount, cFee)
	C.free(unsafe.Pointer(cAddress))
	return bool(result)
}

func attemptBmm(criticalHash string, prevMainBlockHash string, amount uint64) {
	cCriticalHash := C.CString(criticalHash)
	cPrevMainBlockHash := C.CString(prevMainBlockHash)
	C.attempt_bmm(cCriticalHash, cPrevMainBlockHash, C.ulong(amount))
	C.free(unsafe.Pointer(cCriticalHash))
	C.free(unsafe.Pointer(cPrevMainBlockHash))
}

func AttemptBmm(header *types.Header, amount uint64) {
	attemptBmm(header.Hash().Hex()[2:], header.PrevMainBlockHash.Hex()[2:], amount)
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
