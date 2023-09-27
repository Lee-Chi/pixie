package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	client *mongo.Client
}

type Options struct {
	User          string
	Password      string
	MaxPoolSize   uint64
	ReplicaSet    string
	UseApmMonitor bool
}

func NewClient(ctx context.Context, domain, user, password string) (*Client, error) {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/?retryWrites=true&w=majority",
		user,
		password,
		domain,
	)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (mongo *Client) Destroy() error {
	return mongo.client.Disconnect(context.Background())
}

func (m *Client) Database(name string) *Database {
	opts := &options.DatabaseOptions{}
	// opts.SetReadPreference(readpref.Primary()).
	// 	SetReadConcern(readconcern.Linearizable()).
	// 	SetWriteConcern(writeconcern.New(writeconcern.WMajority(), writeconcern.J(true)))
	return &Database{
		database:    m.client.Database(name, opts),
		readOnly:    false,
		upgradeOnly: false,
	}
}

func (m *Client) DatabaseReadOnly(name string) *Database {
	opts := &options.DatabaseOptions{}
	// opts.SetReadPreference(readpref.SecondaryPreferred()).
	// SetReadConcern(readconcern.Local()).
	// SetWriteConcern(writeconcern.New(writeconcern.W(0)))

	return &Database{
		database:    m.client.Database(name, opts),
		readOnly:    true,
		upgradeOnly: false,
	}
}

func (m *Client) DatabaseUpgradeOnly(name string) *Database {
	opts := &options.DatabaseOptions{}
	// opts.SetReadPreference(readpref.SecondaryPreferred()).
	// SetReadConcern(readconcern.Local()).
	// SetWriteConcern(writeconcern.New(writeconcern.W(0)))

	return &Database{
		database:    m.client.Database(name, opts),
		readOnly:    false,
		upgradeOnly: true,
	}
}

type Database struct {
	database *mongo.Database

	readOnly    bool
	upgradeOnly bool
}

func (d *Database) Collection(name string) *Collection {
	return &Collection{
		collection:  d.database.Collection(name),
		name:        name,
		readOnly:    d.readOnly,
		upgradeOnly: d.upgradeOnly,
	}
}

type Collection struct {
	collection *mongo.Collection

	name string

	readOnly    bool
	upgradeOnly bool
}

func (c *Collection) CreateIndex(ctx context.Context, index Index, opts ...*options.CreateIndexesOptions) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}

	if _, err := c.collection.Indexes().CreateOne(ctx, index.ToIndexModel(), opts...); err != nil {
		return err
	}

	return nil
}

func (c *Collection) CreateIndexes(ctx context.Context, indexes Indexes, opts ...*options.CreateIndexesOptions) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}

	if _, err := c.collection.Indexes().CreateMany(ctx, indexes.ToIndexModels(), opts...); err != nil {
		return err
	}

	return nil
}

func (c *Collection) DropOneIndex(ctx context.Context, name string) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}

	if _, err := c.collection.Indexes().DropOne(ctx, name); err != nil {
		return err
	}

	return nil
}

func (c *Collection) DropAllIndexes(ctx context.Context) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}

	if _, err := c.collection.Indexes().DropAll(ctx); err != nil {
		return err
	}

	return nil
}

