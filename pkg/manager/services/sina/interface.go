/**
 * @Time: 2023/10/20 16:25
 * @Author: jzechen
 * @File: interface.go
 * @Software: GoLand collector
 */

package sina

import (
	"context"
	"github.com/jzechen/toresa/pkg/manager/dto"
)

type Interface interface {
	// TODO: add collect sinaWeibo handler method here
	Hello(ctx context.Context, req *dto.NullRsp) (*dto.NullRsp, error)
}
