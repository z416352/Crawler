package apiservice

import (
	"Crawler/pkg/logger"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
)

var Coincap_time_frame_cases = map[string]int{
	"m1":  1,
	"m5":  5,
	"m15": 15,
	"m30": 30,
	"h1":  60,
	"h2":  120,
	"h4":  240,
	"h8":  480,
	"h12": 720,
	"d1":  1440,
	"w1":  10080,
}

type CoincapApi_opt struct {
	exchange  string
	interval  string
	baseId    string
	quoteId   string
	startTime *int64
	endTime   *int64
}

type CoincapApi_Kline struct {
	OpenDateTime string `json:"openDateTime" bson:"_id"`
	Open         string `json:"open" bson:"open"`
	High         string `json:"high" bson:"high"`
	Low          string `json:"low" bson:"low"`
	Close        string `json:"close" bson:"close"`
	Volume       string `json:"volume" bson:"volume"`
	OpenTime     int64  `json:"period" bson:"openTime"`
}

func (k *CoincapApi_opt) Exchange(exchange string) *CoincapApi_opt {
	k.exchange = exchange
	return k
}

func (k *CoincapApi_opt) Interval(interval string) *CoincapApi_opt {
	k.interval = interval
	return k
}

func (k *CoincapApi_opt) BaseId(baseId string) *CoincapApi_opt {
	k.baseId = baseId
	return k
}

func (k *CoincapApi_opt) QuoteId(quoteId string) *CoincapApi_opt {
	k.quoteId = quoteId
	return k
}

func (k *CoincapApi_opt) StartTime(start int64) *CoincapApi_opt {
	k.startTime = &start
	return k
}

func (k *CoincapApi_opt) EndTime(end int64) *CoincapApi_opt {
	k.endTime = &end
	return k
}

func (k *CoincapApi_opt) GetURL(baseURL string, resource string) string {
	p := url.Values{}

	p.Add("exchange", k.exchange)
	p.Add("interval", k.interval)
	p.Add("baseId", k.baseId)
	p.Add("quoteId", k.quoteId)

	if k.startTime != nil {
		p.Add("start", strconv.FormatInt(*k.startTime, 10))
	}

	if k.endTime != nil {
		p.Add("end", strconv.FormatInt(*k.endTime, 10))
	}

	u, _ := url.ParseRequestURI(baseURL + resource)
	// u.Path = resource
	u.RawQuery = p.Encode()

	return fmt.Sprintf("%v", u)
}

func (k *CoincapApi_opt) Response2json(body []byte) (res []*CoincapApi_Kline, err error) {
	j, err := simplejson.NewJson(body)
	if err != nil {
		panic(err)
	}

	num := len(j.MustArray())
	res = make([]*CoincapApi_Kline, num)
	for i := 0; i < num; i++ {
		item := j.GetIndex(i)
		if len(item.MustArray()) < 11 {
			logger.CrawlerLog.Errorf("invalid kline response")
			return []*CoincapApi_Kline{}, err
		}

		res[i] = &CoincapApi_Kline{
			OpenDateTime: time.UnixMilli(item.GetIndex(5).MustInt64()).Format("2006-01-02 15:04:05"),
			Open:         item.GetIndex(0).MustString(),
			High:         item.GetIndex(1).MustString(),
			Low:          item.GetIndex(2).MustString(),
			Close:        item.GetIndex(3).MustString(),
			Volume:       item.GetIndex(4).MustString(),
			OpenTime:     item.GetIndex(5).MustInt64(),
		}
	}

	return res, nil
}
