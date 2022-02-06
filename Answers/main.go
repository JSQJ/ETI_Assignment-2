package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type answer struct {
	AnswerID   string `json:"answerid"`
	QuestionID string `json:"questionid"`
	Content    string `json:"content"`
}

type answercomments struct {
	CommentID int    `json:"commentid"`
	AnswerID  string `json:"answerid"`
	Comment   string `json:"comment"`
	StudentID string `json:"studentid"`
}
type answerratings struct {
	RatingID  int    `json:"ratingid"`
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
	return comments
}

func EditComment(db *sql.DB, comment answercomments) bool {
	query := fmt.Sprintf("Update AnswerComments SET Comments = '%s' WHERE ID = '%d' ",
		comment.Comment, comment.CommentID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func CreateComments(db *sql.DB, Comment answercomments) bool {
	query := fmt.Sprintf("INSERT INTO AnswerComments(AnswerID, Comments, StudentID) VALUES ('%s', '%s', '%s')",
		Comment.AnswerID,
		Comment.Comment,
		Comment.StudentID)

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
				//To update comment details
				EditComment(db, comment)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("comment updated successfully"))
				return
			}
		}
	}
}

func GetRating(db *sql.DB, AnswerID string) answerratings {
	query := fmt.Sprintf("Select * FROM AnswerRatings WHERE AnswerID = '%s'", AnswerID)
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
	return ratings
}

func EditRatings(db *sql.DB, ratings answerratings) bool {
	if ratings.AnswerID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE AnswerRatings SET Rating = '%d' WHERE AnswerID = '%s';",
		ratings.AnsRating, ratings.AnswerID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func IncreaseRatings(db *sql.DB, ratings answerratings) bool {
	AnswerID := ratings.AnswerID
	Rating := (getRatingByID(db, AnswerID) + 1)
	query := fmt.Sprintf("Update AnswerRatings SET Rating = '%d'", Rating)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}
func IncreaseAnswerRatings(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment2_db") //connect to database
	if err != nil {
		fmt.Println(err)
	}
	//PUT is for getting increasing rating
	if r.Method == "PUT" {
		var rating answerratings
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &rating)

			//To update rating
			IncreaseRatings(db, rating)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("comment updated successfully"))
			return
		}
	}
}
func DecreaseRatings(db *sql.DB, ratings answerratings) bool {
	AnswerID := ratings.AnswerID
	Rating := (getRatingByID(db, AnswerID) - 1)
	query := fmt.Sprintf("Update AnswerRatings SET Rating = '%d'", Rating)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}
func DecreaseAnswerRatings(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment2_db") //connect to database
	if err != nil {
		fmt.Println(err)
	}
	//PUT is for decreasing rating
	if r.Method == "PUT" {
		var rating answerratings
		reqBody, err := ioutil.ReadAll(r.Body)

		if err == nil {
			json.Unmarshal(reqBody, &rating)

			//To update rating
			DecreaseRatings(db, rating)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("comment updated successfully"))
			return
		}
	}
}
func getRatingByID(db *sql.DB, AnswerID string) int {
	query := fmt.Sprintf("SELECT Rating FROM AnswerRatings WHERE AnswerID= '%s'", AnswerID)
	results, err := db.Query(query) //Run Query
	var ID int
	if err != nil {
		panic(err.Error())
	}
	if results.Next() {
		results.Scan(&ID)
	}
	return ID
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
		Rating := GetRating(db, params["AnswerID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Rating.AnswerID))
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if r.Header.Get("Content-type") == "application/json" {
		if err != nil {
			fmt.Println(err)
		}
		//---PUT is for creating or updating
		// existing rating---
		if r.Method == "PUT" {
			var rating answerratings
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &rating)

				if rating.AnswerID == "" {
					w.WriteHeader(
						http.StatusUnprocessableEntity)
					w.Write([]byte(
						"422 - Please supply answer " +
							" information " +
							"in JSON format"))
					return
				} else {
					EditRatings(db, rating)
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("comment updated successfully"))
					return
				}
			}
		}

	}
}

func GetAnswersWithComments(db *sql.DB) answer {
	query := "Select Answers.*  FROM AnswerComments INNER JOIN Answers ON AnswerComments.AnswerID = Answers.ID GROUP BY AnswerComments.AnswerID"

	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var answers answer
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&answers.AnswerID, &answers.QuestionID,
			&answers.Content)
		if err != nil {
			panic(err.Error())
		}
	}

	return answers
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
	//GET is for getting comments
	if r.Method == "GET" {
		//GET commented questions
		json.NewEncoder(w).Encode(GetAnswersWithComments(db))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func getAllRatedAnswers(db *sql.DB) answer {
	query := "Select Answers.*  FROM AnswerRatings INNER JOIN Answers ON AnswerRatings.AnswerID = Answers.ID WHERE AnswerRatings.Rating > 0"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var answers answer
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&answers.AnswerID, &answers.QuestionID,
			&answers.Content)
		if err != nil {
			panic(err.Error())
		}
	}
	return answers
}

func getRatedAnswers(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Invalid key"))
		return
	}
	db, err := sql.Open("mysql", "root:12N28c02@tcp(127.0.0.1:3306)/assignment2_db") //connect to database
	if err != nil {
		fmt.Println(err)
	}
	//GET is for getting passenger trips
	if r.Method == "GET" {
		//GET rated questions
		json.NewEncoder(w).Encode(getAllRatedAnswers(db))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func main() {
	//API part
	router := mux.NewRouter()
	//Web Front-end CORS
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/api/v1/Answer", home)                                                       //Test API
	router.HandleFunc("/api/v1/Answer/Comments/{AnswerID}", comments).Methods("GET", "PUT", "POST") //API Manipulation
	router.HandleFunc("/api/v1/Answer/Ratings/{AnswerID}", ratings).Methods("Get", "PUT")
	router.HandleFunc("/api/v1/Answer/Ratings/Increase", IncreaseAnswerRatings).Methods("PUT")
	router.HandleFunc("/api/v1/Answer/Ratings/Decrease", DecreaseAnswerRatings).Methods("PUT")
	router.HandleFunc("/api/v1/CommentedAnswers", getCommentedAnswers).Methods("GET")
	router.HandleFunc("/api/v1/RatedAnswers", getRatedAnswers).Methods("GET")
	fmt.Println("Listening at port 9082")
	log.Fatal(http.ListenAndServe(":9082", handlers.CORS(headers, methods, origins)(router)))
}
