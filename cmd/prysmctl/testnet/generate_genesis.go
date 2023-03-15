package testnet

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/altair"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/capella"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/core/execution"
	"github.com/prysmaticlabs/prysm/v3/beacon-chain/state"
	state_native "github.com/prysmaticlabs/prysm/v3/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v3/cmd/flags"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/container/trie"
	"github.com/prysmaticlabs/prysm/v3/io/file"
	ethpb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v3/runtime/interop"
	"github.com/prysmaticlabs/prysm/v3/runtime/version"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	generateGenesisStateFlags = struct {
		DepositJsonFile   string
		ChainConfigFile   string
		ConfigName        string
		NumValidators     uint64
		GenesisTime       uint64
		OutputSSZ         string
		OutputJSON        string
		OutputYaml        string
		ForkName          string
		OverrideEth1Data  bool
		ExecutionEndpoint string
	}{}
	log           = logrus.WithField("prefix", "genesis")
	outputSSZFlag = &cli.StringFlag{
		Name:        "output-ssz",
		Destination: &generateGenesisStateFlags.OutputSSZ,
		Usage:       "Output filename of the SSZ marshaling of the generated genesis state",
		Value:       "",
	}
	outputYamlFlag = &cli.StringFlag{
		Name:        "output-yaml",
		Destination: &generateGenesisStateFlags.OutputYaml,
		Usage:       "Output filename of the YAML marshaling of the generated genesis state",
		Value:       "",
	}
	outputJsonFlag = &cli.StringFlag{
		Name:        "output-json",
		Destination: &generateGenesisStateFlags.OutputJSON,
		Usage:       "Output filename of the JSON marshaling of the generated genesis state",
		Value:       "",
	}
	generateGenesisStateCmd = &cli.Command{
		Name:  "generate-genesis",
		Usage: "Generate a beacon chain genesis state",
		Action: func(cliCtx *cli.Context) error {
			if err := cliActionGenerateGenesisState(cliCtx); err != nil {
				log.WithError(err).Fatal("Could not generate beacon chain genesis state")
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "chain-config-file",
				Destination: &generateGenesisStateFlags.ChainConfigFile,
				Usage:       "The path to a YAML file with chain config values",
			},
			&cli.StringFlag{
				Name:        "deposit-json-file",
				Destination: &generateGenesisStateFlags.DepositJsonFile,
				Usage:       "Path to deposit_data.json file generated by the staking-deposit-cli tool for optionally specifying validators in genesis state",
			},
			&cli.StringFlag{
				Name:        "config-name",
				Usage:       "Config kind to be used for generating the genesis state. Default: mainnet. Options include mainnet, interop, minimal, prater, sepolia. --chain-config-file will override this flag.",
				Destination: &generateGenesisStateFlags.ConfigName,
				Value:       params.MainnetName,
			},
			&cli.Uint64Flag{
				Name:        "num-validators",
				Usage:       "Number of validators to deterministically generate in the genesis state",
				Destination: &generateGenesisStateFlags.NumValidators,
				Required:    true,
			},
			&cli.Uint64Flag{
				Name:        "genesis-time",
				Destination: &generateGenesisStateFlags.GenesisTime,
				Usage:       "Unix timestamp seconds used as the genesis time in the genesis state. If unset, defaults to now()",
			},
			&cli.BoolFlag{
				Name:        "override-eth1data",
				Destination: &generateGenesisStateFlags.OverrideEth1Data,
				Usage:       "Overrides Eth1Data with values from execution client. If unset, defaults to false",
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "execution-endpoint",
				Destination: &generateGenesisStateFlags.ExecutionEndpoint,
				Usage:       "Endpoint to preferred execution client. If unset, defaults to Geth",
				Value:       "http://localhost:8545",
			},
			flags.EnumValue{
				Name:        "fork",
				Usage:       fmt.Sprintf("Name of the BeaconState schema to use in output encoding [%s]", strings.Join(versionNames(), ",")),
				Enum:        versionNames(),
				Value:       versionNames()[0],
				Destination: &generateGenesisStateFlags.ForkName,
			}.GenericFlag(),
			outputSSZFlag,
			outputYamlFlag,
			outputJsonFlag,
		},
	}
)

