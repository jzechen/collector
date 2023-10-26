/**
 * @Time: 2023/10/20 16:26
 * @Author: jzechen
 * @File: sina.go
 * @Software: GoLand collector
 */

package sina

import (
	"context"
	"fmt"
	"github.com/jzechen/toresa/pkg/manager/dto"
)

type Handler struct {
	ctx context.Context
}

func NewSinaHandler(ctx context.Context) *Handler {
	sinaHandler := &Handler{
		ctx: ctx,
	}

	return sinaHandler
}

// TODO: add handler implement here
func (hd *Handler) Hello(ctx context.Context, req *dto.NullRsp) (*dto.NullRsp, error) {
	fmt.Print("Hello")
	return &dto.NullRsp{}, nil
}
