package endpoints

type Request interface {
	validate() error
}

// TicRequest collects the request parameters for the Tic method.
type TicRequest struct {
}

func (r TicRequest) validate() error {
	return nil // TBA
}

// TacRequest collects the request parameters for the Tac method.
type TacRequest struct {
}

func (r TacRequest) validate() error {
	return nil // TBA
}
