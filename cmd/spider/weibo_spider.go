package main

import (
	"WeiboSpiderGo/pkg/config"
	scrapy_rules2 "WeiboSpiderGo/pkg/scrapy_rules"
	"WeiboSpiderGo/pkg/utils"
	"fmt"
)

var uidLi = utils.GetTargetUidList()

func scrapyInfomation() {
	getInfoC := scrapy_rules2.GetDefaultCollector()
	getMoreInfoC := scrapy_rules2.GetDefaultCollector()
	scrapy_rules2.SetMoreInfoCallback(getMoreInfoC)

	scrapy_rules2.SetInfoCallback(getInfoC, getMoreInfoC)

	for _, uid := range uidLi {
		url := fmt.Sprintf("%s/%s/info", scrapy_rules2.BaseUrl, uid)
		getInfoC.Visit(url)
	}
	getInfoC.Wait()
	getMoreInfoC.Wait()
}

func scrapyTweet() {
	getTweetsC := scrapy_rules2.GetDefaultCollector()
	getContentSubC := scrapy_rules2.GetDefaultCollector()
	scrapy_rules2.SetFullContentCallback(getContentSubC)
	getCommentSubC := scrapy_rules2.GetDefaultCollector()
	scrapy_rules2.SetCommentCallback(getCommentSubC)

	scrapy_rules2.SetTweetCallback(getTweetsC, getContentSubC, getCommentSubC)

	for _, uid := range uidLi {
		url := scrapy_rules2.GetTweetUrl(uid)
		getTweetsC.Visit(url)
	}
	getTweetsC.Wait()
	getContentSubC.Wait()
	getCommentSubC.Wait()
}

func scrapyFollow() {
	getFollowC := scrapy_rules2.GetDefaultCollector()
	scrapy_rules2.SetFollowCallback(getFollowC)
	//read files
	for _, uid := range uidLi {
		url := scrapy_rules2.GetFollowUrl(uid)
		getFollowC.Visit(url)
	}
	getFollowC.Wait()
}

func scrapyFans() {
	getFansC := scrapy_rules2.GetDefaultCollector()
	scrapy_rules2.SetFansCallback(getFansC)

	for _, uid := range uidLi {
		url := scrapy_rules2.GetFansUrl(uid)
		getFansC.Visit(url)
	}
	getFansC.Wait()
}

func main() {
	if config.Conf.GetBool("SCRAPY_TYPE.Info") {
		scrapyInfomation()
	}
	if config.Conf.GetBool("SCRAPY_TYPE.Follow") {
		scrapyFollow()
	}
	//修复去重问题
	if config.Conf.GetBool("SCRAPY_TYPE.Fans") {
		scrapyFans()
	}
	if config.Conf.GetBool("SCRAPY_TYPE.Tweet.Main") {
		scrapyTweet()
	}
}
