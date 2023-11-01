package apiservice

import (
	"os"
	"strings"

	"github.com/z416352/Crawler/pkg/logger"
)

// The symbols and timeframes you want to get.
// You can set in GetCrawlTarget() func
type BinanceCrawlTarget struct {
	Symbol_list    []string
	TimeFrame_list []string
}

func (t *BinanceCrawlTarget) GetCrawlTarget() *BinanceCrawlTarget {
	// target := BinanceCrawlTarget{
	// 	Symbol_list:    []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"},
	// 	TimeFrame_list: []string{"15m", "1h"},
	// }

	currenciesEnv, ok := os.LookupEnv("currencies")
	if !ok {
		logger.CrawlerLog.Panicf("currencies env not found")
	}

	timeframesEnv, ok := os.LookupEnv("timeframes")
	if !ok {
		logger.CrawlerLog.Panicf("timeframes env not found")
	}

	target := BinanceCrawlTarget{
		Symbol_list:    format_env(currenciesEnv),
		TimeFrame_list: format_env(timeframesEnv),
	}

	return &target
}

func format_env(strEnv string) []string {
	str_arr := strings.Split(strEnv, "\n")
	res := []string{}

	for _, str := range str_arr {
		str = strings.TrimSpace(str)
		str = strings.TrimPrefix(str, "-")
		str = strings.TrimSpace(str)

		if str != "" {
			res = append(res, str)
		}
	}

	return res
}
