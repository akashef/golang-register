package middleware

import (
	// package to encode and decode the json into struct and vice versa
	"database/sql"
	"fmt"
	"go-postgres/models" // models package where User schema is defined

	// used to access the request and response object of the api
	// used to read the environment variable
	// package used to covert string into int type
	// used to get the params from the route

	// package used to read the .env file
	_ "github.com/lib/pq" // postgres golang driver
)

// create connection with postgres db

// CreateUser create a user in the postgres db

func InsertUser(user models.User, db *sql.DB) (int64, []string) {

	user.Db = db
	// close the db connection
	defer db.Close()

	var strError = user.Validate()
	if len(strError) > 0 {
		return 0, strError
	}

	returnResult := user.CheckExistValues()
	if len(returnResult) > 0 {
		return 0, returnResult
	}

	id := user.InsertUser()
	//sendMail(user)   check email account before un-comment it
	return id, []string{"User created successfully"}
}

func sendMail(user models.User) {
	sender := NewMail()
	body := fmt.Sprintf("New Accout: Name: %s, Email: %s, phone: %s", user.Username, user.Useremail, user.Userphone)
	m := NewMessage("New Regestration Account", body)
	m.To = []string{"info@futiracoin.com"}
	m.AttachFile("/images/" + user.FilePath)
	fmt.Println(sender.Send(m))
}
