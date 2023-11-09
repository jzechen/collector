/**
 * @Time: 2023/10/20 16:26
 * @Author: jzechen
 * @File: sina.go
 * @Software: GoLand collector
 */

package sina

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/jzechen/toresa/pkg/manager/config"
	"github.com/jzechen/toresa/pkg/manager/dto"
	"github.com/jzechen/toresa/pkg/manager/mdb"
	"github.com/jzechen/toresa/pkg/manager/utils/browser"
	"github.com/tebeka/selenium"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog/v2"
)

type Handler struct {
	conf    *config.CollectorManager
	account *mongo.Collection
}

func NewSinaHandler(cfg *config.CollectorManager, mongo mdb.Interface) *Handler {
	ac, err := mongo.GetCollection(cfg.Mongo.Database, "account")
	if err != nil {
		klog.Fatal(err)
	}
	sinaHandler := &Handler{
		conf:    cfg,
		account: ac,
	}

	return sinaHandler
}

// TODO: add handler implement here
func (hd *Handler) Hello(ctx context.Context, req *dto.NullRsp) (*dto.NullRsp, error) {
	fmt.Print("Hello")
	return &dto.NullRsp{}, nil
}

func (hd *Handler) Login(ctx context.Context, req *dto.LoginReq) (*dto.NullRsp, error) {
	cookie, err := getCookieStr(ctx, &hd.conf.Drive, req.UserID, req.Password)
	if err != nil {
		return nil, err
	}
	err = hd.saveToMgo(req.UserID, req.Password, cookie)
	if err != nil {
		return nil, err
	}

	return &dto.NullRsp{}, nil
}

func getCookieStr(ctx context.Context, dc *config.DriveConfig, userName, password string) (string, error) {
	ctx, cancel := browser.GetChromeAllocateFunc()(ctx, dc)
	defer cancel()

	// Connect to the WebDriver instance running locally.
	//caps := selenium.Capabilities{"browserName": "chrome"}

	// Disable image loading to speed up rendering
	//imagCaps := map[string]interface{}{
	//	"profile.managed_default_content_settings.images": 2,
	//}

	if err := chromedp.Run(ctx,
		// navigate to sina weibo login page
		chromedp.Navigate("https://passport.weibo.cn/signin/login?entry=mweibo&r=https://weibo.cn/"),
		// wait for footer element is visible (ie, page is loaded)
		chromedp.WaitVisible(`body > footer`),
		// input the userName and password
		// TODO: Handle verifyCodeImage code input
		chromedp.SetValue("#loginName", userName),
		chromedp.SetValue("#loginPassword", password),
		// click and submit login request
		chromedp.Click("#loginAction", chromedp.NodeVisible),
	); err != nil {
		return "", err
	}
	// TODO: get cookie

	//cobra.CheckErr(drive.Wait(func(wdtemp selenium.WebDriver) (b bool, e error) {
	//	tit, err := wdtemp.Title()
	//	if err != nil {
	//		return false, nil
	//	}
	//	if tit != "我的首页" {
	//		return false, nil
	//	}
	//	return true, nil
	//}))
	//cookie, err := drive.GetCookies()
	//var cookieSlice []string
	//for _, c := range cookie {
	//	cookieSlice = append(cookieSlice, c.Name+"="+c.Value)
	//}
	return "cookie", nil
}

func isDisplayed(by, elementName string) func(selenium.WebDriver) (bool, error) {
	return func(wd selenium.WebDriver) (bool, error) {
		el, err := wd.FindElement(by, elementName)
		if err != nil {
			return false, nil
		}
		enabled, err := el.IsDisplayed()
		if err != nil {
			return false, nil
		}

		if !enabled {
			return false, nil
		}

		return true, nil
	}
}

func (hd *Handler) saveToMgo(userName, password, cookie string) error {
	if cookie == "" {
		return errors.New("null cookie")
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{
					Key:   "password",
					Value: password,
				},
				{
					Key:   "cookie",
					Value: cookie,
				},
				{
					Key:   "status",
					Value: "success",
				},
			},
		},
	}

	// 设置更新选项
	opts := options.Update().SetUpsert(true)

	// 执行upsert操作
	filter := bson.D{{"_id", userName}}
	updateResult, err := hd.account.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	klog.V(2).InfoS("upsert a document in the account collection", "_id", updateResult.UpsertedID, "matchCount", updateResult.MatchedCount, "modifiedCount", updateResult.ModifiedCount, "upsertCount", updateResult.UpsertedCount)
	return nil
}
