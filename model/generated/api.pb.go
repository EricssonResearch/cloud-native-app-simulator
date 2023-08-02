// protoc --go_out=model --go_opt=module=application-model model/api.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: model/api.proto

package generated

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

type CPUTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of all services that executed CPU tasks and their execution times
	Services map[string]float32 `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"fixed32,2,opt,name=value,proto3"`
}

func (x *CPUTaskResponse) Reset() {
	*x = CPUTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CPUTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CPUTaskResponse) ProtoMessage() {}

func (x *CPUTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CPUTaskResponse.ProtoReflect.Descriptor instead.
func (*CPUTaskResponse) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{0}
}

func (x *CPUTaskResponse) GetServices() map[string]float32 {
	if x != nil {
		return x.Services
	}
	return nil
}

type ServiceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Protocol used to contact this service
	Protocol string `protobuf:"bytes,1,opt,name=protocol,proto3" json:"protocol,omitempty"`
	// Response status code
	// HTTP: 200 OK, 404 Not Found, etc
	// gRPC: OK, INVALID_ARGUMENT, etc
	Status string `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *ServiceResponse) Reset() {
	*x = ServiceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServiceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceResponse) ProtoMessage() {}

func (x *ServiceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServiceResponse.ProtoReflect.Descriptor instead.
func (*ServiceResponse) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{1}
}

func (x *ServiceResponse) GetProtocol() string {
	if x != nil {
		return x.Protocol
	}
	return ""
}

func (x *ServiceResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type NetworkTaskResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of all services that executed or responded to network tasks
	Services []string `protobuf:"bytes,1,rep,name=services,proto3" json:"services,omitempty"`
	// Responses from services
	Responses map[string]*ServiceResponse `protobuf:"bytes,2,rep,name=responses,proto3" json:"responses,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Random payload
	Payload string `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *NetworkTaskResponse) Reset() {
	*x = NetworkTaskResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NetworkTaskResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkTaskResponse) ProtoMessage() {}

func (x *NetworkTaskResponse) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkTaskResponse.ProtoReflect.Descriptor instead.
func (*NetworkTaskResponse) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{2}
}

func (x *NetworkTaskResponse) GetServices() []string {
	if x != nil {
		return x.Services
	}
	return nil
}

func (x *NetworkTaskResponse) GetResponses() map[string]*ServiceResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

func (x *NetworkTaskResponse) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type TaskResponses struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CpuTask     *CPUTaskResponse     `protobuf:"bytes,1,opt,name=cpu_task,json=cpuTask,proto3" json:"cpu_task,omitempty"`
	NetworkTask *NetworkTaskResponse `protobuf:"bytes,2,opt,name=network_task,json=networkTask,proto3" json:"network_task,omitempty"`
}

func (x *TaskResponses) Reset() {
	*x = TaskResponses{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskResponses) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskResponses) ProtoMessage() {}

func (x *TaskResponses) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskResponses.ProtoReflect.Descriptor instead.
func (*TaskResponses) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{3}
}

func (x *TaskResponses) GetCpuTask() *CPUTaskResponse {
	if x != nil {
		return x.CpuTask
	}
	return nil
}

func (x *TaskResponses) GetNetworkTask() *NetworkTaskResponse {
	if x != nil {
		return x.NetworkTask
	}
	return nil
}

type Request struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Random payload
	Payload string `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *Request) Reset() {
	*x = Request{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Request) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Request) ProtoMessage() {}

func (x *Request) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Request.ProtoReflect.Descriptor instead.
func (*Request) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{4}
}

func (x *Request) GetPayload() string {
	if x != nil {
		return x.Payload
	}
	return ""
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of called endpoint
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	// Task responses
	Tasks *TaskResponses `protobuf:"bytes,2,opt,name=tasks,proto3" json:"tasks,omitempty"`
	// Error message
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_api_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_model_api_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_model_api_proto_rawDescGZIP(), []int{5}
}

