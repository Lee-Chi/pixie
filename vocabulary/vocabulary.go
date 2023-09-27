package vocabulary

import (
	"context"
	"encoding/json"
	"fmt"
	"pixie/db"
	"pixie/db/model"
	"pixie/db/mongo"
	"pixie/utils/http"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Field_UserId       mongo.F = mongo.Field("user_id")
	Field_Word         mongo.F = mongo.Field("word")
	Field_VocabularyId mongo.F = mongo.Field("vocabulary_id")
)

type Definition struct {
	Text         string `json:"text"`
	PartOfSpeech string `json:"part_of_speech"`
}
type Vocabulary struct {
	Id          primitive.ObjectID `json:"id"`
	Word        string             `json:"word"`
	Definitions []Definition       `json:"definitions"`
}

func (v Vocabulary) Marshal() string {
	lines := []string{
		fmt.Sprintf("#%s", v.Id.Hex()),
		fmt.Sprintf("|- %s", v.Word),
	}
	for _, def := range v.Definitions {
		lines = append(lines, fmt.Sprintf("|- %s, %s", def.PartOfSpeech, def.Text))
	}

	return strings.Join(lines, "\n")
}

func FromModel(id primitive.ObjectID, md model.Vocabulary) Vocabulary {
	vocabulary := Vocabulary{
		Id:          id,
		Word:        md.Word,
		Definitions: []Definition{},
	}
	for _, def := range md.Definitions {
		vocabulary.Definitions = append(vocabulary.Definitions, Definition{
			Text:         def.Text,
			PartOfSpeech: def.PartOfSpeech,
		})
	}

	return vocabulary
}

func Commit(ctx context.Context, userId string, vocabulary Vocabulary) error {
	var md model.Vocabulary
	if err := db.Pixie().Collection(model.CVocabulary).FindOneOrZero(
		ctx,
		Field_Word.Equal(vocabulary.Word),
		&md,
	); err != nil {
		return err
	}

	if md.Word == "" {
		md.Word = vocabulary.Word
		md.Contributor = userId
		md.Definitions = []model.Definition{}
		for _, def := range vocabulary.Definitions {
			md.Definitions = append(md.Definitions, model.Definition{
				Text:         def.Text,
				PartOfSpeech: def.PartOfSpeech,
			})
		}

		if err := db.Pixie().Collection(model.CVocabulary).Insert(
			ctx,
			md,
		); err != nil {
			return err
		}

		return nil
	}

	waitAppends := []model.Definition{}
	for _, commitDef := range vocabulary.Definitions {
		needAppend := true
		for _, mdDef := range md.Definitions {
			if commitDef.PartOfSpeech == mdDef.PartOfSpeech {
				needAppend = false
			}
		}

		if !needAppend {
			continue
		}

		waitAppends = append(waitAppends, model.Definition{
			Text:         commitDef.Text,
			PartOfSpeech: commitDef.PartOfSpeech,
		})
	}

	if len(waitAppends) > 0 {
		md.Contributor = strings.Join([]string{
			md.Contributor,
			userId,
		}, ",")
		md.Definitions = append(md.Definitions, waitAppends...)
	}

	if err := db.Pixie().Collection(model.CVocabulary).Update(
		ctx,
		Field_Word.Equal(vocabulary.Word),
		md,
	); err != nil {
		return err
	}

	return nil
}

func At(ctx context.Context, userId string, vocabularyId primitive.ObjectID) (Vocabulary, error) {
	if vocabularyId.IsZero() {
		user := model.VocabularyUser{}
		if err := db.Pixie().Collection(model.CVocabularyUser).FindOneOrZero(
			ctx,
			Field_UserId.Equal(userId),
			&user,
		); err != nil {
			return Vocabulary{}, err
		}

		vocabularyId = user.CurrentVocabularyId
	}

	voc := struct {
		Id               primitive.ObjectID `bson:"_id"`
		model.Vocabulary `bson:"-,inline"`
	}{}
	if err := db.Pixie().Collection(model.CVocabulary).FindOneByID(
		ctx,
		vocabularyId,
		&voc,
	); err != nil {
		return Vocabulary{}, err
	}

	vocabulary := FromModel(voc.Id, voc.Vocabulary)

	return vocabulary, nil
}

