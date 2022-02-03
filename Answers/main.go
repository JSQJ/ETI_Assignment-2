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
type answerratings struct {
	RatingID  string `json:"ratingid"`
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
	//fmt.Println(&comments.CommentID, &comments.AnswerID,
		//&comments.Comment, &comments.StudentID)
	return comments
}

func EditComment(db *sql.DB, comment answercomments) bool {
	if comment.CommentID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE AnswerComments SET Comments = '%s' WHERE ID = '%d';",
		comment.Comment,comment.CommentID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func CreateComments(db *sql.DB, Comment answercomments) bool {
	query := fmt.Sprintf("INSERT INTO AnswerComments(ID, AnswerID, Comments, StudentID) VALUES (%d,'%s', '%s', '%s')",
		getlastid(db)+1,
		Comment.AnswerID,
		Comment.Comment,
		Comment.StudentID)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func getlastid(db *sql.DB) int { //Gets the last id of passengers
	query1 := "SELECT COUNT(*) FROM AnswerComments"
	query2 := "SELECT ID FROM AnswerComments ORDER BY ID DESC LIMIT 1"
	var commentCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		results.Scan(&commentCount)
	}
	if commentCount > 0 {
		results, err := db.Query(query2) //Run Query
		var ID int
		if err != nil {
			panic(err.Error())
		}
		if results.Next() {
			results.Scan(&ID)
		}
		return ID
	} else {
		return 0
	}

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
		//GET Comment using answer id
		Comment := GetComments(db, params["AnswerID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Comment.AnswerID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if err != nil {
			fmt.Println(err)
		}
		// POST is for creating new comment
		if r.Method == "POST" {

			var newComment answercomments
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newComment)

				if newComment.Comment == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please enter the required information " +
							"in JSON format"))
					return
				} else {
					CreateComments(db, newComment)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Comment created successfully"))
					return
				}
			}
		}
		//---PUT is for creating or updating
		// existing comments---
		if r.Method == "PUT" {
			var comment answercomments
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &comment)

				if comment.CommentID == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply driver " +
							" information " +
							"in JSON format"))
					return
				} else {
					//To update comment details
					EditComment(db, comment)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("comment updated successfully"))
					return
				}
			}
		}

	}
}

func GetAnswersWithComments(db *sql.DB, AnswerID string) answercomments {
	query := fmt.Sprintf("Select * FROM AnswerComments WHERE AnswerID= '%s' AND COUNT(", AnswerID)
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

func GetRatings(db *sql.DB, AnswerID string) answerratings {
	query := fmt.Sprintf("Select * FROM AnswerRatings WHERE AnswerID= '%s'", AnswerID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var ratings answerratings
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&ratings.RatingID, &ratings.AnswerID,
			&ratings.AnsRating)
		if err != nil {
			panic(err.Error())
		}
	}
	//fmt.Println(&ratings.RatingID, &ratings.AnswerID,
		//&ratings.AnsRating)
	return ratings
}

func EditRatings(db *sql.DB, ratings answerratings) bool {
	if ratings.AnswerID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE AnswerRatings SET Rating = '%d' WHERE ID = '%s';",
	ratings.AnsRating)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func CreateRatings(db *sql.DB, ratings answerratings) bool {
	query := fmt.Sprintf("INSERT INTO AnswerRatings(ID, AnswerID, Rating) VALUES (%d,'%s',%d)",
		getlastid(db)+1,
		ratings.AnswerID,
		ratings.AnsRating
		)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func getlastid(db *sql.DB) int { //Gets the last id of passengers
	query1 := "SELECT COUNT(*) FROM AnswerComments"
	query2 := "SELECT ID FROM AnswerComments ORDER BY ID DESC LIMIT 1"
	var commentCount int
	results, err := db.Query(query1) //Run Query
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		results.Scan(&commentCount)
	}
	if commentCount > 0 {
		results, err := db.Query(query2) //Run Query
		var ID int
		if err != nil {
			panic(err.Error())
		}
		if results.Next() {
			results.Scan(&ID)
		}
		return ID
	} else {
		return 0
	}

}

func ratings(w http.ResponseWriter, r *http.Request) {
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
		//GET Ratings using Answer ID
		Comment := GetComments(db, params["AnswerID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Comment.AnswerID))
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

			var newComment answercomments
			reqBody, err := ioutil.ReadAll(r.Body)
			if err == nil {
				// convert JSON to object
				json.Unmarshal(reqBody, &newComment)

				if newComment.Comment == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please enter the required information " +
							"in JSON format"))
					return
				} else {
					// check if driver already exists by email; add only if
					// driver does not exist
					CreateComments(db, newComment)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("Comment created successfully"))
					return
				}
			}
		}
		//---PUT is for creating or updating
		// existing course---
		if r.Method == "PUT" {
			//---PUT is for creating or updating
			// existing driver---
			var comment answercomments
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &comment)

				if comment.CommentID == "" {
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
					//To update driver details
					EditComment(db, comment)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("comment updated successfully"))
					return
				}
			}
		}

	}
}
func getCommentedAnswers(w http.ResponseWriter, r *http.Request) {
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
	//GET is for getting comments
	if r.Method == "GET" {
		//GET comments using answerid
		Comment := GetComments(db, params["AnswerID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Comment.AnswerID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
}

func getRatedAnswers(w http.ResponseWriter, r *http.Request) {
	//GET is for getting passenger trips
	if r.Method == "GET" {
		params := mux.Vars(r)
		passengerID := params["passengerid"]
		var passengerTrips []Trip = getTrips(passengerID)
		json.NewEncoder(w).Encode(passengerTrips)
		fmt.Println(passengerTrips)
	}
}