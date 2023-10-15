package apiservice

// The klines responses from Binance API can be stored in this struct.
type BinanceAPI_Kline struct {
	OpenDateTime             string `json:"openDateTime" bson:"_id"`
	OpenTime                 int64  `json:"openTime" bson:"OpenTime"`
	Open                     string `json:"open" bson:"Open"`
	High                     string `json:"high" bson:"High"`
	Low                      string `json:"low" bson:"Low"`
	Close                    string `json:"close" bson:"Close"`
	Volume                   string `json:"volume" bson:"Volume"`
	CloseTime                int64  `json:"closeTime" bson:"CloseTime"`
	QuoteAssetVolume         string `json:"quoteAssetVolume" bson:"QuoteAssetVolume"`
	TradeNum                 int64  `json:"tradeNum" bson:"TradeNum"`
	TakerBuyBaseAssetVolume  string `json:"takerBuyBaseAssetVolume" bson:"TakerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string `json:"takerBuyQuoteAssetVolume" bson:"TakerBuyQuoteAssetVolume"`
}
