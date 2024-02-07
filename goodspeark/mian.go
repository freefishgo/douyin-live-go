package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/utils"
	"strconv"
	"time"
)

func main() {

	//Flag("headless", true)(a)
	//// Like in Puppeteer.
	//Flag("hide-scrollbars", true)(a)
	//Flag("mute-audio", true)(a)
	//
	path, _ := launcher.LookPath()
	//Set("--window-size", "1920,1080")
	//u := launcher.New().Delete("--headless").Bin(path).MustLaunch()
	u := launcher.New().Headless(false).Leakless(false).Set("hide-scrollbars").Set("mute-audio").Bin(path).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	// 创建一个浏览器窗口
	//browser := rod.New().MustConnect()

	// 在浏览器中创建一个新页面
	page := browser.MustPage()
	page.MustSetViewport(0, 0, 1, false)

	// 导航到百度首页
	page.MustNavigate("https://buyin.jinritemai.com/dashboard/live/control")
	time.Sleep(time.Second * 5)
	page.WaitLoad()

	img, _ := page.Screenshot(true, nil)
	// 截图保存为文件
	//if err := page.Screenshot("screenshot.png"); err != nil {
	//	log.Fatal(err)
	//}
	_ = utils.OutputFile("screenshot.png", img)
	return

	//page.Emulate(rod.EmulateiPhoneX)

	// 获取搜索框元素并输入关键字
	searchInput := page.MustElement("//*[@id=\"app\"]/div/div/div[3]/div[1]/div[2]/div[2]/div/div/span/div/textarea")
	searchInput.MustInput("你好1")
	//searchInput.
	for i := 0; i < 100; i++ {
		searchInput.MustInput(strconv.Itoa(i))
		time.Sleep(time.Second)
		// 清空输入框内容
		searchInput.MustType(input.Backspace)
		//if err := searchInput.Press(keyboard.Clear); err != nil {
		//	log.Fatal(err)
		//}
	}

	// 获取搜索按钮元素并点击
	searchButton := page.MustElement("#su")
	searchButton.MustClick()

	// 等待搜索结果加载完成
	page.WaitLoad()
	img, _ = page.Screenshot(true, nil)
	// 截图保存为文件
	//if err := page.Screenshot("screenshot.png"); err != nil {
	//	log.Fatal(err)
	//}
	_ = utils.OutputFile("screenshot.png", img)
	time.Sleep(time.Hour)

	// 关闭浏览器
	//browser.MustClose()
}
