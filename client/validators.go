package rpc

// import (
// 	"context"
// 	"net/http"

// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/tendermint/tendermint/crypto/secp256k1"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/types/rest"
// )

// type ValidatorMinimal struct {
// 	Address     string `json:"address"`
// 	VotingPower int64  `json:"voting_power"`
// }

// type ValidatorsMinimal struct {
// 	BlockHeight int64              `json:"block_height"`
// 	Validators  []ValidatorMinimal `json:"validators"`
// }

// func GetEVMValidators(cliCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		node, err := cliCtx.GetNode()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		// Get top 100 validators for now
// 		var page int = 1
// 		var perPage int = 100
// 		validators, err := node.Validators(context.Background(), nil, &page, &perPage)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		validatorsMinimal := ValidatorsMinimal{}
// 		validatorsMinimal.BlockHeight = validators.BlockHeight
// 		validatorsMinimal.Validators = []ValidatorMinimal{}

// 		for _, validator := range validators.Validators {
// 			pubKeyBytes, ok := validator.PubKey.(secp256k1.PubKey)
// 			if !ok {
// 				rest.WriteErrorResponse(w, http.StatusInternalServerError, "fail to cast pubkey")
// 				return
// 			}

// 			if pubkey, err := crypto.DecompressPubkey(pubKeyBytes[:]); err != nil {
// 				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 				return
// 			} else {
// 				validatorsMinimal.Validators = append(
// 					validatorsMinimal.Validators,
// 					ValidatorMinimal{
// 						Address:     crypto.PubkeyToAddress(*pubkey).String(),
// 						VotingPower: validator.VotingPower,
// 					},
// 				)
// 			}
// 		}

// 		rest.PostProcessResponseBare(w, cliCtx, validatorsMinimal)
// 	}
// }
