package sutando

import (
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type testStruct struct {
	StructName   string
	StructAge    int
	RealBirthday time.Time
	Ignore       int `bson:"-"`
	primitive.M
	Inner     testSubStruct
	Inner2    testSubStruct
	ArrTest   []int
	MapTest   map[string]int
	StructMap map[string]testSubStruct
	StructArr []testSubStruct
	FitValue  decimal.Decimal
	Random    int64
}

type testSubStruct struct {
	JustMyName string `bson:"name"`
	OhTheAge   int
}

func mockData() testStruct {
	return testStruct{
		StructName:   "Yanun",
		StructAge:    27,
		RealBirthday: time.Date(1995, time.March, 23, 0, 0, 0, 0, time.Local),
		Ignore:       10,
		M:            primitive.M{"Hello": 123},
		Inner: testSubStruct{
			JustMyName: "inner",
			OhTheAge:   10,
		},
		Inner2: testSubStruct{
			JustMyName: "inner2",
			OhTheAge:   10,
		},
		ArrTest: []int{1, 2, 3, 4, 5},
		MapTest: map[string]int{"1": 2, "3": 4, "5": 6},
		StructMap: map[string]testSubStruct{
			"1": {
				JustMyName: "inner",
				OhTheAge:   50,
			},
			"2": {
				JustMyName: "inner2",
				OhTheAge:   10,
			},
		},
		StructArr: []testSubStruct{
			{
				JustMyName: "inner",
				OhTheAge:   50,
			},
			{
				JustMyName: "inner2",
				OhTheAge:   10,
			},
			{
				JustMyName: "inner3",
				OhTheAge:   1,
			},
		},
		FitValue: decimal.RequireFromString("0.03"),
		Random:   time.Now().UnixMilli(),
	}
}

// Setting 一件買幣總設定值
type Setting struct {
	ID               primitive.ObjectID         `bson:"_id" json:"id"`
	MonitorSettings  map[string]MonitorSetting  `bson:"monitorSettings"`  // 監控相關設定
	CurrencySettings map[string]CurrencySetting `bson:"currencySettings"` /* 幣種：設定值 */
	ExchangeSettings map[string]ExchangeSetting `bson:"exchangeSettings"` /* 交易所：設定值（目前只有買賣開關） */
}

type MonitorSetting struct {
	BalanceAlarmThreshold    decimal.Decimal `bson:"balanceAlarmThreshold"`    // comment:"USD警示下限"
	BalanceStopRateThreshold decimal.Decimal `bson:"balanceStopRateThreshold"` // comment:"USD停止報價下限"
}

type CurrencySetting struct {
	Currency          string `bson:"currency"`
	RateSetting       `bson:"rateSetting"`
	SupplementSetting `bson:"supplementSetting"`
}
type RateSetting struct {
	SellSetting     TradeSetting    `bson:"sellSetting"`
	BuySetting      TradeSetting    `bson:"buySetting"`
	InverseRate     decimal.Decimal `bson:"inverseRate"`     // comment:"逆價差趴數"
	EnableFluctRate bool            `bson:"enableFluctRate"` // comment:"是否啟動浮動趴數"
}

type TradeSetting struct {
	ProfitRate  decimal.Decimal `comment:"利潤趴數"`
	USDAmount   decimal.Decimal `comment:"抓價格用的 USD 深度"`
	CryptAmount decimal.Decimal `comment:"抓價格用的 虛擬幣 深度"`
}
type SupplementSetting struct {
	BuyOrderMode          BuyOrderMode    `bson:"buyOrderMode"`          // `comment:"1) 買單以顆數為主下限價 2) 買單以總價為主下市價單"`
	RiskPercentage        decimal.Decimal `bson:"riskPercentage"`        // `comment:"風控趴數"`
	BitoProMinAmountOrder decimal.Decimal `bson:"bitoProMinAmountOrder"` // comment: 補水位訂單量小於該值後在 bitopro 站內下單
}
type (
	BuyOrderMode uint8
)
type ExchangeSetting struct {
	BuyEnable  bool `bson:"buyEnable"`
	SellEnable bool `bson:"sellEnable"`
}

const (
	BuyOrderWithAmount = 1
)

var _defaultSettings = Setting{
	MonitorSettings: map[string]MonitorSetting{
		"USD": {
			BalanceAlarmThreshold:    decimal.NewFromInt(100),
			BalanceStopRateThreshold: decimal.NewFromInt(0),
		},
		"TWD": {
			BalanceAlarmThreshold:    decimal.NewFromInt(5000),
			BalanceStopRateThreshold: decimal.NewFromInt(0),
		},
	},
	CurrencySettings: map[string]CurrencySetting{
		"BTC": {
			Currency: "BTC",
			RateSetting: RateSetting{
				SellSetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("0.99"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				BuySetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("1.01"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				InverseRate:     decimal.RequireFromString("0.01"),
				EnableFluctRate: true,
			},
			SupplementSetting: SupplementSetting{
				BuyOrderMode:          BuyOrderWithAmount,
				RiskPercentage:        decimal.RequireFromString("0.01"),
				BitoProMinAmountOrder: decimal.RequireFromString("0.5"),
			},
		},
		"ETH": {
			Currency: "ETH",
			RateSetting: RateSetting{
				SellSetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("0.99"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				BuySetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("1.01"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				InverseRate:     decimal.RequireFromString("0.01"),
				EnableFluctRate: true,
			},
			SupplementSetting: SupplementSetting{
				BuyOrderMode:          BuyOrderWithAmount,
				RiskPercentage:        decimal.RequireFromString("0.01"),
				BitoProMinAmountOrder: decimal.RequireFromString("5"),
			},
		},
		"USDT": {
			Currency: "USDT",
			RateSetting: RateSetting{
				SellSetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("0.995"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				BuySetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("1.005"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				InverseRate:     decimal.RequireFromString("0.001"),
				EnableFluctRate: true,
			},
			SupplementSetting: SupplementSetting{
				BuyOrderMode:          BuyOrderWithAmount,
				RiskPercentage:        decimal.RequireFromString("0.01"),
				BitoProMinAmountOrder: decimal.RequireFromString("3000"),
			},
		},
		"USDC": {
			Currency: "USDC",
			RateSetting: RateSetting{
				SellSetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("0.995"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				BuySetting: TradeSetting{
					ProfitRate: decimal.RequireFromString("1.005"),
					USDAmount:  decimal.RequireFromString("10000"),
				},
				InverseRate:     decimal.RequireFromString("0.001"),
				EnableFluctRate: true,
			},
			SupplementSetting: SupplementSetting{
				BuyOrderMode:          1,
				RiskPercentage:        decimal.RequireFromString("0.01"),
				BitoProMinAmountOrder: decimal.RequireFromString("1000"),
			},
		},
	},
	ExchangeSettings: map[string]ExchangeSetting{
		"Ftx": {
			BuyEnable:  false,
			SellEnable: false,
		},
		"Bitfinex": {
			BuyEnable:  false,
			SellEnable: false,
		},
		"Liquid": {
			BuyEnable:  false,
			SellEnable: false,
		},
		"OkCoin": {
			BuyEnable:  true,
			SellEnable: true,
		},
		"Bitopro": {
			BuyEnable:  true,
			SellEnable: true,
		},
		"FtxOtc": {
			BuyEnable:  false,
			SellEnable: false,
		},
		"Binance": {
			BuyEnable:  false,
			SellEnable: false,
		},
	},
}
