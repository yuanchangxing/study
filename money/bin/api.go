package bin

import (
	"context"
	"fmt"
	"github.com/yuanchangxing/study/xlog"
	"math/rand"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type IBinConfiger interface {
	GetApiKey() string
	GetApiSecret() string
}

var (
	client *binance.Client
	conf   IBinConfiger
)

func Init(config IBinConfiger) {
	conf = config
}

//const (
//	apiKey           = "你的API密钥"
//	apiSecret        = "你的API私钥"
//	symbol           = "AITECH/USDT" // 替换为目标Alpha代币交易对
//	quantity         = 10.0          // 购买数量
//	priceLimit       = 0.1           // 限价单价格
//	maxTradesPerHour = 10            // 每小时最大交易次数
//	checkIntervalMin = 5             // 最小检查间隔（秒）
//	checkIntervalMax = 10            // 最大检查间隔（秒）
//)

func DoSom() error {
	// 初始化币安客户端
	client = binance.NewClient(conf.GetApiKey(), conf.GetApiSecret())

	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	// 检查账户余额
	balance, err := getWalletBalance(client, "SUI")
	if err != nil {
		return fmt.Errorf("get wallet balance err: %v", err)
	}

	xlog.Infof("balance is %v", balance)
	//fmt.Printf("USDT余额: %.2f\n", balance)

	//if balance < quantity*priceLimit {
	//	//fmt.Println("余额不足，无法购买！")
	//	//return
	//}

	//// 交易计数和时间跟踪
	//tradeCount := 0
	//startTime := time.Now()
	//
	//// 主循环：监控价格并下单
	//for {
	//	// 检查交易频率限制
	//	if tradeCount >= maxTradesPerHour && time.Since(startTime).Hours() < 1 {
	//		fmt.Println("达到小时交易上限，暂停1小时...")
	//		time.Sleep(time.Hour)
	//		tradeCount = 0
	//		startTime = time.Now()
	//	}
	//
	//	// 获取市场价格
	//	currentPrice, err := getMarketPrice(client, symbol)
	//	if err != nil {
	//		fmt.Printf("获取市场价格失败: %v\n", err)
	//		time.Sleep(time.Second * time.Duration(rand.Intn(checkIntervalMax-checkIntervalMin)+checkIntervalMin))
	//		continue
	//	}
	//	fmt.Printf("当前市场价格: %.4f\n", currentPrice)
	//
	//	// 检查是否达到目标价格
	//	if currentPrice <= priceLimit {
	//		fmt.Printf("价格达到目标: %.4f，开始购买...\n", currentPrice)
	//		err := placeBuyOrder(client, symbol, quantity, priceLimit)
	//		if err != nil {
	//			fmt.Printf("下买单失败: %v\n", err)
	//		} else {
	//			fmt.Println("购买成功，退出程序！")
	//			tradeCount++
	//			break
	//		}
	//	} else {
	//		fmt.Printf("价格 %.4f 高于目标价 %.4f，继续监控...\n", currentPrice, priceLimit)
	//	}
	//
	//	// 随机间隔，避免固定频率
	//	time.Sleep(time.Second * time.Duration(rand.Intn(checkIntervalMax-checkIntervalMin)+checkIntervalMin))
	//}
	return nil
}

func getWalletBalance(client *binance.Client, asset string) (float64, error) {
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, b := range account.Balances {
		free, _ := strconv.ParseFloat(b.Free, 64)
		if free > 0 {
			xlog.Infof("Asset : %s  balance %v", b.Asset, b.Free)
		}

		if b.Asset == asset {
			xlog.Infof("balance is %v", b)
			free, _ := strconv.ParseFloat(b.Free, 64)
			return free, nil
		}
	}
	return 0, fmt.Errorf("notfound %s balance", asset)
}

func getMarketPrice(client *binance.Client, symbol string) (float64, error) {
	price, err := client.NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return 0, err
	}
	if len(price) == 0 {
		return 0, fmt.Errorf("未获取到 %s 价格", symbol)
	}
	return strconv.ParseFloat(price[0].Price, 64)
}

func placeBuyOrder(client *binance.Client, symbol string, quantity, price float64) error {
	_, err := client.NewCreateOrderService().
		Symbol(symbol).
		Side(binance.SideTypeBuy).
		Type(binance.OrderTypeLimit).
		TimeInForce(binance.TimeInForceTypeGTC).
		Quantity(fmt.Sprintf("%.2f", quantity)).
		Price(fmt.Sprintf("%.4f", price)).
		Do(context.Background())
	return err
}
