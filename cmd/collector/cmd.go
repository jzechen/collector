/**
 * @Time: 2023/10/20 15:39
 * @Author: jzechen
 * @File: cmd.go
 * @Software: GoLand collector
 */

package main

import (
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/component-base/logs"
	"os"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := NewCollectorManager()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
