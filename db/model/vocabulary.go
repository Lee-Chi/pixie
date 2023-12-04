package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CVocabulary string = "vocabulary"

type Vocabulary struct {
	Word        string       `bson:"word"`
	Definitions []Definition `bson:"definitions"`
	Level       int          `bson:"level"`
	Contributor string       `bson:"contributor"`
}

type Definition struct {
	Text         string `bson:"text"`
	PartOfSpeech string `bson:"part_of_speech"`
}

const CVocabularyBookmark string = "vocabulary_bookmark"

type VocabularyBookmark struct {
	UserId       primitive.ObjectID `bson:"user_id"`
	VocabularyId primitive.ObjectID `bson:"vocabulary_id"`
	CreatedAt    time.Time          `bson:"created_at"`
}

const CVocabularyUser string = "vocabulary_user"

type VocabularyUser struct {
	UserId              string             `bson:"user_id"`
	CurrentVocabularyId primitive.ObjectID `bson:"current_vocabulary_id"`
}
