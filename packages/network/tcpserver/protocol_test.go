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

package tcpserver

import (
	"bytes"
	"fmt"

	"testing"

	"reflect"

	"github.com/AplaProject/go-apla/packages/converter"
	"github.com/stretchr/testify/require"
)

func TestReadRequest(t *testing.T) {
	type testStruct struct {
		Id   uint32
		Data []byte
	}

	request := &bytes.Buffer{}
	request.Write(converter.DecToBin(10, 4))
	request.Write(converter.DecToBin(len("test"), 4))
	request.Write([]byte("test"))

	test := &testStruct{}
	err := ReadRequest(test, request)
	if err != nil {
		t.Errorf("read request return err: %s", err)
	}
	if test.Id != 10 {
		t.Errorf("bad id value")
	}
	if string(test.Data) != "test" {
		t.Errorf("bad data value: %+v", string(test.Data))
	}
}

func TestReadRequestTag(t *testing.T) {
	type testStruct2 struct {
		Id   uint32
		Data []byte `size:"4"`
	}

	request := &bytes.Buffer{}
	request.Write(converter.DecToBin(10, 4))
	request.Write([]byte("test"))

	test := &testStruct2{}
	err := ReadRequest(test, request)
	if err != nil {
		t.Errorf("read request return err: %s", err)
	}
	if test.Id != 10 {
		t.Errorf("bad id value")
	}
	if string(test.Data) != "test" {
		t.Errorf("bad data value: %+v", string(test.Data))
	}
}

func TestSendRequest(t *testing.T) {
	type testStruct2 struct {
		Id   uint32
		Id2  int64
		Test []byte
		Text []byte `size:"4"`
	}

	test := testStruct2{
		Id:   15,
		Id2:  0x1BCDEF0010203040,
		Test: []byte("test"),
		Text: []byte("text"),
	}

	bin := bytes.Buffer{}
	err := SendRequest(&test, &bin)
	if err != nil {
		t.Fatalf("send request failed: %s", err)
	}

	test2 := testStruct2{}
	err = ReadRequest(&test2, &bin)
	if err != nil {
		t.Fatalf("read request failed: %s", err)
	}

	if !reflect.DeepEqual(test, test2) {
		t.Errorf("different values: %+v and %+v", test, test2)
	}
}

func TestRequestType(t *testing.T) {
	source := RequestType{
		Type: uint16(RequestTypeNotFullNode),
	}

	buf := bytes.Buffer{}
	require.NoError(t, source.Write(&buf))

	target := RequestType{}
	require.NoError(t, target.Read(&buf))
	require.Equal(t, source.Type, target.Type)
}

func TestGetBodiesRequest(t *testing.T) {
	source := GetBodiesRequest{
		BlockID:      33,
		ReverseOrder: true,
	}

	buf := bytes.Buffer{}
	require.NoError(t, source.Write(&buf))

	target := GetBodiesRequest{}
	require.NoError(t, target.Read(&buf))
	fmt.Printf("%+v %+v\n", source, target)
	require.Equal(t, source, target)
}
