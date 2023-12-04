package model

import "pixie/db/mongo"

type UserMemberFields struct {
	Id       mongo.F
	Account  mongo.F
	Password mongo.F
}

var Field_UserMember UserMemberFields = UserMemberFields{
	Id:       mongo.Field_ID,
	Account:  mongo.Field("account"),
	Password: mongo.Field("password"),
}

type UserSessionFields struct {
	Id     mongo.F
	UserId mongo.F
	Device mongo.F
}

var Field_UserSession UserSessionFields = UserSessionFields{
	Id:     mongo.Field_ID,
	UserId: mongo.Field("user_id"),
	Device: mongo.Field("device"),
}
