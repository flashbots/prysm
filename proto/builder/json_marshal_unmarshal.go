package builder

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
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

type builderPayloadAttributesV2JSON struct {
	Timestamp   hexutil.Uint64   `json:"timestamp"`
	PrevRandao  hexutil.Bytes    `json:"prevRandao"`
	Slot        types.Slot       `json:"slot"`
	BlockHash   *common.Hash     `json:"blockHash"`
	Withdrawals []*v1.Withdrawal `json:"withdrawals"`
}

// MarshalJSON --
func (p *BuilderPayloadAttributesV2) MarshalJSON() ([]byte, error) {
	bHash := common.BytesToHash(p.BlockHash)
	if p.Withdrawals == nil {
		p.Withdrawals = make([]*v1.Withdrawal, 0)
	}
	return json.Marshal(builderPayloadAttributesV2JSON{
		Timestamp:   hexutil.Uint64(p.Timestamp),
		PrevRandao:  p.PrevRandao,
		Slot:        p.Slot,
		BlockHash:   &bHash,
		Withdrawals: p.Withdrawals,
	})
}

// UnmarshalJSON --
func (p *BuilderPayloadAttributesV2) UnmarshalJSON(enc []byte) error {
	dec := builderPayloadAttributesV2JSON{}
	if err := json.Unmarshal(enc, &dec); err != nil {
		return err
	}
	*p = BuilderPayloadAttributesV2{}
	p.Timestamp = uint64(dec.Timestamp)
	p.PrevRandao = dec.PrevRandao
	p.Slot = dec.Slot
	p.BlockHash = dec.BlockHash.Bytes()
	if p.Withdrawals == nil {
		p.Withdrawals = make([]*v1.Withdrawal, 0)
	}
	p.Withdrawals = dec.Withdrawals

	return nil
}
