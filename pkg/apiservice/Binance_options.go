package apiservice

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/z416352/Crawler/pkg/logger"

	"github.com/bitly/go-simplejson"
)

const (
	BinanceAPI_BASE_URL              = "https://api.binance.com/"
	BinanceAPI_PATH_CANDLESTICK_DATA = "api/v3/klines"
	BinanceAPI_PATH_EXCHANGEINFO     = "api/v3/exchangeInfo"
	BinanceAPI_PATH_PRICE            = "api/v3/ticker/price"
	BinanceAPI_Limit_Results         = 500 // MAX = 500
)

const format = "2006-01-02 15:04:05"

var Binance_TimeframeCases = map[string]int{
	"1m":  1,
	"3m":  3,
	"5m":  5,
	"15m": 15,
	"30m": 30,
	"1h":  60,
	"2h":  120,
	"4h":  240,
	"6h":  360,
	"8h":  480,
	"12h": 720,
	"1d":  1440,
	"3d":  4320,
	"1w":  10080,
	"1M":  43200, // 30 days
}

// Binance API options
type BinanceAPI_opt struct {
	symbol    string
	interval  string
	limit     *int
	startTime *int64
	endTime   *int64
}

func (k *BinanceAPI_opt) Symbol(symbol string) *BinanceAPI_opt {
	k.symbol = symbol
	return k
}

func (k *BinanceAPI_opt) Interval(interval string) *BinanceAPI_opt {
	k.interval = interval
	return k
}

func (k *BinanceAPI_opt) Limit(limit int) *BinanceAPI_opt {
	k.limit = &limit
	return k
}

func (k *BinanceAPI_opt) StartTime(start int64) *BinanceAPI_opt {
	k.startTime = &start
	return k
}

func (k *BinanceAPI_opt) EndTime(end int64) *BinanceAPI_opt {
	k.endTime = &end
	return k
}

func (k *BinanceAPI_opt) GetURL(baseURL string, resource string) string {
	p := url.Values{}

	p.Add("symbol", k.symbol)
	p.Add("interval", k.interval)

	if k.startTime != nil {
		p.Add("startTime", strconv.FormatInt(*k.startTime, 10))
	}

	if k.endTime != nil {
		p.Add("endTime", strconv.FormatInt(*k.endTime, 10))
	}

	if k.limit != nil {
		p.Add("limit", strconv.Itoa(*k.limit))
	}

	u, _ := url.ParseRequestURI(baseURL + resource)
	// u.Path = resource
	u.RawQuery = p.Encode()

	return fmt.Sprintf("%v", u)
}

// Convert responses from Binance to JSON because the responses are in an array type.
func (k *BinanceAPI_opt) Response2json(body []byte) (res []*BinanceAPI_Kline, err error) {
	j, err := simplejson.NewJson(body)
	if err != nil {
		panic(err)
	}

	num := len(j.MustArray())
	res = make([]*BinanceAPI_Kline, num)
	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		if len(item.MustArray()) < 11 {
			logger.CrawlerLog.Errorf("invalid kline response")
			return []*BinanceAPI_Kline{}, err
		}
		utcTime := time.Unix(0, item.GetIndex(0).MustInt64()*int64(time.Millisecond))
		// Convert time to UTC+8 time zone
		loc, _ := time.LoadLocation("Asia/Taipei")
		localTime := utcTime.In(loc)

		res[i] = &BinanceAPI_Kline{
			// OpenDateTime:             time.UnixMilli(item.GetIndex(0).MustInt64()).Format("2006-01-02 15:04:05"),
			OpenDateTime:             localTime.Format(format),
			OpenTime:                 item.GetIndex(0).MustInt64(),
			Open:                     item.GetIndex(1).MustString(),
			High:                     item.GetIndex(2).MustString(),
			Low:                      item.GetIndex(3).MustString(),
			Close:                    item.GetIndex(4).MustString(),
			Volume:                   item.GetIndex(5).MustString(),
			CloseTime:                item.GetIndex(6).MustInt64(),
			QuoteAssetVolume:         item.GetIndex(7).MustString(),
			TradeNum:                 item.GetIndex(8).MustInt64(),
			TakerBuyBaseAssetVolume:  item.GetIndex(9).MustString(),
			TakerBuyQuoteAssetVolume: item.GetIndex(10).MustString(),
		}
	}

	return res, nil
}

// The result from the API is an array list, not a JSON type.
func (k *BinanceAPI_opt) StartCrawler() ([]*BinanceAPI_Kline, error) {
	if k.endTime == nil || k.startTime == nil || k.interval == "" || k.symbol == "" {
		return nil, fmt.Errorf("")
	}

	logger.CrawlerLog.Tracef("--------------------- Get Binance API url ---------------------")

	urlStr := k.GetURL(BinanceAPI_BASE_URL, BinanceAPI_PATH_CANDLESTICK_DATA)
	// endTime = utils.ShiftTimeframe(timeframe_mins, endTime, 1)

	logger.CrawlerLog.Debugf("Binance API: %s", urlStr)

	/* Get K lines from Binance */
	logger.CrawlerLog.Tracef("--------------------- Get K lines from Binance ---------------------")

	resp, err := http.Get(urlStr)
	if err != nil {
		logger.CrawlerLog.Panic(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	kline, err := k.Response2json(body)
	if err != nil {
		logger.CrawlerLog.Error(err)
	}

	return kline, nil
}
