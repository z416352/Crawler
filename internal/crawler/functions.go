package crawler

import (
	req_svc "Crawler/internal/request_services"
	api "Crawler/pkg/apiservice"
	"slices"
	"time"
)

var target *api.BinanceCrawlTarget

func init(){
	target = new(api.BinanceCrawlTarget).GetCrawlTarget()
}

// Counting the number of intervals between the start time and end time.
func countInterval(sTime, eTime int64, timeframe string) int {
	s := time.UnixMilli(sTime)
	e := time.UnixMilli(eTime)
	mins := api.Binance_TimeframeCases[timeframe]

	gap := e.Sub(s)
	nums := int(gap.Minutes()) / mins

	return nums
}

// List the missing timeframe.
func ListNonExistTimeframe() map[string][]string {
	nonExistTimeframeMap := make(map[string][]string)

	for _, symbol := range target.Symbol_list {
		existTimeframeList := req_svc.Get_DBExistTimeframes(symbol)

		for _, tf := range target.TimeFrame_list {
			if !slices.Contains(existTimeframeList, tf) {
				nonExistTimeframeMap[symbol] = append(nonExistTimeframeMap[symbol], tf)
			}
		}
	}

	return nonExistTimeframeMap
}
