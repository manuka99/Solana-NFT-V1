package NFTv1

import (
	"context"
	"errors"
	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func Transfer(mintPK common.PublicKey, fromWalletSecret []byte, toWalletSecret []byte) (*string, error) {

	var fromWallet, _ = types.AccountFromBytes(fromWalletSecret)
	var toWallet, _ = types.AccountFromBytes(toWalletSecret)

	fromTokenAccount, err := GetTokenAccount(mintPK, fromWallet.PublicKey)
	if err != nil {
		return nil, err
	}

	toTokenAccount, err := GetTokenAccount(mintPK, toWallet.PublicKey)
	if err != nil {
		toTokenAccount, _, err = CreateTokenAccount(mintPK, toWalletSecret)
		if err != nil {
			return nil, err
		}
	}

	c := client.NewClient(rpc.TestnetRPCEndpoint)

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, err
	}
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        toWallet.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				tokenprog.TransferChecked(tokenprog.TransferCheckedParam{
					From:     *fromTokenAccount,
					To:       *toTokenAccount,
					Mint:     mintPK,
					Auth:     fromWallet.PublicKey,
					Signers:  []common.PublicKey{},
					Amount:   1,
					Decimals: 0,
				}),
			},
		}),
		Signers: []types.Account{toWallet, fromWallet},
	})
	if err != nil {
		return nil, err
	}

	txhash, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil, err
	}

	wsClient, err := ws.Connect(context.Background(), rp.TestNet_WS)
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(txhash),
		rp.CommitmentFinalized,
	)

	if err != nil {
		panic(err)
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &txhash, nil
		}
		if got.Value.Err != nil {
			return nil, errors.New("transaction confirmation failed")
		} else {
			return &txhash, nil
		}
	}

}
