package wrapper

import (
	"context"
	"github.com/micro/go-micro/server"
)

func TokenVerify(wrap server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context,req server.Request,resp interface{}) error {
		req.Body()
	}
}
