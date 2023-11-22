package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/x/protocolpool/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx context.Context, data *types.GenesisState) error {
	for _, cf := range data.ContinuousFund {
		recipientAddress, err := k.authKeeper.AddressCodec().StringToBytes(cf.Recipient)
		if err != nil {
			return fmt.Errorf("failed to decode recipient address: %w", err)
		}
		err = k.ContinuousFund.Set(ctx, recipientAddress, *cf)
		if err != nil {
			return fmt.Errorf("failed to set continuous fund: %w", err)
		}
	}
	for _, budget := range data.Budget {
		recipientAddress, err := k.authKeeper.AddressCodec().StringToBytes(budget.RecipientAddress)
		if err != nil {
			return fmt.Errorf("failed to decode recipient address: %w", err)
		}
		err = k.BudgetProposal.Set(ctx, recipientAddress, *budget)
	}
	return nil
}

func (k Keeper) ExportGenesis(ctx context.Context) *types.GenesisState {
	var cf []*types.ContinuousFund
	err := k.ContinuousFund.Walk(ctx, nil, func(key sdk.AccAddress, value types.ContinuousFund) (stop bool, err error) {
		cf = append(cf, &types.ContinuousFund{
			Recipient:  key.String(),
			Percentage: value.Percentage,
			Cap:        value.Cap,
			Expiry:     value.Expiry,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	var budget []*types.Budget
	err = k.BudgetProposal.Walk(ctx, nil, func(key sdk.AccAddress, value types.Budget) (stop bool, err error) {
		budget = append(budget, &types.Budget{
			RecipientAddress: key.String(),
			TotalBudget:      value.TotalBudget,
			ClaimedAmount:    value.ClaimedAmount,
			StartTime:        value.StartTime,
			NextClaimFrom:    value.NextClaimFrom,
			Tranches:         value.Tranches,
			TranchesLeft:     value.TranchesLeft,
			Period:           value.Period,
		})
		return false, nil
	})
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(cf, budget)
}
