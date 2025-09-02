package bin

import (
	"context"
	binance_connector "github.com/binance/binance-connector-go"
)

func WalletBalance() (res []*binance_connector.CoinInfo, err error) {
	client := binance_connector.NewClient(conf.GetApiKey(), conf.GetApiSecret())
	return client.NewGetAllCoinsInfoService().Do(context.Background())
}