func Next(ctx context.Context, userId string) (Vocabulary, error) {
	user := model.VocabularyUser{}
	if err := db.Pixie().Collection(model.CVocabularyUser).FindOneOrZero(
		ctx,
		Field_UserId.Equal(userId),
		&user,
	); err != nil {
		return Vocabulary{}, err
	}

	currentId := user.CurrentVocabularyId

	vocs := []model.Vocabulary{}
	if err := db.Pixie().Collection(model.CVocabulary).FindSortSkipLimit(
		ctx,
		mongo.Field_ID.Greater(currentId),
		mongo.Field_ID.Asc(),
		0,
		1,
		&vocs,
	); err != nil {
		return Vocabulary{}, err
	}

	if len(vocs) == 0 {
		currentId = primitive.NilObjectID
	}

	voc := struct {
		Id               primitive.ObjectID `bson:"_id"`
		model.Vocabulary `bson:"-,inline"`
	}{}
	if err := db.Pixie().Collection(model.CVocabulary).First(
		ctx,
		mongo.Field_ID.Greater(currentId),
		mongo.Field_ID.Asc(),
		&voc,
	); err != nil {
		return Vocabulary{}, err
	}

	user.UserId = userId
	user.CurrentVocabularyId = voc.Id

	if err := db.Pixie().Collection(model.CVocabularyUser).Upsert(
		ctx,
		Field_UserId.Equal(userId),
		user,
	); err != nil {
		return Vocabulary{}, err
	}

	vocabulary := FromModel(voc.Id, voc.Vocabulary)

	return vocabulary, nil
}

func Previous(ctx context.Context, userId string) (Vocabulary, error) {
	user := model.VocabularyUser{}
	if err := db.Pixie().Collection(model.CVocabularyUser).FindOneOrZero(
		ctx,
		Field_UserId.Equal(userId),
		&user,
	); err != nil {
		return Vocabulary{}, err
	}

	currentId := user.CurrentVocabularyId
	if currentId.IsZero() {
		currentId = primitive.NewObjectIDFromTimestamp(time.Now())
	}

	vocs := []model.Vocabulary{}
	if err := db.Pixie().Collection(model.CVocabulary).FindSortSkipLimit(
		ctx,
		mongo.Field_ID.Less(currentId),
		mongo.Field_ID.Desc(),
		0,
		1,
		&vocs,
	); err != nil {
		return Vocabulary{}, err
	}

	if len(vocs) == 0 {
		currentId = primitive.NewObjectIDFromTimestamp(time.Now())
	}

	voc := struct {
		Id               primitive.ObjectID `bson:"_id"`
		model.Vocabulary `bson:"-,inline"`
	}{}
	if err := db.Pixie().Collection(model.CVocabulary).First(
		ctx,
		mongo.Field_ID.Less(currentId),
		mongo.Field_ID.Desc(),
		&voc,
	); err != nil {
		return Vocabulary{}, err
	}

	user.UserId = userId
	user.CurrentVocabularyId = voc.Id

	if err := db.Pixie().Collection(model.CVocabularyUser).Upsert(
		ctx,
		Field_UserId.Equal(userId),
		user,
	); err != nil {
		return Vocabulary{}, err
	}

	vocabulary := FromModel(voc.Id, voc.Vocabulary)

	return vocabulary, nil
}

