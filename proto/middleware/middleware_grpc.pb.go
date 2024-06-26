// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: middleware.proto

package middleware

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MiddlewareClient is the client API for Middleware service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MiddlewareClient interface {
	CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*Empty, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error)
	Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Empty, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Empty, error)
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
	Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*Empty, error)
	ChooseReplica(ctx context.Context, in *ChooseReplicaRequest, opts ...grpc.CallOption) (*Empty, error)
	AddReplica(ctx context.Context, in *AddReplicaRequest, opts ...grpc.CallOption) (*Empty, error)
	GetKeygroupInfo(ctx context.Context, in *GetKeygroupInfoRequest, opts ...grpc.CallOption) (*GetKeygroupInfoResponse, error)
	RemoveReplica(ctx context.Context, in *RemoveReplicaRequest, opts ...grpc.CallOption) (*Empty, error)
	GetReplica(ctx context.Context, in *GetReplicaRequest, opts ...grpc.CallOption) (*GetReplicaResponse, error)
	GetAllReplica(ctx context.Context, in *GetAllReplicaRequest, opts ...grpc.CallOption) (*GetAllReplicaResponse, error)
	GetKeygroupTriggers(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error)
	AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*Empty, error)
	AddUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*Empty, error)
}

type middlewareClient struct {
	cc grpc.ClientConnInterface
}

func NewMiddlewareClient(cc grpc.ClientConnInterface) MiddlewareClient {
	return &middlewareClient{cc}
}

