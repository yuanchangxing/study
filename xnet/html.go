package xnet

import (
	"context"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

func GetHttpHTML(url string) (tmp [][]string) {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true), //设置成无浏览器弹出模式
		chromedp.Flag("blink-settings", "imageEnable=false"),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36"),
	}
	c, _ := chromedp.NewExecAllocator(context.Background(), options...)
	chromeCtx, cancel := chromedp.NewContext(c, chromedp.WithLogf(log.Printf))
	_ = chromedp.Run(chromeCtx, make([]chromedp.Action, 0, 1)...)

	timeOutCtx, cancel := context.WithTimeout(chromeCtx, 60*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(timeOutCtx,
		chromedp.Navigate(url),
		////需要爬取的网页的url
		////chromedp.WaitVisible(`#content > div > section.fp-tournament-award-badge-carousel_awardBadgeCarouselSection__w_Ys5 > div > div > div.col-12.fp-tournament-award-badge-carousel_awardCarouselColumn__fQJLf.g-0 > div > div > div > div > div > div > div.slick-slide.slick-active.slick-current > div > div > div`),
		//chromedp.WaitVisible(selector),
		////等待某个特定的元素出现
		//chromedp.OuterHTML(`document.querySelector("body")`, &htmlContent, chromedp.ByJSPath),
		////生成最终的html文件并保存在htmlContent文件中
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(htmlContent)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))

	awardTmp := [][]string{}
	doc.Find(`div[class="fp-tournament-award-badge_awardContent__dUtoO"]`).Each(func(i int, selection *goquery.Selection) {
		award := selection.Find(`p[class="fp-tournament-award-badge_awardName__JpsZZ"]`).Text()
		//goquery通过Find()查找到我们选择的位置，Each()的功能与遍历相似返回所有的结果，Text()返回文本内容
		name := selection.Find(`h4[class=" fp-tournament-award-badge_awardWinner__P_z2d"]`).Text()
		country := selection.Find(`p[class="fp-tournament-award-badge_awardWinnerCountry__EmjVU"]`).Text()

		awardTmp = append(awardTmp, []string{award, name, country})
	})
	return awardTmp
}
