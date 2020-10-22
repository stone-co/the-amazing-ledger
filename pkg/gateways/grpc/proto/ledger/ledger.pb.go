// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.7.1
// source: pkg/gateways/grpc/proto/ledger/ledger.proto

package ledger

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Operation int32

const (
	Operation_DEBIT  Operation = 0
	Operation_CREDIT Operation = 1
)

// Enum value maps for Operation.
var (
	Operation_name = map[int32]string{
		0: "DEBIT",
		1: "CREDIT",
	}
	Operation_value = map[string]int32{
		"DEBIT":  0,
		"CREDIT": 1,
	}
)

func (x Operation) Enum() *Operation {
	p := new(Operation)
	*p = x
	return p
}

func (x Operation) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Operation) Descriptor() protoreflect.EnumDescriptor {
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_enumTypes[0].Descriptor()
}

func (Operation) Type() protoreflect.EnumType {
	return &file_pkg_gateways_grpc_proto_ledger_ledger_proto_enumTypes[0]
}

func (x Operation) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Operation.Descriptor instead.
func (Operation) EnumDescriptor() ([]byte, []int) {
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescGZIP(), []int{0}
}

type SaveTransactionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id      string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Entries []*Entry `protobuf:"bytes,2,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *SaveTransactionRequest) Reset() {
	*x = SaveTransactionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveTransactionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveTransactionRequest) ProtoMessage() {}

func (x *SaveTransactionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveTransactionRequest.ProtoReflect.Descriptor instead.
func (*SaveTransactionRequest) Descriptor() ([]byte, []int) {
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescGZIP(), []int{0}
}

func (x *SaveTransactionRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *SaveTransactionRequest) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	AccountId       string    `protobuf:"bytes,2,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	ExpectedVersion uint64    `protobuf:"varint,3,opt,name=expected_version,json=expectedVersion,proto3" json:"expected_version,omitempty"`
	Operation       Operation `protobuf:"varint,4,opt,name=operation,proto3,enum=ledger.Operation" json:"operation,omitempty"`
	Amount          int32     `protobuf:"varint,5,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescGZIP(), []int{1}
}

func (x *Entry) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Entry) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *Entry) GetExpectedVersion() uint64 {
	if x != nil {
		return x.ExpectedVersion
	}
	return 0
}

func (x *Entry) GetOperation() Operation {
	if x != nil {
		return x.Operation
	}
	return Operation_DEBIT
}

func (x *Entry) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type SaveTransactionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Error string `protobuf:"bytes,1,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *SaveTransactionResponse) Reset() {
	*x = SaveTransactionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SaveTransactionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SaveTransactionResponse) ProtoMessage() {}

func (x *SaveTransactionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SaveTransactionResponse.ProtoReflect.Descriptor instead.
func (*SaveTransactionResponse) Descriptor() ([]byte, []int) {
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescGZIP(), []int{2}
}

func (x *SaveTransactionResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_pkg_gateways_grpc_proto_ledger_ledger_proto protoreflect.FileDescriptor

var file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79, 0x73, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72,
	0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x22, 0x51, 0x0a, 0x16, 0x53, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x27, 0x0a, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x0d, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0xaa, 0x01, 0x0a, 0x05, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x29, 0x0a, 0x10, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0f, 0x65, 0x78, 0x70,
	0x65, 0x63, 0x74, 0x65, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2f, 0x0a, 0x09,
	0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x11, 0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x2f, 0x0a, 0x17, 0x53, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2a, 0x22, 0x0a, 0x09, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x45, 0x42, 0x49, 0x54, 0x10, 0x00, 0x12, 0x0a,
	0x0a, 0x06, 0x43, 0x52, 0x45, 0x44, 0x49, 0x54, 0x10, 0x01, 0x32, 0x65, 0x0a, 0x0d, 0x4c, 0x65,
	0x64, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x54, 0x0a, 0x0f, 0x53,
	0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1e,
	0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f,
	0x2e, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2e, 0x53, 0x61, 0x76, 0x65, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x73, 0x74, 0x6f, 0x6e, 0x65, 0x2d, 0x63, 0x6f, 0x2f, 0x74, 0x68, 0x65, 0x2d, 0x61, 0x6d, 0x61,
	0x7a, 0x69, 0x6e, 0x67, 0x2d, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescOnce sync.Once
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescData = file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDesc
)

func file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescGZIP() []byte {
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescOnce.Do(func() {
		file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescData)
	})
	return file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDescData
}

