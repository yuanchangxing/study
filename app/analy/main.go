package main

import (
	"context"
	"encoding/json"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/yuanchangxing/study/tools"
	"github.com/yuanchangxing/study/xlog"
	"net/http"
)

type Holder struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

func getTopHolders(tokenMint string) ([]Holder, error) {
	resp, err := http.Get("https://api.solscan.io/token/holders?tokenAddress=" + tokenMint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Data []Holder `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

var (
	solCli *rpc.Client
	//rpcAddress = "https://api.mainnet-beta.solana.com"
	//rpcAddress = "https://api.devnet.solana.com"
	rpcAddress = "https://go.getblock.io/9dbc6cbf8bd54e17938eb7ef5ffd63a5"
)

var limit = tools.NewRateLimit(10, 30)

func InitAndSearchBalance(pubAddress string) error {
	solCli = rpc.New(rpcAddress)
	ctx := context.Background()

	// 查询特定地址的余额（SOL）
	address := solana.MustPublicKeyFromBase58(pubAddress)
	balance, err := solCli.GetBalance(ctx, address, rpc.CommitmentFinalized)
	if err != nil {
		return err
	}
	xlog.Infof("address: %s  balance :%d", address, balance.Value)
	return nil
}

func main() {
	myAddress := "24damVGETUTZvLdZZpiUdSX8qsWDKSoKz3uWd7vGJ9GM"
	err := InitAndSearchBalance(myAddress)
	if err != nil {
		xlog.Errorf(err.Error())
		return
	}

	pub, err := solana.PublicKeyFromBase58(myAddress)
	if err != nil {
		xlog.Errorf(err.Error())
		return
	}

	err = history(context.TODO(), pub)
	if err != nil {
		xlog.Errorf(err.Error())
	}
}

func history(ctx context.Context, address solana.PublicKey) error {
	var num = 100
	signatures, err := solCli.GetSignaturesForAddressWithOpts(ctx, address, &rpc.GetSignaturesForAddressOpts{
		Limit: &num,
	})
	if err != nil {
		return err
	}
	xlog.Logger.Info("signatures:", signatures)
	xlog.Infof("signatures: %+v", len(signatures))
	for _, sig := range signatures {
		limit.Wait()
		tx, err := solCli.GetTransaction(ctx, sig.Signature, &rpc.GetTransactionOpts{})
		if err != nil {
			xlog.Infof("获取交易 %s 失败: %v", sig.Signature, err)
			continue
		}
		// 解析交易，检查是否为代币转账
		xlog.Infof("交易 %s: %v", sig.Signature, tx)
	}
	return nil
}
