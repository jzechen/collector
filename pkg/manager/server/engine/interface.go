/**
 * @Time: 2023/10/20 16:14
 * @Author: jzechen
 * @File: interface.go
 * @Software: GoLand collector
 */

package engine

import "net/http"

type Interface interface {
	// CreateHandler register router handler in the engine and return http server handler
	CreateHandler() http.Handler
}
