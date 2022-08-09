package execution

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethRPC "github.com/ethereum/go-ethereum/rpc"
	"github.com/holiman/uint256"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/execution/types"
	fieldparams "github.com/prysmaticlabs/prysm/v3/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v3/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v3/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	pb "github.com/prysmaticlabs/prysm/v3/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v3/runtime/version"
	"github.com/prysmaticlabs/prysm/v3/time/slots"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

const (
	// PayloadAttributes request string for JSON-RPC
	PayloadAttributesMethod = "builder_payloadAttributes"
	// NewPayloadMethod v1 request string for JSON-RPC.
	NewPayloadMethod = "engine_newPayloadV1"
	// NewPayloadMethodV2 v2 request string for JSON-RPC.
	NewPayloadMethodV2 = "engine_newPayloadV2"
	// ForkchoiceUpdatedMethod v1 request string for JSON-RPC.
	ForkchoiceUpdatedMethod = "engine_forkchoiceUpdatedV1"
	// ForkchoiceUpdatedMethodV2 v2 request string for JSON-RPC.
	ForkchoiceUpdatedMethodV2 = "engine_forkchoiceUpdatedV2"
	// GetPayloadMethod v1 request string for JSON-RPC.
	GetPayloadMethod = "engine_getPayloadV1"
	// GetPayloadMethodV2 v2 request string for JSON-RPC.
	GetPayloadMethodV2 = "engine_getPayloadV2"
	// ExchangeTransitionConfigurationMethod v1 request string for JSON-RPC.
	ExchangeTransitionConfigurationMethod = "engine_exchangeTransitionConfigurationV1"
	// ExecutionBlockByHashMethod request string for JSON-RPC.
	ExecutionBlockByHashMethod = "eth_getBlockByHash"
	// ExecutionBlockByNumberMethod request string for JSON-RPC.
	ExecutionBlockByNumberMethod = "eth_getBlockByNumber"
	// Defines the seconds before timing out engine endpoints with non-block execution semantics.
	defaultEngineTimeout = time.Second
)

// ForkchoiceUpdatedResponse is the response kind received by the
// engine_forkchoiceUpdatedV1 endpoint.
type ForkchoiceUpdatedResponse struct {
	Status    *pb.PayloadStatus  `json:"payloadStatus"`
	PayloadId *pb.PayloadIDBytes `json:"payloadId"`
}

// PayloadAttributesResponse is the response kind received by the
// builder_payloadAttributes endpoint.
type PayloadAttributesResponse struct {
	Status *pb.PayloadStatus `json:"payloadStatus"`
} 

// ExecutionPayloadReconstructor defines a service that can reconstruct a full beacon
// block with an execution payload from a signed beacon block and a connection
// to an execution client's engine API.
type ExecutionPayloadReconstructor interface {
	ReconstructFullBlock(
		ctx context.Context, blindedBlock interfaces.ReadOnlySignedBeaconBlock,
	) (interfaces.SignedBeaconBlock, error)
	ReconstructFullBellatrixBlockBatch(
		ctx context.Context, blindedBlocks []interfaces.ReadOnlySignedBeaconBlock,
	) ([]interfaces.SignedBeaconBlock, error)
}

// EngineCaller defines a client that can interact with an Ethereum
// execution node's engine service via JSON-RPC.
type EngineCaller interface {
	NewPayload(ctx context.Context, payload interfaces.ExecutionData) ([]byte, error)
	ForkchoiceUpdated(
		ctx context.Context, state *pb.ForkchoiceState, attrs payloadattribute.Attributer,
	) (*pb.PayloadIDBytes, []byte, error)
	GetPayload(ctx context.Context, payloadId [8]byte, slot primitives.Slot) (interfaces.ExecutionData, error)
	ExchangeTransitionConfiguration(
		ctx context.Context, cfg *pb.TransitionConfiguration,
	) error
	ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool) (*pb.ExecutionBlock, error)
	GetTerminalBlockHash(ctx context.Context) ([]byte, bool, error)
}

var EmptyBlockHash = errors.New("Block hash is empty 0x0000...")

