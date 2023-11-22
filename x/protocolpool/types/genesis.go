package types

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewGenesisState(cf []*ContinuousFund, budget []*Budget) *GenesisState {
	return &GenesisState{
		ContinuousFund: cf,
		Budget:         budget,
	}
}

func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		ContinuousFund: []*ContinuousFund{},
		Budget:         []*Budget{},
	}
}

// ValidateGenesis validates the genesis state of protocolpool genesis input
func ValidateGenesis(gs *GenesisState) error {
	for _, cf := range gs.ContinuousFund {
		if err := validateContinuousFund(*cf); err != nil {
			return err
		}
	}
	for _, bp := range gs.Budget {
		if err := validateBudget(*bp); err != nil {
			return err
		}
	}
	return nil
}

func validateBudget(bp Budget) error {
	if bp.RecipientAddress == "" {
		return fmt.Errorf("recipient cannot be empty")
	}

	if bp.TotalBudget.IsZero() {
		return fmt.Errorf("invalid budget proposal: total budget cannot be zero")
	}

	amount := sdk.NewCoins(*bp.TotalBudget)
	if amount == nil {
		return fmt.Errorf("amount cannot be nil")
	}
	if err := amount.Validate(); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, amount.String())
	}

	if bp.Tranches == 0 {
		return fmt.Errorf("invalid budget proposal: tranches must be greater than zero")
	}

	if bp.Period == nil || *bp.Period == 0 {
		return fmt.Errorf("invalid budget proposal: period length should be greater than zero")
	}
	return nil
}

func validateContinuousFund(cf ContinuousFund) error {
	if cf.Recipient == "" {
		return fmt.Errorf("recipient cannot be empty")
	}

	cap := sdk.NewCoins(*cf.Cap)
	if cap == nil {
		return fmt.Errorf("amount cannot be nil")
	}
	if err := cap.Validate(); err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, cap.String())
	}

	// Validate percentage
	if cf.Percentage.IsZero() || cf.Percentage.IsNil() {
		return fmt.Errorf("percentage cannot be zero or empty")
	}
	if cf.Percentage.IsNegative() {
		return fmt.Errorf("percentage cannot be negative")
	}
	if cf.Percentage.GT(math.LegacyOneDec()) {
		return fmt.Errorf("percentage cannot be greater than one")
	}
	return nil
}
