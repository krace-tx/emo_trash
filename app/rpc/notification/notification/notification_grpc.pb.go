// goctl rpc protoc notification.proto --go_out=. --go-grpc_out=. --zrpc_out=.  -m

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: notification.proto

package notification

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Notification_GetNotificationList_FullMethodName        = "/notification.Notification/GetNotificationList"
	Notification_MarkNotificationAsRead_FullMethodName     = "/notification.Notification/MarkNotificationAsRead"
	Notification_MarkAllNotificationsAsRead_FullMethodName = "/notification.Notification/MarkAllNotificationsAsRead"
	Notification_DeleteNotification_FullMethodName         = "/notification.Notification/DeleteNotification"
	Notification_DeleteAllNotifications_FullMethodName     = "/notification.Notification/DeleteAllNotifications"
	Notification_PushNotification_FullMethodName           = "/notification.Notification/PushNotification"
	Notification_GetUnreadNotificationCount_FullMethodName = "/notification.Notification/GetUnreadNotificationCount"
)

// NotificationClient is the client API for Notification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 通知管理服务接口，用于用户通知的获取和管理
type NotificationClient interface {
	// 获取用户的通知列表
	GetNotificationList(ctx context.Context, in *GetNotificationListRequest, opts ...grpc.CallOption) (*GetNotificationListResponse, error)
	// 标记通知为已读
	MarkNotificationAsRead(ctx context.Context, in *MarkNotificationAsReadRequest, opts ...grpc.CallOption) (*MarkNotificationAsReadResponse, error)
	// 标记所有通知为已读
	MarkAllNotificationsAsRead(ctx context.Context, in *MarkAllNotificationsAsReadRequest, opts ...grpc.CallOption) (*MarkAllNotificationsAsReadResponse, error)
	// 删除通知
	DeleteNotification(ctx context.Context, in *DeleteNotificationRequest, opts ...grpc.CallOption) (*DeleteNotificationResponse, error)
	// 删除所有通知
	DeleteAllNotifications(ctx context.Context, in *DeleteAllNotificationsRequest, opts ...grpc.CallOption) (*DeleteAllNotificationsResponse, error)
	// 推送通知
	PushNotification(ctx context.Context, in *PushNotificationRequest, opts ...grpc.CallOption) (*PushNotificationResponse, error)
	// 获取未读通知数量
	GetUnreadNotificationCount(ctx context.Context, in *GetUnreadNotificationCountRequest, opts ...grpc.CallOption) (*GetUnreadNotificationCountResponse, error)
}

type notificationClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationClient(cc grpc.ClientConnInterface) NotificationClient {
	return &notificationClient{cc}
}