// NewPayload calls the engine_newPayloadVX method via JSON-RPC.
func (s *Service) NewPayload(ctx context.Context, payload interfaces.ExecutionData) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.NewPayload")
	defer span.End()
	start := time.Now()
	defer func() {
		newPayloadLatency.Observe(float64(time.Since(start).Milliseconds()))
	}()

	d := time.Now().Add(time.Duration(params.BeaconConfig().ExecutionEngineTimeoutValue) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	result := &pb.PayloadStatus{}

	switch payload.Proto().(type) {
	case *pb.ExecutionPayload:
		payloadPb, ok := payload.Proto().(*pb.ExecutionPayload)
		if !ok {
			return nil, errors.New("execution data must be a Bellatrix or Capella execution payload")
		}
		err := s.rpcClient.CallContext(ctx, result, NewPayloadMethod, payloadPb)
		if err != nil {
			return nil, handleRPCError(err)
		}
	case *pb.ExecutionPayloadCapella:
		payloadPb, ok := payload.Proto().(*pb.ExecutionPayloadCapella)
		if !ok {
			return nil, errors.New("execution data must be a Capella execution payload")
		}
		err := s.rpcClient.CallContext(ctx, result, NewPayloadMethodV2, payloadPb)
		if err != nil {
			return nil, handleRPCError(err)
		}
	default:
		return nil, errors.New("unknown execution data type")
	}

	switch result.Status {
	case pb.PayloadStatus_INVALID_BLOCK_HASH:
		return nil, ErrInvalidBlockHashPayloadStatus
	case pb.PayloadStatus_ACCEPTED, pb.PayloadStatus_SYNCING:
		return nil, ErrAcceptedSyncingPayloadStatus
	case pb.PayloadStatus_INVALID:
		return result.LatestValidHash, ErrInvalidPayloadStatus
	case pb.PayloadStatus_VALID:
		return result.LatestValidHash, nil
	default:
		return nil, ErrUnknownPayloadStatus
	}
}

