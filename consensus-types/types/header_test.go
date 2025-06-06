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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	karalabessz "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockHeader_Equals(t *testing.T) {
	t.Parallel()
	var (
		slot            = math.Slot(100)
		valIdx          = math.ValidatorIndex(200)
		parentBlockRoot = common.Root{1}
		stateRoot       = common.Root{2}
		bodyRoot        = common.Root{3}

		lhs = types.NewBeaconBlockHeader(
			slot, valIdx, parentBlockRoot, stateRoot, bodyRoot,
		)
	)

	tests := []struct {
		name string
		rhs  *types.BeaconBlockHeader
		want bool
	}{
		{
			name: "equal",
			rhs: types.NewBeaconBlockHeader(
				slot, valIdx, parentBlockRoot, stateRoot, bodyRoot,
			),
			want: true,
		},
		{
			name: "slot differs",
			rhs: types.NewBeaconBlockHeader(
				2*slot, valIdx, parentBlockRoot, stateRoot, bodyRoot,
			),
			want: false,
		},
		{
			name: "validator index differs",
			rhs: types.NewBeaconBlockHeader(
				slot, 2*valIdx, parentBlockRoot, stateRoot, bodyRoot,
			),
			want: false,
		},
		{
			name: "parent block root differs",
			rhs: types.NewBeaconBlockHeader(
				slot, valIdx, common.Root{0xff}, stateRoot, bodyRoot,
			),
			want: false,
		},
		{
			name: "state root differs",
			rhs: types.NewBeaconBlockHeader(
				slot, valIdx, parentBlockRoot, common.Root{0xff}, bodyRoot,
			),
			want: false,
		},
		{
			name: "body root differs",
			rhs: types.NewBeaconBlockHeader(
				slot, valIdx, parentBlockRoot, stateRoot, common.Root{0xff},
			),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got1 := lhs.Equals(tt.rhs)
			require.Equal(t, tt.want, got1)

			// check commutativity as well
			got2 := tt.rhs.Equals(lhs)
			require.Equal(t, got1, got2)

			// copies stays equal/disequal
			rhsCopy := &types.BeaconBlockHeader{}
			*rhsCopy = *tt.rhs
			got3 := rhsCopy.Equals(lhs)
			require.Equal(t, got1, got3)
		})
	}
}

func TestBeaconBlockHeader_Serialization(t *testing.T) {
	t.Parallel()
	original := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	unmarshalled := new(types.BeaconBlockHeader)
	err = ssz.Unmarshal(data, unmarshalled)
	require.NoError(t, err)
	require.Equal(t, original, unmarshalled)

	var buf []byte
	buf, err = original.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, data, buf)
}

func TestBeaconBlockHeader_SizeSSZ(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	size := karalabessz.Size(header)
	require.Equal(t, uint32(112), size)
}

func TestBeaconBlockHeader_HashTreeRoot(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	_ = header.HashTreeRoot()
}

func TestBeaconBlockHeader_GetTree(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	tree, err := header.GetTree()

	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestBeaconBlockHeader_SetStateRoot(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newStateRoot := common.Root{}
	header.SetStateRoot(newStateRoot)

	require.Equal(t, newStateRoot, header.GetStateRoot())
}

func TestBeaconBlockHeader_SetSlot(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newSlot := math.Slot(101)
	header.SetSlot(newSlot)

	require.Equal(t, newSlot, header.GetSlot())
}

func TestBeaconBlockHeader_SetProposerIndex(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newProposerIndex := math.ValidatorIndex(201)
	header.SetProposerIndex(newProposerIndex)
	require.Equal(t, newProposerIndex, header.GetProposerIndex())
}

func TestBeaconBlockHeader_SetParentBlockRoot(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newParentBlockRoot := common.Root{}
	header.SetParentBlockRoot(newParentBlockRoot)

	require.Equal(t, newParentBlockRoot, header.GetParentBlockRoot())
}

func TestBeaconBlockHeader_SetBodyRoot(t *testing.T) {
	t.Parallel()
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newBodyRoot := common.Root{}
	header.SetBodyRoot(newBodyRoot)

	require.Equal(t, newBodyRoot, header.GetBodyRoot())
}

func TestBeaconBlockHeader_UnmarshalSSZ_ErrSize(t *testing.T) {
	t.Parallel()
	buf := make([]byte, 100) // Incorrect size

	unmarshalled := new(types.BeaconBlockHeader)
	err := ssz.Unmarshal(buf, unmarshalled)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
