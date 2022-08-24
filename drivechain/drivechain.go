package drivechain

/*
#cgo LDFLAGS: ./drivechain/target/debug/libdrivechain_eth.a -ldl -lm
#include "./bindings.h"
*/
import "C"
import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"unsafe"
)

const THIS_SIDECHAIN = 1

// db_path: *const libc::c_char,
// this_sidechain: usize,
// rpcuser: *const libc::c_char,
// rpcpassword: *const libc::c_char,
func Init(dbPath, rpcUser, rpcPassword string) {
	log.Info("drivechain initialized")
	C.init(C.CString(dbPath), C.ulong(THIS_SIDECHAIN), C.CString(rpcUser), C.CString(rpcPassword))
}

func AttemptBmm(criticalHash common.Hash, amount uint64) {
	C.attempt_bmm(C.CString(criticalHash.Hex()[2:]), C.ulong(amount))
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

func GetPrevMainBlockHash(mainBlockHash string) string {
	var cPrevMainBlockHash = C.get_prev_main_block_hash(C.CString(mainBlockHash))
	var prevMainBlockHash = C.GoString(cPrevMainBlockHash)
	C.free(unsafe.Pointer(cPrevMainBlockHash))
	return prevMainBlockHash
}

func VerifyBmm(mainBlockHash common.Hash, criticalHash common.Hash) bool {
	return bool(C.verify_bmm(C.CString(mainBlockHash.Hex()[2:]), C.CString(criticalHash.Hex()[2:])))
}
