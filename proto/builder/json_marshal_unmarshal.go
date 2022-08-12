package builder

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
)

type builderPayloadAttributesJSON struct {
	Timestamp  hexutil.Uint64 `json:"timestamp"`
	PrevRandao hexutil.Bytes  `json:"prevRandao"`
	Slot       types.Slot     `json:"slot"`
	BlockHash  *common.Hash   `json:"blockHash"`
}

// MarshalJSON --
func (p *BuilderPayloadAttributes) MarshalJSON() ([]byte, error) {
	bHash := common.BytesToHash(p.BlockHash)
	return json.Marshal(builderPayloadAttributesJSON{
		Timestamp:  hexutil.Uint64(p.Timestamp),
		PrevRandao: p.PrevRandao,
		Slot:       p.Slot,
		BlockHash:  &bHash,
	})
}

// UnmarshalJSON --
func (p *BuilderPayloadAttributes) UnmarshalJSON(enc []byte) error {
	dec := builderPayloadAttributesJSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*p = BuilderPayloadAttributes{}
	p.Timestamp = uint64(dec.Timestamp)
	p.PrevRandao = dec.PrevRandao
	p.Slot = dec.Slot
	p.BlockHash = dec.BlockHash.Bytes()
	return nil
}
