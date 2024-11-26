package main

import (
	"BasicGolangLeaderBoard/internal/auth"
	"BasicGolangLeaderBoard/internal/handler"
	"flag"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

type UserScore struct {
	UserId string `dynamodbav:"UserId"`
	Score  int    `dynamodbav:"Score"`
}

type TopScoreList struct {
}

func main() {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	table_name := flag.String("table-name", "foo", "table-name")
	secretKey := flag.String("secret-key", "foo", "secret-key")
	adminSecretKey := flag.String("admin-key", "foo", "admin-key")
	flag.Parse()

	log.Println(*table_name)
	log.Println(*secretKey)
	log.Println(*adminSecretKey)
	handler := &handler.HandlerHelper{
		DynamodbHelper: handler.DynamodbHelper{
			TableName: *table_name,
			Svc:       dynamodb.New(sess),
		},
		AuthHelper: auth.CreateAuthHelper(*secretKey, *adminSecretKey),
	}

	router := gin.Default()
	router.POST("/PostNewScore", handler.GetTopThreeUserScores)
	router.GET("/GetTopScores", handler.GetTopThreeUserScores)
	router.POST("/CheckUsernameAvailability", handler.CheckUsernameAvailability)
	router.POST("/Auth", handler.Auth)

	router.Run("localhost:50051")
}
