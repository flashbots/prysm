package builder_test

import (
	"encoding/json"
	"testing"

	fieldparams "github.com/prysmaticlabs/prysm/config/fieldparams"
	types "github.com/prysmaticlabs/prysm/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/encoding/bytesutil"
	builder "github.com/prysmaticlabs/prysm/proto/builder"
	"github.com/prysmaticlabs/prysm/testing/require"
)

func TestJsonMarshalUnmarshal(t *testing.T) {
	t.Run("builder payload attributes", func(t *testing.T) {
		random := bytesutil.PadTo([]byte("random"), fieldparams.RootLength)
		blockHash := bytesutil.PadTo([]byte("blockHash"), fieldparams.RootLength)
		jsonPayload := &builder.BuilderPayloadAttributes{
			Timestamp:  1,
			PrevRandao: random,
			Slot: 1,
			BlockHash: blockHash,
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