// ForkchoiceUpdated calls the engine_forkchoiceUpdatedV1 method via JSON-RPC.
func (s *Service) ForkchoiceUpdated(
	ctx context.Context, state *pb.ForkchoiceState, attrs payloadattribute.Attributer,
) (*pb.PayloadIDBytes, []byte, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.ForkchoiceUpdated")
	defer span.End()
	start := time.Now()
	defer func() {
		forkchoiceUpdatedLatency.Observe(float64(time.Since(start).Milliseconds()))
	}()

	d := time.Now().Add(time.Duration(params.BeaconConfig().ExecutionEngineTimeoutValue) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	result := &ForkchoiceUpdatedResponse{}

	if attrs == nil {
		return nil, nil, errors.New("nil payload attributer")
	}
	switch attrs.Version() {
	case version.Bellatrix:
		a, err := attrs.PbV1()
		if err != nil {
			return nil, nil, err
		}
		err = s.rpcClient.CallContext(ctx, result, ForkchoiceUpdatedMethod, state, a)
		if err != nil {
			return nil, nil, handleRPCError(err)
		}
	case version.Capella:
		a, err := attrs.PbV2()
		if err != nil {
			return nil, nil, err
		}
		err = s.rpcClient.CallContext(ctx, result, ForkchoiceUpdatedMethodV2, state, a)
		if err != nil {
			return nil, nil, handleRPCError(err)
		}
	default:
		return nil, nil, fmt.Errorf("unknown payload attribute version: %v", attrs.Version())
	}

	if result.Status == nil {
		return nil, nil, ErrNilResponse
	}
	resp := result.Status
	switch resp.Status {
	case pb.PayloadStatus_SYNCING:
		return nil, nil, ErrAcceptedSyncingPayloadStatus
	case pb.PayloadStatus_INVALID:
		return nil, resp.LatestValidHash, ErrInvalidPayloadStatus
	case pb.PayloadStatus_VALID:
		return result.PayloadId, resp.LatestValidHash, nil
	default:
		return nil, nil, ErrUnknownPayloadStatus
	}
}

// GetPayload calls the engine_getPayloadVX method via JSON-RPC.
func (s *Service) GetPayload(ctx context.Context, payloadId [8]byte, slot primitives.Slot) (interfaces.ExecutionData, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.GetPayload")
	defer span.End()
	start := time.Now()
	defer func() {
		getPayloadLatency.Observe(float64(time.Since(start).Milliseconds()))
	}()

	d := time.Now().Add(defaultEngineTimeout)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()

	if slots.ToEpoch(slot) >= params.BeaconConfig().CapellaForkEpoch {
		result := &pb.ExecutionPayloadCapellaWithValue{}
		err := s.rpcClient.CallContext(ctx, result, GetPayloadMethodV2, pb.PayloadIDBytes(payloadId))
		if err != nil {
			return nil, handleRPCError(err)
		}

		return blocks.WrappedExecutionPayloadCapella(result.Payload, big.NewInt(0).SetBytes(bytesutil.ReverseByteOrder(result.Value)))
	}

	result := &pb.ExecutionPayload{}
	err := s.rpcClient.CallContext(ctx, result, GetPayloadMethod, pb.PayloadIDBytes(payloadId))
	if err != nil {
		return nil, handleRPCError(err)
	}
	return blocks.WrappedExecutionPayload(result)
}

// ExchangeTransitionConfiguration calls the engine_exchangeTransitionConfigurationV1 method via JSON-RPC.
func (s *Service) ExchangeTransitionConfiguration(
	ctx context.Context, cfg *pb.TransitionConfiguration,
) error {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.ExchangeTransitionConfiguration")
	defer span.End()

	// We set terminal block number to 0 as the parameter is not set on the consensus layer.
	zeroBigNum := big.NewInt(0)
	cfg.TerminalBlockNumber = zeroBigNum.Bytes()
	d := time.Now().Add(defaultEngineTimeout)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	result := &pb.TransitionConfiguration{}
	if err := s.rpcClient.CallContext(ctx, result, ExchangeTransitionConfigurationMethod, cfg); err != nil {
		return handleRPCError(err)
	}

	// We surface an error to the user if local configuration settings mismatch
	// according to the response from the execution node.
	cfgTerminalHash := params.BeaconConfig().TerminalBlockHash[:]
	if !bytes.Equal(cfgTerminalHash, result.TerminalBlockHash) {
		return errors.Wrapf(
			ErrConfigMismatch,
			"got %#x from execution node, wanted %#x",
			result.TerminalBlockHash,
			cfgTerminalHash,
		)
	}
	ttdCfg := params.BeaconConfig().TerminalTotalDifficulty
	ttdResult, err := hexutil.DecodeBig(result.TerminalTotalDifficulty)
	if err != nil {
		return errors.Wrap(err, "could not decode received terminal total difficulty")
	}
	if ttdResult.String() != ttdCfg {
		return errors.Wrapf(
			ErrConfigMismatch,
			"got %s from execution node, wanted %s",
			ttdResult.String(),
			ttdCfg,
		)
	}
	return nil
}

// GetTerminalBlockHash returns the valid terminal block hash based on total difficulty.
//
// Spec code:
// def get_pow_block_at_terminal_total_difficulty(pow_chain: Dict[Hash32, PowBlock]) -> Optional[PowBlock]:
//
//	# `pow_chain` abstractly represents all blocks in the PoW chain
//	for block in pow_chain:
//	    parent = pow_chain[block.parent_hash]
//	    block_reached_ttd = block.total_difficulty >= TERMINAL_TOTAL_DIFFICULTY
//	    parent_reached_ttd = parent.total_difficulty >= TERMINAL_TOTAL_DIFFICULTY
//	    if block_reached_ttd and not parent_reached_ttd:
//	        return block
//
//	return None
func (s *Service) GetTerminalBlockHash(ctx context.Context, transitionTime uint64) ([]byte, bool, error) {
	ttd := new(big.Int)
	ttd.SetString(params.BeaconConfig().TerminalTotalDifficulty, 10)
	terminalTotalDifficulty, overflows := uint256.FromBig(ttd)
	if overflows {
		return nil, false, errors.New("could not convert terminal total difficulty to uint256")
	}
	blk, err := s.LatestExecutionBlock(ctx)
	if err != nil {
		return nil, false, errors.Wrap(err, "could not get latest execution block")
	}
	if blk == nil {
		return nil, false, errors.New("latest execution block is nil")
	}

	for {
		if ctx.Err() != nil {
			return nil, false, ctx.Err()
		}
		currentTotalDifficulty, err := tDStringToUint256(blk.TotalDifficulty)
		if err != nil {
			return nil, false, errors.Wrap(err, "could not convert total difficulty to uint256")
		}
		blockReachedTTD := currentTotalDifficulty.Cmp(terminalTotalDifficulty) >= 0

		parentHash := blk.ParentHash
		if parentHash == params.BeaconConfig().ZeroHash {
			return nil, false, nil
		}
		parentBlk, err := s.ExecutionBlockByHash(ctx, parentHash, false /* no txs */)
		if err != nil {
			return nil, false, errors.Wrap(err, "could not get parent execution block")
		}
		if parentBlk == nil {
			return nil, false, errors.New("parent execution block is nil")
		}

		if blockReachedTTD {
			parentTotalDifficulty, err := tDStringToUint256(parentBlk.TotalDifficulty)
			if err != nil {
				return nil, false, errors.Wrap(err, "could not convert total difficulty to uint256")
			}

			// If terminal block has time same timestamp or greater than transition time,
			// then the node violates the invariant that a block's timestamp must be
			// greater than its parent's timestamp. Execution layer will reject
			// a fcu call with such payload attributes. It's best that we return `None` in this a case.
			parentReachedTTD := parentTotalDifficulty.Cmp(terminalTotalDifficulty) >= 0
			if !parentReachedTTD {
				if blk.Time >= transitionTime {
					return nil, false, nil
				}

				log.WithFields(logrus.Fields{
					"number":   blk.Number,
					"hash":     fmt.Sprintf("%#x", bytesutil.Trunc(blk.Hash[:])),
					"td":       blk.TotalDifficulty,
					"parentTd": parentBlk.TotalDifficulty,
					"ttd":      terminalTotalDifficulty,
				}).Info("Retrieved terminal block hash")
				return blk.Hash[:], true, nil
			}
		} else {
			return nil, false, nil
		}
		blk = parentBlk
	}
}

// PayloadAttributes sends payload attributes to a block builder to trigger building of a block
func (s *Service) PayloadAttributes(ctx context.Context, attrs *pb.BuilderPayloadAttributes) ([]byte, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.builder-api-client.PayloadAttributes")
	defer span.End()
	start := time.Now()
	defer func() {
		payloadAttributesLatency.Observe(float64(time.Since(start).Milliseconds()))
	}()
	d := time.Now().Add(time.Duration(params.BeaconConfig().ExecutionEngineTimeoutValue) * time.Second)
	ctx, cancel := context.WithDeadline(ctx, d)
	defer cancel()
	result := &PayloadAttributesResponse{}
	err := s.rpcClient.CallContext(ctx, result, PayloadAttributesMethod, attrs)
	if err != nil {
		return nil, handleRPCError(err)
	}

	if result.Status == nil {
		return nil, ErrNilResponse
	}
	return nil, nil
}

// LatestExecutionBlock fetches the latest execution engine block by calling
// eth_blockByNumber via JSON-RPC.
func (s *Service) LatestExecutionBlock(ctx context.Context) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.LatestExecutionBlock")
	defer span.End()

	result := &pb.ExecutionBlock{}
	err := s.rpcClient.CallContext(
		ctx,
		result,
		ExecutionBlockByNumberMethod,
		"latest",
		false, /* no full transaction objects */
	)
	return result, handleRPCError(err)
}