func versionNames() []string {
	enum := version.All()
	names := make([]string, len(enum))
	for i := range enum {
		names[i] = version.String(enum[i])
	}
	return names
}

// Represents a json object of hex string and uint64 values for
// validators on Ethereum. This file can be generated using the official staking-deposit-cli.
type depositDataJSON struct {
	PubKey                string `json:"pubkey"`
	Amount                uint64 `json:"amount"`
	WithdrawalCredentials string `json:"withdrawal_credentials"`
	DepositDataRoot       string `json:"deposit_data_root"`
	Signature             string `json:"signature"`
}

func cliActionGenerateGenesisState(cliCtx *cli.Context) error {
	if generateGenesisStateFlags.GenesisTime == 0 {
		log.Info("No genesis time specified, defaulting to now()")
	}
	outputJson := generateGenesisStateFlags.OutputJSON
	outputYaml := generateGenesisStateFlags.OutputYaml
	outputSSZ := generateGenesisStateFlags.OutputSSZ
	noOutputFlag := outputSSZ == "" && outputJson == "" && outputYaml == ""
	if noOutputFlag {
		return fmt.Errorf(
			"no %s, %s, %s flag(s) specified. At least one is required",
			outputJsonFlag.Name,
			outputYamlFlag.Name,
			outputSSZFlag.Name,
		)
	}
	if err := setGlobalParams(); err != nil {
		return fmt.Errorf("could not set config params: %v", err)
	}
	genesisState, err := generateGenesis(cliCtx.Context)
	if err != nil {
		return fmt.Errorf("could not generate genesis state: %v", err)
	}
	st, err := upgradeStateToForkName(cliCtx.Context, genesisState, generateGenesisStateFlags.ForkName)
	if err != nil {
		return err
	}

	if outputJson != "" {
		if err := writeToOutputFile(outputJson, st, json.Marshal); err != nil {
			return err
		}
	}
	if outputYaml != "" {
		if err := writeToOutputFile(outputYaml, st, yaml.Marshal); err != nil {
			return err
		}
	}
	if outputSSZ != "" {
		type MinimumSSZMarshal interface {
			MarshalSSZ() ([]byte, error)
		}
		marshalFn := func(o interface{}) ([]byte, error) {
			marshaler, ok := o.(MinimumSSZMarshal)
			if !ok {
				return nil, errors.New("not a marshaler")
			}
			return marshaler.MarshalSSZ()
		}
		if err := writeToOutputFile(outputSSZ, st, marshalFn); err != nil {
			return err
		}
	}
	log.Info("Command completed")
	return nil
}

func upgradeStateToForkName(ctx context.Context, pbst *ethpb.BeaconState, name string) (state.BeaconState, error) {
	st, err := state_native.InitializeFromProtoUnsafePhase0(pbst)
	if err != nil {
		return nil, err
	}
	if name == "" || name == "phase0" {
		return st, nil
	}

	st, err = altair.UpgradeToAltair(ctx, st)
	if err != nil {
		return nil, err
	}
	if name == "altair" {
		return st, nil
	}

	st, err = execution.UpgradeToBellatrix(st)
	if err != nil {
		return nil, err
	}
	if name == "bellatrix" {
		return st, nil
	}

	st, err = capella.UpgradeToCapella(st)
	if err != nil {
		return nil, err
	}
	if name == "capella" {
		return st, nil
	}

	return nil, fmt.Errorf("unrecognized fork name '%s'", name)
}

func setGlobalParams() error {
	chainConfigFile := generateGenesisStateFlags.ChainConfigFile
	if chainConfigFile != "" {
		log.Infof("Specified a chain config file: %s", chainConfigFile)
		return params.LoadChainConfigFile(chainConfigFile, nil)
	}
	cfg, err := params.ByName(generateGenesisStateFlags.ConfigName)
	if err != nil {
		return fmt.Errorf("unable to find config using name %s: %v", generateGenesisStateFlags.ConfigName, err)
	}
	return params.SetActive(cfg.Copy())
}

