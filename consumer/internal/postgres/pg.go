package postgres

import (
	"database/sql"
	"fmt"
	"os"

	// DB adapter
	_ "github.com/lib/pq"
	"github.com/wednesday-solutions/go-template-consumer/pkg/utl/db"
)

// Connect ...
func Connect() (*sql.DB, error) {
	// unfurl the database connection environment variables
	dbDetails, err := db.GetDbDetails()
	if err != nil {
		fmt.Print("error loading database environment variables")

		os.Exit(1)
	}
	psqlInfo := fmt.Sprintf("dbname=%s host=%s user=%s password=%s port=%s sslmode=%s",
		dbDetails.Dbname,
		dbDetails.Host,
		dbDetails.Username,
		dbDetails.Password,
		fmt.Sprintf("%d", dbDetails.Port),
		"disable")
	fmt.Println("Connecting to DB\n", psqlInfo)
	return sql.Open("postgres", psqlInfo)
}