var file_pkg_gateways_grpc_proto_ledger_ledger_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pkg_gateways_grpc_proto_ledger_ledger_proto_goTypes = []interface{}{
	(Operation)(0),                  // 0: ledger.Operation
	(*SaveTransactionRequest)(nil),  // 1: ledger.SaveTransactionRequest
	(*Entry)(nil),                   // 2: ledger.Entry
	(*SaveTransactionResponse)(nil), // 3: ledger.SaveTransactionResponse
}
var file_pkg_gateways_grpc_proto_ledger_ledger_proto_depIdxs = []int32{
	2, // 0: ledger.SaveTransactionRequest.entries:type_name -> ledger.Entry
	0, // 1: ledger.Entry.operation:type_name -> ledger.Operation
	1, // 2: ledger.LedgerService.SaveTransaction:input_type -> ledger.SaveTransactionRequest
	3, // 3: ledger.LedgerService.SaveTransaction:output_type -> ledger.SaveTransactionResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pkg_gateways_grpc_proto_ledger_ledger_proto_init() }
func file_pkg_gateways_grpc_proto_ledger_ledger_proto_init() {
	if File_pkg_gateways_grpc_proto_ledger_ledger_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveTransactionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SaveTransactionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_gateways_grpc_proto_ledger_ledger_proto_goTypes,
		DependencyIndexes: file_pkg_gateways_grpc_proto_ledger_ledger_proto_depIdxs,
		EnumInfos:         file_pkg_gateways_grpc_proto_ledger_ledger_proto_enumTypes,
		MessageInfos:      file_pkg_gateways_grpc_proto_ledger_ledger_proto_msgTypes,
	}.Build()
	File_pkg_gateways_grpc_proto_ledger_ledger_proto = out.File
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_rawDesc = nil
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_goTypes = nil
	file_pkg_gateways_grpc_proto_ledger_ledger_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// LedgerServiceClient is the client API for LedgerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LedgerServiceClient interface {
	SaveTransaction(ctx context.Context, in *SaveTransactionRequest, opts ...grpc.CallOption) (*SaveTransactionResponse, error)
}

type ledgerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLedgerServiceClient(cc grpc.ClientConnInterface) LedgerServiceClient {
	return &ledgerServiceClient{cc}
}

func (c *ledgerServiceClient) SaveTransaction(ctx context.Context, in *SaveTransactionRequest, opts ...grpc.CallOption) (*SaveTransactionResponse, error) {
	out := new(SaveTransactionResponse)
	err := c.cc.Invoke(ctx, "/ledger.LedgerService/SaveTransaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LedgerServiceServer is the server API for LedgerService service.
type LedgerServiceServer interface {
	SaveTransaction(context.Context, *SaveTransactionRequest) (*SaveTransactionResponse, error)
}

// UnimplementedLedgerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedLedgerServiceServer struct {
}

func (*UnimplementedLedgerServiceServer) SaveTransaction(context.Context, *SaveTransactionRequest) (*SaveTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveTransaction not implemented")
}

func RegisterLedgerServiceServer(s *grpc.Server, srv LedgerServiceServer) {
	s.RegisterService(&_LedgerService_serviceDesc, srv)
}

func _LedgerService_SaveTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveTransactionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LedgerServiceServer).SaveTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ledger.LedgerService/SaveTransaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LedgerServiceServer).SaveTransaction(ctx, req.(*SaveTransactionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _LedgerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ledger.LedgerService",
	HandlerType: (*LedgerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveTransaction",
			Handler:    _LedgerService_SaveTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/gateways/grpc/proto/ledger/ledger.proto",
}
