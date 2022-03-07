package main

import (
	"fmt"
	"solana-nft/NFTv1"
)

func main() {

	var fromWalletSecret = []byte{10, 75, 10, 90, 145, 78, 142, 248, 104, 3, 36, 7, 69, 207, 109, 98, 82, 58, 146, 202, 44, 188, 70, 70, 64, 173, 35, 130, 18, 133, 107, 236, 231, 43, 70, 165, 182, 191, 162, 242, 126, 119, 49, 3, 231, 43, 249, 47, 228, 225, 70, 91, 254, 22, 160, 42, 20, 186, 184, 196, 240, 151, 157, 207}

	var toWalletSecret = []byte{47, 163, 68, 180, 12, 82, 124, 0, 101, 163, 250, 17, 181, 250, 63, 165, 179, 85, 112, 117, 245, 102, 63, 181, 48, 68, 190, 193, 178, 112, 227, 57, 17, 239, 150, 83, 192, 134, 121, 241, 161, 240, 133, 128, 9, 112, 247, 2, 71, 181, 138, 177, 227, 201, 12, 225, 164, 158, 122, 91, 176, 169, 10, 147}

	fmt.Println("\n..................... BEGIN MINTING NFT ...................")
	mintPK, ownerPK, mintedTxHash, _ := NFTv1.Mint(fromWalletSecret)
	fmt.Println("\nMINTED PK", mintPK)
	fmt.Println("OWNER PK", ownerPK)
	fmt.Println("TX HASH", *mintedTxHash)
	fmt.Println("\n..................... END MINTING NFT ...................")

	fmt.Println("\n\n................ BEGIN NFT TRANSFER .................")
	transferTXHash, _ := NFTv1.Transfer(*mintPK, fromWalletSecret, toWalletSecret)
	fmt.Println("\nTX HASH", *transferTXHash)
	fmt.Println("\n................ END NFT TRANSFER .................")
}

// 9aE476sH92Vz7DMPyq5WLPkrKWivxeuTKEFKd2sZZcde
