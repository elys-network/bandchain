package rpc

// import (
// 	"net/http"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/types/rest"
// )

// func GetChainIDFn(cliCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		rest.PostProcessResponseBare(w, cliCtx, map[string]string{"chain_id": cliCtx.ChainID})
// 	}
// }

// func GetGenesisHandlerFn(cliCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		node, err := cliCtx.GetNode()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		genesis, err := node.Genesis(nil)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}
// 		rest.PostProcessResponseBare(w, cliCtx, genesis.Genesis)
// 	}
// }
