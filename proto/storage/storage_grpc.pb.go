// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package storage

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

// DatabaseClient is the client API for Database service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseClient interface {
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error)
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error)
	ReadAll(ctx context.Context, in *ReadAllRequest, opts ...grpc.CallOption) (*ReadAllResponse, error)
	IDs(ctx context.Context, in *IDsRequest, opts ...grpc.CallOption) (*IDsResponse, error)
	Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error)
	CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*CreateKeygroupResponse, error)
	DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*DeleteKeygroupResponse, error)
	ExistsKeygroup(ctx context.Context, in *ExistsKeygroupRequest, opts ...grpc.CallOption) (*ExistsKeygroupResponse, error)
	AddKeygroupTrigger(ctx context.Context, in *AddKeygroupTriggerRequest, opts ...grpc.CallOption) (*AddKeygroupTriggerResponse, error)
	DeleteKeygroupTrigger(ctx context.Context, in *DeleteKeygroupTriggerRequest, opts ...grpc.CallOption) (*DeleteKeygroupTriggerResponse, error)
	GetKeygroupTrigger(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error)
}

type databaseClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseClient(cc grpc.ClientConnInterface) DatabaseClient {
	return &databaseClient{cc}
}

func (c *databaseClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*UpdateResponse, error) {
	out := new(UpdateResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Scan(ctx context.Context, in *ScanRequest, opts ...grpc.CallOption) (*ScanResponse, error) {
	out := new(ScanResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Scan", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ReadAll(ctx context.Context, in *ReadAllRequest, opts ...grpc.CallOption) (*ReadAllResponse, error) {
	out := new(ReadAllResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/ReadAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) IDs(ctx context.Context, in *IDsRequest, opts ...grpc.CallOption) (*IDsResponse, error) {
	out := new(IDsResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/IDs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Exists(ctx context.Context, in *ExistsRequest, opts ...grpc.CallOption) (*ExistsResponse, error) {
	out := new(ExistsResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) CreateKeygroup(ctx context.Context, in *CreateKeygroupRequest, opts ...grpc.CallOption) (*CreateKeygroupResponse, error) {
	out := new(CreateKeygroupResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/CreateKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteKeygroup(ctx context.Context, in *DeleteKeygroupRequest, opts ...grpc.CallOption) (*DeleteKeygroupResponse, error) {
	out := new(DeleteKeygroupResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/DeleteKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ExistsKeygroup(ctx context.Context, in *ExistsKeygroupRequest, opts ...grpc.CallOption) (*ExistsKeygroupResponse, error) {
	out := new(ExistsKeygroupResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/ExistsKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) AddKeygroupTrigger(ctx context.Context, in *AddKeygroupTriggerRequest, opts ...grpc.CallOption) (*AddKeygroupTriggerResponse, error) {
	out := new(AddKeygroupTriggerResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/AddKeygroupTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteKeygroupTrigger(ctx context.Context, in *DeleteKeygroupTriggerRequest, opts ...grpc.CallOption) (*DeleteKeygroupTriggerResponse, error) {
	out := new(DeleteKeygroupTriggerResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/DeleteKeygroupTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetKeygroupTrigger(ctx context.Context, in *GetKeygroupTriggerRequest, opts ...grpc.CallOption) (*GetKeygroupTriggerResponse, error) {
	out := new(GetKeygroupTriggerResponse)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/GetKeygroupTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServer is the server API for Database service.
// All implementations should embed UnimplementedDatabaseServer
// for forward compatibility
type DatabaseServer interface {
	Update(context.Context, *UpdateRequest) (*UpdateResponse, error)
	Delete(context.Context, *DeleteRequest) (*DeleteResponse, error)
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Scan(context.Context, *ScanRequest) (*ScanResponse, error)
	ReadAll(context.Context, *ReadAllRequest) (*ReadAllResponse, error)
	IDs(context.Context, *IDsRequest) (*IDsResponse, error)
	Exists(context.Context, *ExistsRequest) (*ExistsResponse, error)
	CreateKeygroup(context.Context, *CreateKeygroupRequest) (*CreateKeygroupResponse, error)
	DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*DeleteKeygroupResponse, error)
	ExistsKeygroup(context.Context, *ExistsKeygroupRequest) (*ExistsKeygroupResponse, error)
	AddKeygroupTrigger(context.Context, *AddKeygroupTriggerRequest) (*AddKeygroupTriggerResponse, error)
	DeleteKeygroupTrigger(context.Context, *DeleteKeygroupTriggerRequest) (*DeleteKeygroupTriggerResponse, error)
	GetKeygroupTrigger(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error)
}

// UnimplementedDatabaseServer should be embedded to have forward compatible implementations.
type UnimplementedDatabaseServer struct {
}

func (UnimplementedDatabaseServer) Update(context.Context, *UpdateRequest) (*UpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDatabaseServer) Delete(context.Context, *DeleteRequest) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDatabaseServer) Append(context.Context, *AppendRequest) (*AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedDatabaseServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedDatabaseServer) Scan(context.Context, *ScanRequest) (*ScanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Scan not implemented")
}
func (UnimplementedDatabaseServer) ReadAll(context.Context, *ReadAllRequest) (*ReadAllResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAll not implemented")
}
func (UnimplementedDatabaseServer) IDs(context.Context, *IDsRequest) (*IDsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IDs not implemented")
}
func (UnimplementedDatabaseServer) Exists(context.Context, *ExistsRequest) (*ExistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedDatabaseServer) CreateKeygroup(context.Context, *CreateKeygroupRequest) (*CreateKeygroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeygroup not implemented")
}
func (UnimplementedDatabaseServer) DeleteKeygroup(context.Context, *DeleteKeygroupRequest) (*DeleteKeygroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroup not implemented")
}
func (UnimplementedDatabaseServer) ExistsKeygroup(context.Context, *ExistsKeygroupRequest) (*ExistsKeygroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExistsKeygroup not implemented")
}
func (UnimplementedDatabaseServer) AddKeygroupTrigger(context.Context, *AddKeygroupTriggerRequest) (*AddKeygroupTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddKeygroupTrigger not implemented")
}
func (UnimplementedDatabaseServer) DeleteKeygroupTrigger(context.Context, *DeleteKeygroupTriggerRequest) (*DeleteKeygroupTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroupTrigger not implemented")
}
func (UnimplementedDatabaseServer) GetKeygroupTrigger(context.Context, *GetKeygroupTriggerRequest) (*GetKeygroupTriggerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKeygroupTrigger not implemented")
}

// UnsafeDatabaseServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServer will
// result in compilation errors.
type UnsafeDatabaseServer interface {
	mustEmbedUnimplementedDatabaseServer()
}

func RegisterDatabaseServer(s grpc.ServiceRegistrar, srv DatabaseServer) {
	s.RegisterService(&Database_ServiceDesc, srv)
}

func _Database_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Scan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ScanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Scan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Scan",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Scan(ctx, req.(*ScanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ReadAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ReadAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/ReadAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ReadAll(ctx, req.(*ReadAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_IDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).IDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/IDs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).IDs(ctx, req.(*IDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).Exists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/Exists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).Exists(ctx, req.(*ExistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_CreateKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).CreateKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/CreateKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).CreateKeygroup(ctx, req.(*CreateKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).DeleteKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/DeleteKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).DeleteKeygroup(ctx, req.(*DeleteKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ExistsKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExistsKeygroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).ExistsKeygroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/ExistsKeygroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).ExistsKeygroup(ctx, req.(*ExistsKeygroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_AddKeygroupTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddKeygroupTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).AddKeygroupTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/AddKeygroupTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).AddKeygroupTrigger(ctx, req.(*AddKeygroupTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteKeygroupTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteKeygroupTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).DeleteKeygroupTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/DeleteKeygroupTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).DeleteKeygroupTrigger(ctx, req.(*DeleteKeygroupTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetKeygroupTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeygroupTriggerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServer).GetKeygroupTrigger(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mcc.fred.storage.Database/GetKeygroupTrigger",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServer).GetKeygroupTrigger(ctx, req.(*GetKeygroupTriggerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Database_ServiceDesc is the grpc.ServiceDesc for Database service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Database_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mcc.fred.storage.Database",
	HandlerType: (*DatabaseServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _Database_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Database_Delete_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _Database_Append_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _Database_Read_Handler,
		},
		{
			MethodName: "Scan",
			Handler:    _Database_Scan_Handler,
		},
		{
			MethodName: "ReadAll",
			Handler:    _Database_ReadAll_Handler,
		},
		{
			MethodName: "IDs",
			Handler:    _Database_IDs_Handler,
		},
		{
			MethodName: "Exists",
			Handler:    _Database_Exists_Handler,
		},
		{
			MethodName: "CreateKeygroup",
			Handler:    _Database_CreateKeygroup_Handler,
		},
		{
			MethodName: "DeleteKeygroup",
			Handler:    _Database_DeleteKeygroup_Handler,
		},
		{
			MethodName: "ExistsKeygroup",
			Handler:    _Database_ExistsKeygroup_Handler,
		},
		{
			MethodName: "AddKeygroupTrigger",
			Handler:    _Database_AddKeygroupTrigger_Handler,
		},
		{
			MethodName: "DeleteKeygroupTrigger",
			Handler:    _Database_DeleteKeygroupTrigger_Handler,
		},
		{
			MethodName: "GetKeygroupTrigger",
			Handler:    _Database_GetKeygroupTrigger_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "storage.proto",
}
