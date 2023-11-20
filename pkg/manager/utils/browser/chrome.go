/**
 * @Time: 2023/11/8 19:33
 * @Author: jzechen
 * @File: chrome.go
 * @Software: GoLand toresa
 */

package browser

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/jzechen/toresa/pkg/manager/config"
	. "github.com/jzechen/toresa/pkg/manager/contants"
)

type Chrome struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewChromeSession(ctx context.Context, dc *config.DriveConfig) Interface {
	ctx, cancel := GetChromeAllocateFunc()(ctx, dc)
	c := &Chrome{
		ctx:    ctx,
		cancel: cancel,
	}
	return c
}

func (c *Chrome) Run(actions ...chromedp.Action) error {
	return chromedp.Run(c.ctx, actions...)
}

func (c *Chrome) Close() {
	c.cancel()
}

type AllocateFunc func(ctx context.Context, dc *config.DriveConfig) (context.Context, context.CancelFunc)

func GetChromeAllocateFunc() AllocateFunc {
	return func(ctx context.Context, dc *config.DriveConfig) (context.Context, context.CancelFunc) {
		var allocCtx context.Context

		if dc.Type == DefaultDriveType {
			// create allocator context for use with remote debug mode chrome.
			// note: the remote chrome browse should start with "--headless" flag.
			// see: https://stackoverflow.com/questions/40538197/chrome-remote-debugging-from-another-machine
			allocCtx, _ = chromedp.NewRemoteAllocator(ctx, dc.Addr)
		} else {
			opts := append(chromedp.DefaultExecAllocatorOptions[:],
				// By default, Chrome will bypass localhost.
				// The test server is bound to localhost, so we should add the
				// following flag to use the proxy for localhost URLs.
				chromedp.Flag("", "<-loopback>"),
				// set Chrome headless mode and run under Linux, you need to set this parameter, otherwise an error will be reported.
				chromedp.Headless,
				// mock user-agent, inorder to beat Anti-spider
				chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"),
			)
			allocCtx, _ = chromedp.NewExecAllocator(ctx, opts...)
		}

		// create context
		return chromedp.NewContext(allocCtx)
	}
}
