/**
 * @Time: 2023/10/20 15:53
 * @Author: jzechen
 * @File: log.go
 * @Software: GoLand collector
 */

package flag

import (
	"flag"
	"fmt"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var packageFlags = flag.NewFlagSet("logging", flag.ContinueOnError)

func init() {
	initFlag(packageFlags)
}

func AddLogFlags(fs *pflag.FlagSet, name string) {
	// Add all supported flags.
	packageFlags.VisitAll(func(f *flag.Flag) {
		pf := pflag.PFlagFromGoFlag(f)
		if fs.Lookup(pf.Name) == nil {
			fs.AddFlag(pf)
		}
	})
	fs.BoolP("help", "h", false, fmt.Sprintf("help for %s", name))
}

func initFlag(fs *flag.FlagSet) {
	var allFlags flag.FlagSet
	klog.InitFlags(&allFlags)
	if fs == nil {
		fs = flag.CommandLine
	}

	allFlags.VisitAll(func(f *flag.Flag) {
		//switch f.Name {
		//case "v", "vmodule":
		fs.Var(f.Value, f.Name, f.Usage)
		//}
	})
}
