package errorx

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidInfo  = "无效数据"
	NotFoundInfo = "请求实体没有找到"
	InternalInfo = "内部错误"
)

func HandlerRpcError(code codes.Code, filed, msg string) error {
	var errorStatus *status.Status

	switch code {
	case codes.InvalidArgument:
		errorStatus = status.New(codes.InvalidArgument, InvalidInfo)
	case codes.NotFound:
		errorStatus = status.New(codes.NotFound, NotFoundInfo)
	case codes.Internal:
		errorStatus = status.New(codes.Internal, InternalInfo)
	default:
	}

	ds, err := errorStatus.WithDetails(
		&errdetails.BadRequest_FieldViolation{
			Field:       filed,
			Description: msg,
		},
	)

	if err != nil {
		return ds.Err()
	}
	return errorStatus.Err()
}
