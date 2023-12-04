package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Field_ID F = Field("_id")
)

type F struct {
	key string
}

func Field(key string) F {
	return F{key: key}
}
func (f F) Less(value interface{}) C {
	return C{f.key: C{"$lt": value}}
}
func (f F) LessEqual(value interface{}) C {
	return C{f.key: C{"$lte": value}}
}
func (f F) Greater(value interface{}) C {
	return C{f.key: C{"$gt": value}}
}
func (f F) GreaterEqual(value interface{}) C {
	return C{f.key: C{"$gte": value}}
}
func (f F) In(value interface{}) C {
	return C{f.key: C{"$in": value}}
}
func (f F) Equal(value interface{}) C {
	return C{f.key: value}
}
func (f F) NotEqual(value interface{}) C {
	return C{f.key: C{"$ne": value}}
}
func (f F) Exist(value interface{}) C {
	return C{f.key: C{"$exists": value}}
}
func (f F) Between(min interface{}, max interface{}) C {
	return C{f.key: C{"$gte": min, "$lte": max}}
}
func (f F) BetweenGreaterMin(min interface{}, max interface{}) C {
	return C{f.key: C{"$gt": min, "$lte": max}}
}
func (f F) BetweenLessMax(min interface{}, max interface{}) C {
	return C{f.key: C{"$gte": min, "$lt": max}}
}

func (f F) Asc() S {
	return S{bson.E{Key: f.key, Value: 1}}
}
func (f F) Desc() S {
	return S{bson.E{Key: f.key, Value: -1}}
}
func (f F) SortBy(value int) S {
	return S{bson.E{Key: f.key, Value: value}}
}

func (f F) Set(value interface{}) U {
	return U{f.key: value}
}

func (f F) Project() P {
	return P{f.key: 1}
}

func (f F) Ignore() P {
	return P{f.key: 0}
}

func (f F) Reference() string {
	return fmt.Sprintf("$%s", f.key)
}

func (f F) Index(value interface{}) primitive.E {
	return primitive.E{Key: f.key, Value: value}
}

func Condition(conditions ...C) C {
	cond := C{}
	for _, condition := range conditions {
		for key, value := range condition {
			cond[key] = value
		}
	}

	return cond
}

func Or(conditions []C) C {
	if len(conditions) == 0 {
		return C{}
	}

	return C{
		"$or": conditions,
	}
}

func And(conditions []C) C {
	if len(conditions) == 0 {
		return C{}
	}

	return C{
		"$and": conditions,
	}
}

type C bson.M

func (self C) Set(condition C) C {
	for key, value := range condition {
		self[key] = value
	}
	return self
}

func (self C) Or(conditions ...C) C {
	if len(conditions) == 1 && len(conditions[0]) == 0 {
		return self
	}

	conds := []C{}
	if len(self) > 0 {
		conds = append(conds, self)
	}

	for _, condition := range conditions {
		if len(condition) > 0 {
			conds = append(conds, condition)
		}
	}

	if len(conds) == 1 {
		return conds[0]
	}

	return C{"$or": conds}
}

func (self C) And(conditions ...C) C {
	if len(conditions) == 1 && len(conditions[0]) == 0 {
		return self
	}

	conds := []C{}
	if len(self) > 0 {
		conds = append(conds, self)
	}

	for _, condition := range conditions {
		if len(condition) > 0 {
			conds = append(conds, condition)
		}
	}

	if len(conds) == 1 {
		return conds[0]
	}

	return C{"$and": conds}
}

type S bson.D

func Sort(sorts ...S) S {
	s := S{}
	for _, sort := range sorts {
		for _, v := range sort {
			s = append(s, v)
		}
	}
	return s
}

func (self S) Append(sort S) S {
	for _, s := range sort {
		self = append(self, s)
	}

	return self
}

type U bson.M

func Update(updates ...U) U {
	u := U{}
	for _, update := range updates {
		for k, v := range update {
			u[k] = v
		}
	}

	return u
}

func (u U) Merge(update U) U {
	for k, v := range update {
		u[k] = v
	}

	return u
}

type P bson.M

func Project(projects ...P) P {
	p := P{}
	for _, project := range projects {
		for k, v := range project {
			p[k] = v
		}
	}

	return p
}

type Index mongo.IndexModel

func (index Index) ToIndexModel() mongo.IndexModel {
	return mongo.IndexModel(index)
}

type Indexes []mongo.IndexModel

func (indexes Indexes) ToIndexModels() []mongo.IndexModel {
	return []mongo.IndexModel(indexes)
}
