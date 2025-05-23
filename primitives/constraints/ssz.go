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

package constraints

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/karalabe/ssz"
)

// SSZMarshaler is an interface for objects that can be
// marshaled to SSZ format.
type SSZMarshaler interface {
	// MarshalSSZ marshals the object into SSZ format.
	MarshalSSZ() ([]byte, error)
}

// SSZUnmarshaler is an interface for objects that can be unmarshaled from SSZ format.
type SSZUnmarshaler interface {
	ssz.Object
	ValidateAfterDecodingSSZ() error // once unmarshalled we will check whether type syntax is correct
}

// SSZMarshallable is an interface that combines SSZMarshaler and SSZUnmarshaler.
type SSZMarshallable interface {
	SSZMarshaler
	SSZUnmarshaler
}

// Versionable is a constraint that requires a type to have a GetForkVersion method.
type Versionable interface {
	GetForkVersion() common.Version
}

// SSZVersionable is an interface that combines SSZMarshallable and Versionable.
type SSZVersionedMarshallable interface {
	Versionable
	SSZMarshallable
}

// SSZRootable is an interface for objects that can compute their hash tree root.
type SSZRootable interface {
	// HashTreeRoot computes the hash tree root of the object.
	HashTreeRoot() common.Root
}

// SSZMarshallableRootable is an interface that combines
// sszMarshaler, sszUnmarshaler, and SSZRootable.
type SSZMarshallableRootable interface {
	SSZMarshallable
	SSZRootable
}

// SSZVersionedMarshallableRootable is an interface that combines
// SSZVersionedMarshallable and SSZRootable.
type SSZVersionedMarshallableRootable interface {
	SSZVersionedMarshallable
	SSZRootable
}
