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
	Update(ctx context.Context, in *UpdateItem, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error)
	Append(ctx context.Context, in *AppendItem, opts ...grpc.CallOption) (*Key, error)
	Read(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Val, error)
	ReadAll(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_ReadAllClient, error)
	IDs(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_IDsClient, error)
	Exists(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error)
	CreateKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error)
	DeleteKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error)
	ExistsKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error)
	AddKeygroupTrigger(ctx context.Context, in *KeygroupTrigger, opts ...grpc.CallOption) (*Response, error)
	DeleteKeygroupTrigger(ctx context.Context, in *KeygroupTrigger, opts ...grpc.CallOption) (*Response, error)
	GetKeygroupTrigger(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_GetKeygroupTriggerClient, error)
}

type databaseClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseClient(cc grpc.ClientConnInterface) DatabaseClient {
	return &databaseClient{cc}
}

func (c *databaseClient) Update(ctx context.Context, in *UpdateItem, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Delete(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Append(ctx context.Context, in *AppendItem, opts ...grpc.CallOption) (*Key, error) {
	out := new(Key)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) Read(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Val, error) {
	out := new(Val)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ReadAll(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_ReadAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &Database_ServiceDesc.Streams[0], "/mcc.fred.storage.Database/ReadAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &databaseReadAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Database_ReadAllClient interface {
	Recv() (*Item, error)
	grpc.ClientStream
}

type databaseReadAllClient struct {
	grpc.ClientStream
}

func (x *databaseReadAllClient) Recv() (*Item, error) {
	m := new(Item)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *databaseClient) IDs(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_IDsClient, error) {
	stream, err := c.cc.NewStream(ctx, &Database_ServiceDesc.Streams[1], "/mcc.fred.storage.Database/IDs", opts...)
	if err != nil {
		return nil, err
	}
	x := &databaseIDsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Database_IDsClient interface {
	Recv() (*Key, error)
	grpc.ClientStream
}

type databaseIDsClient struct {
	grpc.ClientStream
}

func (x *databaseIDsClient) Recv() (*Key, error) {
	m := new(Key)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *databaseClient) Exists(ctx context.Context, in *Key, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/Exists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) CreateKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/CreateKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/DeleteKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) ExistsKeygroup(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/ExistsKeygroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) AddKeygroupTrigger(ctx context.Context, in *KeygroupTrigger, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/AddKeygroupTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) DeleteKeygroupTrigger(ctx context.Context, in *KeygroupTrigger, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/mcc.fred.storage.Database/DeleteKeygroupTrigger", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseClient) GetKeygroupTrigger(ctx context.Context, in *Keygroup, opts ...grpc.CallOption) (Database_GetKeygroupTriggerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Database_ServiceDesc.Streams[2], "/mcc.fred.storage.Database/GetKeygroupTrigger", opts...)
	if err != nil {
		return nil, err
	}
	x := &databaseGetKeygroupTriggerClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Database_GetKeygroupTriggerClient interface {
	Recv() (*Trigger, error)
	grpc.ClientStream
}

type databaseGetKeygroupTriggerClient struct {
	grpc.ClientStream
}

func (x *databaseGetKeygroupTriggerClient) Recv() (*Trigger, error) {
	m := new(Trigger)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DatabaseServer is the server API for Database service.
// All implementations should embed UnimplementedDatabaseServer
// for forward compatibility
type DatabaseServer interface {
	Update(context.Context, *UpdateItem) (*Response, error)
	Delete(context.Context, *Key) (*Response, error)
	Append(context.Context, *AppendItem) (*Key, error)
	Read(context.Context, *Key) (*Val, error)
	ReadAll(*Keygroup, Database_ReadAllServer) error
	IDs(*Keygroup, Database_IDsServer) error
	Exists(context.Context, *Key) (*Response, error)
	CreateKeygroup(context.Context, *Keygroup) (*Response, error)
	DeleteKeygroup(context.Context, *Keygroup) (*Response, error)
	ExistsKeygroup(context.Context, *Keygroup) (*Response, error)
	AddKeygroupTrigger(context.Context, *KeygroupTrigger) (*Response, error)
	DeleteKeygroupTrigger(context.Context, *KeygroupTrigger) (*Response, error)
	GetKeygroupTrigger(*Keygroup, Database_GetKeygroupTriggerServer) error
}

// UnimplementedDatabaseServer should be embedded to have forward compatible implementations.
type UnimplementedDatabaseServer struct {
}

func (UnimplementedDatabaseServer) Update(context.Context, *UpdateItem) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDatabaseServer) Delete(context.Context, *Key) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedDatabaseServer) Append(context.Context, *AppendItem) (*Key, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedDatabaseServer) Read(context.Context, *Key) (*Val, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedDatabaseServer) ReadAll(*Keygroup, Database_ReadAllServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadAll not implemented")
}
func (UnimplementedDatabaseServer) IDs(*Keygroup, Database_IDsServer) error {
	return status.Errorf(codes.Unimplemented, "method IDs not implemented")
}
func (UnimplementedDatabaseServer) Exists(context.Context, *Key) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Exists not implemented")
}
func (UnimplementedDatabaseServer) CreateKeygroup(context.Context, *Keygroup) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateKeygroup not implemented")
}
func (UnimplementedDatabaseServer) DeleteKeygroup(context.Context, *Keygroup) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroup not implemented")
}
func (UnimplementedDatabaseServer) ExistsKeygroup(context.Context, *Keygroup) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExistsKeygroup not implemented")
}
func (UnimplementedDatabaseServer) AddKeygroupTrigger(context.Context, *KeygroupTrigger) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddKeygroupTrigger not implemented")
}
func (UnimplementedDatabaseServer) DeleteKeygroupTrigger(context.Context, *KeygroupTrigger) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteKeygroupTrigger not implemented")
}
func (UnimplementedDatabaseServer) GetKeygroupTrigger(*Keygroup, Database_GetKeygroupTriggerServer) error {
	return status.Errorf(codes.Unimplemented, "method GetKeygroupTrigger not implemented")
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
	in := new(UpdateItem)
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
		return srv.(DatabaseServer).Update(ctx, req.(*UpdateItem))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
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
		return srv.(DatabaseServer).Delete(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendItem)
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
		return srv.(DatabaseServer).Append(ctx, req.(*AppendItem))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
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
		return srv.(DatabaseServer).Read(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ReadAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Keygroup)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DatabaseServer).ReadAll(m, &databaseReadAllServer{stream})
}

