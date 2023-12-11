package db

import (
	"context"
	"github.com/claudiu/gocron"
	"github.com/why-xn/alap-backend/pkg/config"
	"github.com/why-xn/alap-backend/pkg/core/log"
	"github.com/why-xn/alap-backend/pkg/enum"
	"github.com/why-xn/alap-backend/pkg/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"sync"
)

type DbManagerInterface interface {
	initConnection()
	InsertSingleDocument(collectionName string, document interface{}) (primitive.ObjectID, error)
	InsertMultipleDocument(collectionName string, documents []interface{}) ([]interface{}, error)
	FindOne(collectionName string, filter interface{}, objType reflect.Type) interface{}
	FindOneByObjId(collectionName string, objId primitive.ObjectID, objType reflect.Type) interface{}
	FindOneByStrId(collectionName string, strId string, objType reflect.Type) interface{}
	FindAll(collectionName string, objType reflect.Type, filter interface{}, sortParam *types.SortParam, start int64, limit int64) []DbFindObj
	UpdateOneByObjId(collectionName string, objId primitive.ObjectID, document interface{}) error
	UpdateOneByStrId(collectionName string, strId string, document interface{}) error
	DeleteOneByObjId(collectionName string, objId primitive.ObjectID) error
	DeleteOneByStrId(collectionName string, strId string) error
	RestoreOneByObjId(collectionName string, objId primitive.ObjectID) error
	PermanentDeleteOneByObjId(collectionName string, objId primitive.ObjectID) error
	PermanentDeleteOneByStrId(collectionName string, strId string) error
	DropCollection(collectionName string) error
}

type dbManager struct {
	ctx context.Context
	db  *mongo.Database
}

// Implementing Singleton
var singletonDbManager *dbManager
var onceDbManager sync.Once

func GetDbManager() *dbManager {
	onceDbManager.Do(func() {
		singletonDbManager = &dbManager{}
		singletonDbManager.initConnection()
	})
	return singletonDbManager
}

func (dm *dbManager) initConnection() {
	log.Logger.Infow("Initializing DB Connection")
	// Base context.
	ctx := context.Background()
	dm.ctx = ctx
	clientOpts := options.Client().ApplyURI(config.DatabaseConnectionString)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Logger.Fatalw("Failed to establish DB Connection", "err", err.Error())
		return
	}

	// Check db connection
	checkDbConnection(ctx, client)

	// Keep checking db connection every 10s
	go func() {
		s := gocron.NewScheduler()
		s.Every(10).Seconds().Do(checkDbConnection, ctx, client)
		s.Start()
	}()

	db := client.Database(config.DatabaseName)
	dm.db = db

	/*err = dm.DropCollection(collection.Rfid)
	if err != nil {
		log.Logger.Fatalw("Failed to drop collection", "err", err.Error())
	}
	err = dm.DropCollection(collection.RfidHistory)
	if err != nil {
		log.Logger.Fatalw("Failed to drop collection", "err", err.Error())
	}*/

	log.Logger.Infow("DB Connection established")
}

func checkDbConnection(ctx context.Context, client *mongo.Client) {
	err := client.Ping(ctx, nil)
	if err != nil {
		log.Logger.Fatalw("Failed to verify DB Connection", err.Error())
	}
}

func (dm *dbManager) InsertSingleDocument(collectionName string, document interface{}) (primitive.ObjectID, error) {
	coll := dm.db.Collection(collectionName)

	result, err := coll.InsertOne(dm.ctx, document)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	// ID of the inserted document.
	objectID := result.InsertedID.(primitive.ObjectID)
	return objectID, nil
}

func (dm *dbManager) InsertMultipleDocument(collectionName string, documents []interface{}) ([]interface{}, error) {
	coll := dm.db.Collection(collectionName)
	results, err := coll.InsertMany(dm.ctx, documents)
	if err != nil {
		return nil, err
	}

	return results.InsertedIDs, nil
}

func (dm *dbManager) FindOne(collectionName string, filter interface{}, objType reflect.Type) interface{} {
	coll := dm.db.Collection(collectionName)

	findResult := coll.FindOne(dm.ctx, filter)
	if err := findResult.Err(); err != nil {
		return nil
	}

	objValue := reflect.New(objType)
	obj := objValue.Interface()
	err := findResult.Decode(obj)
	if err != nil {
		log.Logger.Errorw("error occurred while find document decoding", "err", err.Error())
		return nil
	}
	return obj
}

