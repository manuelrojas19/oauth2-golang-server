package main

import (
	"log"

	"github.com/manuelrojas19/go-oauth2-server/database"
	"github.com/manuelrojas19/go-oauth2-server/oauth"
	"github.com/manuelrojas19/go-oauth2-server/store/dao"
)

func main() {
	db, _ := database.InitDatabaseConnection()

	oauthClient := &oauth.OauthClient{
		OauthClientDao: &dao.OauthClientDao{
			Db: db,
		},
	}

	oauthToken := &oauth.OauthToken{
		OauthClient: oauthClient,
		OauthTokenDao: &dao.OauthTokenDao{
			Db: db,
		},
	}

	client, error := oauthClient.FindOauthClient("client149")

	log.Println(client)

	if error != nil {
		log.Fatal(error)
	}

	accessToken, error := oauthToken.GrantAccessToken(client.Key)
	log.Println(accessToken)
	if error != nil {
		log.Fatal(error)
		return
	}
}