func (c *middlewareClient) CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/CreateKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/DeleteKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error) {
	out := new(ScanResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Scan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error) {
	out := new(KeysResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Keys", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) Notify(ctx context.Context, in *NotifyRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/Notify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) ChooseReplica(ctx context.Context, in *ChooseReplicaRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/ChooseReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) AddReplica(ctx context.Context, in *AddReplicaRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/AddReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) GetKeygroupInfo(ctx context.Context, in *GetKeygroupInfoRequest, opts ...grpc.CallOption) (*GetKeygroupInfoResponse, error) {
	out := new(GetKeygroupInfoResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/GetKeygroupInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) RemoveReplica(ctx context.Context, in *RemoveReplicaRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/RemoveReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) GetReplica(ctx context.Context, in *GetReplicaRequest, opts ...grpc.CallOption) (*GetReplicaResponse, error) {
	out := new(GetReplicaResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/GetReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) GetAllReplica(ctx context.Context, in *GetAllReplicaRequest, opts ...grpc.CallOption) (*GetAllReplicaResponse, error) {
	out := new(GetAllReplicaResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/GetAllReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) GetKeygroupTriggers(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error) {
	out := new(GetKeygroupTriggerResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/GetKeygroupTriggers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/AddTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/RemoveTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) AddUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *middlewareClient) RemoveUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.middleware.Middleware/RemoveUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MiddlewareServer is the server API for Middleware service.
// All implementations should embed UnimplementedMiddlewareServer
// for forward compatibility
type MiddlewareServer interface {
	CreateKeygroup(context.Context, *CreateKeygroupRequest) (*Empty, error)
	DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*Empty, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Scan(context.Context, *ScanRequest) (*ScanResponse, error)
	Keys(context.Context, *KeysRequest) (*KeysResponse, error)
	Update(context.Context, *UpdateRequest) (*Empty, error)
	Delete(context.Context, *DeleteRequest) (*Empty, error)
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	Notify(context.Context, *NotifyRequest) (*Empty, error)
	ChooseReplica(context.Context, *ChooseReplicaRequest) (*Empty, error)
	AddReplica(context.Context, *AddReplicaRequest) (*Empty, error)
	GetKeygroupInfo(context.Context, *GetKeygroupInfoRequest) (*GetKeygroupInfoResponse, error)
	RemoveReplica(context.Context, *RemoveReplicaRequest) (*Empty, error)
	GetReplica(context.Context, *GetReplicaRequest) (*GetReplicaResponse, error)
	GetAllReplica(context.Context, *GetAllReplicaRequest) (*GetAllReplicaResponse, error)
	GetKeygroupTriggers(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error)
	AddTrigger(context.Context, *AddTriggerRequest) (*Empty, error)
	RemoveTrigger(context.Context, *RemoveTriggerRequest) (*Empty, error)
	AddUser(context.Context, *UserRequest) (*Empty, error)
	RemoveUser(context.Context, *UserRequest) (*Empty, error)
}

// UnimplementedMiddlewareServer should be embedded to have forward compatible implementations.
type UnimplementedMiddlewareServer struct {
}

func (UnimplementedMiddlewareServer) CreateKeygroup(context.Context, *CreateKeygroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeygroup not implemented")
}
func (UnimplementedMiddlewareServer) DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroup not implemented")
}
func (UnimplementedMiddlewareServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedMiddlewareServer) Scan(context.Context, *ScanRequest) (*ScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Scan not implemented")
}
func (UnimplementedMiddlewareServer) Keys(context.Context, *KeysRequest) (*KeysResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Keys not implemented")
}
func (UnimplementedMiddlewareServer) Update(context.Context, *UpdateRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMiddlewareServer) Delete(context.Context, *DeleteRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedMiddlewareServer) Append(context.Context, *AppendRequest) (*AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedMiddlewareServer) Notify(context.Context, *NotifyRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Notify not implemented")
}
func (UnimplementedMiddlewareServer) ChooseReplica(context.Context, *ChooseReplicaRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChooseReplica not implemented")
}
func (UnimplementedMiddlewareServer) AddReplica(context.Context, *AddReplicaRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddReplica not implemented")
}
func (UnimplementedMiddlewareServer) GetKeygroupInfo(context.Context, *GetKeygroupInfoRequest) (*GetKeygroupInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeygroupInfo not implemented")
}
func (UnimplementedMiddlewareServer) RemoveReplica(context.Context, *RemoveReplicaRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveReplica not implemented")
}
func (UnimplementedMiddlewareServer) GetReplica(context.Context, *GetReplicaRequest) (*GetReplicaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReplica not implemented")
}
func (UnimplementedMiddlewareServer) GetAllReplica(context.Context, *GetAllReplicaRequest) (*GetAllReplicaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllReplica not implemented")
}
func (UnimplementedMiddlewareServer) GetKeygroupTriggers(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeygroupTriggers not implemented")
}
func (UnimplementedMiddlewareServer) AddTrigger(context.Context, *AddTriggerRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTrigger not implemented")
}
func (UnimplementedMiddlewareServer) RemoveTrigger(context.Context, *RemoveTriggerRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTrigger not implemented")
}
func (UnimplementedMiddlewareServer) AddUser(context.Context, *UserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedMiddlewareServer) RemoveUser(context.Context, *UserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}

// UnsafeMiddlewareServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MiddlewareServer will
// result in compilation errors.
type UnsafeMiddlewareServer interface {
	mustEmbedUnimplementedMiddlewareServer()
}

func RegisterMiddlewareServer(s grpc.ServiceRegistrar, srv MiddlewareServer) {
	s.RegisterService(&Middleware_ServiceDesc, srv)
}

func _Middleware_CreateKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).CreateKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/CreateKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).CreateKeygroup(ctx, req.(*CreateKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_DeleteKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).DeleteKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/DeleteKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).DeleteKeygroup(ctx, req.(*DeleteKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Scan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Scan(ctx, req.(*ScanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Keys_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeysRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Keys(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Keys",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Keys(ctx, req.(*KeysRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_Notify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NotifyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).Notify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/Notify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).Notify(ctx, req.(*NotifyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_ChooseReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChooseReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).ChooseReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/ChooseReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).ChooseReplica(ctx, req.(*ChooseReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_AddReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).AddReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/AddReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).AddReplica(ctx, req.(*AddReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_GetKeygroupInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeygroupInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).GetKeygroupInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/GetKeygroupInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).GetKeygroupInfo(ctx, req.(*GetKeygroupInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_RemoveReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).RemoveReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/RemoveReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).RemoveReplica(ctx, req.(*RemoveReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_GetReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).GetReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/GetReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).GetReplica(ctx, req.(*GetReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_GetAllReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).GetAllReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/GetAllReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).GetAllReplica(ctx, req.(*GetAllReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_GetKeygroupTriggers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeygroupTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).GetKeygroupTriggers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/GetKeygroupTriggers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).GetKeygroupTriggers(ctx, req.(*GetKeygroupTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_AddTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).AddTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/AddTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).AddTrigger(ctx, req.(*AddTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_RemoveTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).RemoveTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/RemoveTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).RemoveTrigger(ctx, req.(*RemoveTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).AddUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Middleware_RemoveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MiddlewareServer).RemoveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.middleware.Middleware/RemoveUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MiddlewareServer).RemoveUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Middleware_ServiceDesc is the grpc.ServiceDesc for Middleware service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Middleware_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mcc.fred.middleware.Middleware",
	HandlerType: (*MiddlewareServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateKeygroup",
			Handler:    _Middleware_CreateKeygroup_Handler,
		},
		{
			MethodName: "DeleteKeygroup",
			Handler:    _Middleware_DeleteKeygroup_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _Middleware_Read_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _Middleware_Scan_Handler,
		},
		{
			MethodName: "Keys",
			Handler:    _Middleware_Keys_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Middleware_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Middleware_Delete_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _Middleware_Append_Handler,
		},
		{
			MethodName: "Notify",
			Handler:    _Middleware_Notify_Handler,
		},
		{
			MethodName: "ChooseReplica",
			Handler:    _Middleware_ChooseReplica_Handler,
		},
		{
			MethodName: "AddReplica",
			Handler:    _Middleware_AddReplica_Handler,
		},
		{
			MethodName: "GetKeygroupInfo",
			Handler:    _Middleware_GetKeygroupInfo_Handler,
		},
		{
			MethodName: "RemoveReplica",
			Handler:    _Middleware_RemoveReplica_Handler,
		},
		{
			MethodName: "GetReplica",
			Handler:    _Middleware_GetReplica_Handler,
		},
		{
			MethodName: "GetAllReplica",
			Handler:    _Middleware_GetAllReplica_Handler,
		},
		{
			MethodName: "GetKeygroupTriggers",
			Handler:    _Middleware_GetKeygroupTriggers_Handler,
		},
		{
			MethodName: "AddTrigger",
			Handler:    _Middleware_AddTrigger_Handler,
		},
		{
			MethodName: "RemoveTrigger",
			Handler:    _Middleware_RemoveTrigger_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _Middleware_AddUser_Handler,
		},
		{
			MethodName: "RemoveUser",
			Handler:    _Middleware_RemoveUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "middleware.proto",
}
