package common

import (
	"errors"
	"math/big"
)

const (
	// Uint128Bytes defines the number of bytes for Uint128 type.
	Uint128Bytes = 16

	// Uint128Bits defines the number of bits for Uint128 type.
	Uint128Bits = 128
)

var (
	// ErrUint128Overflow indicates the value is greater than maximum value 2^128.
	ErrUint128Overflow = errors.New("uint128: overflow")

	// ErrUint128Underflow indicates the value is smaller than minimum value 0.
	ErrUint128Underflow = errors.New("uint128: underflow")

	// ErrUint128InvalidBytesSize indicates the bytes size is not equal to Uint128Bytes.
	ErrUint128InvalidBytesSize = errors.New("uint128: invalid bytes")

	// ErrUint128InvalidString indicates the string is not valid when converted to uint128.
	ErrUint128InvalidString = errors.New("uint128: invalid string to uint128")
)

// Uint128 defines uint128 type, based on big.Int.
type Uint128 struct {
	value *big.Int
}

// Validate returns error if u is not a valid uint128, otherwise returns nil.
func (u *Uint128) Validate() error {
	// Sign returns:
	//
	//	-1 if x <  0
	//	 0 if x == 0
	//	+1 if x >  0
	//
	if u.value.Sign() < 0 {
		return ErrUint128Underflow
	}

	// bitlen return length of bit
	if u.value.BitLen() > 128 {
		return ErrUint128Overflow
	}
	return nil
}

// NewUint128 returns a new Uint128 struct with default value.
func NewUint128() *Uint128 {
	return &Uint128{big.NewInt(0)}
}

// NewUint128FromString returns a new Uint128 struct with given value and have a check
func NewUint128FromString(str string) (*Uint128, error) {
	big := new(big.Int)
	_, success := big.SetString(str, 10)
	if !success {
		return nil, ErrUint128InvalidString
	}
	if err := (&Uint128{big}).Validate(); err != nil {
		return nil, err
	}
	return &Uint128{big}, nil
}

// String returns the string representation of x.
func (u *Uint128) String() string {
	return u.value.Text(10)
}

// ToFixedSizeBytes converts Uint128 to Big-Endian fixed size bytes.
func (u *Uint128) ToFixedSizeBytes() ([16]byte, error) {
	var res [16]byte
	if err := u.Validate(); err != nil {
		return res, err
	}
	bs := u.value.Bytes()
	l := len(bs)
	if l == 0 {
		return res, nil
	}
	idx := Uint128Bytes - len(bs)
	if idx < Uint128Bytes {
		copy(res[idx:], bs)
	}
	return res, nil
}

// ToFixedSizeByteSlice converts Uint128 to Big-Endian fixed size byte slice.
func (u *Uint128) ToFixedSizeByteSlice() ([]byte, error) {
	bytes, err := u.ToFixedSizeByteSlice()
	return bytes[:], err
}
