package smartcontract

import (
	"testing"

	"github.com/ahmetson/common-lib/blockchain"
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/common-lib/smartcontract_key"
	"github.com/stretchr/testify/suite"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestSmartcontractSuite struct {
	suite.Suite
	smartcontract Smartcontract
}

func (suite *TestSmartcontractSuite) SetupTest() {
	_, _ = smartcontract_key.New("1", "0xaddress")
	_, _ = blockchain.NewHeader(uint64(1), uint64(23))
	suite.smartcontract = Smartcontract{}
}

func (suite *TestSmartcontractSuite) TestNew() {
	// creating a new smartcontract from empty parameter
	// should fail
	kv := key_value.Empty()
	_, err := New(kv)
	suite.Require().Error(err)

	// creating a new smartcontract with the exact type
	// should be successful
	key, _ := smartcontract_key.New("1", "0xaddress")
	abiId := "base64="
	txKey := blockchain.TransactionKey{
		Id:    "0xtx_id",
		Index: 0,
	}
	header, _ := blockchain.NewHeader(uint64(1), uint64(23))
	deployer := "0xahmetson"
	kv = key_value.Empty().
		Set("key", key).
		Set("abi_id", abiId).
		Set("transaction_key", txKey).
		Set("block_header", header).
		Set("deployer", deployer)

	smartcontract, err := New(kv)
	suite.Require().NoError(err)
	suite.Require().EqualValues(suite.smartcontract, *smartcontract)

	// creating a smartcontract with the missing data
	// should fail.
	// In this case Transaction key's Topic is missing
	txKey = blockchain.TransactionKey{Index: 0}
	kv = key_value.Empty().
		Set("key", key).
		Set("abi_id", abiId).
		Set("transaction_key", txKey).
		Set("block_header", header).
		Set("deployer", deployer)
	_, err = New(kv)
	suite.Require().Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSmartcontract(t *testing.T) {
	suite.Run(t, new(TestSmartcontractSuite))
}
