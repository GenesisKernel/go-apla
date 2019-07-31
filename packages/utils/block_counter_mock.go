// Copyright (C) 2017, 2018, 2019 EGAAS S.A.
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or (at
// your option) any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.

package utils

import mock "github.com/stretchr/testify/mock"

// mockIntervalBlocksCounter is an autogenerated mock type for the intervalBlocksCounter type
type mockIntervalBlocksCounter struct {
	mock.Mock
}

// count provides a mock function with given fields: state
func (_m *mockIntervalBlocksCounter) count(state blockGenerationState) (int, error) {
	ret := _m.Called(state)

	var r0 int
	if rf, ok := ret.Get(0).(func(blockGenerationState) int); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(blockGenerationState) error); ok {
		r1 = rf(state)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
