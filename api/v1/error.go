package v1

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ErrOffsetOutOfRange represents an error found when the offset
// is too large
type ErrOffsetOutOfRange struct {
	Offset uint64
}

// GRPCStatus implements the GRPC status interface
func (e ErrOffsetOutOfRange) GRPCStatus() *status.Status {
	st := status.New(codes.NotFound, fmt.Sprintf("offset out of range: %d", e.Offset))
	msg := fmt.Sprintf("The requested offset is outside the log's range: %d", e.Offset)

	d := &errdetails.LocalizedMessage{
		Locale:  "en-GB",
		Message: msg,
	}
	std, err := st.WithDetails(d)
	if err != nil {
		return st
	}
	return std
}

// Error implements the error interface
func (e ErrOffsetOutOfRange) Error() string {
	return e.GRPCStatus().Err().Error()
}