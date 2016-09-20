package catalogue

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"golang.org/x/net/context"
)

type Endpoints struct {
	ListEndpoint   endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	HealthEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		ListEndpoint:   MakeListEndpoint(s),
		GetEndpoint:    MakeGetEndpoint(s),
		HealthEndpoint: MakeHealthEndpoint(s),
	}
}

func MakeListEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listRequest)
		socks, err := s.List(req.Tags, req.Order, req.PageNum, req.PageSize)
		return listResponse{Socks: socks, Err: err}, nil
	}
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getRequest)
		sock, err := s.Get(req.ID)
		return getResponse{Sock: sock, Err: err}, nil
	}
}

// MakeHealthEndpoint returns current health of the given service.
func MakeHealthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return healthResponse{Status: "OK", Time: time.Now().String()}, nil
	}
}

type listRequest struct {
	Tags     []string `json:"tags"`
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

type healthRequest struct {
	//
}

type healthResponse struct {
	Status string `json:"status"`
	Time   string `json:"time"`
}
