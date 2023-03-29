package oracle

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/bandprotocol/bandchain-packet/obi"
	"github.com/bandprotocol/chain/v2/x/oracle/keeper"
	"github.com/bandprotocol/chain/v2/x/oracle/types"
)

// handleBeginBlock re-calculates and saves the rolling seed value based on block hashes.
func handleBeginBlock(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	currentReqID := k.GetRequestLastExpired(ctx) + 1
	lastReqID := types.RequestID(k.GetRequestCount(ctx))
	for ; currentReqID <= lastReqID; currentReqID++ {
		k.AddPendingRequest(ctx, currentReqID)
	}
	k.SetRequestLastExpired(ctx, lastReqID)
}

type CoinRatesCallData struct {
	Symbols    []string `protobuf:"bytes,1,rep,name=symbols,proto3" json:"symbols,omitempty"`
	Multiplier uint64   `protobuf:"varint,2,opt,name=multiplier,proto3" json:"multiplier,omitempty"`
}

type CoinRatesResult struct {
	Rates []uint64 `protobuf:"varint,1,rep,packed,name=rates,proto3" json:"rates,omitempty"`
}

// handleEndBlock cleans up the state during end block. See comment in the implementation!
func handleEndBlock(ctx sdk.Context, k keeper.Keeper) {
	// Loops through all requests in the resolvable list to resolve all of them!
	for _, reqID := range k.GetPendingResolveList(ctx) {
		request := k.MustGetRequest(ctx, reqID)
		coinRatesData := CoinRatesCallData{}
		ratesResult := CoinRatesResult{Rates: []uint64{10}}
		err := obi.Decode(request.Calldata, &coinRatesData)
		fmt.Println("handleEndBlock.Decode", err)
		if err == nil {
			ratesResult.Rates = []uint64{}
			for index := range coinRatesData.Symbols {
				ratesResult.Rates = append(ratesResult.Rates, (uint64(ctx.BlockTime().Unix())*uint64(index+1))/3432)
			}
		}

		result, err := obi.Encode(ratesResult)
		fmt.Println("handleEndBlock.Encode", err)
		k.ResolveSuccess(ctx, reqID, result, 1000)
	}

	// Once all the requests are resolved, we can clear the list.
	k.SetPendingResolveList(ctx, []types.RequestID{})
	// Lastly, we clean up data requests that are supposed to be expired.
	k.ProcessExpiredRequests(ctx)
	// NOTE: We can remove old requests from state to optimize space, using `k.DeleteRequest`
	// and `k.DeleteReports`. We don't do that now as it is premature optimization at this state.
}
