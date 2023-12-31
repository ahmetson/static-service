package configuration

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	database "github.com/Seascape-Foundation/mysql-seascape-extension"
	"github.com/ahmetson/common-lib/blockchain"
	"github.com/ahmetson/common-lib/smartcontract_key"
	"github.com/ahmetson/common-lib/topic"
	"github.com/ahmetson/service-lib/configuration"
	parameter "github.com/ahmetson/service-lib/identity"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
	"github.com/ahmetson/static-service/abi"
	"github.com/ahmetson/static-service/smartcontract"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

// We won't test the requests.
// The requests are tested in the controllers
// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TestConfigurationDbSuite struct {
	suite.Suite
	dbName        string
	configuration Configuration
	container     *mysql.MySQLContainer
	dbCon         *remote.ClientSocket
	ctx           context.Context
}

func (suite *TestConfigurationDbSuite) SetupTest() {
	// prepare the database creation
	suite.dbName = "test"
	_, filename, _, _ := runtime.Caller(0)
	// configuration depends on smartcontract.
	// smartcontract depends on abi.
	fileDir := filepath.Dir(filename)
	storageAbi := "20230308171023_storage_abi.sql"
	storageSmartcontract := "20230308173919_storage_smartcontract.sql"
	storageConfiguration := "20230308173943_storage_configuration.sql"
	changeGroupType := "20230314150414_storage_configuration_group_type.sql"

	abiSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageAbi)
	smartcontractSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageSmartcontract)
	configurationSqlPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", storageConfiguration)
	changeGroupPath := filepath.Join(fileDir, "..", "..", "_db", "migrations", changeGroupType)

	suite.T().Log("the configuration table path", configurationSqlPath)

	// run the container
	ctx := context.TODO()
	container, err := mysql.RunContainer(ctx,
		mysql.WithDatabase(suite.dbName),
		mysql.WithUsername("root"),
		mysql.WithPassword("tiger"),
		mysql.WithScripts(abiSqlPath, smartcontractSqlPath, configurationSqlPath, changeGroupPath),
	)
	suite.Require().NoError(err)
	suite.container = container
	suite.ctx = ctx

	logger, err := log.New("mysql-suite", false)
	suite.Require().NoError(err)
	appConfig, err := configuration.NewAppConfig(logger)
	suite.Require().NoError(err)

	// Creating a database client
	// after settings the default parameters
	// we should have the username and password

	// Overwrite the default parameters to use test container
	host, err := container.Host(ctx)
	suite.Require().NoError(err)
	ports, err := container.Ports(ctx)
	suite.Require().NoError(err)
	exposedPort := ports["3306/tcp"][0].HostPort

	database.DatabaseConfigurations.Parameters["SDS_DATABASE_HOST"] = host
	database.DatabaseConfigurations.Parameters["SDS_DATABASE_PORT"] = exposedPort
	database.DatabaseConfigurations.Parameters["SDS_DATABASE_NAME"] = suite.dbName

	//go database.Run(appConfig, logger)
	// wait for initiation of the controller
	time.Sleep(time.Second * 1)

	databaseService, err := parameter.Inprocess("database")
	suite.Require().NoError(err)
	client, err := remote.InprocRequestSocket(databaseService.Url(), logger, appConfig)
	suite.Require().NoError(err)

	suite.dbCon = client

	// add the storage abi
	abiId := "base64="
	sampleAbi := abi.Abi{
		Body: "[{}]",
		Id:   abiId,
	}
	err = sampleAbi.Insert(suite.dbCon)
	suite.Require().NoError(err)

	// add the storage smartcontract
	_, _ = smartcontract_key.New("1", "0xaddress")
	_, _ = blockchain.NewHeader(uint64(1), uint64(23))
	sm := smartcontract.Smartcontract{}
	err = sm.Insert(suite.dbCon)
	suite.Require().NoError(err)

	sample := topic.Topic{
		Organization: "seascape",
		Project:      "sds-core",
		NetworkId:    "1",
		Group:        "test-suite",
		Name:         "TestErc20",
	}
	suite.configuration = Configuration{
		Topic: sample,
	}

	suite.T().Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			suite.T().Fatalf("failed to terminate container: %s", err)
		}
		if err := suite.dbCon.Close(); err != nil {
			suite.T().Fatalf("failed to terminate database connection: %s", err)
		}
	})
}

func (suite *TestConfigurationDbSuite) TestConfiguration() {
	var configs []*Configuration

	err := suite.configuration.SelectAll(suite.dbCon, &configs)
	suite.Require().NoError(err)
	suite.Require().Len(configs, 0)

	err = suite.configuration.Insert(suite.dbCon)
	suite.Require().NoError(err)

	err = suite.configuration.SelectAll(suite.dbCon, &configs)
	suite.Require().NoError(err)
	suite.Require().Len(configs, 1)
	suite.Require().EqualValues(suite.configuration, *configs[0])

	// inserting a configuration
	// that links to the non-existing smartcontract
	// should fail
	sample := topic.Topic{
		Organization: "seascape",
		Project:      "sds-core",
		NetworkId:    "1",
		Group:        "test-suite",
		Name:         "TestToken",
	}
	conf := Configuration{
		Topic: sample,
	}
	err = conf.Insert(suite.dbCon)
	suite.Require().Error(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestConfigurationDb(t *testing.T) {
	suite.Run(t, new(TestConfigurationDbSuite))
}
