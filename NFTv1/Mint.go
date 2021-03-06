package NFTv1

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	rp "github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

func Mint(fromWalletSecret []byte) (*common.PublicKey, *common.PublicKey, *string, error) {

	// var feePayer, _ = types.AccountFromBytes(feePayerSecret)
	var fromWallet, _ = types.AccountFromBytes(fromWalletSecret)

	c := client.NewClient(rpc.TestnetRPCEndpoint)

	mint := types.NewAccount()

	ata, _, err := common.FindAssociatedTokenAddress(fromWallet.PublicKey, mint.PublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	tokenMetadataPubkey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	tokenMasterEditionPubkey, err := tokenmeta.GetMasterEdition(mint.PublicKey)
	if err != nil {
		return nil, nil, nil, err
	}

	mintAccountRent, err := c.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		return nil, nil, nil, err
	}

	recentBlockhashResponse, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return nil, nil, nil, err
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, fromWallet},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        fromWallet.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     fromWallet.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals: 0,
					Mint:     mint.PublicKey,
					MintAuth: fromWallet.PublicKey,
				}),
				tokenmeta.CreateMetadataAccount(tokenmeta.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           fromWallet.PublicKey,
					Payer:                   fromWallet.PublicKey,
					UpdateAuthority:         fromWallet.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: tokenmeta.Data{
						Name:                 "Fake SMS #1355",
						Symbol:               "FSMB",
						Uri:                  "https://34c7ef24f4v2aejh75xhxy5z6ars4xv47gpsdrei6fiowptk2nqq.arweave.net/3wXyF1wvK6ARJ_9ue-O58CMuXrz5nyHEiPFQ6z5q02E",
						SellerFeeBasisPoints: 100,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  fromWallet.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				}),
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 fromWallet.PublicKey,
					Owner:                  fromWallet.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				tokenprog.MintTo(tokenprog.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   fromWallet.PublicKey,
					Amount: 1,
				}),
				tokenmeta.CreateMasterEdition(tokenmeta.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: fromWallet.PublicKey,
					MintAuthority:   fromWallet.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           fromWallet.PublicKey,
					MaxSupply:       pointer.Uint64(0),
				}),
			},
		}),
	})
	if err != nil {
		return nil, nil, nil, err
	}

	sign, err := c.SendTransactionWithConfig(context.TODO(), tx, client.SendTransactionConfig{
		SkipPreflight:       false,
		PreflightCommitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		return nil, nil, nil, err
	}

	wsClient, err := ws.Connect(context.Background(), rp.TestNet_WS)
	sub, err := wsClient.SignatureSubscribe(
		solana.MustSignatureFromBase58(sign),
		rp.CommitmentFinalized,
	)

	if err != nil {
		return nil, nil, nil, err
	}
	defer sub.Unsubscribe()

	for {
		got, err := sub.Recv()
		if err != nil {
			return &mint.PublicKey, &fromWallet.PublicKey, &sign, nil
		}
		if got.Value.Err != nil {
			panic(errors.New("transaction confirmation failed"))
		} else {
			return &mint.PublicKey, &fromWallet.PublicKey, &sign, nil
		}
	}

}