type Database_ReadAllServer interface {
	Send(*Item) error
	grpc.ServerStream
}

type databaseReadAllServer struct {
	grpc.ServerStream
}

func (x *databaseReadAllServer) Send(m *Item) error {
	return x.ServerStream.SendMsg(m)
}

func _Database_IDs_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Keygroup)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DatabaseServer).IDs(m, &databaseIDsServer{stream})
}

type Database_IDsServer interface {
	Send(*Key) error
	grpc.ServerStream
}

type databaseIDsServer struct {
	grpc.ServerStream
}

func (x *databaseIDsServer) Send(m *Key) error {
	return x.ServerStream.SendMsg(m)
}

func _Database_Exists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Key)
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
		return srv.(DatabaseServer).Exists(ctx, req.(*Key))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_CreateKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Keygroup)
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
		return srv.(DatabaseServer).CreateKeygroup(ctx, req.(*Keygroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Keygroup)
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
		return srv.(DatabaseServer).DeleteKeygroup(ctx, req.(*Keygroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_ExistsKeygroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Keygroup)
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
		return srv.(DatabaseServer).ExistsKeygroup(ctx, req.(*Keygroup))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_AddKeygroupTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeygroupTrigger)
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
		return srv.(DatabaseServer).AddKeygroupTrigger(ctx, req.(*KeygroupTrigger))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_DeleteKeygroupTrigger_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeygroupTrigger)
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
		return srv.(DatabaseServer).DeleteKeygroupTrigger(ctx, req.(*KeygroupTrigger))
	}
	return interceptor(ctx, in, info, handler)
}

func _Database_GetKeygroupTrigger_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Keygroup)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DatabaseServer).GetKeygroupTrigger(m, &databaseGetKeygroupTriggerServer{stream})
}

type Database_GetKeygroupTriggerServer interface {
	Send(*Trigger) error
	grpc.ServerStream
}

type databaseGetKeygroupTriggerServer struct {
	grpc.ServerStream
}

func (x *databaseGetKeygroupTriggerServer) Send(m *Trigger) error {
	return x.ServerStream.SendMsg(m)
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
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReadAll",
			Handler:       _Database_ReadAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "IDs",
			Handler:       _Database_IDs_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetKeygroupTrigger",
			Handler:       _Database_GetKeygroupTrigger_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "storage.proto",
}
