// Package account handles the user information for authentication that accessed
// to SDS Service.
//
// It has two types of accounts:
//   - Account for developers, it uses CurveZMQ for the authentication.
//   - SmartcontractDeveloper for smartcontract developers.
//     It uses blockchain public/private key for authentication.
package account

import "github.com/blocklords/sds/common/blockchain"

// Account is the developer that accessed to the SDS Service.
type Account struct {
	PublicKey      string               `json:"public_key"`      // Public Key for authentication.
	Organization   string               `json:"organization"`    // Organization
	NonceTimestamp blockchain.Timestamp `json:"nonce_timestamp"` // Nonce since the last usage. Only acceptable for developers
}

type Accounts []*Account

// NewFromPublicKey creates the account from the public key
func NewFromPublicKey(publicKey string) *Account {
	return &Account{
		PublicKey:      publicKey,
		Organization:   "",
		NonceTimestamp: 0,
	}
}

// New Account from the fields
func New(publicKey string, timestamp blockchain.Timestamp, organization string) *Account {
	return &Account{
		PublicKey:      publicKey,
		NonceTimestamp: timestamp,
		Organization:   organization,
	}
}

///////////////////////////////////////////////////////////
//
// Group operations
//
///////////////////////////////////////////////////////////

// NewAccounts converts list of accounts into Accounts
func NewAccounts(newAccounts ...*Account) Accounts {
	accounts := make(Accounts, len(newAccounts))
	copy(accounts, newAccounts)

	return accounts
}

// PublicKeys of the accounts
func (accounts Accounts) PublicKeys() []string {
	publicKeys := make([]string, len(accounts))

	for i, account := range accounts {
		publicKeys[i] = account.PublicKey
	}

	return publicKeys
}

// Add newAccounts to the list of Accounts
//
// Example:
//
//		accounts.Add(acc_1, acc_2).
//	 	Add(acc_3, acc_4)
func (accounts Accounts) Add(newAccounts ...*Account) Accounts {
	for _, account := range newAccounts {
		accounts = append(accounts, account)
	}

	return accounts
}

func (accounts Accounts) Remove(newAccounts ...*Account) Accounts {
	for _, account := range newAccounts {
		for i := range accounts {
			if account.PublicKey == accounts[i].PublicKey {
				accounts = append(accounts[:i], accounts[i+1:]...)
				return accounts
			}
		}
	}

	return accounts
}