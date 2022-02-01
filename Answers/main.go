package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type answercomments struct {
	CommentID string `json:"commentid"`
	AnswerID  string `json:"answerid"`
	Comment   string `json:"comment"`
	StudentID string `json:"studentid"`
}
type answerratingss struct {
	AnswerID  string `json:"answerid"`
	AnsRating int    `json:"answerrating"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Comments and ratings for Answers")
}
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func GetComments(db *sql.DB, AnswerID string) answercomments {
	query := fmt.Sprintf("Select * FROM AnswerComments WHERE AnswerID= '%s'", AnswerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var comments answercomments
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&comments.CommentID, &comments.AnswerID,
			&comments.Comment, &comments.StudentID)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println(&comments.CommentID, &comments.AnswerID,
		&comments.Comment, &comments.StudentID)
	return comments
}

func EditComment(db *sql.DB, driver driverinfo) bool {
	if driver.DriverID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE Driver SET FirstName = '%s', LastName = '%s', MobileNumber= '%s', LicenseNumber = '%s', EmailAddress = '%s' WHERE DriverID = '%s';",
		driver.FirstName, driver.LastName, driver.MobileNumber, driver.LicenseNumber, driver.EmailAddress, driver.DriverID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func comments(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment2_db") //connect to database
	if err != nil {
		fmt.Println(err)
	}
	params := mux.Vars(r)
	if r.Method == "DELETE" {
		println("Can't delete comment")
	} else if r.Method == "GET" {
		//GET Driver using email address
		Comment := GetComments(db, params["AnswerID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else if DriverInformation.EmailAddress == "" { // Check if data is empty
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("Driver does not exists"))
			return
		} else {
			json.NewEncoder(w).Encode(GetDriver(db, DriverInformation.EmailAddress))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if err != nil {
			fmt.Println(err)
		}
		// POST is for creating new course
		if r.Method == "POST" {

			var newDriver driverinfo
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newDriver)

				if newDriver.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please enter the required information " +
							"in JSON format"))
					return
				} else {
					// check if driver already exists by email; add only if
					// driver does not exist
					if !CheckDriver(db, newDriver.EmailAddress) {
						CreateDriver(db, newDriver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Driver created successfully"))
						return
					} else {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("Email is already in use"))
						return
					}
				}
			}
		}
		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			//---PUT is for creating or updating
			// existing driver---
			var driver driverinfo
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &driver)

				if driver.EmailAddress == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply driver " +
							" information " +
							"in JSON format"))
					return
				} else {
					// check if Driver does not exists; update only if
					// driver does exist
					if !CheckDriver(db, driver.EmailAddress) {
						w.WriteHeader(http.StatusUnprocessableEntity)
						w.Write([]byte("There is no exsiting driver with " + driver.EmailAddress))
						return
					} else {
						//To update driver details
						EditDriver(db, driver)
						w.WriteHeader(http.StatusCreated)
						w.Write([]byte("Driver updated successfully"))
						return
					}

				}
			}
		}

	}
}
