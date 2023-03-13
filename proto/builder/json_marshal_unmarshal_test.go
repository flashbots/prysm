package builder_test

import (
	"encoding/json"
	"testing"

	fieldparams "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	types "github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	builder "github.com/prysmaticlabs/prysm/v3/proto/builder"
	v1 "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v3/testing/require"
)

func TestJsonMarshalUnmarshal(t *testing.T) {
	t.Run("builder payload attributes", func(t *testing.T) {
		random := bytesutil.PadTo([]byte("random"), fieldparams.RootLength)
		blockHash := bytesutil.PadTo([]byte("blockHash"), fieldparams.RootLength)
		jsonPayload := &builder.BuilderPayloadAttributes{
			Timestamp:  1,
			PrevRandao: random,
			Slot:       1,
			BlockHash:  blockHash,
		}
		enc, err := json.Marshal(jsonPayload)
		require.NoError(t, err)
		payloadPb := &builder.BuilderPayloadAttributes{}
		require.NoError(t, json.Unmarshal(enc, payloadPb))
		require.DeepEqual(t, uint64(1), payloadPb.Timestamp)
		require.DeepEqual(t, random, payloadPb.PrevRandao)
		require.DeepEqual(t, types.Slot(1), payloadPb.Slot)
		require.DeepEqual(t, blockHash, payloadPb.BlockHash)
	})
}

func TestJsonMarshalUnmarshalV2(t *testing.T) {
	t.Run("builder payload attributes v2", func(t *testing.T) {
		random := bytesutil.PadTo([]byte("random"), fieldparams.RootLength)
		blockHash := bytesutil.PadTo([]byte("blockHash"), fieldparams.RootLength)
		jsonPayload := &builder.BuilderPayloadAttributesV2{
			Timestamp:  1,
			PrevRandao: random,
			Slot:       1,
			BlockHash:  blockHash,
			Withdrawals: []*v1.Withdrawal{{
				Index:          1,
				ValidatorIndex: 1,
				Address:        bytesutil.PadTo([]byte("address"), 20),
				Amount:         1,
			}},
		}
		enc, err := json.Marshal(jsonPayload)
		require.NoError(t, err)
		payloadPb := &builder.BuilderPayloadAttributesV2{}
		require.NoError(t, json.Unmarshal(enc, payloadPb))
		require.DeepEqual(t, uint64(1), payloadPb.Timestamp)
		require.DeepEqual(t, random, payloadPb.PrevRandao)
		require.DeepEqual(t, types.Slot(1), payloadPb.Slot)
		require.DeepEqual(t, blockHash, payloadPb.BlockHash)
		require.Equal(t, 1, len(payloadPb.Withdrawals))
		withdrawal := payloadPb.Withdrawals[0]
		require.Equal(t, uint64(1), withdrawal.Index)
		require.Equal(t, primitives.ValidatorIndex(1), withdrawal.ValidatorIndex)
		require.DeepEqual(t, bytesutil.PadTo([]byte("address"), 20), withdrawal.Address)
		require.Equal(t, uint64(1), withdrawal.Amount)
	})
}
