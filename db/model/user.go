package model

import "go.mongodb.org/mongo-driver/bson/primitive"

const CUserMember string = "user_member"

type UserMember struct {
	Id     primitive.ObjectID `bson:"_id"`
	Member `bson:"-,inline"`
}

type Member struct {
	Account  string `bson:"account"`
	Password string `bson:"password"`
}

const CUserSession string = "user_session"

type UserSession struct {
	Id      primitive.ObjectID `bson:"_id"`
	Session `bson:"-,inline"`
}

type Session struct {
	UserId primitive.ObjectID `bson:"user_id"`
	Device string             `bson:"device"`
}