// ExecutionBlockByHash fetches an execution engine block by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Service) ExecutionBlockByHash(ctx context.Context, hash common.Hash, withTxs bool) (*pb.ExecutionBlock, error) {
	ctx, span := trace.StartSpan(ctx, "powchain.engine-api-client.ExecutionBlockByHash")
	defer span.End()
	result := &pb.ExecutionBlock{}
	err := s.rpcClient.CallContext(ctx, result, ExecutionBlockByHashMethod, hash, withTxs)
	return result, handleRPCError(err)
}

// ExecutionBlocksByHashes fetches a batch of execution engine blocks by hash by calling
// eth_blockByHash via JSON-RPC.
func (s *Service) ExecutionBlocksByHashes(ctx context.Context, hashes []common.Hash, withTxs bool) ([]*pb.ExecutionBlock, error) {
	_, span := trace.StartSpan(ctx, "powchain.engine-api-client.ExecutionBlocksByHashes")
	defer span.End()
	numOfHashes := len(hashes)
	elems := make([]gethRPC.BatchElem, 0, numOfHashes)
	execBlks := make([]*pb.ExecutionBlock, 0, numOfHashes)
	if numOfHashes == 0 {
		return execBlks, nil
	}
	for _, h := range hashes {
		blk := &pb.ExecutionBlock{}
		newH := h
		elems = append(elems, gethRPC.BatchElem{
			Method: ExecutionBlockByHashMethod,
			Args:   []interface{}{newH, withTxs},
			Result: blk,
			Error:  error(nil),
		})
		execBlks = append(execBlks, blk)
	}
	ioErr := s.rpcClient.BatchCall(elems)
	if ioErr != nil {
		return nil, ioErr
	}
	for _, e := range elems {
		if e.Error != nil {
			return nil, handleRPCError(e.Error)
		}
	}
	return execBlks, nil
}

