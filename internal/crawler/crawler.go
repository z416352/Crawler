package crawler

import (
	"github.com/z416352/Crawler/internal/utils"
	"github.com/z416352/Crawler/pkg/apiservice"
	api "github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// https://docs.coincap.io/
const (
	CoincapAPI_BASE_URL              = "http://api.coincap.io/v2/"
	CoincapAPI_PATH_CANDLESTICK_DATA = "candles"
	CoincapAPI_PATH_EXCHANGEINFO     = "exchanges"
	CoincapAPI_PATH_PRICE            = "assets" // current price
)

const format = "2006-01-02 15:04:05"

type Crawler struct {
	Symbol    string
	Timeframe string
	Data_nums int
}

// Coincap API not work now.
func (c *Crawler) Coincap_Crawler(time_frame, baseId, quoteId string) []*api.CoincapApi_Kline {
	timeframe_mins, exist := api.Coincap_time_frame_cases[time_frame]
	if !exist {
		logger.CrawlerLog.Errorf("Invalid interval. Your interval is '%s'", time_frame)
		os.Exit(1)
	}
	logger.CrawlerLog.Tracef("BaseId = %s, QuoteId = %s,  Interval = %s", baseId, quoteId, time_frame)

	/* Get Coincap API url */
	logger.CrawlerLog.Infof("Get Coincap API url")
	k := new(api.CoincapApi_opt)
	k.Exchange("binance")
	k.Interval(time_frame)
	k.BaseId(baseId)
	k.QuoteId(quoteId)

	// Get c.Data_nums klines
	if c.Data_nums != 0 {
		endTime := utils.NewestKlineTime(timeframe_mins)
		startTime := utils.ShiftTimeframe(timeframe_mins, endTime, c.Data_nums)
		k.StartTime(startTime)
		k.EndTime(endTime)
	}
	urlStr := k.GetURL(CoincapAPI_BASE_URL, CoincapAPI_PATH_CANDLESTICK_DATA)
	logger.CrawlerLog.Debugf("Coincap API url = %s", urlStr)

	/* Get K lines from Coincap */
	logger.CrawlerLog.Infof("Get K lines from Coincap")
	resp, err := http.Get(urlStr)
	if err != nil {
		logger.CrawlerLog.Panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	Kline, err := k.Response2json(body)
	if err != nil {
		logger.CrawlerLog.Error(err)
	}

	logger.CrawlerLog.Infof(`Get "%d" K-lines.`, len(Kline))

	return Kline
}

// timeRange => endTime, startTime. startTime is must earlier than endTime.
func (c *Crawler) Binance_Crawler(timeRange ...int64) ([]*api.BinanceAPI_Kline, error) {
	/* Check parameters for validity  */
	if err := c.isVaild(); err != nil {
		return nil, err
	}

	var startTime, endTime int64
	var klinesNum int
	switch len(timeRange) {
	case 2:
		endTime = timeRange[0]
		startTime = timeRange[1]
		klinesNum = countInterval(startTime, endTime, c.Timeframe)
	case 1:
		endTime = timeRange[0]
		klinesNum = 1
	default:
		return nil, fmt.Errorf("timeRange can only have one or two variables")
	}

	timeframe_mins := api.Binance_TimeframeCases[c.Timeframe]

	/* Setting params and crawl data from Binance API */
	var klines []*api.BinanceAPI_Kline
	k := new(api.BinanceAPI_opt)
	k.Symbol(c.Symbol)
	k.Interval(c.Timeframe)

	logger.CrawlerLog.Infof("====================== Binance Crawler start ======================")
	logger.CrawlerLog.Infof("Symbol = %s, Interval = %s", c.Symbol, c.Timeframe)
	logger.CrawlerLog.Infof("klinesNum: %d", klinesNum)

	check := klinesNum
	for klinesNum != 0 {
		if klinesNum > api.BinanceAPI_Limit_Results {
			c.Data_nums = api.BinanceAPI_Limit_Results
			klinesNum -= api.BinanceAPI_Limit_Results
		} else {
			c.Data_nums = klinesNum
			klinesNum = 0
		}
		startTime = utils.ShiftTimeframe(timeframe_mins, endTime, c.Data_nums-1)
		logger.CrawlerLog.Debugf("StartTime: '%v', EndTime: '%v'", time.UnixMilli(startTime).Format(format), time.UnixMilli(endTime).Format(format))

		k.EndTime(endTime)
		k.StartTime(startTime)
		endTime = utils.ShiftTimeframe(timeframe_mins, startTime, 1)

		kline, err := k.StartCrawler()
		if err != nil {
			logger.CrawlerLog.Errorf("Err: %v", err)
		}

		// logger.CrawlerLog.Debugf("From '%s' to '%s'. Get '%d' klines.", kline[0].OpenDateTime, kline[len(kline)-1].OpenDateTime, len(kline))
		// logger.CrawlerLog.Debugf("-------------")

		klines = append(kline, klines...)
	}

	if len(klines) != check {
		return nil, fmt.Errorf("different from expected number of klines. len(klines): '%d', klinesNum: '%d'", len(klines), check)
	}
	if len(klines) > 0 {
		logger.CrawlerLog.Infof("From '%s' to '%s'. Number of k lines: '%d'", klines[0].OpenDateTime, klines[len(klines)-1].OpenDateTime, len(klines))
	}

	logger.CrawlerLog.Infof("======================= Binance Crawler end =======================\n")
	return klines, nil
}

func (c *Crawler) isVaild() error {
	if c.Symbol == "" || c.Timeframe == "" {
		return fmt.Errorf("invalid 'Symbol' and 'TimeFrame' shouldn't empty")
	}

	if c.Data_nums > apiservice.BinanceAPI_Limit_Results {
		return fmt.Errorf("invalid, 'Data_nums' is too large. Max is %d", apiservice.BinanceAPI_Limit_Results)
	}

	if _, exist := api.Binance_TimeframeCases[c.Timeframe]; !exist {
		return fmt.Errorf("invalid interval. Your interval is '%s'", c.Timeframe)
	}

	return nil
}
