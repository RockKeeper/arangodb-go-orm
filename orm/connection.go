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

func (dc *DatabaseConnection) DB(database string) *DatabaseConnection {

	dc.currentDatabase = dc.GetDB(database)
	dc.currentCollection = nil

	return dc
}

func (dc *DatabaseConnection) GetDB(database string) driver.Database {

	db, err := dc.httpClient.Database(dc.currentContext, database)
	if err != nil {
		panic(err)
	}

	return db
}

func (dc *DatabaseConnection) Collection(collectionName string) *Collection {
	collection, err := dc.currentDatabase.Collection(dc.currentContext, collectionName)
	if err != nil {
		// handle error
	}

	currentCollection := &Collection{
		Name:               collectionName,
		collection:         collection,
		databaseConnection: dc,
	}

	return currentCollection
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
