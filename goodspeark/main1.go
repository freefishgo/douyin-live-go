package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func main2() {
	ops := append(chromedp.DefaultExecAllocatorOptions[3:], chromedp.NoFirstRun, chromedp.NoDefaultBrowserCheck, chromedp.NoSandbox)
	ops = append(chromedp.DefaultExecAllocatorOptions[2:], chromedp.NoSandbox, chromedp.Flag("headless", true))
	ctx1, c1 := chromedp.NewExecAllocator(context.Background(), ops...)
	defer c1()
	//// 创建一个上下文和取消函数
	browserCtx, cancel := chromedp.NewContext(ctx1)
	defer cancel()

	// 导航到百度首页
	if err := chromedp.Run(browserCtx, chromedp.Navigate("https://buyin.jinritemai.com/dashboard/live/control")); err != nil {
		log.Fatal(err)
	}

	// 等待页面加载完成
	time.Sleep(2 * time.Second)

	// 获取当前页面的截图
	if err := captureScreenshot(browserCtx, "screenshot.png"); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Minute)
}

// captureScreenshot 获取当前页面的截图
func captureScreenshot(ctx context.Context, filename string) error {
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
		return err
	}

	// 将截图保存为文件
	if err := ioutil.WriteFile(filename, buf, 0644); err != nil {
		return err
	}

	fmt.Printf("Screenshot saved to %s\n", filename)
	return nil
}
