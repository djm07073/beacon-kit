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

package constants

const (
	MaxProposerSlashings     = 16
	MaxAttesterSlashings     = 2
	MaxAttestations          = 128
	MaxVoluntaryExits        = 16
	MaxBlsToExecutionChanges = 16

	// MaxTxsPerPayload is the maximum number of transactions in a execution payload.
	MaxTxsPerPayload uint64 = 1048576

	// MaxWithdrawalsPerPayload is the maximum number of withdrawals in a execution payload.
	MaxWithdrawalsPerPayload uint64 = 16

	// MaxWithdrawalRequestsPerPayload is the maximum number of withdrawal requests in a execution
	// payload.
	MaxWithdrawalRequestsPerPayload = 16

	// MaxConsolidationRequestsPerPayload is the maximum number of consolidation requests in a
	// execution payload.
	MaxConsolidationRequestsPerPayload = 2

	// MaxDepositRequestsPerPayload is the maximum number of deposit requests in a execution payload.
	MaxDepositRequestsPerPayload = 8192
)