func (c *notificationClient) GetNotificationList(ctx context.Context, in *GetNotificationListRequest, opts ...grpc.CallOption) (*GetNotificationListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetNotificationListResponse)
	err := c.cc.Invoke(ctx, Notification_GetNotificationList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) MarkNotificationAsRead(ctx context.Context, in *MarkNotificationAsReadRequest, opts ...grpc.CallOption) (*MarkNotificationAsReadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MarkNotificationAsReadResponse)
	err := c.cc.Invoke(ctx, Notification_MarkNotificationAsRead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) MarkAllNotificationsAsRead(ctx context.Context, in *MarkAllNotificationsAsReadRequest, opts ...grpc.CallOption) (*MarkAllNotificationsAsReadResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MarkAllNotificationsAsReadResponse)
	err := c.cc.Invoke(ctx, Notification_MarkAllNotificationsAsRead_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) DeleteNotification(ctx context.Context, in *DeleteNotificationRequest, opts ...grpc.CallOption) (*DeleteNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteNotificationResponse)
	err := c.cc.Invoke(ctx, Notification_DeleteNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) DeleteAllNotifications(ctx context.Context, in *DeleteAllNotificationsRequest, opts ...grpc.CallOption) (*DeleteAllNotificationsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAllNotificationsResponse)
	err := c.cc.Invoke(ctx, Notification_DeleteAllNotifications_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) PushNotification(ctx context.Context, in *PushNotificationRequest, opts ...grpc.CallOption) (*PushNotificationResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PushNotificationResponse)
	err := c.cc.Invoke(ctx, Notification_PushNotification_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) GetUnreadNotificationCount(ctx context.Context, in *GetUnreadNotificationCountRequest, opts ...grpc.CallOption) (*GetUnreadNotificationCountResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUnreadNotificationCountResponse)
	err := c.cc.Invoke(ctx, Notification_GetUnreadNotificationCount_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServer is the server API for Notification service.
// All implementations must embed UnimplementedNotificationServer
// for forward compatibility.
//
// 通知管理服务接口，用于用户通知的获取和管理
type NotificationServer interface {
	// 获取用户的通知列表
	GetNotificationList(context.Context, *GetNotificationListRequest) (*GetNotificationListResponse, error)
	// 标记通知为已读
	MarkNotificationAsRead(context.Context, *MarkNotificationAsReadRequest) (*MarkNotificationAsReadResponse, error)
	// 标记所有通知为已读
	MarkAllNotificationsAsRead(context.Context, *MarkAllNotificationsAsReadRequest) (*MarkAllNotificationsAsReadResponse, error)
	// 删除通知
	DeleteNotification(context.Context, *DeleteNotificationRequest) (*DeleteNotificationResponse, error)
	// 删除所有通知
	DeleteAllNotifications(context.Context, *DeleteAllNotificationsRequest) (*DeleteAllNotificationsResponse, error)
	// 推送通知
	PushNotification(context.Context, *PushNotificationRequest) (*PushNotificationResponse, error)
	// 获取未读通知数量
	GetUnreadNotificationCount(context.Context, *GetUnreadNotificationCountRequest) (*GetUnreadNotificationCountResponse, error)
	mustEmbedUnimplementedNotificationServer()
}

// UnimplementedNotificationServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedNotificationServer struct{}

func (UnimplementedNotificationServer) GetNotificationList(context.Context, *GetNotificationListRequest) (*GetNotificationListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNotificationList not implemented")
}
func (UnimplementedNotificationServer) MarkNotificationAsRead(context.Context, *MarkNotificationAsReadRequest) (*MarkNotificationAsReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkNotificationAsRead not implemented")
}
func (UnimplementedNotificationServer) MarkAllNotificationsAsRead(context.Context, *MarkAllNotificationsAsReadRequest) (*MarkAllNotificationsAsReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkAllNotificationsAsRead not implemented")
}
func (UnimplementedNotificationServer) DeleteNotification(context.Context, *DeleteNotificationRequest) (*DeleteNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteNotification not implemented")
}
func (UnimplementedNotificationServer) DeleteAllNotifications(context.Context, *DeleteAllNotificationsRequest) (*DeleteAllNotificationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllNotifications not implemented")
}
func (UnimplementedNotificationServer) PushNotification(context.Context, *PushNotificationRequest) (*PushNotificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushNotification not implemented")
}
func (UnimplementedNotificationServer) GetUnreadNotificationCount(context.Context, *GetUnreadNotificationCountRequest) (*GetUnreadNotificationCountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUnreadNotificationCount not implemented")
}
func (UnimplementedNotificationServer) mustEmbedUnimplementedNotificationServer() {}
func (UnimplementedNotificationServer) testEmbeddedByValue()                      {}

// UnsafeNotificationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServer will
// result in compilation errors.
type UnsafeNotificationServer interface {
	mustEmbedUnimplementedNotificationServer()
}

func RegisterNotificationServer(s grpc.ServiceRegistrar, srv NotificationServer) {
	// If the following call pancis, it indicates UnimplementedNotificationServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Notification_ServiceDesc, srv)
}

func _Notification_GetNotificationList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNotificationListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).GetNotificationList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_GetNotificationList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).GetNotificationList(ctx, req.(*GetNotificationListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_MarkNotificationAsRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkNotificationAsReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).MarkNotificationAsRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_MarkNotificationAsRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).MarkNotificationAsRead(ctx, req.(*MarkNotificationAsReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_MarkAllNotificationsAsRead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MarkAllNotificationsAsReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).MarkAllNotificationsAsRead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_MarkAllNotificationsAsRead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).MarkAllNotificationsAsRead(ctx, req.(*MarkAllNotificationsAsReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_DeleteNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).DeleteNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_DeleteNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).DeleteNotification(ctx, req.(*DeleteNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_DeleteAllNotifications_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllNotificationsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).DeleteAllNotifications(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_DeleteAllNotifications_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).DeleteAllNotifications(ctx, req.(*DeleteAllNotificationsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_PushNotification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PushNotificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).PushNotification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_PushNotification_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).PushNotification(ctx, req.(*PushNotificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_GetUnreadNotificationCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUnreadNotificationCountRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).GetUnreadNotificationCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Notification_GetUnreadNotificationCount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).GetUnreadNotificationCount(ctx, req.(*GetUnreadNotificationCountRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Notification_ServiceDesc is the grpc.ServiceDesc for Notification service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notification_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "notification.Notification",
	HandlerType: (*NotificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetNotificationList",
			Handler:    _Notification_GetNotificationList_Handler,
		},
		{
			MethodName: "MarkNotificationAsRead",
			Handler:    _Notification_MarkNotificationAsRead_Handler,
		},
		{
			MethodName: "MarkAllNotificationsAsRead",
			Handler:    _Notification_MarkAllNotificationsAsRead_Handler,
		},
		{
			MethodName: "DeleteNotification",
			Handler:    _Notification_DeleteNotification_Handler,
		},
		{
			MethodName: "DeleteAllNotifications",
			Handler:    _Notification_DeleteAllNotifications_Handler,
		},
		{
			MethodName: "PushNotification",
			Handler:    _Notification_PushNotification_Handler,
		},
		{
			MethodName: "GetUnreadNotificationCount",
			Handler:    _Notification_GetUnreadNotificationCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "notification.proto",
}
