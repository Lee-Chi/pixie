package main

import (
	"context"
	"pixie/core/user"
	"pixie/db"
	"pixie/db/model"
	"pixie/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func upgrade() {
	ctx := context.Background()

	if err := db.Pixie().Collection(model.CUserMember).CreateIndex(
		ctx,
		mongo.Index{
			Keys: bson.D{
				model.Field_UserMember.Account.Index(1),
			},
			Options: options.Index().SetUnique(true),
		},
	); err != nil {
		panic(err)
	}

	if _, err := user.CreateMember(ctx, "leo", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"); err != nil {
		panic(err)
	}
}
