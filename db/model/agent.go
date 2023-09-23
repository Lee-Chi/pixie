package model

const CAgent string = "agent"

type Agent struct {
	UserId string `bson:"user_id"`
	Pixie  Pixie  `bson:"pixie"`
}

type Pixie struct {
	Name    string `bson:"name"`
	Payload string `bson:"payload"`
}
