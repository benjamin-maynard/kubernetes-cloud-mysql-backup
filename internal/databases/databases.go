package databases

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"filippo.io/age"
	"filippo.io/age/agessh"
	"github.com/benjamin-maynard/kubernetes-cloud-mysql-backup/internal/config"
)

// ListDatabases returns a comma seperated string of databases on the given MySQL Host
// it ignores common known databases
func ListDatabases(config config.DatabaseBackupConfig) (string, error) {

	// Build query
	excludedDatabaseList := []string{"'mysql'", "'sys'", "'tmp'", "'information_schema'", "'performance_schema'"}
	excludedDatabaseSQLStatement := fmt.Sprintf("SELECT GROUP_CONCAT(schema_name) FROM information_schema.schemata WHERE schema_name NOT IN (%s)", strings.Join(excludedDatabaseList, ","))

	// Create a new slice for storing the command
	executeCmd := []string{
		"mysql",
		fmt.Sprintf("--user=%s", config.TargetDatabaseUser),           // Database User
		fmt.Sprintf("--host=%s", config.TargetDatabaseHost),           // Database Host
		fmt.Sprintf("--port=%v", config.TargetDatabasePort),           // Database Port
		fmt.Sprintf("--password='%s'", config.TargetDatabasePassword), // Database Password
		"--no-auto-rehash",
		"--skip-column-names",
		fmt.Sprintf("--execute=\"%s\"", excludedDatabaseSQLStatement), // SQL Statement to List Databases
	}

	// Execute the command to list databases
	result, err := executeShellCmd(executeCmd)
	if err != nil {
		return "", err
	}

	// Split the databases into a slice
	return result, nil

}

// DatabaseList refers to a collection of Databases
type DatabaseList []*Database

// NewDatabaseList returns a slice of Databases, populated from a comma seperated
// string of database names
func NewDatabaseList(dbList string) DatabaseList {
	list := DatabaseList{}
	for _, database := range strings.Split(dbList, ",") {
		list = append(list, &Database{DatabaseName: database})
	}
	return list
}

// ProcessBackups perform all backup and upload activities for databases
func (d DatabaseList) ProcessBackups(config config.AppConfig) {

	// Loop through all databases
	for _, db := range d {

		// Perform the actual backup
		err := db.dumpToFile(config.DatabaseBackupConfig)
		// If we got an error, log and move on to the next database
		if err != nil {
			log.Printf("error dumping database '%s', got error: %v\n", db.DatabaseName, err)
			db.BackupError = err.Error()
			continue
		}

		// Compress the backup if neccesary
		if config.BackupCompress {
			err = db.compress(config.DatabaseBackupConfig)
			// If we got an error, log and move on to the next database
			if err != nil {
				log.Printf("error compressing database '%s', got error: %v\n", db.DatabaseName, err)
				db.BackupError = err.Error()
				continue
			}
		}

		// Encrypt the backups if neccesary
		if config.AgePublicKey != "" {
			err := db.encrypt(config.DatabaseBackupConfig)
			// If we got an error, log and move on to the next database
			if err != nil {
				log.Printf("error encryption database '%s', got error: %v\n", db.DatabaseName, err)
				db.BackupError = err.Error()
				continue
			}
		}
	}

}

// Database stores information on a database throughout the backup lifecycle
type Database struct {
	// The name of the MySQL Database
	DatabaseName string
	// The path on the local system where the backup file is stored
	LocalBackupPath string
	// If there was an error backing up the database at any point
	// a string representation of the error is stored here
	BackupError string
	// Path where the backup was uploaded, may be in a range of
	// formats depending on the destination
	WrittenDestination string
}

// buildInitialFilePath builds the file path where the backup will be stored
// the file may later be renamed, for example when compressed or encrypted
func (d Database) buildInitialFilePath(config config.DatabaseBackupConfig) string {

	// Create the file path where the backup will be stored
	var timestamp string
	if config.BackupTimestamp != "" {
		timestamp = fmt.Sprintf("%s-", time.Now().Format(config.BackupTimestamp))
	}
	return fmt.Sprintf("/dumps/%s%s.sql", timestamp, d.DatabaseName)

}