func (dm *dbManager) FindOneByObjId(collectionName string, objId primitive.ObjectID, objType reflect.Type) interface{} {
	coll := dm.db.Collection(collectionName)

	findResult := coll.FindOne(dm.ctx, bson.M{"_id": objId, "status": enum.StatusValid})
	if err := findResult.Err(); err != nil {
		return nil
	}

	objValue := reflect.New(objType)
	obj := objValue.Interface()
	err := findResult.Decode(obj)
	if err != nil {
		log.Logger.Errorw("error occurred while find document decoding", "err", err.Error())
		return nil
	}
	return obj
}

func (dm *dbManager) FindOneByStrId(collectionName string, strId string, objType reflect.Type) interface{} {
	objId, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		log.Logger.Errorw("error occurred while converting str id to obj id", "err", err.Error())
		return nil
	}
	return dm.FindOneByObjId(collectionName, objId, objType)
}

func (dm *dbManager) FindAll(collectionName string, objType reflect.Type, filter interface{}, sortParam *types.SortParam, skip int64, limit int64) []DbFindObj {
	// Pass these options to the Find method
	findOptions := options.Find()

	if sortParam != nil {
		findOptions.SetSort(bson.D{{sortParam.SortBy, sortParam.Type}})
	}

	if skip > 0 {
		findOptions.SetSkip(skip)
	}

	if limit > -1 {
		findOptions.SetLimit(limit)
	}

	coll := dm.db.Collection(collectionName)

	reflect.New(objType)

	var results []DbFindObj

	// Passing bson.D{{}} as the filter matches all documents in the collection
	if filter == nil {
		filter = bson.D{{}}
	}
	cur, err := coll.Find(dm.ctx, filter, findOptions)
	if err != nil {
		log.Logger.Errorw("error occurred while finding", "err", err.Error())
		//log.Println("[ERROR]", err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		elemValue := reflect.New(objType)
		elem := elemValue.Interface()

		err := cur.Decode(elem)
		if err != nil {
			log.Logger.Errorw("error occurred while decoding", "err", err.Error())
			break
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Logger.Errorw("error occurred in cursor", "err", err.Error())
	} else {
		// Close the cursor once finished
		cur.Close(context.TODO())
	}

	return results
}

func (dm *dbManager) UpdateOneByObjId(collectionName string, objId primitive.ObjectID, document interface{}) error {
	coll := dm.db.Collection(collectionName)

	filter := bson.M{"_id": objId, "status": enum.StatusValid}

	update := bson.M{"$set": document}

	// Call the driver's UpdateOne() method and pass filter and update to it
	_, err := coll.UpdateOne(
		dm.ctx,
		filter,
		update,
	)

	return err
}

func (dm *dbManager) UpdateOneByStrId(collectionName string, strId string, document interface{}) error {
	objId, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		log.Logger.Errorw("error occurred while converting str id to obj id", "err", err.Error())
		return nil
	}
	return dm.UpdateOneByObjId(collectionName, objId, document)
}

func (dm *dbManager) DeleteOneByObjId(collectionName string, objId primitive.ObjectID) error {
	coll := dm.db.Collection(collectionName)

	filter := bson.M{"_id": objId, "status": enum.StatusValid}

	update := bson.M{"$set": bson.M{"status": enum.StatusDeleted}}

	// Call the driver's UpdateOne() method and pass filter and update to it
	_, err := coll.UpdateOne(
		dm.ctx,
		filter,
		update,
	)

	return err
}

func (dm *dbManager) DeleteOneByStrId(collectionName string, strId string) error {
	objId, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		log.Logger.Errorw("error occurred while converting str id to obj id", "err", err.Error())
		return nil
	}
	return dm.DeleteOneByObjId(collectionName, objId)
}

func (dm *dbManager) RestoreOneByObjId(collectionName string, objId primitive.ObjectID) error {
	coll := dm.db.Collection(collectionName)

	filter := bson.M{"_id": objId, "status": enum.StatusDeleted}

	update := bson.M{"$set": bson.M{"status": enum.StatusValid}}

	// Call the driver's UpdateOne() method and pass filter and update to it
	_, err := coll.UpdateOne(
		dm.ctx,
		filter,
		update,
	)

	return err
}

func (dm *dbManager) PermanentDeleteOneByObjId(collectionName string, objId primitive.ObjectID) error {
	coll := dm.db.Collection(collectionName)

	filter := bson.M{"_id": objId}

	_, err := coll.DeleteOne(dm.ctx, filter)

	return err
}

func (dm *dbManager) PermanentDeleteOneByStrId(collectionName string, strId string) error {
	objId, err := primitive.ObjectIDFromHex(strId)
	if err != nil {
		log.Logger.Errorw("error occurred while converting str id to obj id", "err", err.Error())
		return nil
	}
	return dm.PermanentDeleteOneByObjId(collectionName, objId)
}

func (dm *dbManager) DropCollection(collectionName string) error {
	coll := dm.db.Collection(collectionName)

	err := coll.Drop(dm.ctx)

	return err
}