func Toggle(ctx context.Context, userId string, vocabularyId primitive.ObjectID) error {
	if vocabularyId.IsZero() {
		var user model.VocabularyUser
		if err := db.Pixie().Collection(model.CVocabularyUser).FindOne(
			ctx,
			Field_UserId.Equal(userId),
			&user,
		); err != nil {
			return err
		}

		vocabularyId = user.CurrentVocabularyId
	}

	bookmark := model.VocabularyBookmark{
		UserId:       userId,
		VocabularyId: vocabularyId,
		CreatedAt:    time.Now(),
	}

	condition := Field_UserId.Equal(userId).And(Field_VocabularyId.Equal(vocabularyId))

	count, err := db.Pixie().Collection(model.CVocabularyBookmark).Count(
		ctx,
		condition,
	)
	if err != nil {
		return err
	}

	if count == 0 {
		if err := db.Pixie().Collection(model.CVocabularyBookmark).Insert(
			ctx,
			bookmark,
		); err != nil {
			return err
		}

		return nil
	}

	if err := db.Pixie().Collection(model.CVocabularyBookmark).Delete(
		ctx,
		condition,
	); err != nil {
		return err
	}

	return nil
}

func BrowseBookmark(ctx context.Context, userId string, page int64) ([]Vocabulary, error) {
	bookmarks := []model.VocabularyBookmark{}

	var (
		limit int64 = 3
		skip  int64 = (page - 1) * 10
	)
	if err := db.Pixie().Collection(model.CVocabularyBookmark).FindSortSkipLimit(
		ctx,
		Field_UserId.Equal(userId),
		mongo.Field_ID.Desc(),
		skip,
		limit,
		&bookmarks,
	); err != nil {
		return nil, err
	}

	vocabularies := []Vocabulary{}
	if len(bookmarks) > 0 {
		ids := []primitive.ObjectID{}
		for _, bookmark := range bookmarks {
			ids = append(ids, bookmark.VocabularyId)
		}

		vocs := []struct {
			Id               primitive.ObjectID `bson:"_id"`
			model.Vocabulary `bson:"-,inline"`
		}{}
		if err := db.Pixie().Collection(model.CVocabulary).Find(
			ctx,
			mongo.Field_ID.In(ids),
			&vocs,
		); err != nil {
			return nil, err
		}

		cache := map[string]model.Vocabulary{}
		for _, voc := range vocs {
			cache[voc.Id.Hex()] = voc.Vocabulary
		}

		for _, bookmark := range bookmarks {
			if found, ok := cache[bookmark.VocabularyId.Hex()]; ok {
				vocabularies = append(vocabularies, FromModel(bookmark.VocabularyId, found))
			}
		}
	}

	return vocabularies, nil
}

func Upload() error {
	resources := []string{
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/國一.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/國二.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/國三.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/1級.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/2級.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/3級.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/4級.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/5級.json",
		"https://raw.githubusercontent.com/AppPeterPan/TaiwanSchoolEnglishVocabulary/main/6級.json",
	}

	ctx := context.Background()

	for i, resource := range resources {
		vocs, err := fetch(resource)
		if err != nil {
			return err
		}
		if err := upload(ctx, vocs, i+1); err != nil {
			return err
		}
	}

	return nil
}

func fetch(url string) ([]Vocabulary, error) {
	data, err := http.NewRequest().Get(url)
	if err != nil {
		return nil, err
	}

	vocabularies := []Vocabulary{}
	if err := json.Unmarshal(data, &vocabularies); err != nil {
		return nil, err
	}

	return vocabularies, nil
}

func upload(ctx context.Context, vocabularies []Vocabulary, level int) error {
	inserts := []interface{}{}
	for _, vocabulary := range vocabularies {
		voc := model.Vocabulary{
			Word:        vocabulary.Word,
			Definitions: []model.Definition{},
			Level:       level,
		}
		for _, definition := range vocabulary.Definitions {
			voc.Definitions = append(voc.Definitions, model.Definition{
				Text:         definition.Text,
				PartOfSpeech: definition.PartOfSpeech,
			})
		}
		inserts = append(inserts, voc)
	}
	if err := db.Pixie().Collection(model.CVocabulary).InsertMany(
		ctx,
		inserts,
	); err != nil {
		return err
	}

	return nil
}
