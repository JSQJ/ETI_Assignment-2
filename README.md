# ETI_Assignment-2
![Copy of Untitled Diagram drawio](https://user-images.githubusercontent.com/78250532/152691013-392e6289-a9d4-405a-aff6-4a25ced2bf98.png)

There are a total of 2 microservices required for this package:
* Question Comments and Ratings Microservice
* Answer Comments and Ratings Microservice

There are a total of 2 databases directly involved in this package:
* AnswerCommentsAndRatings
* QuestionCommentsAndRatings

Design Considerations:
I designed the architecture to be independent and loosely coupled. For the requirements of my package, the two main microservices are question comments and ratings and answer comments and ratings are dependent on the other microservices for questions and answers. 

The two main MS, question comments and ratings and answer comments and ratings have their respective databases.To Support this feature, the microservices depend on the question and answers microservices to get data and is called by the the UI for questions and answers. 

The questions comment and ratings microservice contains the following features:
   * Viewing Comments / Ratings for questions
   * Create Comments / Ratings for quetions
   * Updating Comments / Ratings for quetions
   * View all questions with Comments / Ratings

The answers comment and ratings microservice contains the following features:
   * Viewing Comments / Ratings for answers
   * Create Comments / Ratings for answers
   * Updating Comments / Ratings for answers
   * View all answers with Comments / Ratings
   
DockerHub images link
   * https://hub.docker.com/repository/docker/jsqj/question_repo
   * https://hub.docker.com/repository/docker/jsqj/answer_repo
  
Prerequisites
   * Please ensure that GOLANG and MYSQL is installed on your system, and is fully operational
   * Please do also ensure that your SQL user login is as such:
   * Username: root
   * Password: 12N28c02
  
  