// Aggregate 聚合查詢
func (c *Collection) Aggregate(ctx context.Context, condition interface{}, result interface{}) error {
	cursor, err := c.collection.Aggregate(ctx, condition)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// Find 查詢多筆資料
func (c *Collection) Find(ctx context.Context, condition C, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	cursor, err := c.collection.Find(ctx, condition)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindProj 查詢多筆指定欄位資料
func (c *Collection) FindProj(ctx context.Context, condition C, projection P, result interface{}) error {
	if condition == nil {
		condition = C{}
	}
	if projection == nil {
		projection = P{}
	}
	cursor, err := c.collection.Find(ctx, condition, options.Find().SetProjection(projection))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindSort 查詢多筆資料
func (c *Collection) FindSort(ctx context.Context, condition C, sort S, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	cursor, err := c.collection.Find(ctx, condition, options.Find().SetSort(sort))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, result)
}

// FindSortSkipLimit 查詢指定筆數資料
func (c *Collection) FindSortSkipLimit(ctx context.Context, condition C, sort S, skip int64, limit int64, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	cursor, err := c.collection.Find(ctx, condition, options.Find().SetSort(sort).SetSkip(skip).SetLimit(limit))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindProjSort 查詢多筆指定欄位資料
func (c *Collection) FindProjSort(ctx context.Context, condition C, projection P, sort interface{}, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	cursor, err := c.collection.Find(ctx, condition, options.Find().SetProjection(projection).SetSort(sort))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindProjSortSkipLimit 查詢指定筆數查詢指定欄位資料
func (c *Collection) FindProjSortSkipLimit(ctx context.Context, condition C, projection P, sort S, skip int64, limit int64, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	cursor, err := c.collection.Find(ctx, condition, options.Find().SetProjection(projection).SetSort(sort).SetSkip(skip).SetLimit(limit))
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

// FindOneByID 查詢指定ID資料
func (c *Collection) FindOneByID(ctx context.Context, id primitive.ObjectID, result interface{}) error {
	sr := c.collection.FindOne(ctx, bson.M{"_id": id})
	if err := sr.Err(); err != nil {
		return err
	}
	if err := sr.Decode(result); err != nil {
		return err
	}

	return nil
}

// FindOne 查詢一筆資料
func (c *Collection) FindOne(ctx context.Context, condition C, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	sr := c.collection.FindOne(ctx, condition)
	if err := sr.Err(); err != nil {
		return err
	}
	if err := sr.Decode(result); err != nil {
		return err
	}

	return nil
}

func (c *Collection) FindOneOrZero(ctx context.Context, condition C, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	sr := c.collection.FindOne(ctx, condition)
	if err := sr.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}

		return err
	}
	if err := sr.Decode(result); err != nil {
		return err
	}

	return nil
}

func (c *Collection) First(ctx context.Context, condition C, sort S, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	if sort == nil {
		sort = S{}
	}

	sr := c.collection.FindOne(ctx, condition, options.FindOne().SetSort(sort).SetSkip(0))
	if err := sr.Err(); err != nil {
		return err
	}
	if err := sr.Decode(result); err != nil {
		return err
	}

	return nil
}

// FindOneProj 查詢一筆指定欄位資料
func (c *Collection) FindOneProj(ctx context.Context, condition C, projection P, result interface{}) error {
	if condition == nil {
		condition = C{}
	}

	sr := c.collection.FindOne(ctx, condition, options.FindOne().SetProjection(projection))
	if err := sr.Err(); err != nil {
		return err
	}
	if err := sr.Decode(result); err != nil {
		return err
	}

	return nil
}

// Count 統計查詢的資料數量
func (c *Collection) Count(ctx context.Context, condition C) (int64, error) {
	if condition == nil {
		condition = C{}
	}

	count, err := c.collection.CountDocuments(ctx, condition)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Upsert ...
func (c *Collection) Upsert(ctx context.Context, condition C, upsert interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	options := options.UpdateOptions{}
	if _, err := c.collection.UpdateMany(ctx, condition, bson.M{"$set": upsert}, options.SetUpsert(true)); err != nil {
		return err
	}

	return nil
}

// UpdateByID 更新指定id資料
func (c *Collection) UpdateByID(ctx context.Context, id interface{}, update interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if _, err := c.collection.UpdateByID(ctx, id, bson.M{"$set": update}); err != nil {
		return err
	}

	return nil
}

// Update 更新資料
func (c *Collection) Update(ctx context.Context, condition C, update interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.UpdateOne(ctx, condition, bson.M{"$set": update}); err != nil {
		return err
	}

	return nil
}

// UpdateOne ...
func (c *Collection) UpdateOne(ctx context.Context, condition C, update interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	result, err := c.collection.UpdateOne(ctx, condition, bson.M{"$set": update})
	if err != nil {
		return err
	}

	if result.ModifiedCount != 1 {
		return fmt.Errorf("no one be modified")
	}

	return nil
}

// UpdateMany 更新多筆資料
func (c *Collection) UpdateMany(ctx context.Context, condition C, update interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.UpdateMany(ctx, condition, bson.M{"$set": update}); err != nil {
		return err
	}

	return nil
}

func (c *Collection) Insert(ctx context.Context, insert interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if _, err := c.collection.InsertOne(ctx, insert); err != nil {
		return err
	}

	return nil
}

// InsertOne 新增一筆資料
func (c *Collection) InsertOne(ctx context.Context, insert interface{}) (primitive.ObjectID, error) {
	if c.readOnly {
		return primitive.NilObjectID, fmt.Errorf("use read only connect")
	}

	result, err := c.collection.InsertOne(ctx, insert)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// InsertMany 新增多筆資料
func (c *Collection) InsertMany(ctx context.Context, inserts []interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if _, err := c.collection.InsertMany(ctx, inserts); err != nil {
		return err
	}

	return nil
}

// DeleteByID 移除指定id資料
func (c *Collection) DeleteByID(ctx context.Context, id interface{}) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if _, err := c.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}

// Delete 移除一筆資料
func (c *Collection) Delete(ctx context.Context, condition C) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.DeleteOne(ctx, condition); err != nil {
		return err
	}

	return nil
}

// DeleteOne
func (c *Collection) DeleteOne(ctx context.Context, condition C) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	result, err := c.collection.DeleteOne(ctx, condition)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no document")
	}

	return nil
}

// DeleteMany 移除資料
func (c *Collection) DeleteMany(ctx context.Context, condition C) error {
	if c.readOnly {
		return fmt.Errorf("use read only connect")
	}

	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.DeleteMany(ctx, condition); err != nil {
		return err
	}

	return nil
}

// Rename ...
func (c *Collection) Rename(ctx context.Context, condition C, update U) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}
	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.UpdateMany(
		ctx,
		condition,
		U{"$rename": update},
	); err != nil {
		return err
	}

	return nil
}

// Replace ...
func (c *Collection) Replace(ctx context.Context, condition C, replacement interface{}) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}
	if condition == nil {
		condition = C{}
	}

	if _, err := c.collection.ReplaceOne(
		ctx,
		condition,
		replacement,
	); err != nil {
		return err
	}

	return nil
}

func (c *Collection) Drop(ctx context.Context) error {
	if !c.upgradeOnly {
		return fmt.Errorf("need to use upgrade only connection")
	}

	if err := c.collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}
