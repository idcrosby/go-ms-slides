package catalogue

import (
	"strings"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware decorates a service.
type Middleware func(Service) Service

// START OMIT
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}
// END OMIT

// START 2-OMIT
func (mw loggingMiddleware) List(order string, pageNum, pageSize int) (socks []Sock, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "List",
			"order", order,
			"pageNum", pageNum,
			"pageSize", pageSize,
			"result", len(socks),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.List(tags, order, pageNum, pageSize)
}
// END 2-OMIT

func (mw loggingMiddleware) Get(id string) (s Sock, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "Get",
			"id", id,
			"sock", s.ID,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.Get(id)
}

