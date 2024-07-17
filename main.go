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

	_, error := oauthClient.CreateOauthClient("client", "secret", "uri")
	client, _ := oauthClient.FindOauthClient("client2")
	log.Println(client)
	if error != nil {
		log.Fatal(error)
	}
}
