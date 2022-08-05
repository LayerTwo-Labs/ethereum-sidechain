package drivechain

/*
#cgo LDFLAGS: ./drivechain/target/debug/libdrivechain_eth.a -ldl -lm
#include "./bindings.h"
*/
import "C"
import (
	"github.com/ethereum/go-ethereum/common"
	"unsafe"
)

func AttemptBmm(criticalHash string, blockData string, amount uint64) {
	C.attempt_bmm(C.CString(criticalHash), C.CString(blockData), C.ulong(amount))
}

type Block struct {
	Data          string
	Time          int64
	MainBlockHash string
}

func ConfirmBmm() *Block {
	var cBlock = C.confirm_bmm()
	if cBlock == nil {
		return nil
	}
	var block = Block{
		Data:          C.GoString(cBlock.data),
		Time:          int64(cBlock.time),
		MainBlockHash: C.GoString(cBlock.main_block_hash),
	}
	C.free(unsafe.Pointer(cBlock.data))
	C.free(unsafe.Pointer(cBlock.main_block_hash))
	C.free(unsafe.Pointer(cBlock))
	return &block
}

func GetPrevMainBlockHash(mainBlockHash string) string {
	var cPrevMainBlockHash = C.get_prev_main_block_hash(C.CString(mainBlockHash))
	var prevMainBlockHash = C.GoString(cPrevMainBlockHash)
	C.free(unsafe.Pointer(cPrevMainBlockHash))
	return prevMainBlockHash
}

func VerifyBmm(mainBlockHash common.Hash, criticalHash common.Hash) bool {
	return bool(C.verify_bmm(C.CString(mainBlockHash.Hex()), C.CString(criticalHash.Hex())))
}
