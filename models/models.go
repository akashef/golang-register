package models

import (
	"database/sql"
	"encoding/base64"
	"log"
	"net/mail"
	"os"
	"regexp"

	"github.com/google/uuid"
)

// User schema of the user table
type User struct {
	Db                  *sql.DB
	Username            string `json:"username"`
	Useremail           string `json:"useremail"`
	Userphone           string `json:"userphone"`
	Useraddress         string `json:"useraddress"`
	Authenticationtype  int64  `json:"authenticationtype"`
	Authenticationvalue string `json:"authenticationvalue"`
	UserNationalId      string `json:"usernationalid"`
	PassportImage       string `json:"passportImage"`
	FilePath            string
}

func (user *User) Validate() []string {
	var strArray []string
	if len(user.Username) < 5 {
		strArray = append(strArray, "Please Enter Valid username, more than 5 charachter")
	}
	if len(user.Useremail) == 0 && len(user.Userphone) == 0 {
		strArray = append(strArray, "Please Enter Email or Mobile Number")
	}
	if 1 > user.Authenticationtype && user.Authenticationtype > 3 {
		strArray = append(strArray, "Please Enter Valid Authenticationtype")
		strArray = append(strArray, "1 for finger print")
		strArray = append(strArray, "2 for eye-detection")
		strArray = append(strArray, "3 for password")
	}

	if len(user.Useremail) > 0 {
		_, err := mail.ParseAddress(user.Useremail)
		if err != nil {
			strArray = append(strArray, "Please Enter Valid Email")
		}
	}
	//save PassportImage
	if user.PassportImage != "" {
		dec, err := base64.StdEncoding.DecodeString(user.PassportImage)
		if err != nil {
			panic(err)
		}

		//filename
		fileName := uuid.New()
		fileNameString := fileName.String()

		re, err := regexp.Compile(`[^\w]`)
		if err != nil {
			log.Fatal(err)
		}
		fileNameString = re.ReplaceAllString(fileNameString, "")

		f, err := os.Create("images/" + fileNameString + ".jpg")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
		user.FilePath = fileNameString + ".jpg"
	}

	return strArray
}

func (user *User) CheckExistValues() []string {
	var returnReutls []string
	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE username=$1`
	// execute the sql statement
	row, _ := user.Db.Query(sqlStatement, user.Username)
	if row.Next() {
		returnReutls = append(returnReutls, "username already exist")
	}

	// create the select sql query
	sqlStatement = `SELECT * FROM users WHERE useremail=$1`
	// execute the sql statement
	row, _ = user.Db.Query(sqlStatement, user.Useremail)
	if row.Next() {
		returnReutls = append(returnReutls, "Useremail already exist")
	}

	// create the select sql query
	sqlStatement = `SELECT * FROM users WHERE userphone=$1`
	// execute the sql statement
	row, _ = user.Db.Query(sqlStatement, user.Userphone)
	if row.Next() {
		returnReutls = append(returnReutls, "Userphone already exist")
	}
	// return empty user on error
	return returnReutls
}

func (user *User) InsertUser() int64 {

	sqlStatement := `INSERT INTO users 
		(
		username,
		useremail,
		userphone,
		useraddress,
		authenticationtype,
		authenticationvalue,
		usernationalid,
		filepath
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING userid`

	// the inserted id will store in this id
	var id int64

	// execute the sql statement
	// Scan function will save the insert id in the id
	err := user.Db.QueryRow(sqlStatement,
		user.Username,
		user.Useremail,
		user.Userphone,
		user.Useraddress,
		user.Authenticationtype,
		user.Authenticationvalue,
		user.UserNationalId,
		user.FilePath).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	return id
}
