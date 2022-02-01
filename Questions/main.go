package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type questioncomments struct {
	CommentID string `json:"commentid"`
	AnswerID string `json:"answerid"`
	QuestionID  string `json:"questionid"`
	Comment   string `json:"comment"`
	StudentID string `json:"studentid"`
}
type questionratingss struct {
	AnswerID  string `json:"answerid"`
	AnsRating int    `json:"answerrating"`
}
