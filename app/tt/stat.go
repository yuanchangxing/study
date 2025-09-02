package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

// Transaction 定义代币交易结构体
type Transaction struct {
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	TokenSymbol     string `json:"tokenSymbol"`
	Timestamp       string `json:"timeStamp"`
	IsError         string `json:"isError"`
	ContractAddress string `json:"contractAddress"`
}

// APIResponse 定义 BSCScan API 返回的结构体
type APIResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Transaction `json:"result"`
}

// USDTStats USDT 交易统计信息
type USDTStats struct {
	TotalTx         int
	SuccessfulTx    int
	FailedTx        int
	TotalUSDTVolume float64
	EarliestTxTime  time.Time
	LatestTxTime    time.Time
	ToAddresses     map[string]int
	FromAddresses   map[string]int
}

const (
	myAddress = "0xf82dc96c681452e5bad265676f49aa28553b04b0"
)

func main() {
	// 替换为你的 BSCScan API Key、要查询的地址和 USDT 合约地址
	apiKey := "T3YQ58XDEUVX6XV46JW894VY8C4UI5W3AI"
	//address := "0xf82dc96c681452e5bad265676f49aa28553b04b0"
	usdtContract := "0x55d398326f99059ff775485246999027b3197955" // BSC USDT 合约地址

	// 获取 USDT 交易数据
	transactions, err := getUSDTTransactions(myAddress, usdtContract, apiKey)
	if err != nil {
		fmt.Printf("Error fetching transactions: %v\n", err)
		return
	}

	// 统计 USDT 交易信息
	stats := analyzeUSDTTransactions(transactions)

	// 打印统计结果
	printUSDTStats(stats)
}

// getUSDTTransactions 从 BSCScan API 获取 USDT 交易数据
func getUSDTTransactions(address, contractAddress, apiKey string) ([]Transaction, error) {
	url := fmt.Sprintf(
		"https://api.bscscan.com/api?module=account&action=tokentx&contractaddress=%s&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
		contractAddress, address, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//log.Println(string(body))

	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		return nil, err
	}

	if apiResp.Status != "1" {
		return nil, fmt.Errorf("API error: %s", apiResp.Message)
	}

	return apiResp.Result, nil
}

// analyzeUSDTTransactions 分析 USDT 交易数据
func analyzeUSDTTransactions(txs []Transaction) USDTStats {
	stats := USDTStats{
		ToAddresses:   make(map[string]int),
		FromAddresses: make(map[string]int),
	}

	for _, tx := range txs {
		// 确认交易是 USDT
		//if tx.TokenSymbol != "USDT" {
		//	continue
		//}

		// 统计交易总数
		stats.TotalTx++

		// 统计成功/失败交易
		if tx.IsError == "0" {
			stats.SuccessfulTx++
		} else {
			stats.FailedTx++
		}

		// 统计 USDT 交易量（转换为 USDT 单位，18 位小数）
		value, err := parseBigInt(tx.Value)
		if err != nil {
			fmt.Printf("Warning: failed to parse value %s for tx %s: %v\n", tx.Value, tx.Hash, err)
			continue
		}
		// Convert to USDT (divide by 10^18)
		usdtValue := new(big.Float).Quo(new(big.Float).SetInt(value), big.NewFloat(1e18))
		usdtFloat, _ := usdtValue.Float64()
		stats.TotalUSDTVolume += usdtFloat
		log.Println(tx.Value, value)
		// 统计交易时间
		//timestamp, _ := parseBigInt(tx.Timestamp)
		//txTime := time.Unix(int64(timestamp), 0)
		//if stats.EarliestTxTime.IsZero() || txTime.Before(stats.EarliestTxTime) {
		//	stats.EarliestTxTime = txTime
		//}
		//if txTime.After(stats.LatestTxTime) {
		//	stats.LatestTxTime = txTime
		//}

		// 统计交易对手地址
		stats.ToAddresses[tx.To]++
		stats.FromAddresses[tx.From]++
	}

	return stats
}

// parseBigInt 将字符串转换为 uint64
func parseBigInt(s string) (*big.Int, error) {
	value, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse big int: %s", s)
	}
	return value, nil
}

// printUSDTStats 打印 USDT 统计结果
func printUSDTStats(stats USDTStats) {
	fmt.Printf("USDT 交易统计信息:\n")
	fmt.Printf("总交易数: %d\n", stats.TotalTx)
	fmt.Printf("成功交易: %d\n", stats.SuccessfulTx)
	fmt.Printf("失败交易: %d\n", stats.FailedTx)
	fmt.Printf("总 USDT 交易量: %.2f USDT\n", stats.TotalUSDTVolume)
	fmt.Printf("最早交易时间: %s\n", stats.EarliestTxTime.Format(time.RFC3339))
	fmt.Printf("最晚交易时间: %s\n", stats.LatestTxTime.Format(time.RFC3339))
	fmt.Printf("唯一接收地址数: %d\n", len(stats.ToAddresses))
	fmt.Printf("唯一发送地址数: %d\n", len(stats.FromAddresses))

	// 打印最常见的交易对手地址
	fmt.Println("\n最常见的5个接收地址:")
	for addr, count := range getTopAddresses(stats.ToAddresses, 5) {
		fmt.Printf("%s: %d 次\n", addr, count)
	}
}

// getTopAddresses 获取交易次数最多的地址
func getTopAddresses(addresses map[string]int, limit int) map[string]int {
	type addrCount struct {
		addr  string
		count int
	}

	var counts []addrCount
	for addr, count := range addresses {
		counts = append(counts, addrCount{addr, count})
	}

	// 按交易次数排序
	for i := 0; i < len(counts)-1; i++ {
		for j := i + 1; j < len(counts); j++ {
			if counts[i].count < counts[j].count {
				counts[i], counts[j] = counts[j], counts[i]
			}
		}
	}

	result := make(map[string]int)
	for i := 0; i < len(counts) && i < limit; i++ {
		result[counts[i].addr] = counts[i].count
	}
	return result
}
