package errorx

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	InvalidInfo  = "Invalid information received"
	NotFoundInfo = "Request entity not found"
)

func HandlerRpcError(code codes.Code, filed, msg string) error {
	var errorStatus *status.Status

	switch code {
	case codes.InvalidArgument:
		errorStatus = status.New(codes.InvalidArgument, InvalidInfo)
	case codes.NotFound:
		errorStatus = status.New(codes.NotFound, NotFoundInfo)
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