// HeaderByHash returns the relevant header details for the provided block hash.
func (s *Service) HeaderByHash(ctx context.Context, hash common.Hash) (*types.HeaderInfo, error) {
	var hdr *types.HeaderInfo
	err := s.rpcClient.CallContext(ctx, &hdr, ExecutionBlockByHashMethod, hash, false /* no transactions */)
	if err == nil && hdr == nil {
		err = ethereum.NotFound
	}
	return hdr, err
}

// HeaderByNumber returns the relevant header details for the provided block number.
func (s *Service) HeaderByNumber(ctx context.Context, number *big.Int) (*types.HeaderInfo, error) {
	var hdr *types.HeaderInfo
	err := s.rpcClient.CallContext(ctx, &hdr, ExecutionBlockByNumberMethod, toBlockNumArg(number), false /* no transactions */)
	if err == nil && hdr == nil {
		err = ethereum.NotFound
	}
	return hdr, err
}

// ReconstructFullBlock takes in a blinded beacon block and reconstructs
// a beacon block with a full execution payload via the engine API.
func (s *Service) ReconstructFullBlock(
	ctx context.Context, blindedBlock interfaces.ReadOnlySignedBeaconBlock,
) (interfaces.SignedBeaconBlock, error) {
	if err := blocks.BeaconBlockIsNil(blindedBlock); err != nil {
		return nil, errors.Wrap(err, "cannot reconstruct bellatrix block from nil data")
	}
	if !blindedBlock.Block().IsBlinded() {
		return nil, errors.New("can only reconstruct block from blinded block format")
	}
	header, err := blindedBlock.Block().Body().Execution()
	if err != nil {
		return nil, err
	}
	if header.IsNil() {
		return nil, errors.New("execution payload header in blinded block was nil")
	}

	// If the payload header has a block hash of 0x0, it means we are pre-merge and should
	// simply return the block with an empty execution payload.
	if bytes.Equal(header.BlockHash(), params.BeaconConfig().ZeroHash[:]) {
		payload := buildEmptyExecutionPayload()
		return blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlock, payload)
	}

	executionBlockHash := common.BytesToHash(header.BlockHash())
	executionBlock, err := s.ExecutionBlockByHash(ctx, executionBlockHash, true /* with txs */)
	if err != nil {
		return nil, fmt.Errorf("could not fetch execution block with txs by hash %#x: %v", executionBlockHash, err)
	}
	if executionBlock == nil {
		return nil, fmt.Errorf("received nil execution block for request by hash %#x", executionBlockHash)
	}
	if bytes.Equal(executionBlock.Hash.Bytes(), []byte{}) {
		return nil, EmptyBlockHash
	}

	executionBlock.Version = blindedBlock.Version()
	payload, err := fullPayloadFromExecutionBlock(header, executionBlock)
	if err != nil {
		return nil, err
	}
	fullBlock, err := blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlock, payload.Proto())
	if err != nil {
		return nil, err
	}
	reconstructedExecutionPayloadCount.Add(1)
	return fullBlock, nil
}

