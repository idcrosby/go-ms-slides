package catalogue

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

// START OMIT
type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		ListEndpoint:   MakeListEndpoint(s),
		GetEndpoint:    MakeGetEndpoint(s),
	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		socks, err := s.List(req.Order, req.PageNum, req.PageSize)
		return listResponse{Socks: socks, Err: err}, nil
	}
}
// END OMIT
// START 2-OMIT
func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		sock, err := s.Get(req.ID)
		return getResponse{Sock: sock, Err: err}, nil
	}
}

type listRequest struct {
	Order    string   `json:"order"`
	PageNum  int      `json:"pageNum"`
	PageSize int      `json:"pageSize"`
}
type listResponse struct {
	Socks []Sock `json:"sock"`
	Err   error  `json:"err"`
}

type getRequest struct {
	ID string `json:"id"`
}
type getResponse struct {
	Sock Sock  `json:"sock"`
	Err  error `json:"err"`
}
// END 2-OMIT