func (x *Response) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *Response) GetTasks() *TaskResponses {
	if x != nil {
		return x.Tasks
	}
	return nil
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_model_api_proto protoreflect.FileDescriptor

var file_model_api_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x22, 0x94, 0x01, 0x0a,
	0x0f, 0x43, 0x50, 0x55, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x44, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x43,
	0x50, 0x55, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x1a, 0x3b, 0x0a, 0x0d, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a,
	0x02, 0x38, 0x01, 0x22, 0x45, 0x0a, 0x0f, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63,
	0x6f, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0xf2, 0x01, 0x0a, 0x13, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x12, 0x4b,
	0x0a, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x2d, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x4e, 0x65,
	0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x09, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61,
	0x79, 0x6c, 0x6f, 0x61, 0x64, 0x1a, 0x58, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x30, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72,
	0x61, 0x74, 0x65, 0x64, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0x89, 0x01, 0x0a, 0x0d, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x73, 0x12, 0x35, 0x0a, 0x08, 0x63, 0x70, 0x75, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e,
	0x43, 0x50, 0x55, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x07, 0x63, 0x70, 0x75, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x41, 0x0a, 0x0c, 0x6e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x5f, 0x74, 0x61, 0x73, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e,
	0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f,
	0x72, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x0b,
	0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x54, 0x61, 0x73, 0x6b, 0x22, 0x23, 0x0a, 0x07, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x22, 0x70, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x2e, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61,
	0x74, 0x65, 0x64, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x73, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x42, 0x1d, 0x5a, 0x1b, 0x61, 0x70, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2d, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65,
	0x64, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_model_api_proto_rawDescOnce sync.Once
	file_model_api_proto_rawDescData = file_model_api_proto_rawDesc
)

func file_model_api_proto_rawDescGZIP() []byte {
	file_model_api_proto_rawDescOnce.Do(func() {
		file_model_api_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_api_proto_rawDescData)
	})
	return file_model_api_proto_rawDescData
}

var file_model_api_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_model_api_proto_goTypes = []interface{}{
	(*CPUTaskResponse)(nil),     // 0: generated.CPUTaskResponse
	(*ServiceResponse)(nil),     // 1: generated.ServiceResponse
	(*NetworkTaskResponse)(nil), // 2: generated.NetworkTaskResponse
	(*TaskResponses)(nil),       // 3: generated.TaskResponses
	(*Request)(nil),             // 4: generated.Request
	(*Response)(nil),            // 5: generated.Response
	nil,                         // 6: generated.CPUTaskResponse.ServicesEntry
	nil,                         // 7: generated.NetworkTaskResponse.ResponsesEntry
}
var file_model_api_proto_depIdxs = []int32{
	6, // 0: generated.CPUTaskResponse.services:type_name -> generated.CPUTaskResponse.ServicesEntry
	7, // 1: generated.NetworkTaskResponse.responses:type_name -> generated.NetworkTaskResponse.ResponsesEntry
	0, // 2: generated.TaskResponses.cpu_task:type_name -> generated.CPUTaskResponse
	2, // 3: generated.TaskResponses.network_task:type_name -> generated.NetworkTaskResponse
	3, // 4: generated.Response.tasks:type_name -> generated.TaskResponses
	1, // 5: generated.NetworkTaskResponse.ResponsesEntry.value:type_name -> generated.ServiceResponse
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_model_api_proto_init() }
func file_model_api_proto_init() {
	if File_model_api_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_api_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CPUTaskResponse); i {
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
		file_model_api_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServiceResponse); i {
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
		file_model_api_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NetworkTaskResponse); i {
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
		file_model_api_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskResponses); i {
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
		file_model_api_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Request); i {
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
		file_model_api_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
			RawDescriptor: file_model_api_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_api_proto_goTypes,
		DependencyIndexes: file_model_api_proto_depIdxs,
		MessageInfos:      file_model_api_proto_msgTypes,
	}.Build()
	File_model_api_proto = out.File
	file_model_api_proto_rawDesc = nil
	file_model_api_proto_goTypes = nil
	file_model_api_proto_depIdxs = nil
}
