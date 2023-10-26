/**
 * @Time: 2023/10/20 15:44
 * @Author: jzechen
 * @File: collector-manager.go
 * @Software: GoLand collector
 */

package cmd

import (
	"context"
	"fmt"
	"github.com/jzechen/toresa/cmd/manager/cmd/options"
	"github.com/jzechen/toresa/pkg/common/apiserver"
	logFlag "github.com/jzechen/toresa/pkg/common/flag"
	"github.com/jzechen/toresa/pkg/manager/config"
	"github.com/jzechen/toresa/pkg/manager/server"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

const (
	applicationName = "collector-manager"
	helpTextShort   = "A server application used to manage the collect jobs in the collector project"
	helpTextLong    = `The Collector-manager is a module that manages the collect job in the collector project.

      Find more information at:
            https://github.com/jzechen/toresa`
)

func NewCollectorManager() *cobra.Command {
	opt := options.NewCollectorManagerOptions()

	cmd := &cobra.Command{
		Use:   applicationName,
		Short: helpTextShort,
		Long:  helpTextLong,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Flags().VisitAll(func(flag *pflag.Flag) {
				klog.V(2).Infof("FLAG: --%s=%q", flag.Name, flag.Value)
			})

			err := opt.Validate()
			if err != nil {
				klog.Fatalf("validate options: %v", err)
			}

			serverConfig := CreateManagerConfig(opt)
			klog.V(4).Infof("load config: %+v", *serverConfig)

			if err := Run(apiserver.SetupSignalContext(), serverConfig); err != nil {
				klog.Exitf("run failed with %v", err)
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	opt.AddFlags(cmd.Flags())
	logFlag.AddLogFlags(cmd.Flags(), cmd.Name())

	return cmd
}

func CreateManagerConfig(managerOptions *options.CollectorManagerOptions) *config.CollectorManager {
	cfg, conf, err := config.BuildConfig(managerOptions)
	cobra.CheckErr(err)

	config.SetDefaultConfig(cfg)
	klog.Info("init config successfully")
	cfg.OriginConf = conf

	return cfg
}

func Run(ctx context.Context, cfg *config.CollectorManager) error {
	klog.V(2).Info("manager server prepare to run")
	s, err := server.NewCollectorManagerServer(ctx, cfg)
	if err != nil {
		return err
	}

	klog.V(2).Info("start the manager server")
	s.Run()

	return nil
}