// ReconstructFullBellatrixBlockBatch takes in a batch of blinded beacon blocks and reconstructs
// them with a full execution payload for each block via the engine API.
func (s *Service) ReconstructFullBellatrixBlockBatch(
	ctx context.Context, blindedBlocks []interfaces.ReadOnlySignedBeaconBlock,
) ([]interfaces.SignedBeaconBlock, error) {
	if len(blindedBlocks) == 0 {
		return []interfaces.SignedBeaconBlock{}, nil
	}
	executionHashes := []common.Hash{}
	validExecPayloads := []int{}
	zeroExecPayloads := []int{}
	for i, b := range blindedBlocks {
		if err := blocks.BeaconBlockIsNil(b); err != nil {
			return nil, errors.Wrap(err, "cannot reconstruct bellatrix block from nil data")
		}
		if !b.Block().IsBlinded() {
			return nil, errors.New("can only reconstruct block from blinded block format")
		}
		header, err := b.Block().Body().Execution()
		if err != nil {
			return nil, err
		}
		if header.IsNil() {
			return nil, errors.New("execution payload header in blinded block was nil")
		}
		// Determine if the block is pre-merge or post-merge. Depending on the result,
		// we will ask the execution engine for the full payload.
		if bytes.Equal(header.BlockHash(), params.BeaconConfig().ZeroHash[:]) {
			zeroExecPayloads = append(zeroExecPayloads, i)
		} else {
			executionBlockHash := common.BytesToHash(header.BlockHash())
			validExecPayloads = append(validExecPayloads, i)
			executionHashes = append(executionHashes, executionBlockHash)
		}
	}
	execBlocks, err := s.ExecutionBlocksByHashes(ctx, executionHashes, true /* with txs*/)
	if err != nil {
		return nil, fmt.Errorf("could not fetch execution blocks with txs by hash %#x: %v", executionHashes, err)
	}

	// For each valid payload, we reconstruct the full block from it with the
	// blinded block.
	fullBlocks := make([]interfaces.SignedBeaconBlock, len(blindedBlocks))
	for sliceIdx, realIdx := range validExecPayloads {
		b := execBlocks[sliceIdx]
		if b == nil {
			return nil, fmt.Errorf("received nil execution block for request by hash %#x", executionHashes[sliceIdx])
		}
		header, err := blindedBlocks[realIdx].Block().Body().Execution()
		if err != nil {
			return nil, err
		}
		payload, err := fullPayloadFromExecutionBlock(header, b)
		if err != nil {
			return nil, err
		}
		fullBlock, err := blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlocks[realIdx], payload.Proto())
		if err != nil {
			return nil, err
		}
		fullBlocks[realIdx] = fullBlock
	}
	// For blocks that are pre-merge we simply reconstruct them via an empty
	// execution payload.
	for _, realIdx := range zeroExecPayloads {
		payload := buildEmptyExecutionPayload()
		fullBlock, err := blocks.BuildSignedBeaconBlockFromExecutionPayload(blindedBlocks[realIdx], payload)
		if err != nil {
			return nil, err
		}
		fullBlocks[realIdx] = fullBlock
	}
	reconstructedExecutionPayloadCount.Add(float64(len(blindedBlocks)))
	return fullBlocks, nil
}

func fullPayloadFromExecutionBlock(
	header interfaces.ExecutionData, block *pb.ExecutionBlock,
) (interfaces.ExecutionData, error) {
	if header.IsNil() || block == nil {
		return nil, errors.New("execution block and header cannot be nil")
	}
	blockHash := block.Hash
	if !bytes.Equal(header.BlockHash(), blockHash[:]) {
		return nil, fmt.Errorf(
			"block hash field in execution header %#x does not match execution block hash %#x",
			header.BlockHash(),
			blockHash,
		)
	}
	blockTransactions := block.Transactions
	txs := make([][]byte, len(blockTransactions))
	for i, tx := range blockTransactions {
		txBin, err := tx.MarshalBinary()
		if err != nil {
			return nil, err
		}
		txs[i] = txBin
	}

	if block.Version == version.Bellatrix {
		return blocks.WrappedExecutionPayload(&pb.ExecutionPayload{
			ParentHash:    header.ParentHash(),
			FeeRecipient:  header.FeeRecipient(),
			StateRoot:     header.StateRoot(),
			ReceiptsRoot:  header.ReceiptsRoot(),
			LogsBloom:     header.LogsBloom(),
			PrevRandao:    header.PrevRandao(),
			BlockNumber:   header.BlockNumber(),
			GasLimit:      header.GasLimit(),
			GasUsed:       header.GasUsed(),
			Timestamp:     header.Timestamp(),
			ExtraData:     header.ExtraData(),
			BaseFeePerGas: header.BaseFeePerGas(),
			BlockHash:     blockHash[:],
			Transactions:  txs,
		})
	}
	return blocks.WrappedExecutionPayloadCapella(&pb.ExecutionPayloadCapella{
		ParentHash:    header.ParentHash(),
		FeeRecipient:  header.FeeRecipient(),
		StateRoot:     header.StateRoot(),
		ReceiptsRoot:  header.ReceiptsRoot(),
		LogsBloom:     header.LogsBloom(),
		PrevRandao:    header.PrevRandao(),
		BlockNumber:   header.BlockNumber(),
		GasLimit:      header.GasLimit(),
		GasUsed:       header.GasUsed(),
		Timestamp:     header.Timestamp(),
		ExtraData:     header.ExtraData(),
		BaseFeePerGas: header.BaseFeePerGas(),
		BlockHash:     blockHash[:],
		Transactions:  txs,
		Withdrawals:   block.Withdrawals,
	}, big.NewInt(0)) // We can't get the block value and don't care about the block value for this instance
}

