package endpoints

type Request interface {
	validate() error
}

// SumRequest collects the request parameters for the Sum method.
type SumRequest struct {
	A int64 `json:"a"` // a
	B int64 `json:"b"` // b
}

func (r SumRequest) validate() error {
	return nil // TBA
}