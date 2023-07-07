package handler

import (
	"github.com/Seascape-Foundation/sds-service-lib/log"
	"github.com/Seascape-Foundation/sds-service-lib/remote"
	"github.com/Seascape-Foundation/static-seascape-service/smartcontract"

	"github.com/Seascape-Foundation/sds-common-lib/data_type/database"
	"github.com/Seascape-Foundation/sds-common-lib/smartcontract_key"

	"github.com/Seascape-Foundation/sds-service-lib/communication/command"
	"github.com/Seascape-Foundation/sds-service-lib/communication/message"
)

type SetSmartcontractRequest = smartcontract.Smartcontract
type SetSmartcontractReply = smartcontract.Smartcontract
type GetSmartcontractRequest = smartcontract_key.Key
type GetSmartcontractReply = smartcontract.Smartcontract

// SmartcontractRegister Register a new smartcontract. It means we are adding smartcontract parameters into
// smartcontract database table.
// Requires abi_id parameter. First call abi_register method first.
func SmartcontractRegister(request message.Request, _ log.Logger, clients remote.Clients) message.Reply {
	var sm SetSmartcontractRequest
	err := request.Parameters.Interface(&sm)
	if err != nil {
		return message.Fail("failed to parse data")
	}
	if err := sm.Validate(); err != nil {
		return message.Fail("failed to validate: " + err.Error())
	}

	var reply = sm
	replyMessage, err := command.Reply(&reply)
	if err != nil {
		return message.Fail("failed to reply")
	}

	dbCon := remote.GetClient(clients, "database")

	var crud database.Crud = &sm
	if err = crud.Insert(dbCon); err != nil {
		return message.Fail("Smartcontract saving in the database failed: " + err.Error())
	}

	return replyMessage
}

// SmartcontractGet Returns configuration and smartcontract information related to the configuration
var SmartcontractGet = func(request message.Request, _ log.Logger, clients remote.Clients) message.Reply {
	var key GetSmartcontractRequest
	err := request.Parameters.Interface(&key)
	if err != nil {
		return message.Fail("failed to parse data")
	}
	if err := key.Validate(); err != nil {
		return message.Fail("key.Validate: " + err.Error())
	}

	dbCon := remote.GetClient(clients, "database")

	var selected = smartcontract.Smartcontract{}
	var crud database.Crud = &selected
	err = crud.Select(dbCon, &selected)
	if err != nil {
		return message.Fail("failed to get configuration from the database: " + err.Error())
	}

	replyMessage, err := command.Reply(&selected)
	if err != nil {
		return message.Fail("failed to reply")
	}

	return replyMessage
}
