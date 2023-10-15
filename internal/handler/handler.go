package handler

import (
	api "Crawler/pkg/apiservice"
	"Crawler/pkg/logger"
	"time"

	"github.com/robfig/cron/v3"
)

// Cron 語法工具
// https://tool.lu/crontab/
// https://crontab.guru/
var target *api.BinanceCrawlTarget

func init() {
	target = new(api.BinanceCrawlTarget).GetCrawlTarget()

	logger.CrawlerLog.Infof("=========================================================")
	logger.CrawlerLog.Infof("There are '%d' symbols and '%d' timeframes.", len(target.Symbol_list), len(target.TimeFrame_list))
	logger.CrawlerLog.Infof("Symbol: %v", target.Symbol_list)
	logger.CrawlerLog.Infof("Timeframe: %v", target.TimeFrame_list)
	logger.CrawlerLog.Infof("=========================================================\n")
}

func Crawler_handler() {
	c := cron.New()
	initialDB_Data()

	updataToNewest()

	// "*/15 * * * *" -> At every 15th minute.
	_, err := c.AddFunc("*/15 * * * *", func() {
		time.Sleep(time.Second * 1)
		crawl_All(target)
	})
	if err != nil {
		logger.CrawlerLog.Panicf("Can't add cron rules.")
	}

	c.Start()
}
