package bmm

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/drivechain"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
	"math/big"
)

var (
	maxUncles = 0 // With blind merge mining there are no uncle blocks.
)

// Bmm is a blind merge mining consensus engine.
type Bmm struct {
}

func (bmm *Bmm) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (bmm *Bmm) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	if !drivechain.VerifyBmm(header.MainBlockHash, header.Root) {
		return errors.New("invalid bmm")
	}
	return nil
}

func (bmm *Bmm) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort, results := make(chan struct{}), make(chan error, len(headers))
	return abort, results
}

func (bmm *Bmm) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return errors.New("uncle blocks are impossible with a blind merge mining consensus engine")
}

func (bmm *Bmm) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return errors.New("unimplemented")
}

func (bmm *Bmm) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	txs []*types.Transaction, uncles []*types.Header) {
}

func (bmm *Bmm) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return nil, errors.New("unimplemented")
}

func (bmm *Bmm) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return errors.New("unimplemented")
}

func (bmm *Bmm) SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	enc := []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra,
	}
	if header.BaseFee != nil {
		enc = append(enc, header.BaseFee)
	}
	rlp.Encode(hasher, enc)
	hasher.Sum(hash[:0])
	return hash
}

func (bmm *Bmm) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	// There is no concept of difficulty for blind merge mining.
	return nil
}

func (bmm *Bmm) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return nil
}

func (bmm *Bmm) Close() error {
	return errors.New("unimplemented")
}
