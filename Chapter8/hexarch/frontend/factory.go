package frontend

import (
	"fmt"
)

func NewFrontEnd(s string) (FrontEnd, error) {
	switch s {
	case "rest":
		return &restFrontEnd{}, nil
	case "grpc":
		return &grpcFrontEnd{}, nil
	case "":
		return nil, fmt.Errorf("frontend type not defined")
	default:
		return nil, fmt.Errorf("no such frontend %s", s)

	}
}
