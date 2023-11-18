package drivechain

/*
#cgo LDFLAGS: ./drivechain/target/debug/libdrivechain_eth.a -ldl -lWs2_32 -lbcrypt -lntdll -lkernel32 -luserenv -lm
#include "./bindings.h"
*/
import "C"

import (
	"strings"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
)

//-LC:/Users/torke/dev/dlfcn-win32

func newDeposits(deposits []Deposit) C.Deposits {
	depositsMemory := C.malloc(C.size_t(len(deposits)) * C.size_t(unsafe.Sizeof(C.Deposit{})))
	depositsSlice := (*[1<<30 - 1]C.Deposit)(depositsMemory)
	for i, deposit := range deposits {
		depositsSlice[i] = C.Deposit{
			address: C.CString(strings.ToLower(deposit.Amount.String())),
			amount:  C.ulonglong(deposit.Amount.Uint64()),
		}
	}
	return C.Deposits{
		ptr: &depositsSlice[0],
		len: C.ulonglong(len(deposits)),
	}
}

func newRefundsFromHash(refunds []common.Hash) C.Refunds {
	refundsMemory := C.malloc(C.size_t(len(refunds)) * C.size_t(unsafe.Sizeof(C.Refund{})))
	refundsSlice := (*[1<<30 - 1]C.Refund)(refundsMemory)
	for i, id := range refunds {
		cRefund := C.Refund{
			id: C.CString(id.Hex()),
		}
		refundsSlice[i] = cRefund
	}
	return C.Refunds{
		ptr: &refundsSlice[0],
		len: C.ulonglong(len(refunds)),
	}
}

func newRefunds(refunds []Refund) C.Refunds {
	refundsMemory := C.malloc(C.size_t(len(refunds)) * C.size_t(unsafe.Sizeof(C.Refund{})))
	refundsSlice := (*[1<<30 - 1]C.Refund)(refundsMemory)
	for i, r := range refunds {
		cRefund := C.Refund{
			id:     C.CString(r.Id.Hex()),
			amount: C.ulonglong(r.Amount.Uint64()),
		}
		refundsSlice[i] = cRefund
	}
	return C.Refunds{
		ptr: &refundsSlice[0],
		len: C.ulonglong(len(refunds)),
	}
}

func newWithdrawalsFromHash(withdrawals []common.Hash) C.Withdrawals {
	withdrawalsMemory := C.malloc(C.size_t(len(withdrawals)) * C.size_t(unsafe.Sizeof(C.Withdrawal{})))
	withdrawalsSlice := (*[1<<30 - 1]C.Withdrawal)(withdrawalsMemory)
	for i, id := range withdrawals {
		cWithdrawal := C.Withdrawal{
			id: C.CString(id.Hex()),
		}
		withdrawalsSlice[i] = cWithdrawal
	}

	return C.Withdrawals{
		ptr: &withdrawalsSlice[0],
		len: C.ulonglong(len(withdrawals)),
	}
}

func newWithdrawals(withdrawals map[common.Hash]Withdrawal) C.Withdrawals {
	withdrawalsMemory := C.malloc(C.size_t(len(withdrawals)) * C.size_t(unsafe.Sizeof(C.Withdrawal{})))
	withdrawalsSlice := (*[1<<30 - 1]C.Withdrawal)(withdrawalsMemory)
	{
		i := 0
		for id, w := range withdrawals {
			cWithdrawal := C.Withdrawal{
				id:      C.CString(id.Hex()),
				address: w.Address,
				amount:  C.ulonglong(w.Amount.Uint64()),
				fee:     C.ulonglong(w.Fee.Uint64()),
			}
			withdrawalsSlice[i] = cWithdrawal
			i += 1
		}
	}
	return C.Withdrawals{
		ptr: &withdrawalsSlice[0],
		len: C.ulonglong(len(withdrawals)),
	}
}

func createDeposit(address common.Address, amount uint64, fee uint64) bool {
	cAddress := C.CString(strings.ToLower(address.Hex()))
	cAmount := C.ulonglong(amount)
	cFee := C.ulonglong(fee)
	result := C.create_deposit(cAddress, cAmount, cFee)
	C.free(unsafe.Pointer(cAddress))
	return bool(result)
}

func attemptBmm(criticalHash string, prevMainBlockHash string, amount uint64) {
	cCriticalHash := C.CString(criticalHash)
	cPrevMainBlockHash := C.CString(prevMainBlockHash)
	C.attempt_bmm(cCriticalHash, cPrevMainBlockHash, C.ulonglong(amount))
	C.free(unsafe.Pointer(cCriticalHash))
	C.free(unsafe.Pointer(cPrevMainBlockHash))
}

func initBmmEngine(dbPath, host, rpcUser, rpcPassword string, port uint16) {
	cDbPath := C.CString(dbPath)
	cHost := C.CString(host)
	cRpcUser := C.CString(rpcUser)
	cRpcPassword := C.CString(rpcPassword)

	C.init(cDbPath, C.ulonglong(THIS_SIDECHAIN), cHost, C.ushort(port), cRpcUser, cRpcPassword)
	C.free(unsafe.Pointer(cDbPath))
	C.free(unsafe.Pointer(cRpcUser))
	C.free(unsafe.Pointer(cRpcPassword))
}