// buildDumpCommand builds the command that will be executed to
// perform the database dump
func (d Database) buildDumpCommand(config config.DatabaseBackupConfig) []string {

	// Create a new slice for storing the command
	executeCmd := []string{
		"mysqldump",
		fmt.Sprintf("--user=%s", config.TargetDatabaseUser), // Database User
		fmt.Sprintf("--host=%s", config.TargetDatabaseHost), // Database Host
		fmt.Sprintf("--port=%v", config.TargetDatabasePort), // Database Port
		"--password=$TARGET_DATABASE_PASSWORD",              // Database password, pulled directly from the environment to prevent password leakage
		config.BackupAdditionalParams,                       // Additional user parameters
		d.DatabaseName,                                      // Database Name
	}

	// If the BackupCreateDatabaseStatement variable is set, then we need to add the --databases flag
	// to the MySQL command so that mysqldump will add the CREATE and USE statements
	if config.BackupCreateDatabaseStatement {
		executeCmd = append(executeCmd[:len(executeCmd)-1], append([]string{"--databases"}, executeCmd[len(executeCmd)-1:]...)...)
	}

	return executeCmd
}

// dumpToFile performs a MySQL dump of the database returning an error if it fails
// it also updates the Database.
func (d *Database) dumpToFile(config config.DatabaseBackupConfig) error {

	// Build initial database dump command
	dumpCmd := d.buildDumpCommand(config)

	// Get target file path and append backup path to command
	path := d.buildInitialFilePath(config)
	dumpCmd = append(dumpCmd, fmt.Sprintf("> %s", path))

	// Execute the command to perform the dump
	_, err := executeShellCmd(dumpCmd)
	if err != nil {
		return err
	}

	// Update the Database.LocalBackupPath with the path to the backup file
	d.LocalBackupPath = path

	return nil

}

// compress performs gzip compression on the database local backup file
func (d *Database) compress(config config.DatabaseBackupConfig) error {

	// Create the target file
	archivePath := fmt.Sprintf("%s.gz", d.LocalBackupPath)
	archive, err := os.Create(archivePath)
	if err != nil {
		return fmt.Errorf("error creating empty backup archive '%s', got error: '%v'", archivePath, err)
	}
	defer archive.Close()

	// Open the uncompressed database
	backup, err := os.Open(d.LocalBackupPath)
	if err != nil {
		return fmt.Errorf("error loading database backup '%s', got error: '%v'", d.LocalBackupPath, err)
	}
	defer backup.Close()

	// Create gzip writer and write the file
	writer := gzip.NewWriter(archive)
	defer writer.Close()
	_, err = io.Copy(writer, backup)
	if err != nil {
		return fmt.Errorf("error compressing backup '%s', got error: '%v'", d.LocalBackupPath, err)
	}

	// Delete the original backup
	err = os.Remove(d.LocalBackupPath)
	if err != nil {
		return fmt.Errorf("error deleting original backup '%s', got error: '%v'", d.LocalBackupPath, err)
	}

	// Update the path to the compressed path
	d.LocalBackupPath = archivePath

	return nil

}

// encrypt encrypts the local database file using the key provided in the config
func (d *Database) encrypt(config config.DatabaseBackupConfig) error {

	// Open the existing database backup
	backup, err := os.Open(d.LocalBackupPath)
	if err != nil {
		return fmt.Errorf("error loading database backup '%s', got error: '%v'", d.LocalBackupPath, err)
	}
	defer backup.Close()

	// Create the target file
	encryptedPath := fmt.Sprintf("%s.age", d.LocalBackupPath)
	encrypted, err := os.Create(encryptedPath)
	if err != nil {
		return fmt.Errorf("error creating empty encryption archive '%s', got error: '%v'", encryptedPath, err)
	}
	defer encrypted.Close()

	// Parse the public key for encrption
	recipient, err := agessh.ParseRecipient(config.AgePublicKey)
	if err != nil {
		return fmt.Errorf("error parsing public key, got error: '%v'", err)
	}

	// Createt the encryption writer
	w, err := age.Encrypt(encrypted, recipient)
	if err != nil {
		return fmt.Errorf("error encrypted file '%s', got error: '%v'", encryptedPath, err)
	}
	defer w.Close()

	// Encrypt the file
	_, err = io.Copy(w, backup)
	if err != nil {
		return fmt.Errorf("error encrpyption file '%s', got error: '%v'", encryptedPath, err)
	}

	// Delete the original backup
	err = os.Remove(d.LocalBackupPath)
	if err != nil {
		return fmt.Errorf("error deleting original backup '%s', got error: '%v'", d.LocalBackupPath, err)
	}

	// Update the path to the compressed path
	d.LocalBackupPath = encryptedPath

	return nil

}
