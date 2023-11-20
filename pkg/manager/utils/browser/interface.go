/**
 * @Time: 2023/11/20 10:06
 * @Author: jzechen
 * @File: interface.go
 * @Software: GoLand toresa
 */

package browser

import "github.com/chromedp/chromedp"

type Interface interface {
	// Run the actions using the browser that is created before.
	Run(actions ...chromedp.Action) error
	// Close the browser session.
	Close()
}
