package main

import (
	"fmt"
	"log"

	"github.com/benjamin-maynard/kubernetes-cloud-mysql-backup/internal/config"
	"github.com/benjamin-maynard/kubernetes-cloud-mysql-backup/internal/databases"
)

func main() {

	// Initialise the application configuration
	config, err := config.NewConfigFromEnvironment()
	if err != nil {
		log.Fatalf("error loading application configuration: %v", err)
	}

	// Build the list of databases if we are backing up all databases
	if config.TargetAllDatabases {
		config.TargetDatabaseNames, err = databases.ListDatabases(config.DatabaseBackupConfig)
		if err != nil {
			log.Fatalf("error loading databases to backup: %v", err)
		}
	}

	// Build a new database lists from our config
	dbList := databases.NewDatabaseList(config.TargetDatabaseNames)

	// Perform the backup activities
	dbList.ProcessBackups(config)

	for _, val := range dbList {
		fmt.Println(val.LocalBackupPath)
	}

}
