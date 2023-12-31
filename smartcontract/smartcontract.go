// Package smartcontract defines the smartcontract data and the link to the abi
package smartcontract

import (
	"fmt"
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/common-lib/topic"
)

// Smartcontract The storage smartcontract
// It keeps the read-only parameters such as
// associated ABI, deployer, address, block parameter as well as the transaction
// where it was deployed.
//
// The Database interaction depends on the sds/storage/abi
type Smartcontract struct {
	// Topic keeps the information about the smartcontract
	// The topic string should be unique for each smartcontract
	// At least the following parameters are required to be:
	// - organization (team, developer)
	// - project  (dapp name)
	// - group (classification, for example: nft, token)
	// - network id (network where the smartcontract was deployed)
	// - name (smartcontract name. recommended that it matches to the file name)
	Topic         topic.Topic `json:"topic"`
	TransactionId string      `json:"transaction_id"`
	Owner         string      `json:"owner,omitempty"`
	Verifier      string      `json:"verifier,omitempty"`
	// Specific parameters of the smartcontract based on the network.
	//
	// For example, Ethereum's data:
	// - abiId
	// - address
	//
	// Sui blockchain's:
	// - packageId
	// - moduleId
	// - resourceId <optional>
	// - resourceType
	Specific key_value.KeyValue `json:"specific"`
}

func (sm *Smartcontract) Validate() error {
	if len(sm.TransactionId) == 0 {
		return fmt.Errorf("missing TransactionId")
	}

	return nil
}

func (sm *Smartcontract) Id() topic.Id {
	return sm.Topic.Id().Only("org", "net", "name")
}
