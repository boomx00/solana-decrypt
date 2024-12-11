package main

import (
	"context"
	"fmt"
	"go-prac/functions"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

func main() {
	config := functions.LoadConfig("config.json")
	client := rpc.New(config.RPC)
	txSig := solana.MustSignatureFromBase58("3w6TGrCSZyD6ReKxkUocoMoDmX3WjkozA6HcENhUvva5R8bPA4PRwz3np4VQsDQ5K47GovRpmsjHJZCp7Tkb2VsV")

	// get transaction meta data
	var version uint64 = 0
	output, err := client.GetTransaction(
		context.TODO(),
		txSig,
		&rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			MaxSupportedTransactionVersion: &version,
			Commitment:                     rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		fmt.Println(err)
	}
	functions.DecodeTokenTransfers(output)

	// fmt.Println(output.Meta.InnerInstructions[0])
	// spew.Dump(output.Meta)
	// // get transaction instructions
	// txInst, err := solana.TransactionFromDecoder(bin.NewBinDecoder(output.Transaction.GetBinary()))
	// txMessage := txInst.Message
	// spew.Dump(txMessage)
	// for i, inst := range txMessage.Instructions {
	// 	fmt.Printf("\nInstruction %d:\n", i)
	// 	fmt.Printf("Program ID: %s\n", txMessage.AccountKeys[inst.ProgramIDIndex])

	// }
}
