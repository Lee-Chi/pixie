package model

import "pixie/db/mongo"

type VocabularyFields struct {
	Word mongo.F
}

var Field_Vocabulary VocabularyFields = VocabularyFields{
	Word: mongo.Field("word"),
}

type VocabularyBookmarkFields struct {
	UserId       mongo.F
	VocabularyId mongo.F
}

var Field_VocabularyBookmark VocabularyBookmarkFields = VocabularyBookmarkFields{
	UserId:       mongo.Field("user_id"),
	VocabularyId: mongo.Field("vocabulary_id"),
}

type VocabularyUserFields struct {
	UserId mongo.F
}

var Field_VocabularyUser VocabularyUserFields = VocabularyUserFields{
	UserId: mongo.Field("user_id"),
}