// Handles errors received from the RPC server according to the specification.
func handleRPCError(err error) error {
	if err == nil {
		return nil
	}
	if isTimeout(err) {
		return ErrHTTPTimeout
	}
	e, ok := err.(gethRPC.Error)
	if !ok {
		if strings.Contains(err.Error(), "401 Unauthorized") {
			log.Error("HTTP authentication to your execution client is not working. Please ensure " +
				"you are setting a correct value for the --jwt-secret flag in Prysm, or use an IPC connection if on " +
				"the same machine. Please see our documentation for more information on authenticating connections " +
				"here https://docs.prylabs.network/docs/execution-node/authentication")
			return fmt.Errorf("could not authenticate connection to execution client: %v", err)
		}
		return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
	}
	switch e.ErrorCode() {
	case -32700:
		errParseCount.Inc()
		return ErrParse
	case -32600:
		errInvalidRequestCount.Inc()
		return ErrInvalidRequest
	case -32601:
		errMethodNotFoundCount.Inc()
		return ErrMethodNotFound
	case -32602:
		errInvalidParamsCount.Inc()
		return ErrInvalidParams
	case -32603:
		errInternalCount.Inc()
		return ErrInternal
	case -38001:
		errUnknownPayloadCount.Inc()
		return ErrUnknownPayload
	case -38002:
		errInvalidForkchoiceStateCount.Inc()
		return ErrInvalidForkchoiceState
	case -38003:
		errInvalidPayloadAttributesCount.Inc()
		return ErrInvalidPayloadAttributes
	case -32000:
		errServerErrorCount.Inc()
		// Only -32000 status codes are data errors in the RPC specification.
		errWithData, ok := err.(gethRPC.DataError)
		if !ok {
			return errors.Wrapf(err, "got an unexpected error in JSON-RPC response")
		}
		return errors.Wrapf(ErrServer, "%v", errWithData.Error())
	default:
		return err
	}
}

// ErrHTTPTimeout returns true if the error is a http.Client timeout error.
var ErrHTTPTimeout = errors.New("timeout from http.Client")

type httpTimeoutError interface {
	Error() string
	Timeout() bool
}

func isTimeout(e error) bool {
	t, ok := e.(httpTimeoutError)
	return ok && t.Timeout()
}

func tDStringToUint256(td string) (*uint256.Int, error) {
	b, err := hexutil.DecodeBig(td)
	if err != nil {
		return nil, err
	}
	i, overflows := uint256.FromBig(b)
	if overflows {
		return nil, errors.New("total difficulty overflowed")
	}
	return i, nil
}

func buildEmptyExecutionPayload() *pb.ExecutionPayload {
	return &pb.ExecutionPayload{
		ParentHash:    make([]byte, fieldparams.RootLength),
		FeeRecipient:  make([]byte, fieldparams.FeeRecipientLength),
		StateRoot:     make([]byte, fieldparams.RootLength),
		ReceiptsRoot:  make([]byte, fieldparams.RootLength),
		LogsBloom:     make([]byte, fieldparams.LogsBloomLength),
		PrevRandao:    make([]byte, fieldparams.RootLength),
		BaseFeePerGas: make([]byte, fieldparams.RootLength),
		BlockHash:     make([]byte, fieldparams.RootLength),
		Transactions:  make([][]byte, 0),
		ExtraData:     make([]byte, 0),
	}
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	pending := big.NewInt(-1)
	if number.Cmp(pending) == 0 {
		return "pending"
	}
	finalized := big.NewInt(int64(gethRPC.FinalizedBlockNumber))
	if number.Cmp(finalized) == 0 {
		return "finalized"
	}
	safe := big.NewInt(int64(gethRPC.SafeBlockNumber))
	if number.Cmp(safe) == 0 {
		return "safe"
	}
	return hexutil.EncodeBig(number)
}
