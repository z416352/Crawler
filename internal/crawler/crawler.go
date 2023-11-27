package crawler

import (
	"fmt"

	"github.com/z416352/Crawler/internal/utils"
	"github.com/z416352/Crawler/pkg/apiservice"
	api "github.com/z416352/Crawler/pkg/apiservice"
	"github.com/z416352/Crawler/pkg/logger"
)

const format = "2006-01-02 15:04:05"

type Crawler struct {
	Symbol    string
	Timeframe string
	Data_nums int
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
		logger.CrawlerLog.Debugf("StartTime: '%v', EndTime: '%v'", utils.ConvertToUTC8(startTime), utils.ConvertToUTC8(endTime))

		k.EndTime(endTime)
		k.StartTime(startTime)

		// next loop end time
		endTime = utils.ShiftTimeframe(timeframe_mins, startTime, 1)

		kline, err := k.StartCrawler()
		if err != nil {
			logger.CrawlerLog.Errorf("Err: %v", err)
		}


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
