package apiservice

// The struct is used to POST to the Database API
// POST func -> Post_InsertACrawlerData()
type BinanceTypeRequestDetail struct {
	Symbol    string
	Timeframe string
	Klines    []*BinanceAPI_Kline
}
