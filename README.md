## PG CONNECTION

This is a library designed to simplify the management of database connections using Gorm in PostgreSQL. It focuses on the use of contexts to handle transactions, thereby simplifying the Gorm handling.

Esta es una librería diseñada para facilitar la gestión de conexiones a bases de datos mediante Gorm en PostgreSQL. Se centra en el uso de contextos para administrar las transacciones, simplificando así el manejo de Gorm.

## Open Connection

```go
    package main

import (
	"github.com/user0608/pg-connection"
)

func main() {
	    // Open a global database connection using the provided configuration parameters.
        _, err := connection.OpenWithConfigs(connection.DBConfigParams{
            DBHost:     "192.16.0.15",
            DBPort:     "5432",
            DBName:     "",
            DBUsername: "",
            DBPassword: "",
            DBLogLevel: "info",
        })

        if err != nil {
            // Handle any error that occurred during database connection setup.           
            panic(err)
        }

        // Now the global database connection is open and can be used throughout the application.
        // It's important to handle errors appropriately, such as logging or terminating the application,
        // depending on your application's requirements.
    }

```

## Get Connection

```go
    // This is how we obtain a transaction using a context.
    tx := connection.Conn(context.Background())
```

## Transaction

```go
        // WithTx initiates a transaction using the provided context.
        // The transaction is propagated through the context, meaning that all database operations
        // using the same context will be within the transaction.
        connection.WithTx(context.Background(), func(ctx context.Context) error {
            // Retrieves the transaction from the context.
            tx := connection.Conn(ctx)
            // 'tx' is of type *gorm.DB, representing the database transaction.

            // Here, you can perform operations within the transaction, such as queries,
            // updates, or inserts into the database.

            // To apply changes to the database, simply return nil.
            // If an error occurs during operations within the transaction,
            // you can return that error, and the transaction will be automatically rolled back.

            // For example:
            // if err := tx.Model(&Model{}).Create(&data).Error; err != nil {
            //     return err
            // }

            // If everything goes well, the transaction will be committed automatically when
            // the anonymous function completes without errors.
            return nil
        })
```

## Environment variables
```go
package main

import (
	"log"
	"github.com/user0608/pg-connection"
)

func main() {
	// If using environment variables, this function will use the following variables
	// to open the PostgreSQL connection: PG_HOST, PG_PORT, PG_USER, PG_PASSWORD, PG_DATABASE.
	_, err := connection.Open()
	if err != nil {
		// Handle any error that occurred during database connection setup.
		// For example, log the error or terminate the application.
		log.Fatalln(err)
	}

	// Now the database connection is open and can be used throughout the application.
	// It's important to handle errors appropriately, such as logging or terminating the application,
	// depending on your application's requirements.
}

```


## Manager
```go
    // It is also possible to use a manager to handle the database connection.
    // This is useful for dependency injection and allows managing multiple database connections.
    func main() {
        // Create a new connection manager instance.
        manager, err := connection.NewConnection(connection.DBConfigParams{})
        if err != nil {
            // Handle any error that occurred during connection manager creation.
            // For example, log the error or terminate the application.
            log.Fatalln(err)
        }

        // Obtain a connection from the manager.
        // This can be useful for regular database operations.
        regularConn := manager.Conn(context.Background())

        // Use the connection for non-transactional database operations.

        // Example: Fetch data
        // data, err := regularConn.Model(&Model{}).Find(&result).Error
        // Handle error if needed.

        // Use the manager to perform operations within a transaction.
        manager.WithTx(context.Background(), func(ctx context.Context) error {
            // Obtain a transactional connection from the manager.
            tx := manager.Conn(ctx)

            // Perform transactional operations using 'tx'.
            // Example: tx.Model(&Model{}).Create(&data)

            // Return an error if the transaction should be rolled back, or nil to commit the transaction.

            return nil
        })
    }
```