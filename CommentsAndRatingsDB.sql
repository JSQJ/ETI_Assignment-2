CREATE database assignment2_db;

USE assignment2_db;

CREATE TABLE IF NOT EXISTS QuestionComments (
 ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
 QuestionID VARCHAR(30)  , 
 Comments VARCHAR(255), 
 StudentID VARCHAR(10) NOT NULL,
 FOREIGN KEY (QuestionID) REFERENCES Questions(ID) ON DELETE CASCADE
 );
 
 CREATE TABLE IF NOT EXISTS QuestionRatings (
 ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
 QuestionID VARCHAR(30)  , 
 Rating int,
 FOREIGN KEY (QuestionID) REFERENCES Questions(ID) ON DELETE CASCADE
 );
 
CREATE TABLE IF NOT EXISTS AnswerComments (
 ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
 AnswerID VARCHAR(30)  , 
 Comments VARCHAR(255), 
 StudentID VARCHAR(10) NOT NULL,
 FOREIGN KEY (AnswerID) REFERENCES Answers(ID) ON DELETE CASCADE
 );
 
 CREATE TABLE IF NOT EXISTS AnswerRatings (
 ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
 AnswerID VARCHAR(30)  , 
 Rating int,
 FOREIGN KEY (AnswerID) REFERENCES Answers(ID) ON DELETE CASCADE
 );
 
delimiter #

create trigger InsertAnswerRatingsTrigger after insert on Answers
for each row
begin
  insert into AnswerRatings (AnswerID, Rating) values (new.ID, 0);
end#

create trigger InsertQuestionsRatingsTrigger after insert on Questions
for each row
begin
  insert into QuestionRatings (QuestionID, Rating) values (new.ID, 0);
end#

delimiter ;


