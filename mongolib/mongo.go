package mongolib

import (
	"book-restapi/book"
	"book-restapi/configs"
	"context"
	"sync"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatastore struct {
	Config  *configs.Mongo
	DB      *mongo.Database
	Session *mongo.Client
}

func NewDatasore(mongoConfigs *configs.Mongo) (*MongoDatastore, error) {
	mongoDatastore := new(MongoDatastore)
	mongoDatastore.Config = mongoConfigs
	db, session, err := connect(mongoConfigs)

	if db != nil && session != nil {
		mongoDatastore.DB = db
		mongoDatastore.Session = session
	}

	return mongoDatastore, err
}

func connect(mongoConfigs *configs.Mongo) (a *mongo.Database, b *mongo.Client, c error) {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	var err error

	connectOnce.Do(func() {
		db, session, err = connectToMongo(mongoConfigs)
	})

	return db, session, err
}

func connectToMongo(mongoConfigs *configs.Mongo) (a *mongo.Database, b *mongo.Client, c error) {
	connStr := mongoConfigs.URI
	client, err := mongo.NewClient(options.Client().ApplyURI(connStr))

	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		return nil, nil, err
	}

	var DB = client.Database(mongoConfigs.Database)

	return DB, client, nil
}

// get data from database
func GetBooks(ctx context.Context, collection *mongo.Collection) ([]book.Book, error) {
	cursor, err := collection.Find(ctx, bson.D{})

	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	var books []book.Book
	if err = cursor.All(ctx, &books); err != nil {
		return nil, err
	}

	return books, nil
}

// get specific data from database

func GetBook(ctx context.Context, collection *mongo.Collection, id string) (book.Book, error) {
	var book book.Book

	err := collection.FindOne(ctx, bson.M{"book_id": id}).Decode(&book)

	if err != nil {
		return book, err
	}

	return book, nil
}

// insert data into database

func InsertBook(ctx context.Context, collection *mongo.Collection, book book.Book) (interface{}, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	book.ID = id.String()

	res, err := collection.InsertOne(ctx, book)

	if err != nil {
		return nil, err
	}

	return res.InsertedID, nil

}

// update data - updating title of the book
func UpdateBook(ctx context.Context, collection *mongo.Collection, id string, title string) (interface{}, error) {

	result, err := collection.UpdateOne(
		ctx,
		bson.M{"book_id": id},
		bson.D{
			{"$set", bson.D{{"title", title}}},
		},
	)

	if err != nil {
		return nil, err
	}

	if result.MatchedCount != 0 {
		return result.MatchedCount, nil
	}

	if result.UpsertedCount != 0 {
		return result.UpsertedCount, nil
	}

	return result.MatchedCount, nil
}

// delete data
func DeleteBook(ctx context.Context, collection *mongo.Collection, title string) (interface{}, error) {
	options := options.Delete().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})

	res, err := collection.DeleteOne(ctx, bson.M{"title": title}, options)

	if err != nil {
		return nil, err
	}

	return res.DeletedCount, nil
}
