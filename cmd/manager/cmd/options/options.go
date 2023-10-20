/**
 * @Time: 2023/10/20 15:45
 * @Author: jzechen
 * @File: options.go
 * @Software: GoLand collector
 */

package options

import (
	"github.com/spf13/pflag"
	utilErrors "k8s.io/apimachinery/pkg/util/errors"
)

type CollectorManagerOptions struct {
	// the path of the config file
	CfgFile string `json:"cfgFile,omitempty"`
}

func (o *CollectorManagerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.CfgFile, "config", "c", "", "server config file")
}

func NewCollectorManagerOptions() *CollectorManagerOptions {
	option := &CollectorManagerOptions{}
	return option
}

func (o *CollectorManagerOptions) Validate() error {
	var ErrList []error

	return utilErrors.NewAggregate(ErrList)
}
