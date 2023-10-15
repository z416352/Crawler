package handler

import (
	"github.com/z416352/Crawler/internal/crawler"
	"github.com/z416352/Crawler/internal/request_services"
	"github.com/z416352/Crawler/internal/utils"
	api "github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"
	"time"
)

const (
	initialDataCount_1d  = 365
	initialDataCount_4h  = 1000
	initialDataCount_1h  = 3000
	initialDataCount_15m = 4000
)

// Crawl all data and send the data to the insert API using the POST method.
func crawl_All(target *api.BinanceCrawlTarget) {
	now := time.Now()

	for _, symbol := range target.Symbol_list {
		for _, timeframe := range target.TimeFrame_list {
			switch timeframe {
			case "15m":
				if now.Minute()%15 != 0 {
					continue
				}
			case "1h":
				if now.Minute() != 0 {
					continue
				}
			case "4h":
				if now.Minute() != 0 || now.Hour()%4 != 0 {
					continue
				}
			case "1d":
				if now.Hour() != 0 || now.Minute() != 0 {
					continue
				}
			}

			// Crawler start
			c := crawler.Crawler{
				Timeframe: timeframe,
				Symbol:    symbol,
			}
			newestDataTime := utils.NewestKlineTime(api.Binance_TimeframeCases[timeframe])
			klines, err := c.Binance_Crawler(newestDataTime)
			if err != nil {
				logger.CrawlerLog.Panic(err)
			}

			kline_detail := api.BinanceTypeRequestDetail{
				Symbol:    c.Symbol,
				Timeframe: c.Timeframe,
				Klines:    klines,
			}

			// Send data to the insert API by POST method
			request_services.Post_InsertACrawlerData(&kline_detail, symbol, timeframe)

			time.Sleep(time.Second)
		}
	}
}

// If the database exists, this function will update data with the newest kline.
func updataToNewest() {
	logger.CrawlerLog.Info("====================== UpdataToNewest start ======================")
	logger.CrawlerLog.Infof("Symbols: %v", target.Symbol_list)
	logger.CrawlerLog.Infof("TimeFrame: %v", target.TimeFrame_list)
	logger.CrawlerLog.Info("==================================================================")

	for _, symbol := range target.Symbol_list {
		for _, timeframe := range target.TimeFrame_list {
			kline := request_services.Get_GetNewestData(symbol, timeframe)
			logger.CrawlerLog.Debugf("[%s][%s]: The lastest date of data: '%v'", symbol, timeframe, kline.OpenDateTime)

			startTime := kline.OpenTime
			endTime := utils.NewestKlineTime(api.Binance_TimeframeCases[timeframe])
			if startTime == endTime {
				continue
			}

			c := crawler.Crawler{
				Timeframe: timeframe,
				Symbol:    symbol,
			}
			klines, err := c.Binance_Crawler(endTime, startTime)

			if err != nil {
				logger.CrawlerLog.Panic(err)
			}

			kline_detail := api.BinanceTypeRequestDetail{
				Symbol:    c.Symbol,
				Timeframe: c.Timeframe,
				Klines:    klines,
			}

			// Send data to the insert API by POST method
			request_services.Post_InsertACrawlerData(&kline_detail, symbol, timeframe)

			time.Sleep(time.Second)
		}
	}

	logger.CrawlerLog.Info("====================== UpdataToNewest end ======================\n")
}

// Initialize database data when the database doesn't exist.
// It will crawl a certain amount of data based on different timeframes and insert it into the database.
func initialDB_Data() {
	logger.CrawlerLog.Infof("====================== initial DB Data start ======================")
	logger.CrawlerLog.Infof("These DBs and timeframes need to initial data.")

	missingTimeframe := crawler.ListNonExistTimeframe()
	for symbol, tf_list := range missingTimeframe {
		logger.CrawlerLog.Infof("Symbol: %s, Timeframe: %v", symbol, tf_list)
	}
	logger.CrawlerLog.Infof("===================================================================")

	for symbol, tf_list := range missingTimeframe {
		for _, tf := range tf_list {
			var dataCount int
			switch tf {
			case "1d":
				dataCount = initialDataCount_1d
			case "4h":
				dataCount = initialDataCount_4h
			case "1h":
				dataCount = initialDataCount_1h
			case "15m":
				dataCount = initialDataCount_15m
			default:
				dataCount = 500
			}

			tf_mins := api.Binance_TimeframeCases[tf]

			c := crawler.Crawler{
				Timeframe: tf,
				Symbol:    symbol,
			}

			endTime := utils.NewestKlineTime(tf_mins)
			startTime := utils.ShiftTimeframe(tf_mins, endTime, dataCount)

			klines, err := c.Binance_Crawler(endTime, startTime)
			if err != nil {
				logger.CrawlerLog.Panic(err)
			}

			kline_detail := api.BinanceTypeRequestDetail{
				Symbol:    c.Symbol,
				Timeframe: c.Timeframe,
				Klines:    klines,
			}

			// Send data to the insert API by POST method
			request_services.Post_InsertACrawlerData(&kline_detail, symbol, tf)

			time.Sleep(time.Second)
		}
	}

	logger.CrawlerLog.Infof("======================= initial DB Data end =======================\n")
}
