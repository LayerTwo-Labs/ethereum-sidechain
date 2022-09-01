package bmm

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/drivechain"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"
	"math/big"
	"path/filepath"
	"time"
)

var (
	maxUncles = 0 // With blind merge mining and 10 minutes block time there are no uncle blocks.
)

// Bmm is a blind merge mining consensus engine.
type Bmm struct {
}

func New(dataDir string) Bmm {
	drivechain.Init(filepath.Join(dataDir, "drivechain"), "user", "password")
	return Bmm{}
}

func (bmm *Bmm) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

// FIXME: Figure out why VerifyHeader is never called in dev mode.
// FIXME: Add non PoW checks from ethash consensus engine.
func (bmm *Bmm) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	log.Info("verifying ", header.PrevMainBlockHash)
	if !drivechain.VerifyBmm(header.PrevMainBlockHash, header.Hash()) {
		return errors.New("invalid bmm")
	}
	return nil
}

func (bmm *Bmm) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	log.Info("verifying ", headers)
	abort, results := make(chan struct{}), make(chan error, len(headers))
	for i := 0; i < len(headers); i++ {
		err := bmm.VerifyHeader(chain, headers[i], seals[i])
		results <- err
	}
	return abort, results
}

func (bmm *Bmm) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return errors.New("uncle blocks are impossible with a blind merge mining consensus engine")
}

func (bmm *Bmm) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// NOTE: Probably PrevMainBlockHash should be set here.
	header.Difficulty = big.NewInt(1)
	return nil
}

func (bmm *Bmm) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	txs []*types.Transaction, uncles []*types.Header) {
	deposits := drivechain.GetDepositOutputs()
	for _, deposit := range deposits {
		state.AddBalance(deposit.Address, deposit.Amount)
	}
	// Accumulate any block and uncle rewards and commit the final state root
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
}

func (bmm *Bmm) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Finalize block
	bmm.Finalize(chain, header, state, txs, uncles)
	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil)), nil
}

func (bmm *Bmm) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// FIXME: Make it possible for the miner to change the amount.
	amount := uint64(10000)
	header := block.Header()
	drivechain.AttemptBmm(header.Hash(), amount)
	log.Info("attempting to bmm block")

	go func() {
		for true {
			// log.Info("checking if block was bmmed")
			state := drivechain.ConfirmBmm()
			if state == drivechain.Succeded {
				log.Info("block was bmmed")
				select {
				case <-stop:
					break
				case results <- block.WithSeal(header):
				default:
				}
				break
			} else if state == drivechain.Failed {
				log.Info("bmm commitment wasn't inclued in a main:block")
				log.Info("attempting new bmm request")
				drivechain.AttemptBmm(header.Hash(), amount)
			}
			time.Sleep(1 * time.Second)
		}
		log.Info("finished attempting to seal block")
	}()
	return nil
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
	return nil
}
