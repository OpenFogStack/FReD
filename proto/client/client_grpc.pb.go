// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package client

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

// ClientClient is the client API for Client service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientClient interface {
	CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*Empty, error)
	DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*Empty, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
	AddReplica(ctx context.Context, in *AddReplicaRequest, opts ...grpc.CallOption) (*Empty, error)
	GetKeygroupReplica(ctx context.Context, in *GetKeygroupReplicaRequest, opts ...grpc.CallOption) (*GetKeygroupReplicaResponse, error)
	RemoveReplica(ctx context.Context, in *RemoveReplicaRequest, opts ...grpc.CallOption) (*Empty, error)
	GetReplica(ctx context.Context, in *GetReplicaRequest, opts ...grpc.CallOption) (*GetReplicaResponse, error)
	GetAllReplica(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllReplicaResponse, error)
	GetKeygroupTriggers(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error)
	AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*Empty, error)
	AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...grpc.CallOption) (*Empty, error)
}

type clientClient struct {
	cc grpc.ClientConnInterface
}

func NewClientClient(cc grpc.ClientConnInterface) ClientClient {
	return &clientClient{cc}
}

func (c *clientClient) CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/CreateKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/DeleteKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error) {
	out := new(ScanResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/Scan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) AddReplica(ctx context.Context, in *AddReplicaRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/AddReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) GetKeygroupReplica(ctx context.Context, in *GetKeygroupReplicaRequest, opts ...grpc.CallOption) (*GetKeygroupReplicaResponse, error) {
	out := new(GetKeygroupReplicaResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/GetKeygroupReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) RemoveReplica(ctx context.Context, in *RemoveReplicaRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/RemoveReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) GetReplica(ctx context.Context, in *GetReplicaRequest, opts ...grpc.CallOption) (*GetReplicaResponse, error) {
	out := new(GetReplicaResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/GetReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) GetAllReplica(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetAllReplicaResponse, error) {
	out := new(GetAllReplicaResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/GetAllReplica", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) GetKeygroupTriggers(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error) {
	out := new(GetKeygroupTriggerResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/GetKeygroupTriggers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) AddTrigger(ctx context.Context, in *AddTriggerRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/AddTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) RemoveTrigger(ctx context.Context, in *RemoveTriggerRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/RemoveTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) AddUser(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/AddUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientClient) RemoveUser(ctx context.Context, in *RemoveUserRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/mcc.fred.client.Client/RemoveUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientServer is the server API for Client service.
// All implementations should embed UnimplementedClientServer
// for forward compatibility
type ClientServer interface {
	CreateKeygroup(context.Context, *CreateKeygroupRequest) (*Empty, error)
	DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*Empty, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Scan(context.Context, *ScanRequest) (*ScanResponse, error)
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	AddReplica(context.Context, *AddReplicaRequest) (*Empty, error)
	GetKeygroupReplica(context.Context, *GetKeygroupReplicaRequest) (*GetKeygroupReplicaResponse, error)
	RemoveReplica(context.Context, *RemoveReplicaRequest) (*Empty, error)
	GetReplica(context.Context, *GetReplicaRequest) (*GetReplicaResponse, error)
	GetAllReplica(context.Context, *Empty) (*GetAllReplicaResponse, error)
	GetKeygroupTriggers(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error)
	AddTrigger(context.Context, *AddTriggerRequest) (*Empty, error)
	RemoveTrigger(context.Context, *RemoveTriggerRequest) (*Empty, error)
	AddUser(context.Context, *AddUserRequest) (*Empty, error)
	RemoveUser(context.Context, *RemoveUserRequest) (*Empty, error)
}

// UnimplementedClientServer should be embedded to have forward compatible implementations.
type UnimplementedClientServer struct {
}

func (UnimplementedClientServer) CreateKeygroup(context.Context, *CreateKeygroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeygroup not implemented")
}
func (UnimplementedClientServer) DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroup not implemented")
}
func (UnimplementedClientServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedClientServer) Scan(context.Context, *ScanRequest) (*ScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Scan not implemented")
}
func (UnimplementedClientServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedClientServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedClientServer) Append(context.Context, *AppendRequest) (*AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedClientServer) AddReplica(context.Context, *AddReplicaRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddReplica not implemented")
}
func (UnimplementedClientServer) GetKeygroupReplica(context.Context, *GetKeygroupReplicaRequest) (*GetKeygroupReplicaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeygroupReplica not implemented")
}
func (UnimplementedClientServer) RemoveReplica(context.Context, *RemoveReplicaRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveReplica not implemented")
}
func (UnimplementedClientServer) GetReplica(context.Context, *GetReplicaRequest) (*GetReplicaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetReplica not implemented")
}
func (UnimplementedClientServer) GetAllReplica(context.Context, *Empty) (*GetAllReplicaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllReplica not implemented")
}
func (UnimplementedClientServer) GetKeygroupTriggers(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeygroupTriggers not implemented")
}
func (UnimplementedClientServer) AddTrigger(context.Context, *AddTriggerRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddTrigger not implemented")
}
func (UnimplementedClientServer) RemoveTrigger(context.Context, *RemoveTriggerRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveTrigger not implemented")
}
func (UnimplementedClientServer) AddUser(context.Context, *AddUserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUser not implemented")
}
func (UnimplementedClientServer) RemoveUser(context.Context, *RemoveUserRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUser not implemented")
}

// UnsafeClientServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientServer will
// result in compilation errors.
type UnsafeClientServer interface {
	mustEmbedUnimplementedClientServer()
}

func RegisterClientServer(s grpc.ServiceRegistrar, srv ClientServer) {
	s.RegisterService(&Client_ServiceDesc, srv)
}

func _Client_CreateKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).CreateKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/CreateKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).CreateKeygroup(ctx, req.(*CreateKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_DeleteKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).DeleteKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/DeleteKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).DeleteKeygroup(ctx, req.(*DeleteKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/Scan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Scan(ctx, req.(*ScanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_AddReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).AddReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/AddReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).AddReplica(ctx, req.(*AddReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_GetKeygroupReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeygroupReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).GetKeygroupReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/GetKeygroupReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).GetKeygroupReplica(ctx, req.(*GetKeygroupReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_RemoveReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).RemoveReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/RemoveReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).RemoveReplica(ctx, req.(*RemoveReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_GetReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetReplicaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).GetReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/GetReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).GetReplica(ctx, req.(*GetReplicaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_GetAllReplica_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).GetAllReplica(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/GetAllReplica",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).GetAllReplica(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_GetKeygroupTriggers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeygroupTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).GetKeygroupTriggers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/GetKeygroupTriggers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).GetKeygroupTriggers(ctx, req.(*GetKeygroupTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_AddTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).AddTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/AddTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).AddTrigger(ctx, req.(*AddTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_RemoveTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).RemoveTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/RemoveTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).RemoveTrigger(ctx, req.(*RemoveTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_AddUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).AddUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/AddUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).AddUser(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Client_RemoveUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientServer).RemoveUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.client.Client/RemoveUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientServer).RemoveUser(ctx, req.(*RemoveUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Client_ServiceDesc is the grpc.ServiceDesc for Client service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Client_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mcc.fred.client.Client",
	HandlerType: (*ClientServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateKeygroup",
			Handler:    _Client_CreateKeygroup_Handler,
		},
		{
			MethodName: "DeleteKeygroup",
			Handler:    _Client_DeleteKeygroup_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _Client_Read_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _Client_Scan_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Client_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Client_Delete_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _Client_Append_Handler,
		},
		{
			MethodName: "AddReplica",
			Handler:    _Client_AddReplica_Handler,
		},
		{
			MethodName: "GetKeygroupReplica",
			Handler:    _Client_GetKeygroupReplica_Handler,
		},
		{
			MethodName: "RemoveReplica",
			Handler:    _Client_RemoveReplica_Handler,
		},
		{
			MethodName: "GetReplica",
			Handler:    _Client_GetReplica_Handler,
		},
		{
			MethodName: "GetAllReplica",
			Handler:    _Client_GetAllReplica_Handler,
		},
		{
			MethodName: "GetKeygroupTriggers",
			Handler:    _Client_GetKeygroupTriggers_Handler,
		},
		{
			MethodName: "AddTrigger",
			Handler:    _Client_AddTrigger_Handler,
		},
		{
			MethodName: "RemoveTrigger",
			Handler:    _Client_RemoveTrigger_Handler,
		},
		{
			MethodName: "AddUser",
			Handler:    _Client_AddUser_Handler,
		},
		{
			MethodName: "RemoveUser",
			Handler:    _Client_RemoveUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client.proto",
}
