// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package deposit

import (
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/cli/utils/parser"
	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

// Commands creates a new command for deposit related actions.
func Commands(
	chainSpec chain.Spec,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "deposit",
		Short:                      "deposit subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2, //nolint:mnd // from sdk.
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		NewValidateDeposit(chainSpec),
		NewCreateValidator(chainSpec),
	)

	return cmd
}

// NewValidateDeposit creates a new command for validating a deposit message.
//
//nolint:mnd,lll // lots of magic numbers, reads better if long description is one line
func NewValidateDeposit(chainSpec chain.Spec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validates a deposit message for creating a new validator",
		Long:  `Validates a deposit message for creating a new validator. The deposit message includes the public key, withdrawal credentials, and deposit amount. The args taken are in the order of the public key, withdrawal credentials, deposit amount, signature, and genesis validator root.`,
		Args:  cobra.ExactArgs(6),
		RunE:  validateDepositMessage(chainSpec),
	}

	return cmd
}

// validateDepositMessage validates a deposit message for creating a new
// validator.
func validateDepositMessage(chainSpec chain.Spec) func(
	_ *cobra.Command,
	args []string,
) error {
	return func(_ *cobra.Command, args []string) error {
		pubkey, err := parser.ConvertPubkey(args[0])
		if err != nil {
			return err
		}

		credentials, err := parser.ConvertWithdrawalCredentials(args[1])
		if err != nil {
			return err
		}

		amount, err := parser.ConvertAmount(args[2])
		if err != nil {
			return err
		}

		signature, err := parser.ConvertSignature(args[3])
		if err != nil {
			return err
		}

		genesisValidatorRoot, err := parser.ConvertGenesisValidatorRoot(args[4])
		if err != nil {
			return err
		}

		depositMessage := types.DepositMessage{
			Pubkey:      pubkey,
			Credentials: credentials,
			Amount:      amount,
		}

		// All deposits are signed with the genesis version.
		genesisVersion := version.FromUint32[common.Version](constants.GenesisVersion)

		return depositMessage.VerifyCreateValidator(
			types.NewForkData(genesisVersion, genesisValidatorRoot),
			signature,
			chainSpec.DomainTypeDeposit(),
			signer.BLSSigner{}.VerifySignature,
		)
	}
}
