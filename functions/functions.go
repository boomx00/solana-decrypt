package functions

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Config struct {
	RPC string `json:"rpc"`
	WS  string `json:"ws"`
}

func LoadConfig(filename string) *Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Errorf("error : %v", err)
	}

	var config Config

	errUnmarshal := json.Unmarshal(data, &config)
	if errUnmarshal != nil {
		fmt.Errorf("error: %v", errUnmarshal)
	}
	return &config
}
func DecodeTokenTransfers(tx *rpc.GetTransactionResult) {
	decoded, err := tx.Transaction.GetTransaction()
	if err != nil {
		fmt.Printf("Error decoding transaction: %v\n", err)
		return
	}
	accountKeys := decoded.Message.AccountKeys

	fmt.Printf("Total account keys: %d\n", len(accountKeys))

	for _, inner := range tx.Meta.InnerInstructions {
		fmt.Printf("\nInner Instructions for Index %d:\n", inner.Index)

		for i, inst := range inner.Instructions {
			// Bounds check for program ID
			if inst.ProgramIDIndex >= uint16(len(accountKeys)) {
				fmt.Printf("Program ID index %d out of bounds\n", inst.ProgramIDIndex)
				continue
			}

			programID := accountKeys[inst.ProgramIDIndex]

			if programID == solana.TokenProgramID {
				data := inst.Data

				if len(data) > 0 {
					discriminator := data[0]

					if discriminator == 3 && len(data) >= 9 {
						amount := binary.LittleEndian.Uint64(data[1:9])

						fmt.Printf("\nTransfer #%d:\n", i+1)

						// Safely get account addresses with bounds checking
						if len(inst.Accounts) >= 3 {
							for j, accountIdx := range inst.Accounts {
								if accountIdx >= uint16(len(accountKeys)) {
									fmt.Printf("Account index %d out of bounds\n", accountIdx)
									continue
								}

								account := accountKeys[accountIdx]
								switch j {
								case 0:
									fmt.Printf("From: %s\n", account)
								case 1:
									fmt.Printf("To: %s\n", account)
								case 2:
									fmt.Printf("Authority: %s\n", account)
								}
							}
						}

						fmt.Printf("Amount: %d (raw amount, divide by decimals for actual token amount)\n", amount)
					} else {
						fmt.Printf("\nNon-transfer instruction #%d:\n", i+1)
						fmt.Printf("Program: %s\n", programID)
						fmt.Printf("Discriminator: %d\n", discriminator)
						fmt.Printf("Raw data: %v\n", data)

						fmt.Printf("Accounts involved:\n")
						for _, accountIdx := range inst.Accounts {
							if accountIdx < uint16(len(accountKeys)) {
								fmt.Printf("- %s\n", accountKeys[accountIdx])
							} else {
								fmt.Printf("Account index %d out of bounds\n", accountIdx)
							}
						}
					}
				}
			}
		}
	}
}
