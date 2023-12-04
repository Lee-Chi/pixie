package user

import (
	"context"
	"fmt"
	"pixie/db"
	"pixie/db/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMember(ctx context.Context, account string, password string) (*Member, error) {
	md := model.Member{
		Account:  account,
		Password: password,
	}

	id, err := db.Pixie().Collection(model.CUserMember).InsertOne(
		ctx,
		md,
	)
	if err != nil {
		return nil, fmt.Errorf("db.insert(), %v", err)
	}

	return &Member{
		UserMember: model.UserMember{
			Id:     id,
			Member: md,
		},
	}, nil
}

func FindOneMember(ctx context.Context, account string) (*Member, error) {
	md := model.UserMember{}
	if err := db.Pixie().Collection(model.CUserMember).FindOne(
		ctx,
		model.Field_UserMember.Account.Equal(account),
		&md,
	); err != nil {
		return nil, err
	}

	return &Member{
		UserMember: md,
	}, nil
}

type Member struct {
	model.UserMember
}

func (m Member) Login(ctx context.Context, password string, device string) (string, error) {
	if m.Password != password {
		return "", fmt.Errorf("wrong password")
	}

	if err := db.Pixie().Collection(model.CUserSession).DeleteMany(
		ctx,
		model.Field_UserSession.UserId.Equal(m.Id),
	); err != nil {
		return "", fmt.Errorf("db.user_session.deleteMany(), %v", err)
	}

	session := model.Session{
		UserId: m.Id,
		Device: device,
	}

	sessionId, err := db.Pixie().Collection(model.CUserSession).InsertOne(
		ctx,
		session,
	)
	if err != nil {
		return "", fmt.Errorf("db.user_session.insert(), %v", err)
	}

	return sessionId.Hex(), nil
}

type UserSession struct {
	model.UserSession
}

func GetSession(ctx context.Context, userId primitive.ObjectID) (*UserSession, error) {
	md := model.UserSession{}

	if err := db.Pixie().Collection(model.CUserSession).FindOne(
		ctx,
		model.Field_UserSession.UserId.Equal(userId),
		&md,
	); err != nil {
		return nil, fmt.Errorf("db.user_session.find(), %v", err)
	}

	return &UserSession{
		UserSession: md,
	}, nil
}
