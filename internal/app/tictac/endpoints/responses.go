package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/ms-sample/internal/app/tictac/service"
	"github.com/cage1016/ms-sample/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*TicResponse)(nil)

	_ httptransport.StatusCoder = (*TicResponse)(nil)

	_ httptransport.Headerer = (*TacResponse)(nil)

	_ httptransport.StatusCoder = (*TacResponse)(nil)
)

// TicResponse collects the response values for the Tic method.
type TicResponse struct {
	Err error `json:"-"`
}

func (r TicResponse) StatusCode() int {
	return http.StatusNoContent
}

func (r TicResponse) Headers() http.Header {
	return http.Header{}
}

func (r TicResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version}
}

// TacResponse collects the response values for the Tac method.
type TacResponse struct {
	Res int64 `json:"res"`
	Err error `json:"-"`
}

func (r TacResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r TacResponse) Headers() http.Header {
	return http.Header{}
}

func (r TacResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}
