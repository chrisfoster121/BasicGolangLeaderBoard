package handler

import (
	"BasicGolangLeaderBoard/internal/auth"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
)

type DynamodbHelper struct {
	TableName string
	Svc       *dynamodb.DynamoDB
}

type UserScoreAWS struct {
	UserId string `dynamodbav:"UserId"`
	Score  string `dynamodbav:"Score"`
}
type LoginCredentials struct {
	username string
	password string
}
type Username struct {
	Username string
}

type HandlerHelper struct {
	DynamodbHelper DynamodbHelper
	AuthHelper     auth.AuthHelper
}

func (h HandlerHelper) GetTopThreeUserScores(c *gin.Context) {
	token := strings.TrimPrefix(c.Request.Header["Authorization"][0], "Bearer ")
	err := h.AuthHelper.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.DynamodbHelper.Svc.Scan(&dynamodb.ScanInput{
		TableName: &h.DynamodbHelper.TableName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user_scores []UserScoreAWS
	for _, element := range out.Items {
		user_score := &UserScoreAWS{}
		err = dynamodbattribute.UnmarshalMap(element, &user_score)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_scores = append(user_scores, *user_score)
	}

	c.JSON(200, user_scores)
}

func (h HandlerHelper) PostNewScore(c *gin.Context) {
	token := strings.TrimPrefix(c.Request.Header["Authorization"][0], "Bearer ")
	err := h.AuthHelper.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var response map[string]string
	err = c.BindJSON(&response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	av, err := dynamodbattribute.MarshalMap(response)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(av)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(h.DynamodbHelper.TableName),
	}

	_, err = h.DynamodbHelper.Svc.PutItem(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h HandlerHelper) CheckUsernameAvailability(c *gin.Context) {
	token := strings.TrimPrefix(c.Request.Header["Authorization"][0], "Bearer ")
	err := h.AuthHelper.VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var body Username
	err = c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	out, err := h.DynamodbHelper.Svc.Scan(&dynamodb.ScanInput{
		TableName: &h.DynamodbHelper.TableName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user_scores []UserScoreAWS
	for _, element := range out.Items {
		user_score := &UserScoreAWS{}
		err = dynamodbattribute.UnmarshalMap(element, &user_score)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user_scores = append(user_scores, *user_score)
	}

	for _, user := range user_scores {
		if body.Username == user.UserId {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Username taken"})
		}
	}

	c.JSON(http.StatusOK, gin.H{"error": "Username available"})
}

func (h HandlerHelper) Auth(c *gin.Context) {
	token := strings.TrimPrefix(c.Request.Header["Authorization"][0], "Bearer ")
	err := h.AuthHelper.VerifyAdminToken(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var creds LoginCredentials
	err = c.BindJSON(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err = h.AuthHelper.CreateToken(creds.username)
	if err != nil {
		log.Fatalf("Got error creating token: %s", err)
	}
	auth_token := &auth.AuthToken{
		Token: token,
	}

	c.JSON(200, auth_token)
}