func generateGenesis(ctx context.Context) (*ethpb.BeaconState, error) {
	genesisTime := generateGenesisStateFlags.GenesisTime
	numValidators := generateGenesisStateFlags.NumValidators
	depositJsonFile := generateGenesisStateFlags.DepositJsonFile
	overrideEth1Data := generateGenesisStateFlags.OverrideEth1Data
	if depositJsonFile != "" {
		expanded, err := file.ExpandPath(depositJsonFile)
		if err != nil {
			return nil, err
		}
		inputJSON, err := os.Open(expanded) // #nosec G304
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := inputJSON.Close(); err != nil {
				log.WithError(err).Printf("Could not close file %s", depositJsonFile)
			}
		}()
		log.Printf("Generating genesis state from input JSON deposit data %s", depositJsonFile)
		return genesisStateFromJSONValidators(ctx, inputJSON, genesisTime)
	}
	if numValidators == 0 {
		return nil, fmt.Errorf(
			"expected --num-validators > 0 to have been provided",
		)
	}
	// If no JSON input is specified, we create the state deterministically from interop keys.
	genesisState, _, err := interop.GenerateGenesisState(ctx, genesisTime, numValidators)
	if err != nil {
		return nil, err
	}
	if overrideEth1Data {
		log.Print("Overriding Eth1Data with data from execution client")
		conn, err := rpc.Dial(generateGenesisStateFlags.ExecutionEndpoint)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"could not dial %s please make sure you are running your execution client",
				generateGenesisStateFlags.ExecutionEndpoint)
		}
		client := ethclient.NewClient(conn)
		header, err := client.HeaderByNumber(ctx, big.NewInt(0))
		if err != nil {
			return nil, errors.Wrap(err, "could not get header by number")
		}
		t, err := trie.NewTrie(params.BeaconConfig().DepositContractTreeDepth)
		if err != nil {
			return nil, errors.Wrap(err, "could not create deposit tree")
		}
		depositRoot, err := t.HashTreeRoot()
		if err != nil {
			return nil, errors.Wrap(err, "could not get hash tree root")
		}
		genesisState.Eth1Data = &ethpb.Eth1Data{
			DepositRoot:  depositRoot[:],
			DepositCount: 0,
			BlockHash:    header.Hash().Bytes(),
		}
		genesisState.Eth1DepositIndex = 0
	}
	return genesisState, err
}

func genesisStateFromJSONValidators(ctx context.Context, r io.Reader, genesisTime uint64) (*ethpb.BeaconState, error) {
	enc, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var depositJSON []*depositDataJSON
	if err := json.Unmarshal(enc, &depositJSON); err != nil {
		return nil, err
	}
	depositDataList := make([]*ethpb.Deposit_Data, len(depositJSON))
	depositDataRoots := make([][]byte, len(depositJSON))
	for i, val := range depositJSON {
		data, dataRootBytes, err := depositJSONToDepositData(val)
		if err != nil {
			return nil, err
		}
		depositDataList[i] = data
		depositDataRoots[i] = dataRootBytes
	}
	beaconState, _, err := interop.GenerateGenesisStateFromDepositData(ctx, genesisTime, depositDataList, depositDataRoots)
	if err != nil {
		return nil, err
	}
	return beaconState, nil
}

func depositJSONToDepositData(input *depositDataJSON) (depositData *ethpb.Deposit_Data, dataRoot []byte, err error) {
	pubKeyBytes, err := hex.DecodeString(strings.TrimPrefix(input.PubKey, "0x"))
	if err != nil {
		return
	}
	withdrawalbytes, err := hex.DecodeString(strings.TrimPrefix(input.WithdrawalCredentials, "0x"))
	if err != nil {
		return
	}
	signatureBytes, err := hex.DecodeString(strings.TrimPrefix(input.Signature, "0x"))
	if err != nil {
		return
	}
	dataRootBytes, err := hex.DecodeString(strings.TrimPrefix(input.DepositDataRoot, "0x"))
	if err != nil {
		return
	}
	depositData = &ethpb.Deposit_Data{
		PublicKey:             pubKeyBytes,
		WithdrawalCredentials: withdrawalbytes,
		Amount:                input.Amount,
		Signature:             signatureBytes,
	}
	dataRoot = dataRootBytes
	return
}

func writeToOutputFile(
	fPath string,
	data interface{},
	marshalFn func(o interface{}) ([]byte, error),
) error {
	encoded, err := marshalFn(data)
	if err != nil {
		return err
	}
	if err := file.WriteFile(fPath, encoded); err != nil {
		return err
	}
	log.Printf("Done writing genesis state to %s", fPath)
	return nil
}
