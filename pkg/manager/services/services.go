/**
 * @Time: 2023/10/20 16:24
 * @Author: jzechen
 * @File: services.go
 * @Software: GoLand collector
 */

package services

import (
	"github.com/jzechen/collector/pkg/manager/services/sina"
)

type Interface interface {
	sina.Interface
}
