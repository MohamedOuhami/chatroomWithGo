package utils

import (
	"ChatroomWithGo/models"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq" // Import PostgreSQL driver
)

// ConnectToDb connects to the database
func ConnectToDb() (*sql.DB, error) {
	host := "127.0.0.1"
	port := "5432"
	user := "postgres"
	password := "v01d"
	dbname := "chatroomdb"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Open a connection to the database
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Test if the connection is alive
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertNewUser(db *sql.DB, newUser models.UserModel) error {
	userUsername := newUser.Username
	userPassword := newUser.Password

	sqlStatement := `INSERT INTO chatroomUsers (username,password) values ($1,$2)`

	_, err := db.Exec(sqlStatement, userUsername, userPassword)

	if err != nil {
		return err
	}

	return nil

}
func ConnectWithUsername(db *sql.DB, username string, password string) error {
	var dbUsername, dbPassword string
	sqlStatement := `SELECT username, password FROM chatroomUsers WHERE username=$1`

	err := db.QueryRow(sqlStatement, username).Scan(&dbUsername, &dbPassword)


	// Handle query error
	if err != nil {
		return err
	}

	// Trim spaces from the password fields before comparing
	password = strings.TrimSpace(password)
	dbPassword = strings.TrimSpace(dbPassword)

	// Compare the passwords
	if password != dbPassword {
		return fmt.Errorf("Password incorrect")
	}

	// If the password matches
	return nil
}
