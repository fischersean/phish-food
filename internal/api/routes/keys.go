package routes

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigateway"

	"github.com/fischersean/phish-food/internal/api"
	db "github.com/fischersean/phish-food/internal/database"

	"fmt"
	"net/http"
)

func CreateApiKeyAndAddToPlan(conn db.Connection, plan string, name string) (key *apigateway.UsagePlanKey, err error) {

	svc := apigateway.New(conn.Session)
	baseKey, err := svc.CreateApiKey(
		&apigateway.CreateApiKeyInput{
			Description: aws.String("Testing out creating a key"),
			Enabled:     aws.Bool(true),
			Name:        aws.String(name),
		},
	)
	if err != nil {
		return key, err
	}

	key, err = svc.CreateUsagePlanKey(
		&apigateway.CreateUsagePlanKeyInput{
			UsagePlanId: aws.String(plan),
			KeyId:       baseKey.Id,
			KeyType:     aws.String("API_KEY"),
		},
	)

	return key, err
}

func HandleGetKey(w http.ResponseWriter, r *http.Request) {

	tokenHeader := r.Header["Authorization"]
	token, err := api.ParseAndValidateJWT(tokenHeader[0])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	conn := db.SharedConnection

	tokenUsername, ok := token.Get("username")
	if !ok {
		http.Error(w, "No username in JWT", 500)
		return
	}

	// Check to see if user already has a key
	user, err := conn.GetUserRecord(db.UserQueryInput{
		Username: fmt.Sprintf("%s", tokenUsername),
	})
	if err != nil {
		http.Error(w, "Could not query user table", 500)
	}
	if user.ApiKey != "" {
		w.Write([]byte(user.ApiKey))
		return
	}

	// if no key exists, issue a new one and add it to the usage plan
	// TODO: This needs to be an env variable
	key, err := CreateApiKeyAndAddToPlan(conn, "fwzevm", fmt.Sprintf("%s", tokenUsername))
	if err != nil {
		http.Error(w, "Could not issue and assign key", 500)
		return
	}

	user, err = conn.UpdateUserRecord(db.UserUpdateInput{
		Username: fmt.Sprintf("%s", tokenUsername),
		NewUserRecord: db.UserRecord{
			ApiKey:        *key.Value,
			ApiKeyEnabled: true,
		},
	})
	if err != nil {
		http.Error(w, "Could not update user data", 500)
		return
	}

	w.Write([]byte(user.ApiKey))
}
