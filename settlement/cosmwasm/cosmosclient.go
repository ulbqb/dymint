package cosmwasm

import (
	"context"
	"fmt"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/cosmosclient/cosmosclient"
	rollapptypes "github.com/dymensionxyz/dymension/x/rollapp/types"
	sequencertypes "github.com/dymensionxyz/dymension/x/sequencer/types"
	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	tmjson "github.com/tendermint/tendermint/libs/json"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"google.golang.org/grpc"
)

// CosmosClient is an interface for interacting with cosmos client chains.
// It is a wrapper around the cosmos client in order to provide with an interface which can be implemented by
// other clients and can easily be mocked for testing purposes.
// Currently it contains only the methods that are used by the dymension hub client.
type CosmosClient interface {
	Context() sdkclient.Context
	StartEventListener() error
	StopEventListener() error
	EventListenerQuit() <-chan struct{}
	SubscribeToEvents(ctx context.Context, subscriber string, query string, outCapacity ...int) (out <-chan ctypes.ResultEvent, err error)
	BroadcastTx(accountName string, msgs ...sdktypes.Msg) (cosmosclient.Response, error)
	GetRollappClient() rollapptypes.QueryClient
	GetSequencerClient() sequencertypes.QueryClient
	GetAccount(accountName string) (cosmosaccount.Account, error)
}

type cosmosClient struct {
	cosmosclient.Client
	constract string
}

var _ CosmosClient = &cosmosClient{}

// NewCosmosClient creates a new cosmos client
func NewCosmosClient(client cosmosclient.Client, constract string) CosmosClient {
	return &cosmosClient{
		Client:    client,
		constract: constract,
	}
}

func (c *cosmosClient) StartEventListener() error {
	return c.Client.RPC.WSEvents.Start()
}

func (c *cosmosClient) StopEventListener() error {
	return c.Client.RPC.WSEvents.Stop()
}

func (c *cosmosClient) EventListenerQuit() <-chan struct{} {
	return c.Client.RPC.GetWSClient().Quit()
}

func (c *cosmosClient) SubscribeToEvents(ctx context.Context, subscriber string, query string, outCapacity ...int) (out <-chan ctypes.ResultEvent, err error) {
	return c.Client.RPC.WSEvents.Subscribe(ctx, subscriber, query, outCapacity...)
}

func (c *cosmosClient) GetRollappClient() rollapptypes.QueryClient {
	return &RollappClient{
		inner:     wasmtypes.NewQueryClient(c.Context()),
		constract: "",
	}
}

func (c *cosmosClient) GetSequencerClient() sequencertypes.QueryClient {
	return &SequencerClient{
		inner:     wasmtypes.NewQueryClient(c.Context()),
		constract: "",
	}
}

func (c *cosmosClient) GetAccount(accountName string) (cosmosaccount.Account, error) {
	return c.AccountRegistry.GetByName(accountName)
}

type RollappClient struct {
	inner     wasmtypes.QueryClient
	constract string
}

// Parameters queries the parameters of the module.
func (c *RollappClient) Params(ctx context.Context, in *rollapptypes.QueryParamsRequest, opts ...grpc.CallOption) (*rollapptypes.QueryParamsResponse, error) {
	panic("this is not supported")
}

// Queries a Rollapp by index.
func (c *RollappClient) Rollapp(ctx context.Context, in *rollapptypes.QueryGetRollappRequest, opts ...grpc.CallOption) (*rollapptypes.QueryGetRollappResponse, error) {
	panic("this is not supported")
}

// Queries a list of Rollapp items.
func (c *RollappClient) RollappAll(ctx context.Context, in *rollapptypes.QueryAllRollappRequest, opts ...grpc.CallOption) (*rollapptypes.QueryAllRollappResponse, error) {
	panic("this is not supported")
}

// Queries a LatestStateIndex by rollapp-id.
func (c *RollappClient) LatestStateIndex(ctx context.Context, in *rollapptypes.QueryGetLatestStateIndexRequest, opts ...grpc.CallOption) (*rollapptypes.QueryGetLatestStateIndexResponse, error) {
	msgBytes, err := tmjson.Marshal(in)
	if err != nil {
		return nil, err
	}
	queryData := fmt.Sprintf("{\"rollapp\":{\"latestStateIndex\":%s}}", string(msgBytes))
	queryRes, err :=
		c.inner.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
			Address:   c.constract,
			QueryData: []byte(queryData),
		}, opts...)
	if err != nil {
		return nil, err
	}

	var contractRes rollapptypes.QueryGetLatestStateIndexResponse
	if err := tmjson.Unmarshal(queryRes.Data, &contractRes); err != nil {
		return nil, err
	}
	return &contractRes, nil
}

