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

type question struct {
	QuestionID int    `json:"questionid"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Module     string `json:"module"`
}
type questioncomments struct {
	CommentID  int    `json:"commentid"`
	QuestionID string `json:"questionid"`
	Comment    string `json:"comment"`
	StudentID  string `json:"studentid"`
}
type questionratings struct {
	RatingID       int    `json:"ratingid"`
	QuestionID     string `json:"answerid"`
	QuestionRating int    `json:"answerrating"`
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Comments and ratings for Questions")
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

func GetComments(db *sql.DB, QuestionID string) questioncomments {
	query := fmt.Sprintf("Select * FROM QuestionComments WHERE QuestionID = '%s'", QuestionID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var comments questioncomments
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&comments.CommentID, &comments.QuestionID,
			&comments.Comment, &comments.StudentID)
		if err != nil {
			panic(err.Error())
		}
	}
	/*
		fmt.Println(&comments.CommentID, &comments.QuestionID,
			&comments.Comment, &comments.StudentID)*/
	return comments
}

func EditComment(db *sql.DB, comment questioncomments) bool {
	query := fmt.Sprintf("UPDATE QuestionComments SET Comments = '%s' WHERE ID = '%d';",
		comment.Comment, comment.CommentID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func CreateComments(db *sql.DB, Comment questioncomments) bool {
	query := fmt.Sprintf("INSERT INTO QuestionComments(QuestionID, Comments, StudentID) VALUES ('%s', '%s', '%s')",
		Comment.QuestionID,
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
		//GET comment using question id
		Comment := GetComments(db, params["QuestionID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Comment.QuestionID))
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

			var newComment questioncomments
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
		// existing comment---
		if r.Method == "PUT" {
			var comment questioncomments
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

func GetRating(db *sql.DB, QuestionID string) questionratings {
	query := fmt.Sprintf("Select * FROM QuestionRatings WHERE QuestionID = '%s'", QuestionID)
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var ratings questionratings
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&ratings.RatingID, &ratings.QuestionID,
			&ratings.QuestionRating)
		if err != nil {
			panic(err.Error())
		}
	}
	//fmt.Println(&ratings.RatingID, &ratings.AnswerID,
	//&ratings.AnsRating)
	return ratings
}

func EditRatings(db *sql.DB, ratings questionratings) bool {
	if ratings.QuestionID == "" {
		return false
	}
	query := fmt.Sprintf("UPDATE QuestionRatings SET Rating = '%d' WHERE QuestionID = '%s';",
		ratings.QuestionRating, ratings.QuestionID)
	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
	return true
}

func IncreaseRatings(db *sql.DB, ratings questionratings) bool {
	QuestionID := ratings.QuestionID
	Rating := (getRatingByID(db, QuestionID) + 1)
	query := fmt.Sprintf("Update QuestionRatings SET Rating = '%d'", Rating)
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
		var rating questionratings
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
func DecreaseRatings(db *sql.DB, ratings questionratings) bool {
	QuestionID := ratings.QuestionID
	Rating := (getRatingByID(db, QuestionID) - 1)
	query := fmt.Sprintf("Update QuestionRatings SET Rating = '%d'", Rating)
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
		var rating questionratings
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
func getRatingByID(db *sql.DB, QuestionID string) int {
	query := fmt.Sprintf("SELECT Rating FROM QuestionRatings WHERE QuestionID = '%s'", QuestionID)
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
		//GET Ratings using Question ID
		Rating := GetRating(db, params["QuestionID"])
		if err != nil {
			fmt.Printf("There was an error encoding the json. err = %s", err)
		} else {
			json.NewEncoder(w).Encode(GetComments(db, Rating.QuestionID))
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
			var rating questionratings
			reqBody, err := ioutil.ReadAll(r.Body)

			if err == nil {
				json.Unmarshal(reqBody, &rating)

				if rating.QuestionID == "" {
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

func GetQuestionsWithComments(db *sql.DB) question {
	query := "Select Questions.*  FROM QuestionComments INNER JOIN Questions ON QuestionComments.QuestionID = Questions.ID GROUP BY QuestionComments.QuestionID"

	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var questions question
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&questions.QuestionID, &questions.Title,
			&questions.Content, &questions.Module)
		if err != nil {
			panic(err.Error())
		}
	}

	return questions
}

func getCommentedQuestions(w http.ResponseWriter, r *http.Request) {
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
		//GET comments using answerid
		json.NewEncoder(w).Encode(GetQuestionsWithComments(db))
		w.WriteHeader(http.StatusAccepted)
		return
	}
}

func getAllRatedQuestion(db *sql.DB) question {
	query := "Select Questions.*  FROM QuestionRatings INNER JOIN Questions ON QuestionRatings.QuestionID = Questions.ID WHERE QuestionRatings.Rating > 0"
	results, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
	var questions question
	for results.Next() {
		// map this type to the record in the table
		err = results.Scan(&questions.QuestionID, &questions.Title,
			&questions.Content, &questions.Module)
		if err != nil {
			panic(err.Error())
		}
	}
	/*
		fmt.Println(&comments.CommentID, &comments.AnswerID,
			&comments.Comment, &comments.StudentID)*/
	return questions
}

func getRatedQuestion(w http.ResponseWriter, r *http.Request) {
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
		//GET comments using answerid
		json.NewEncoder(w).Encode(getAllRatedQuestion(db))
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
	router.HandleFunc("/api/v1/Question", home)                                                         //Test API
	router.HandleFunc("/api/v1/Question/Comments/{QuestionID}", comments).Methods("GET", "PUT", "POST") //API Manipulation
	router.HandleFunc("/api/v1/Question/Ratings/{QuestionID}", ratings).Methods("Get", "PUT")
	router.HandleFunc("/api/v1/Question/Ratings/Increase", IncreaseAnswerRatings).Methods("PUT")
	router.HandleFunc("/api/v1/Question/Ratings/Decrease", DecreaseAnswerRatings).Methods("PUT")
	router.HandleFunc("/api/v1/CommentedQuestion", getCommentedQuestions).Methods("GET")
	router.HandleFunc("/api/v1/RatedQuestion", getRatedQuestion).Methods("GET")
	fmt.Println("Listening at port 9083")
	log.Fatal(http.ListenAndServe(":9083", handlers.CORS(headers, methods, origins)(router)))
}
