// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tictac.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type TicRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TicRequest) Reset()         { *m = TicRequest{} }
func (m *TicRequest) String() string { return proto.CompactTextString(m) }
func (*TicRequest) ProtoMessage()    {}
func (*TicRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcd35ae87bd1f440, []int{0}
}

func (m *TicRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TicRequest.Unmarshal(m, b)
}
func (m *TicRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TicRequest.Marshal(b, m, deterministic)
}
func (m *TicRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TicRequest.Merge(m, src)
}
func (m *TicRequest) XXX_Size() int {
	return xxx_messageInfo_TicRequest.Size(m)
}
func (m *TicRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TicRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TicRequest proto.InternalMessageInfo

type TicResponse struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TicResponse) Reset()         { *m = TicResponse{} }
func (m *TicResponse) String() string { return proto.CompactTextString(m) }
func (*TicResponse) ProtoMessage()    {}
func (*TicResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcd35ae87bd1f440, []int{1}
}

func (m *TicResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TicResponse.Unmarshal(m, b)
}
func (m *TicResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TicResponse.Marshal(b, m, deterministic)
}
func (m *TicResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TicResponse.Merge(m, src)
}
func (m *TicResponse) XXX_Size() int {
	return xxx_messageInfo_TicResponse.Size(m)
}
func (m *TicResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TicResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TicResponse proto.InternalMessageInfo

func (m *TicResponse) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type TacRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TacRequest) Reset()         { *m = TacRequest{} }
func (m *TacRequest) String() string { return proto.CompactTextString(m) }
func (*TacRequest) ProtoMessage()    {}
func (*TacRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcd35ae87bd1f440, []int{2}
}

func (m *TacRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TacRequest.Unmarshal(m, b)
}
func (m *TacRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TacRequest.Marshal(b, m, deterministic)
}
func (m *TacRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TacRequest.Merge(m, src)
}
func (m *TacRequest) XXX_Size() int {
	return xxx_messageInfo_TacRequest.Size(m)
}
func (m *TacRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_TacRequest.DiscardUnknown(m)
}

var xxx_messageInfo_TacRequest proto.InternalMessageInfo

type TacResponse struct {
	Res                  int64    `protobuf:"varint,1,opt,name=res,proto3" json:"res,omitempty"`
	Err                  string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TacResponse) Reset()         { *m = TacResponse{} }
func (m *TacResponse) String() string { return proto.CompactTextString(m) }
func (*TacResponse) ProtoMessage()    {}
func (*TacResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_dcd35ae87bd1f440, []int{3}
}

func (m *TacResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TacResponse.Unmarshal(m, b)
}
func (m *TacResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TacResponse.Marshal(b, m, deterministic)
}
func (m *TacResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TacResponse.Merge(m, src)
}
func (m *TacResponse) XXX_Size() int {
	return xxx_messageInfo_TacResponse.Size(m)
}
func (m *TacResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_TacResponse.DiscardUnknown(m)
}

var xxx_messageInfo_TacResponse proto.InternalMessageInfo

func (m *TacResponse) GetRes() int64 {
	if m != nil {
		return m.Res
	}
	return 0
}

func (m *TacResponse) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*TicRequest)(nil), "pb.TicRequest")
	proto.RegisterType((*TicResponse)(nil), "pb.TicResponse")
	proto.RegisterType((*TacRequest)(nil), "pb.TacRequest")
	proto.RegisterType((*TacResponse)(nil), "pb.TacResponse")
}

func init() { proto.RegisterFile("tictac.proto", fileDescriptor_dcd35ae87bd1f440) }

var fileDescriptor_dcd35ae87bd1f440 = []byte{
	// 153 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0xc9, 0x4c, 0x2e,
	0x49, 0x4c, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xe2, 0xe1, 0xe2,
	0x0a, 0xc9, 0x4c, 0x0e, 0x4a, 0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x51, 0x92, 0xe7, 0xe2, 0x06, 0xf3,
	0x8a, 0x0b, 0xf2, 0xf3, 0x8a, 0x53, 0x85, 0x04, 0xb8, 0x98, 0x53, 0x8b, 0x8a, 0x24, 0x18, 0x15,
	0x18, 0x35, 0x38, 0x83, 0x40, 0x4c, 0xb0, 0xf2, 0x44, 0xb8, 0x72, 0x43, 0x2e, 0x6e, 0x30, 0x0f,
	0xa1, 0xbc, 0x28, 0xb5, 0x18, 0xac, 0x9c, 0x39, 0x08, 0xc4, 0x84, 0x19, 0xc0, 0x04, 0x37, 0xc0,
	0x28, 0x82, 0x8b, 0x2d, 0x04, 0xec, 0x06, 0x21, 0x35, 0x2e, 0xe6, 0x90, 0xcc, 0x64, 0x21, 0x3e,
	0xbd, 0x82, 0x24, 0x3d, 0x84, 0x13, 0xa4, 0xf8, 0xe1, 0x7c, 0xa8, 0xa9, 0x20, 0x75, 0x89, 0x30,
	0x75, 0x89, 0x68, 0xea, 0x10, 0xb6, 0x27, 0xb1, 0x81, 0x3d, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff,
	0xff, 0x52, 0xb6, 0xee, 0x7a, 0xe4, 0x00, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TictacClient is the client API for Tictac service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TictacClient interface {
	Tic(ctx context.Context, in *TicRequest, opts ...grpc.CallOption) (*TicResponse, error)
	Tac(ctx context.Context, in *TacRequest, opts ...grpc.CallOption) (*TacResponse, error)
}

type tictacClient struct {
	cc *grpc.ClientConn
}

func NewTictacClient(cc *grpc.ClientConn) TictacClient {
	return &tictacClient{cc}
}

func (c *tictacClient) Tic(ctx context.Context, in *TicRequest, opts ...grpc.CallOption) (*TicResponse, error) {
	out := new(TicResponse)
	err := c.cc.Invoke(ctx, "/pb.Tictac/Tic", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tictacClient) Tac(ctx context.Context, in *TacRequest, opts ...grpc.CallOption) (*TacResponse, error) {
	out := new(TacResponse)
	err := c.cc.Invoke(ctx, "/pb.Tictac/Tac", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TictacServer is the server API for Tictac service.
type TictacServer interface {
	Tic(context.Context, *TicRequest) (*TicResponse, error)
	Tac(context.Context, *TacRequest) (*TacResponse, error)
}

// UnimplementedTictacServer can be embedded to have forward compatible implementations.
type UnimplementedTictacServer struct {
}

func (*UnimplementedTictacServer) Tic(ctx context.Context, req *TicRequest) (*TicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tic not implemented")
}
func (*UnimplementedTictacServer) Tac(ctx context.Context, req *TacRequest) (*TacResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Tac not implemented")
}

func RegisterTictacServer(s *grpc.Server, srv TictacServer) {
	s.RegisterService(&_Tictac_serviceDesc, srv)
}

func _Tictac_Tic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TictacServer).Tic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Tictac/Tic",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TictacServer).Tic(ctx, req.(*TicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tictac_Tac_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TacRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TictacServer).Tac(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Tictac/Tac",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TictacServer).Tac(ctx, req.(*TacRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tictac_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Tictac",
	HandlerType: (*TictacServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Tic",
			Handler:    _Tictac_Tic_Handler,
		},
		{
			MethodName: "Tac",
			Handler:    _Tictac_Tac_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tictac.proto",
}
