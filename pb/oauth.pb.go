// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v4.24.3
// source: oauth.proto

package pb

import (
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

type CheckTokenRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Token string `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *CheckTokenRequest) Reset() {
	*x = CheckTokenRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oauth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenRequest) ProtoMessage() {}

func (x *CheckTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_oauth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenRequest.ProtoReflect.Descriptor instead.
func (*CheckTokenRequest) Descriptor() ([]byte, []int) {
	return file_oauth_proto_rawDescGZIP(), []int{0}
}

func (x *CheckTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type UserDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId      string   `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Username    string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Authorities []string `protobuf:"bytes,3,rep,name=authorities,proto3" json:"authorities,omitempty"`
}

func (x *UserDetails) Reset() {
	*x = UserDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oauth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDetails) ProtoMessage() {}

func (x *UserDetails) ProtoReflect() protoreflect.Message {
	mi := &file_oauth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDetails.ProtoReflect.Descriptor instead.
func (*UserDetails) Descriptor() ([]byte, []int) {
	return file_oauth_proto_rawDescGZIP(), []int{1}
}

func (x *UserDetails) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UserDetails) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UserDetails) GetAuthorities() []string {
	if x != nil {
		return x.Authorities
	}
	return nil
}

type ClientDetails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientId                    string   `protobuf:"bytes,1,opt,name=clientId,proto3" json:"clientId,omitempty"`
	AccessTokenValiditySeconds  int32    `protobuf:"varint,2,opt,name=accessTokenValiditySeconds,proto3" json:"accessTokenValiditySeconds,omitempty"`
	RefreshTokenValiditySeconds int32    `protobuf:"varint,3,opt,name=refreshTokenValiditySeconds,proto3" json:"refreshTokenValiditySeconds,omitempty"`
	AuthorizedGrantTypes        []string `protobuf:"bytes,4,rep,name=authorizedGrantTypes,proto3" json:"authorizedGrantTypes,omitempty"`
}

func (x *ClientDetails) Reset() {
	*x = ClientDetails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oauth_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientDetails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientDetails) ProtoMessage() {}

