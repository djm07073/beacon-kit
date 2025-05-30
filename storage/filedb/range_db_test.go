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

package filedb_test

import (
	"reflect"
	"testing"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/errors"
	file "github.com/berachain/beacon-kit/storage/filedb"
	"github.com/berachain/beacon-kit/storage/interfaces/mocks"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// =========================== BASIC OPERATIONS ============================

func TestRangeDB(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupFunc     func(rdb *file.RangeDB) error
		testFunc      func(t *testing.T, rdb *file.RangeDB)
		expectedError bool
	}{
		{
			name: "Get",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				gotValue, err := rdb.Get(1, []byte("testKey"))
				require.NoError(t, err)
				require.Equal(t, []byte("testValue"), gotValue)
			},
		},
		{
			name: "Has",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.True(t, exists)
			},
		},
		{
			name: "Set",
			setupFunc: func(_ *file.RangeDB) error {
				return nil // No setup required
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				err := rdb.Set(1, []byte("testKey"), []byte("testValue"))
				require.NoError(t, err)

				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.True(t, exists)
			},
		},
		{
			name: "Delete",
			setupFunc: func(rdb *file.RangeDB) error {
				return rdb.Set(1, []byte("testKey"), []byte("testValue"))
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()

				exists, err := rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.True(t, exists)

				err = rdb.Delete(1, []byte("testKey"))
				require.NoError(t, err)

				exists, err = rdb.Has(1, []byte("testKey"))
				require.NoError(t, err)
				require.False(t, exists)
			},
		},
		{
			name: "Prune",
			setupFunc: func(rdb *file.RangeDB) error {
				for index := uint64(1); index <= 5; index++ {
					if err := rdb.Set(
						index, []byte("testKey"), []byte("testValue"),
					); err != nil {
						return err
					}
				}
				return nil
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				err := rdb.Prune(1, 4)
				require.NoError(t, err)

				for index := uint64(1); index <= 3; index++ {
					var exists bool
					exists, err = rdb.Has(index, []byte("testKey"))
					require.NoError(t, err, "index %d", index)
					require.False(t, exists, "index %d", index)
				}

				for index := uint64(4); index <= 5; index++ {
					var exists bool
					exists, err = rdb.Has(index, []byte("testKey"))
					require.NoError(t, err, "index %d", index)
					require.True(t, exists, "index %d", index)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rdb := file.NewRangeDB(newTestFDB("/tmp/testdb-1"))

			if tt.setupFunc != nil {
				require.NoError(t, tt.setupFunc(rdb))
			}

			tt.testFunc(t, rdb)
		})
	}
}

func TestExtractIndex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		prefixedKey []byte
		expectedIdx uint64
		expectedErr error
	}{
		{
			name:        "ValidKey",
			prefixedKey: []byte("12345/testKey"),
			expectedIdx: 12345,
			expectedErr: nil,
		},
		{
			name:        "InvalidKeyFormat",
			prefixedKey: []byte("testKey"),
			expectedIdx: 0,
			expectedErr: errors.New("invalid key format"),
		},
		{
			name:        "InvalidIndex",
			prefixedKey: []byte("abc/testKey"),
			expectedIdx: 0,
			expectedErr: errors.New(
				"strconv.ParseUint: parsing \"abc\": invalid syntax",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			idx, err := file.ExtractIndex(tt.prefixedKey)
			require.Equal(t, tt.expectedIdx, idx)
			if tt.expectedErr != nil {
				require.ErrorContains(t, err, tt.expectedErr.Error())
			}
		})
	}
}

// =========================== PRUNING =====================================

func TestRangeDB_DeleteRange_NotSupported(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		db   *mocks.DB
	}{
		{
			name: "DeleteRangeNotSupported",
			db:   new(mocks.DB),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			t.Helper()
			require.Panics(t, func() { _ = file.NewRangeDB(tt.db) })
		})
	}
}

func TestRangeDB_Prune(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		setupFunc     func(rdb *file.RangeDB) error
		start         uint64
		end           uint64
		expectedError bool
		testFunc      func(t *testing.T, rdb *file.RangeDB)
	}{
		{
			name: "PruneWithDeleteRange",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 0, 50)
			},
			start:         2,
			end:           7,
			expectedError: false,
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				requireNotExist(t, rdb, 2, 6)
				requireExist(t, rdb, 0, 1)
				requireExist(t, rdb, 7, 50)
			},
		},
		{
			name: "PruneWithDeleteRange-InvalidRange",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 0, 50)
			},
			start:         7,
			end:           2,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rdb := file.NewRangeDB(newTestFDB("/tmp/testdb-2"))

			if tt.setupFunc != nil {
				require.NoError(t, tt.setupFunc(rdb))
			}
			err := rdb.Prune(tt.start, tt.end)
			if (err != nil) != tt.expectedError {
				t.Fatalf(
					"Prune() error = %v, expectedError %v",
					err,
					tt.expectedError,
				)
			}

			if tt.testFunc != nil {
				tt.testFunc(t, rdb)
			}
		})
	}
}

