// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package evm

import (
	"fmt"

	. "github.com/hyperledger/burrow/binary"
)

// Not goroutine safe
type Stack struct {
	data []Word256
	ptr  int

	gas *uint64
	err *error
}

func NewStack(capacity int, gas *uint64, err *error) *Stack {
	return &Stack{
		data: make([]Word256, capacity),
		ptr:  0,
		gas:  gas,
		err:  err,
	}
}

func (st *Stack) useGas(gasToUse uint64) {
	if *st.gas > gasToUse {
		*st.gas -= gasToUse
	} else {
		st.setErr(ErrInsufficientGas)
	}
}

func (st *Stack) setErr(err error) {
	if *st.err == nil {
		*st.err = err
	}
}

func (st *Stack) Push(d Word256) {
	st.useGas(GasStackOp)
	if st.ptr == cap(st.data) {
		st.setErr(ErrDataStackOverflow)
		return
	}
	st.data[st.ptr] = d
	st.ptr++
}

// currently only called after sha3.Sha3
func (st *Stack) PushBytes(bz []byte) {
	if len(bz) != 32 {
		panic("Invalid bytes size: expected 32")
	}
	st.Push(LeftPadWord256(bz))
}

func (st *Stack) Push64(i int64) {
	st.Push(Int64ToWord256(i))
}

func (st *Stack) PushU64(i uint64) {
	st.Push(Uint64ToWord256(i))
}

func (st *Stack) Pop() Word256 {
	st.useGas(GasStackOp)
	if st.ptr == 0 {
		st.setErr(ErrDataStackUnderflow)
		return Zero256
	}
	st.ptr--
	return st.data[st.ptr]
}

func (st *Stack) PopBytes() []byte {
	return st.Pop().Bytes()
}

func (st *Stack) Pop64() int64 {
	d := st.Pop()
	return Int64FromWord256(d)
}

func (st *Stack) PopU64() uint64 {
	d := st.Pop()
	return Uint64FromWord256(d)
}

func (st *Stack) Len() int {
	return st.ptr
}

func (st *Stack) Swap(n int) {
	st.useGas(GasStackOp)
	if st.ptr < n {
		st.setErr(ErrDataStackUnderflow)
		return
	}
	st.data[st.ptr-n], st.data[st.ptr-1] = st.data[st.ptr-1], st.data[st.ptr-n]
}

func (st *Stack) Dup(n int) {
	st.useGas(GasStackOp)
	if st.ptr < n {
		st.setErr(ErrDataStackUnderflow)
		return
	}
	st.Push(st.data[st.ptr-n])
}

// Not an opcode, costs no gas.
func (st *Stack) Peek() Word256 {
	if st.ptr == 0 {
		st.setErr(ErrDataStackUnderflow)
		return Zero256
	}
	return st.data[st.ptr-1]
}

func (st *Stack) Print(n int) {
	fmt.Println("### stack ###")
	if st.ptr > 0 {
		nn := n
		if st.ptr < n {
			nn = st.ptr
		}
		for j, i := 0, st.ptr-1; i > st.ptr-1-nn; i-- {
			fmt.Printf("%-3d  %X\n", j, st.data[i])
			j += 1
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