// Queries a StateInfo by index.
func (c *RollappClient) StateInfo(ctx context.Context, in *rollapptypes.QueryGetStateInfoRequest, opts ...grpc.CallOption) (*rollapptypes.QueryGetStateInfoResponse, error) {
	msgBytes, err := tmjson.Marshal(in)
	if err != nil {
		return nil, err
	}
	queryData := fmt.Sprintf("{\"rollapp\":{\"stateInfo\":%s}}", string(msgBytes))
	queryRes, err :=
		c.inner.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
			Address:   c.constract,
			QueryData: []byte(queryData),
		}, opts...)
	if err != nil {
		return nil, err
	}

	var contractRes rollapptypes.QueryGetStateInfoResponse
	if err := tmjson.Unmarshal(queryRes.Data, &contractRes); err != nil {
		return nil, err
	}
	return &contractRes, nil
}

// Queries a list of StateInfo items.
func (c *RollappClient) StateInfoAll(ctx context.Context, in *rollapptypes.QueryAllStateInfoRequest, opts ...grpc.CallOption) (*rollapptypes.QueryAllStateInfoResponse, error) {
	panic("this is not supported")
}

type SequencerClient struct {
	inner     wasmtypes.QueryClient
	constract string
}

// Parameters queries the parameters of the module.
func (c *SequencerClient) Params(ctx context.Context, in *sequencertypes.QueryParamsRequest, opts ...grpc.CallOption) (*sequencertypes.QueryParamsResponse, error) {
	panic("this is not supported")
}

// Queries a Sequencer by index.
func (c *SequencerClient) Sequencer(ctx context.Context, in *sequencertypes.QueryGetSequencerRequest, opts ...grpc.CallOption) (*sequencertypes.QueryGetSequencerResponse, error) {
	panic("this is not supported")
}

// Queries a list of Sequencer items.
func (c *SequencerClient) SequencerAll(ctx context.Context, in *sequencertypes.QueryAllSequencerRequest, opts ...grpc.CallOption) (*sequencertypes.QueryAllSequencerResponse, error) {
	panic("this is not supported")
}

// Queries a SequencersByRollapp by index.
func (c *SequencerClient) SequencersByRollapp(ctx context.Context, in *sequencertypes.QueryGetSequencersByRollappRequest, opts ...grpc.CallOption) (*sequencertypes.QueryGetSequencersByRollappResponse, error) {
	msgBytes, err := tmjson.Marshal(in)
	if err != nil {
		return nil, err
	}
	queryData := fmt.Sprintf("{\"sequencer\":{\"sequencersByRollapp\":%s}}", string(msgBytes))
	queryRes, err :=
		c.inner.SmartContractState(ctx, &wasmtypes.QuerySmartContractStateRequest{
			Address:   c.constract,
			QueryData: []byte(queryData),
		}, opts...)
	if err != nil {
		return nil, err
	}

	var contractRes sequencertypes.QueryGetSequencersByRollappResponse
	if err := tmjson.Unmarshal(queryRes.Data, &contractRes); err != nil {
		return nil, err
	}
	return &contractRes, nil
}

// Queries a list of SequencersByRollapp items.
func (c *SequencerClient) SequencersByRollappAll(ctx context.Context, in *sequencertypes.QueryAllSequencersByRollappRequest, opts ...grpc.CallOption) (*sequencertypes.QueryAllSequencersByRollappResponse, error) {
	panic("this is not supported")
}

// Queries a Scheduler by index.
func (c *SequencerClient) Scheduler(ctx context.Context, in *sequencertypes.QueryGetSchedulerRequest, opts ...grpc.CallOption) (*sequencertypes.QueryGetSchedulerResponse, error) {
	panic("this is not supported")
}

// Queries a list of Scheduler items.
func (c *SequencerClient) SchedulerAll(ctx context.Context, in *sequencertypes.QueryAllSchedulerRequest, opts ...grpc.CallOption) (*sequencertypes.QueryAllSchedulerResponse, error) {
	panic("this is not supported")
}
