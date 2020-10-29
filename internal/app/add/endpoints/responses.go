package endpoints

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/cage1016/ms-sample/internal/app/add/service"
	"github.com/cage1016/ms-sample/internal/pkg/responses"
)

var (
	_ httptransport.Headerer = (*SumResponse)(nil)

	_ httptransport.StatusCoder = (*SumResponse)(nil)
)

// SumResponse collects the response values for the Sum method.
type SumResponse struct {
	Res int64 `json:"res"`
	Err error `json:"-"`
}

func (r SumResponse) StatusCode() int {
	return http.StatusOK // TBA
}

func (r SumResponse) Headers() http.Header {
	return http.Header{}
}

func (r SumResponse) Response() interface{} {
	return responses.DataRes{APIVersion: service.Version, Data: r}
}
