/**
 * @Time: 2023/11/1 11:43
 * @Author: jzechen
 * @File: selenium_test.go
 * @Software: GoLand toresa
 */

package envtest

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	// DrivePath the path of Chrome drive
	DrivePath = "/mnt/d/Google/Chrome/Application/chromedriver.exe"
	// ChromePath the path of Chrome browser
	ChromePath = "/mnt/d/Google/Chrome/Application/chrome.exe"
	// ChromePort the --remote-debugging-port of Chrome browser
	ChromePort = 9515
	// InstanceURL the listen url of the Chrome instance
	InstanceURL = "http://127.0.0.1:%d/wd/hub"
)

func TestSeleniumLink(t *testing.T) {
	// set selenium service option
	opts := []selenium.ServiceOption{
		selenium.ChromeDriver(DrivePath),
		selenium.Output(os.Stderr), // redirect output to stderr
	}

	// start Selenium server
	service, err := selenium.NewChromeDriverService(DrivePath, ChromePort, opts...)
	if err != nil {
		t.Fatalf("Error starting the ChromeDriver server: %s", err)
	}
	defer service.Stop()

	// set Chrome option
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Path: ChromePath,
		Args: []string{
			"--headless",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36",
		},
	}
	caps.AddChrome(chromeCaps)

	// connect to the Chrome instance
	wd, err := selenium.NewRemote(caps, fmt.Sprintf(InstanceURL, ChromePort))
	if err != nil {
		t.Fatalf("Failed to connect to the WebDriver: %s", err)
	}
	defer wd.Quit()

	// get html page
	err = wd.Get("https://www.baidu.com")
	if err != nil {
		t.Fatalf("Failed to load webpage: %s", err)
	}

	// wait and exit
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func TestChromeDriver(t *testing.T) {
	// 启动ChromeDriver
	cmd := exec.Command(DrivePath)
	if err := cmd.Start(); err != nil {
		fmt.Printf("无法启动ChromeDriver: %v", err)
		return
	}
	defer cmd.Process.Kill()

	// 启动Chrome浏览器
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless", // 无头模式
		},
	}
	caps.AddChrome(chromeCaps)

	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://172.21.53.193:9222"))
	if err != nil {
		fmt.Printf("无法连接到远程浏览器: %v", err)
		return
	}
	defer wd.Quit()

	// 在此处编写你的Selenium代码
	// 例如：
	wd.Get("https://www.google.com")
	title, err := wd.Title()
	if err != nil {
		fmt.Printf("无法获取页面标题: %v", err)
		return
	}
	fmt.Println("页面标题:", title)
}
