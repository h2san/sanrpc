// Code generated by protoc-gen-go. DO NOT EDIT.
// source: example.proto

package example

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Req struct {
	A                    int64    `protobuf:"varint,1,opt,name=a,proto3" json:"a,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Req) Reset()         { *m = Req{} }
func (m *Req) String() string { return proto.CompactTextString(m) }
func (*Req) ProtoMessage()    {}
func (*Req) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{0}
}

func (m *Req) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Req.Unmarshal(m, b)
}
func (m *Req) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Req.Marshal(b, m, deterministic)
}
func (m *Req) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Req.Merge(m, src)
}
func (m *Req) XXX_Size() int {
	return xxx_messageInfo_Req.Size(m)
}
func (m *Req) XXX_DiscardUnknown() {
	xxx_messageInfo_Req.DiscardUnknown(m)
}

var xxx_messageInfo_Req proto.InternalMessageInfo

func (m *Req) GetA() int64 {
	if m != nil {
		return m.A
	}
	return 0
}

type Resq struct {
	B                    int64    `protobuf:"varint,1,opt,name=b,proto3" json:"b,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Resq) Reset()         { *m = Resq{} }
func (m *Resq) String() string { return proto.CompactTextString(m) }
func (*Resq) ProtoMessage()    {}
func (*Resq) Descriptor() ([]byte, []int) {
	return fileDescriptor_15a1dc8d40dadaa6, []int{1}
}

func (m *Resq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Resq.Unmarshal(m, b)
}
func (m *Resq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Resq.Marshal(b, m, deterministic)
}
func (m *Resq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Resq.Merge(m, src)
}
func (m *Resq) XXX_Size() int {
	return xxx_messageInfo_Resq.Size(m)
}
func (m *Resq) XXX_DiscardUnknown() {
	xxx_messageInfo_Resq.DiscardUnknown(m)
}

var xxx_messageInfo_Resq proto.InternalMessageInfo

func (m *Resq) GetB() int64 {
	if m != nil {
		return m.B
	}
	return 0
}

func init() {
	proto.RegisterType((*Req)(nil), "example.Req")
	proto.RegisterType((*Resq)(nil), "example.Resq")
}

func init() { proto.RegisterFile("example.proto", fileDescriptor_15a1dc8d40dadaa6) }

var fileDescriptor_15a1dc8d40dadaa6 = []byte{
	// 82 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4d, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x84, 0xb9,
	0x98, 0x83, 0x52, 0x0b, 0x85, 0x78, 0xb8, 0x18, 0x13, 0x25, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83,
	0x18, 0x13, 0x95, 0x44, 0xb8, 0x58, 0x82, 0x52, 0x8b, 0xc1, 0xa2, 0x49, 0x30, 0xd1, 0xa4, 0x24,
	0x36, 0xb0, 0x56, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x0c, 0xd9, 0xd2, 0x3b, 0x4b, 0x00,
	0x00, 0x00,
}
