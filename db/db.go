package db

import (
	"context"
	"pixie/db/mongo"
)

// var firebaseClient *firebase.Client

// func Firebase() *firebase.Client {
// 	return firebaseClient
// }

// func Connect(ctx context.Context, config *firebase.Config, opt ...option.ClientOption) error {
// 	firebaseClient = &firebase.Client{}
// 	if err := firebaseClient.Connect(ctx, config, opt...); err != nil {
// 		return err
// 	}

// 	return nil
// }

var mongoClient *mongo.Client
var pixiedb *mongo.Database

func Build(ctx context.Context, domain string, user string, password string) error {
	client, err := mongo.NewClient(ctx, domain, user, password)
	if err != nil {
		return err
	}

	mongoClient = client
	pixiedb = client.DatabaseUpgradeOnly("pixiedb")

	return nil
}

func Destroy(ctx context.Context) error {
	if err := mongoClient.Destroy(); err != nil {
		return err
	}

	return nil
}

func Pixie() *mongo.Database {
	return pixiedb
}
