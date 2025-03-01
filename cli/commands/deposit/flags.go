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

const (
	// privateKey is the flag for the private key to sign the deposit message.
	privateKey = "private-key"

	// overrideNodeKey is the flag for overriding the node key.
	overrideNodeKey = "override-node-key"

	// validatorPrivateKey is the flag for the validator private key.
	valPrivateKey = "validator-private-key"
)

const (
	// overrideNodeKeyShorthand is the shorthand flag for the overrideNodeKey
	// flag.
	overrideNodeKeyShorthand = "o"
)

const (
	// defaultPrivateKey is the default value for the privateKey flag.
	defaultPrivateKey = ""

	// defaultOverrideNodeKey is the default value for the overrideNodeKey flag.
	defaultOverrideNodeKey = false

	// defaultValidatorPrivateKey is the default value for the
	// validatorPrivateKey flag.
	defaultValidatorPrivateKey = ""
)

const (
	// privateKeyFlagMsg is the usage description for the privateKey flag.
	privateKeyMsg = `private key to sign and pay for the deposit message. 
	This is required if the broadcast flag is set.`

	// overrideNodeKeyFlagMsg is the usage description for the overrideNodeKey
	// flag.
	overrideNodeKeyMsg = "override the node private key"

	// valPrivateKeyMsg is the usage description for the
	// valPrivateKey flag.
	valPrivateKeyMsg = `validator private key. This is required if the 
	override-node-key flag is set.`
)
