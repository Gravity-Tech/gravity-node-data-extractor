package tests

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)


func generateEthKeys(n int) []*ecdsa.PrivateKey {
	var result []*ecdsa.PrivateKey

	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		key.
		result = append(result, key)
	}

	return result
}

func generateEthAcccounts(n int) []*core.GenesisAccount {
	keys := generateEthKeys(n)

	var result []*core.GenesisAccount
	for key := range keys {
		result = append(result, core.GenesisAccount{
			Balance:    balanceWithBase(10, 18),
			PrivateKey: hexutil..,
		})
	}
}

func balanceWithBase(balance, dec int64) *big.Int {
	initialBalance := big.NewInt(balance)
	initialBalance.Exp(big.NewInt(10), big.NewInt(dec), nil)
	return initialBalance
}

func TestExtractionWavesSourceLock(t *testing.T) {
	// Generate a new random account and a funded simulator
	keys := generateEthKeys(3)

	nebulaAddress, portUserAddress, portAddress := keys[0], keys[1], keys[2]
	auth := bind.NewKeyedTransactor(key)

	sim := backends.NewSimulatedBackend(core.GenesisAccount{ Balance: initialBalance}, 100_000)

	// Deploy a token contract on the simulated blockchain
	_, _, token, err := DeployMyToken(auth, sim, new(big.Int), "Simulated blockchain tokens", 0, "SBT")
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	// Print the current (non existent) and pending name of the contract
	name, _ := token.Name(nil)
	fmt.Println("Pre-mining name:", name)

	name, _ = token.Name(&bind.CallOpts{Pending: true})
	fmt.Println("Pre-mining pending name:", name)

	// Commit all pending transactions in the simulator and print the names again
	sim.Commit()

	name, _ = token.Name(nil)
	fmt.Println("Post-mining name:", name)

	name, _ = token.Name(&bind.CallOpts{Pending: true})
	fmt.Println("Post-mining pending name:", name)
}
