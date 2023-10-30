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
	"github.com/jzechen/toresa/pkg/manager/config"
	"github.com/jzechen/toresa/pkg/manager/dto"
	"github.com/jzechen/toresa/pkg/manager/mdb"
	"github.com/spf13/cobra"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"k8s.io/klog/v2"
	"strings"
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
	cookie, err := getCookieStr(&hd.conf.Drive, req.UserID, req.Password)
	if err != nil {
		return nil, err
	}
	err = hd.saveToMgo(req.UserID, req.Password, cookie)
	if err != nil {
		return nil, err
	}

	return &dto.NullRsp{}, nil
}

func getCookieStr(dc *config.DriveConfig, userName, password string) (string, error) {
	// Start a Selenium WebDriver server instance
	svc, err := selenium.NewChromeDriverService(dc.Path, dc.Port)
	if err != nil {
		return "", fmt.Errorf("start a chromedriver service failed with %s", err)
	}
	//注意这里，server关闭之后，chrome窗口也会关闭
	defer svc.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{"browserName": "chrome"}
	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			//"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)
	drive, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", dc.Port))
	if err != nil {
		return "", fmt.Errorf("connect to the webDriver failded with %s", err)
	}
	defer drive.Quit()

	err = drive.Get("https://passport.weibo.cn/signin/login?entry=mweibo&r=https://weibo.cn/")
	if err != nil {
		return "", fmt.Errorf("get page failed with %s", err)
	}
	cobra.CheckErr(drive.Wait(isDisplayed(selenium.ByCSSSelector, "#loginName")))
	cobra.CheckErr(drive.Wait(isDisplayed(selenium.ByCSSSelector, "#loginPassword")))
	cobra.CheckErr(drive.Wait(isDisplayed(selenium.ByCSSSelector, "#loginAction")))
	username, err := drive.FindElement(selenium.ByCSSSelector, "#loginName")
	if err != nil {
		return "", fmt.Errorf("get username element failed with %s", err)
	}
	passwordElem, err := drive.FindElement(selenium.ByCSSSelector, "#loginPassword")
	if err != nil {
		return "", fmt.Errorf("get password element failed with %s", err)
	}
	submit, err := drive.FindElement(selenium.ByCSSSelector, "#loginAction")
	if err != nil {
		return "", fmt.Errorf("get login action element failed with %s", err)
	}
	cobra.CheckErr(username.SendKeys(userName))
	cobra.CheckErr(passwordElem.SendKeys(password))
	cobra.CheckErr(submit.Click())
	cobra.CheckErr(drive.Wait(func(wdtemp selenium.WebDriver) (b bool, e error) {
		tit, err := wdtemp.Title()
		if err != nil {
			return false, nil
		}
		if tit != "我的首页" {
			return false, nil
		}
		return true, nil
	}))
	cookie, err := drive.GetCookies()
	var cookieSlice []string
	for _, c := range cookie {
		cookieSlice = append(cookieSlice, c.Name+"="+c.Value)
	}
	return strings.Join(cookieSlice, ";"), nil
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
