/**
 * @Time: 2023/10/20 16:29
 * @Author: jzechen
 * @File: dto.go
 * @Software: GoLand collector
 */

package scrape

type SinaWeiboConfig struct {
	UserName []string    `yaml:"userName"`
	Password []string    `yaml:"password"`
	Info     bool        `yaml:"info"`
	Follow   bool        `yaml:"follow"`
	Fans     bool        `yaml:"fans"`
	Tweet    TweetConfig `yaml:"tweet"`
}

type TweetConfig struct {
	Main    bool `yaml:"main"`
	Comment bool `yaml:"comment"`
}
