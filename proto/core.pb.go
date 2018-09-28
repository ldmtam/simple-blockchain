// Code generated by protoc-gen-go. DO NOT EDIT.
// source: core.proto

package corepb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Transaction struct {
	Hash                 []byte   `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
	From                 []byte   `protobuf:"bytes,2,opt,name=from,proto3" json:"from,omitempty"`
	To                   []byte   `protobuf:"bytes,3,opt,name=to,proto3" json:"to,omitempty"`
	Value                []byte   `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	Nonce                uint64   `protobuf:"varint,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Timestamp            int64    `protobuf:"varint,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature            []byte   `protobuf:"bytes,7,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Transaction) Reset()         { *m = Transaction{} }
func (m *Transaction) String() string { return proto.CompactTextString(m) }
func (*Transaction) ProtoMessage()    {}
func (*Transaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{0}
}

func (m *Transaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Transaction.Unmarshal(m, b)
}
func (m *Transaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Transaction.Marshal(b, m, deterministic)
}
func (m *Transaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Transaction.Merge(m, src)
}
func (m *Transaction) XXX_Size() int {
	return xxx_messageInfo_Transaction.Size(m)
}
func (m *Transaction) XXX_DiscardUnknown() {
	xxx_messageInfo_Transaction.DiscardUnknown(m)
}

var xxx_messageInfo_Transaction proto.InternalMessageInfo

func (m *Transaction) GetHash() []byte {
	if m != nil {
		return m.Hash
	}
	return nil
}

func (m *Transaction) GetFrom() []byte {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *Transaction) GetTo() []byte {
	if m != nil {
		return m.To
	}
	return nil
}

func (m *Transaction) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Transaction) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *Transaction) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *Transaction) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

type Account struct {
	Address              []byte   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Balance              []byte   `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Nonce                uint64   `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7e43720d1edc0fe, []int{1}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

func (m *Account) GetBalance() []byte {
	if m != nil {
		return m.Balance
	}
	return nil
}

func (m *Account) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func init() {
	proto.RegisterType((*Transaction)(nil), "corepb.Transaction")
	proto.RegisterType((*Account)(nil), "corepb.Account")
}

func init() { proto.RegisterFile("core.proto", fileDescriptor_f7e43720d1edc0fe) }

var fileDescriptor_f7e43720d1edc0fe = []byte{
	// 205 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x41, 0x4e, 0x86, 0x30,
	0x10, 0x85, 0x53, 0xe0, 0x87, 0x38, 0x1a, 0x17, 0x8d, 0x8b, 0x59, 0xb8, 0x20, 0xac, 0x58, 0xb9,
	0xf1, 0x04, 0x5e, 0x01, 0xbd, 0xc0, 0x50, 0xaa, 0x90, 0x40, 0x87, 0xb4, 0x83, 0xb7, 0xf2, 0x8e,
	0xa6, 0xad, 0x04, 0x77, 0xef, 0x7d, 0x2f, 0x69, 0xbf, 0x0c, 0x80, 0x61, 0x6f, 0x5f, 0x76, 0xcf,
	0xc2, 0xba, 0x8e, 0x79, 0x1f, 0xbb, 0x1f, 0x05, 0xf7, 0x1f, 0x9e, 0x5c, 0x20, 0x23, 0x0b, 0x3b,
	0xad, 0xa1, 0x9a, 0x29, 0xcc, 0xa8, 0x5a, 0xd5, 0x3f, 0x0c, 0x29, 0x47, 0xf6, 0xe9, 0x79, 0xc3,
	0x22, 0xb3, 0x98, 0xf5, 0x23, 0x14, 0xc2, 0x58, 0x26, 0x52, 0x08, 0xeb, 0x27, 0xb8, 0x7d, 0xd3,
	0x7a, 0x58, 0xac, 0x12, 0xca, 0x25, 0x52, 0xc7, 0xce, 0x58, 0xbc, 0xb5, 0xaa, 0xaf, 0x86, 0x5c,
	0xf4, 0x33, 0xdc, 0xc9, 0xb2, 0xd9, 0x20, 0xb4, 0xed, 0x58, 0xb7, 0xaa, 0x2f, 0x87, 0x0b, 0xc4,
	0x35, 0x2c, 0x5f, 0x8e, 0xe4, 0xf0, 0x16, 0x9b, 0xf4, 0xda, 0x05, 0xba, 0x77, 0x68, 0xde, 0x8c,
	0xe1, 0xc3, 0x89, 0x46, 0x68, 0x68, 0x9a, 0xbc, 0x0d, 0xe1, 0xcf, 0xf6, 0xac, 0x71, 0x19, 0x69,
	0xa5, 0xf8, 0x71, 0x76, 0x3e, 0xeb, 0x25, 0x54, 0xfe, 0x13, 0x1a, 0xeb, 0x74, 0x93, 0xd7, 0xdf,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x4d, 0x8e, 0xd8, 0x9e, 0x21, 0x01, 0x00, 0x00,
}