func (x *ClientDetails) ProtoReflect() protoreflect.Message {
	mi := &file_oauth_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientDetails.ProtoReflect.Descriptor instead.
func (*ClientDetails) Descriptor() ([]byte, []int) {
	return file_oauth_proto_rawDescGZIP(), []int{2}
}

func (x *ClientDetails) GetClientId() string {
	if x != nil {
		return x.ClientId
	}
	return ""
}

func (x *ClientDetails) GetAccessTokenValiditySeconds() int32 {
	if x != nil {
		return x.AccessTokenValiditySeconds
	}
	return 0
}

func (x *ClientDetails) GetRefreshTokenValiditySeconds() int32 {
	if x != nil {
		return x.RefreshTokenValiditySeconds
	}
	return 0
}

func (x *ClientDetails) GetAuthorizedGrantTypes() []string {
	if x != nil {
		return x.AuthorizedGrantTypes
	}
	return nil
}

type CheckTokenResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserDetails   *UserDetails   `protobuf:"bytes,1,opt,name=userDetails,proto3" json:"userDetails,omitempty"`
	ClientDetails *ClientDetails `protobuf:"bytes,2,opt,name=clientDetails,proto3" json:"clientDetails,omitempty"`
	IsValidToken  bool           `protobuf:"varint,3,opt,name=isValidToken,proto3" json:"isValidToken,omitempty"`
	Err           string         `protobuf:"bytes,4,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *CheckTokenResponse) Reset() {
	*x = CheckTokenResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_oauth_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckTokenResponse) ProtoMessage() {}

func (x *CheckTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_oauth_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckTokenResponse.ProtoReflect.Descriptor instead.
func (*CheckTokenResponse) Descriptor() ([]byte, []int) {
	return file_oauth_proto_rawDescGZIP(), []int{3}
}

func (x *CheckTokenResponse) GetUserDetails() *UserDetails {
	if x != nil {
		return x.UserDetails
	}
	return nil
}

func (x *CheckTokenResponse) GetClientDetails() *ClientDetails {
	if x != nil {
		return x.ClientDetails
	}
	return nil
}

func (x *CheckTokenResponse) GetIsValidToken() bool {
	if x != nil {
		return x.IsValidToken
	}
	return false
}

func (x *CheckTokenResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

var File_oauth_proto protoreflect.FileDescriptor

var file_oauth_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6f, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70,
	0x62, 0x22, 0x29, 0x0a, 0x11, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x63, 0x0a, 0x0b,
	0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x22, 0xe1, 0x01, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61,
	0x69, 0x6c, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x3e, 0x0a, 0x1a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x1a, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x56, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x12,
	0x40, 0x0a, 0x1b, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x56,
	0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x1b, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x69, 0x74, 0x79, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x73, 0x12, 0x32, 0x0a, 0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x47,
	0x72, 0x61, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x14, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x64, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x54, 0x79, 0x70, 0x65, 0x73, 0x22, 0xb6, 0x01, 0x0a, 0x12, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x31, 0x0a, 0x0b,
	0x75, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x73, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12,
	0x37, 0x0a, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x52, 0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x69, 0x73, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x69, 0x73, 0x56, 0x61, 0x6c, 0x69, 0x64, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x72, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x32, 0x4d,
	0x0a, 0x0c, 0x4f, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3d,
	0x0a, 0x0a, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x15, 0x2e, 0x70,
	0x62, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x54, 0x6f,
	0x6b, 0x65, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x2a, 0x5a,
	0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x68, 0x65, 0x6e,
	0x67, 0x78, 0x75, 0x66, 0x65, 0x6e, 0x67, 0x31, 0x39, 0x39, 0x34, 0x2f, 0x67, 0x6f, 0x2d, 0x73,
	0x65, 0x63, 0x6b, 0x69, 0x6c, 0x6c, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_oauth_proto_rawDescOnce sync.Once
	file_oauth_proto_rawDescData = file_oauth_proto_rawDesc
)

func file_oauth_proto_rawDescGZIP() []byte {
	file_oauth_proto_rawDescOnce.Do(func() {
		file_oauth_proto_rawDescData = protoimpl.X.CompressGZIP(file_oauth_proto_rawDescData)
	})
	return file_oauth_proto_rawDescData
}

var file_oauth_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_oauth_proto_goTypes = []interface{}{
	(*CheckTokenRequest)(nil),  // 0: pb.CheckTokenRequest
	(*UserDetails)(nil),        // 1: pb.UserDetails
	(*ClientDetails)(nil),      // 2: pb.ClientDetails
	(*CheckTokenResponse)(nil), // 3: pb.CheckTokenResponse
}
var file_oauth_proto_depIdxs = []int32{
	1, // 0: pb.CheckTokenResponse.userDetails:type_name -> pb.UserDetails
	2, // 1: pb.CheckTokenResponse.clientDetails:type_name -> pb.ClientDetails
	0, // 2: pb.OAuthService.CheckToken:input_type -> pb.CheckTokenRequest
	3, // 3: pb.OAuthService.CheckToken:output_type -> pb.CheckTokenResponse
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_oauth_proto_init() }
func file_oauth_proto_init() {
	if File_oauth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_oauth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenRequest); i {
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
		file_oauth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDetails); i {
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
		file_oauth_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientDetails); i {
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
		file_oauth_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckTokenResponse); i {
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
			RawDescriptor: file_oauth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_oauth_proto_goTypes,
		DependencyIndexes: file_oauth_proto_depIdxs,
		MessageInfos:      file_oauth_proto_msgTypes,
	}.Build()
	File_oauth_proto = out.File
	file_oauth_proto_rawDesc = nil
	file_oauth_proto_goTypes = nil
	file_oauth_proto_depIdxs = nil
}
