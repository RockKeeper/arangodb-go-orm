package orm

import (
	"context"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type DatabaseConnectionData struct {
	Host     string
	Username string
	Password string
	Database string
}

type DatabaseConnection struct {
	httpClient        driver.Client
	connectionData    *DatabaseConnectionData
	currentDatabase   driver.Database
	currentCollection *Collection
	currentContext    context.Context
}

func NewDatabaseConnection(connectionData *DatabaseConnectionData) (*DatabaseConnection, error) {

	databaseConnection := &DatabaseConnection{}

	connection, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{connectionData.Host},
	})

	if err != nil {
		log.Fatalln("==> ERROR")
		log.Fatal(err)
		log.Fatalln("--------------------------------")
		return nil, err
	}
	databaseClient, err := driver.NewClient(driver.ClientConfig{
		Connection:     connection,
		Authentication: driver.BasicAuthentication(connectionData.Username, connectionData.Password),
	})
	if err != nil {
		log.Fatalln("==> ERROR")
		log.Fatal(err)
		log.Fatalln("--------------------------------")
		return nil, err
	}

	databaseConnection.httpClient = databaseClient
	databaseConnection.connectionData = connectionData
	databaseConnection.currentDatabase = databaseConnection.GetDB(connectionData.Database)
	databaseConnection.currentContext = context.Background()

	return databaseConnection, nil
}
