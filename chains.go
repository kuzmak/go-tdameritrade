package tdameritrade

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"strconv"
)

// ChainsService handles communication with the chains related methods of
// the TDAmeritrade API.
//
// TDAmeritrade API docs: https://developer.tdameritrade.com/option-chains/apis
type ChainsService struct {
	client *Client
}

// a float64 whose JSON unmarshaller supports NaN and Inf
type Float64WithSpecial float64

func (v *Float64WithSpecial) UnmarshalJSON(b []byte) error {
	// Is it a float?
	var f float64
	err := json.Unmarshal(b, &f)
	if err == nil {
		*v = Float64WithSpecial(f)
		return nil
	}

	// Try decoding a string instead
	var s string
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	// the field was a JSON string, try to parse that string
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*v = Float64WithSpecial(n)
	return nil
}

func (v Float64WithSpecial) MarshalJSON() ([]byte, error) {
	f := float64(v)
	s := strconv.FormatFloat(f, 'f', -1, 64)

	if math.IsInf(f, 0) || math.IsNaN(f) {
		s = `"` + s + `"`
	}

	return []byte(s), nil
}

type Underlying struct {
	Symbol            string  `json:"symbol"`
	Description       string  `json:"description"`
	Change            float64 `json:"change"`
	PercentChange     float64 `json:"percentChange"`
	Close             float64 `json:"close"`
	QuoteTime         int     `json:"quoteTime"`
	TradeTime         int     `json:"tradeTime"`
	Bid               float64 `json:"bid"`
	Ask               float64 `json:"ask"`
	Last              float64 `json:"last"`
	Mark              float64 `json:"mark"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
	BidSize           int     `json:"bidSize"`
	AskSize           int     `json:"askSize"`
	HighPrice         float64 `json:"highPrice"`
	LowPrice          float64 `json:"lowPrice"`
	OpenPrice         float64 `json:"openPrice"`
	TotalVolume       int     `json:"totalVolume"`
	ExchangeName      string  `json:"exchangeName"`
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
	Delayed           bool    `json:"delayed"`
}

type ExpDateOption struct {
	PutCall                string             `json:"putCall"`
	Symbol                 string             `json:"symbol"`
	Description            string             `json:"description"`
	ExchangeName           string             `json:"exchangeName"`
	Bid                    float64            `json:"bid"`
	Ask                    float64            `json:"ask"`
	Last                   float64            `json:"last"`
	Mark                   float64            `json:"mark"`
	BidSize                int                `json:"bidSize"`
	AskSize                int                `json:"askSize"`
	BidAskSize             string             `json:"bidAskSize"`
	LastSize               float64            `json:"lastSize"`
	HighPrice              float64            `json:"highPrice"`
	LowPrice               float64            `json:"lowPrice"`
	OpenPrice              float64            `json:"openPrice"`
	ClosePrice             float64            `json:"closePrice"`
	TotalVolume            int                `json:"totalVolume"`
	TradeDate              string             `json:"tradeDate"`
	TradeTimeInLong        int                `json:"tradeTimeInLong"`
	QuoteTimeInLong        int                `json:"quoteTimeInLong"`
	NetChange              float64            `json:"netChange"`
	Volatility             Float64WithSpecial `json:"volatility"`
	Delta                  Float64WithSpecial `json:"delta"`
	Gamma                  Float64WithSpecial `json:"gamma"`
	Theta                  Float64WithSpecial `json:"theta"`
	Vega                   Float64WithSpecial `json:"vega"`
	Rho                    Float64WithSpecial `json:"rho"`
	OpenInterest           int                `json:"openInterest"`
	TimeValue              float64            `json:"timeValue"`
	TheoreticalOptionValue Float64WithSpecial `json:"theoreticalOptionValue"`
	TheoreticalVolatility  Float64WithSpecial `json:"theoreticalVolatility"`
	OptionDeliverablesList string             `json:"optionDeliverablesList"`
	StrikePrice            float64            `json:"strikePrice"`
	ExpirationDate         int                `json:"expirationDate"`
	DaysToExpiration       int                `json:"daysToExpiration"`
	ExpirationType         string             `json:"expirationType"`
	LastTradingDate        int                `json:"lastTradingDay"`
	Multiplier             float64            `json:"multiplier"`
	SettlementType         string             `json:"settlementType"`
	DeliverableNote        string             `json:"deliverableNote"`
	IsIndexOption          bool               `json:"isIndexOption"`
	PercentChange          float64            `json:"percentChange"`
	MarkChange             float64            `json:"markChange"`
	MarkPercentChange      float64            `json:"markPercentChange"`
	InTheMoney             bool               `json:"inTheMoney"`
	Mini                   bool               `json:"mini"`
	NonStandard            bool               `json:"nonStandard"`
}

// the first string is the exp date.  the second string is the strike price.
type ExpDateMap map[string]map[string][]ExpDateOption

type Chains struct {
	Symbol            string     `json:"symbol"`
	Status            string     `json:"status"`
	Underlying        Underlying `json:"underlying"`
	Strategy          string     `json:"strategy"`
	Interval          float64    `json:"interval"`
	IsDelayed         bool       `json:"isDelayed"`
	IsIndex           bool       `json:"isIndex"`
	InterestRate      float64    `json:"interestRate"`
	UnderlyingPrice   float64    `json:"underlyingPrice"`
	Volatility        float64    `json:"volatility"`
	DaysToExpiration  float64    `json:"daysToExpiration"`
	NumberOfContracts int        `json:"numberOfContracts"`
	CallExpDateMap    ExpDateMap `json:"callExpDateMap"`
	PutExpDateMap     ExpDateMap `json:"putExpDateMap"`
}

// Users must provide the required URL queryValues for this function to work.
// TD Ameritrade url values: https://developer.tdameritrade.com/option-chains/apis/get/marketdata/chains
// Instructions for using url.Values: https://golang.org/pkg/net/url/#Values
func (s *ChainsService) GetChains(ctx context.Context, queryValues url.Values) (*Chains, *Response, error) {
	u := fmt.Sprintf("marketdata/chains?%s", queryValues.Encode())

	req, err := s.client.NewRequest("GET", u, nil)

	if err != nil {
		return nil, nil, err
	}

	chains := new(Chains)

	resp, err := s.client.Do(ctx, req, chains)
	if err != nil {
		return nil, resp, err
	}

	return chains, resp, nil
}
