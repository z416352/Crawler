package apiservice

// The symbols and timeframes you want to get.
// You can set in GetCrawlTarget() func
type BinanceCrawlTarget struct {
	Symbol_list    []string
	TimeFrame_list []string
}

func (t *BinanceCrawlTarget) GetCrawlTarget() *BinanceCrawlTarget {
	target := BinanceCrawlTarget{
		Symbol_list:    []string{"BTCUSDT", "ETHUSDT", "BNBUSDT"},
		TimeFrame_list: []string{"15m", "1h"},
	}

	return &target
}
