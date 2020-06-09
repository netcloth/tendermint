package types

import (
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
)

// ABCIResults wraps the deliver tx results to return a proof
type ABCIResults []*abci.DeterministicResponseDeliverTx

// NewResults creates ABCIResults from the list of ResponseDeliverTx.
func NewResults(responses []*abci.ResponseDeliverTx) ABCIResults {
	res := make(ABCIResults, len(responses))
	for i, d := range responses {
		res[i] = NewResultFromResponse(d)
	}
	return res
}

// NewResultFromResponse creates ABCIResult from ResponseDeliverTx.
func NewResultFromResponse(response *abci.ResponseDeliverTx) *abci.DeterministicResponseDeliverTx {
	return &abci.DeterministicResponseDeliverTx{
		Code:      response.Code,
		Data:      response.Data,
		GasWanted: response.GasWanted,
		GasUsed:   response.GasUsed,
		Events:    response.Events,
	}
}

// Hash returns a merkle hash of all results
func (a ABCIResults) Hash() []byte {
	// NOTE: we copy the impl of the merkle tree for txs -
	// we should be consistent and either do it for both or not.
	return merkle.SimpleHashFromByteSlices(a.toByteSlices())
}

// ProveResult returns a merkle proof of one result from the set
func (a ABCIResults) ProveResult(i int) merkle.SimpleProof {
	_, proofs := merkle.SimpleProofsFromByteSlices(a.toByteSlices())
	return *proofs[i]
}

func (a ABCIResults) toByteSlices() [][]byte {
	l := len(a)
	bzs := make([][]byte, l)
	for i := 0; i < l; i++ {
		bz, err := proto.Marshal(a[i])
		if err != nil {
			panic(err)
		}
		bzs[i] = bz
	}
	return bzs
}
