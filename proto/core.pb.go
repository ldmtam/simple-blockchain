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
	Chainid              uint32   `protobuf:"varint,2,opt,name=chainid,proto3" json:"chainid,omitempty"`
	From                 []byte   `protobuf:"bytes,3,opt,name=from,proto3" json:"from,omitempty"`
	To                   []byte   `protobuf:"bytes,4,opt,name=to,proto3" json:"to,omitempty"`
	Value                []byte   `protobuf:"bytes,5,opt,name=value,proto3" json:"value,omitempty"`
	Fee                  []byte   `protobuf:"bytes,6,opt,name=fee,proto3" json:"fee,omitempty"`
	Nonce                uint64   `protobuf:"varint,7,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Timestamp            int64    `protobuf:"varint,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Signature            []byte   `protobuf:"bytes,9,opt,name=signature,proto3" json:"signature,omitempty"`
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

func (m *Transaction) GetChainid() uint32 {
	if m != nil {
		return m.Chainid
	}
	return 0
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

func (m *Transaction) GetFee() []byte {
	if m != nil {
		return m.Fee
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
	// 232 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x90, 0x4d, 0x6e, 0x84, 0x30,
	0x0c, 0x85, 0x15, 0x60, 0xa0, 0xe3, 0xfe, 0xa8, 0x8a, 0xba, 0xf0, 0xa2, 0x0b, 0x34, 0x2b, 0x56,
	0xdd, 0xf4, 0x04, 0xbd, 0x02, 0xed, 0x05, 0x4c, 0xc8, 0x94, 0x48, 0x43, 0x8c, 0x92, 0xd0, 0xbb,
	0xf6, 0x36, 0x95, 0xc3, 0x8c, 0xe8, 0xee, 0xbd, 0xef, 0x3d, 0x25, 0xb6, 0x01, 0x0c, 0x07, 0xfb,
	0xb6, 0x04, 0x4e, 0xac, 0x6b, 0xd1, 0xcb, 0x70, 0xfa, 0x55, 0x70, 0xff, 0x15, 0xc8, 0x47, 0x32,
	0xc9, 0xb1, 0xd7, 0x1a, 0xaa, 0x89, 0xe2, 0x84, 0xaa, 0x55, 0xdd, 0x43, 0x9f, 0xb5, 0x46, 0x68,
	0xcc, 0x44, 0xce, 0xbb, 0x11, 0x8b, 0x56, 0x75, 0x8f, 0xfd, 0xcd, 0x4a, 0xfb, 0x1c, 0x78, 0xc6,
	0x72, 0x6b, 0x8b, 0xd6, 0x4f, 0x50, 0x24, 0xc6, 0x2a, 0x93, 0x22, 0xb1, 0x7e, 0x81, 0xc3, 0x0f,
	0x5d, 0x56, 0x8b, 0x87, 0x8c, 0x36, 0xa3, 0x9f, 0xa1, 0x3c, 0x5b, 0x8b, 0x75, 0x66, 0x22, 0xa5,
	0xe7, 0xd9, 0x1b, 0x8b, 0x4d, 0xab, 0xba, 0xaa, 0xdf, 0x8c, 0x7e, 0x85, 0x63, 0x72, 0xb3, 0x8d,
	0x89, 0xe6, 0x05, 0xef, 0x5a, 0xd5, 0x95, 0xfd, 0x0e, 0x24, 0x8d, 0xee, 0xdb, 0x53, 0x5a, 0x83,
	0xc5, 0x63, 0x7e, 0x6b, 0x07, 0xa7, 0x4f, 0x68, 0x3e, 0x8c, 0xe1, 0xd5, 0x27, 0x59, 0x81, 0xc6,
	0x31, 0xd8, 0x18, 0xaf, 0x9b, 0xdd, 0xac, 0x24, 0x03, 0x5d, 0x48, 0x3e, 0x2e, 0xb6, 0xe4, 0x6a,
	0xf7, 0x81, 0xca, 0x7f, 0x03, 0x0d, 0x75, 0xbe, 0xdf, 0xfb, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x17, 0xb5, 0x7a, 0x88, 0x4d, 0x01, 0x00, 0x00,
}
