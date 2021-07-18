package bitkub

type BitkubWebsocket struct {
	Stream string `json:"stream" example:"market.ticker.thb_btc"`
}

type TickerMarketWebsocket struct {
	Stream         string  `json:"stream" example:"market.ticker.thb_btc"`
	ID             int     `json:"id" example:"1"`
	Last           float64 `json:"last" example:"1065201"`
	LowestAsk      float64 `json:"lowestAsk" example:"1066424.99"`
	LowestAskSize  float64 `json:"lowestAskSize" example:"2.1746439"`
	HighestBid     float64 `json:"highestBid" example:"1065201"`
	HighestBidSize float64 `json:"highestBidSize" example:"0.00508662"`
	Change         float64 `json:"change" example:"500.95"`
	PercentChange  float64 `json:"percentChange" example:"0.05"`
	BaseVolume     float64 `json:"baseVolume" example:"116.8407195"`
	QuoteVolume    float64 `json:"quoteVolume" example:"125125608.44"`
	IsFrozen       float64 `json:"isFrozen" example:"0"`
	High24Hr       float64 `json:"high24hr" example:"1082999"`
	Low24Hr        float64 `json:"low24hr" example:"1060065.78"`
	Open           float64 `json:"open" example:"1064700.05"`
	Close          float64 `json:"close" example:"1065201"`
}

type TradeMarketWebsocket struct {
	Stream    string  `json:"stream" example:"market.trade.thb_btc"`
	Symbol    string  `json:"sym" example:"THB_ETH"`
	Txn       string  `json:"txn" example:"ETHSELL0000074282"`
	Rate      float64 `json:"rat" example:"5977.00"`
	Amount    float64 `json:"amt" example:"1.556539"`
	Buy       int     `json:"bid" example:"2048451"`
	Sell      int     `json:"sid" example:"2924729"`
	Timestamp int     `json:"ts" example:"1542268567"`
}