// =========================== INVARIANTS ================================.

// invariant: all indexes up to the firstNonNilIndex should be nil.
func TestRangeDB_Invariants(t *testing.T) {
	t.Parallel()
	// we ignore errors for most of the tests below because we want to ensure
	// that the invariants hold in exceptional circumstances.
	tests := []struct {
		name      string
		setupFunc func(rdb *file.RangeDB) error
		testFunc  func(t *testing.T, rdb *file.RangeDB)
	}{
		{
			name: "Populate from empty",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 1, 5)
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				requireNotExist(t, rdb, 0, lastConsequetiveNilIndex(rdb))
			},
		},
		{
			name: "Delete from populated",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 1, 5)
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				_ = rdb.Delete(2, []byte("key"))
				requireNotExist(t, rdb, 0, lastConsequetiveNilIndex(rdb))
			},
		},
		{
			name: "Prune from populated",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 1, 10)
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				_ = rdb.Prune(0, 3)
				requireNotExist(t, rdb, 0, lastConsequetiveNilIndex(rdb))
			},
		},
		{
			name: "Populate, Prune, Set round trip",
			setupFunc: func(rdb *file.RangeDB) error {
				return populateTestDB(rdb, 1, 30)
			},
			testFunc: func(t *testing.T, rdb *file.RangeDB) {
				t.Helper()
				if err := rdb.Prune(0, 25); err != nil {
					t.Fatalf("Prune() error = %v", err)
				}
				_ = populateTestDB(rdb, 5, 10)
				requireNotExist(t, rdb, 0, lastConsequetiveNilIndex(rdb))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rdb := file.NewRangeDB(newTestFDB("/tmp/testdb-3"))

			if tt.setupFunc != nil {
				if err := tt.setupFunc(rdb); err != nil {
					requireNotExist(
						t,
						rdb,
						0,
						lastConsequetiveNilIndex(rdb),
					)
				}
			}
			if tt.testFunc != nil {
				tt.testFunc(t, rdb)
			}
		})
	}
}

// =============================== HELPERS ==================================

// newTestFDB returns a new file DB instance with an in-memory filesystem.
func newTestFDB(path string) *file.DB {
	fs := afero.NewMemMapFs()
	return file.NewDB(
		// don't reuse the same txt file for consecutive unit tests bc file
		// db slow AF
		file.WithRootDirectory(path),
		file.WithFileExtension("txt"),
		file.WithDirectoryPermissions(0700),
		file.WithLogger(log.NewNopLogger()),
		file.WithAferoFS(fs),
	)
}

func getFirstNonNilIndex(rdb *file.RangeDB) uint64 {
	return reflect.ValueOf(rdb).Elem().FieldByName("lowerBoundIndex").Uint()
}

func lastConsequetiveNilIndex(rdb *file.RangeDB) uint64 {
	return uint64(max(int64(getFirstNonNilIndex(rdb))-1, 0))
}

// requireNotExist requires the indexes from `from` to `to` to be empty.
//

func requireNotExist(t *testing.T, rdb *file.RangeDB, from uint64, to uint64) {
	t.Helper()
	for i := from; i <= to; i++ {
		exists, err := rdb.Has(i, []byte("key"))
		require.NoError(t, err)
		require.False(t, exists, "Index %d should have been pruned", i)
	}
}

// requireExist requires the indexes from `from` to `to` not be empty.
func requireExist(t *testing.T, rdb *file.RangeDB, from uint64, to uint64) {
	t.Helper()
	for i := from; i <= to; i++ {
		exists, err := rdb.Has(i, []byte("key"))
		require.NoError(t, err)
		require.True(t, exists, "Index %d should not have been pruned", i)
	}
}

// populateTestDB populates the test DB with indexes from `from` to `to`.
func populateTestDB(rdb *file.RangeDB, from, to uint64) error {
	for i := from; i <= to; i++ {
		if err := rdb.Set(i, []byte("key"), []byte("value")); err != nil {
			return err
		}
	}
	return nil
}
