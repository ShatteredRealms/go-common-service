package srv

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrPermissionDenied = status.New(codes.PermissionDenied, "unauthorized").Err()
)
