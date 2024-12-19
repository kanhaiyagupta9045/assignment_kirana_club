package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kanhaiyagupta9045/kirana_club/internals/db"
	"github.com/kanhaiyagupta9045/kirana_club/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreVisitService struct {
	client *mongo.Client
}

func NewStoreService() *StoreVisitService {
	return &StoreVisitService{
		client: db.DBConnection(),
	}
}

func (sv *StoreVisitService) InsertStoreVisitService(storeVisit models.StoresVisit) (primitive.ObjectID, error) {
	collection := sv.client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, storeVisit)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil

}

func (sv *StoreVisitService) UpdateStoreVisitServiceStatus(id primitive.ObjectID, status, errMssg, failedStoreID string) error {
	collection := sv.client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": status}}

	if status == "completed" {
		update = bson.M{"$set": bson.M{"status": status}}
	} else if status == "failed" {

		if errMssg == "" || failedStoreID == "" {
			return errors.New("error message and failed_store_id missing")
		} else {
			update = bson.M{
				"$set": bson.M{
					"status":          status,
					"error":           errMssg,
					"failed_store_id": failedStoreID,
				},
			}
		}

	} else {
		return mongo.ErrNoDocuments
	}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return err
}

func (sv *StoreVisitService) UpdateVisitInfo(id primitive.ObjectID, visitIndex int, newPerimeters []int64, newImageUUIDs []string) error {
	collection := sv.client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION"))

	log.Printf("Called UpdateVisitInfo: %v\n", id.Hex())
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("visits.%d.perimeters", visitIndex):  newPerimeters,
			fmt.Sprintf("visits.%d.image_uuids", visitIndex): newImageUUIDs,
		},
	}
	_, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (sv *StoreVisitService) GetStatusAndErrorByID(id primitive.ObjectID) (string, string, string, error) {
	collection := sv.client.Database(os.Getenv("DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION"))

	log.Printf("Called GetStatusAndErrorByID: %v\n", id.Hex())

	result := struct {
		Status        string `bson:"status"`
		Error         string `bson:"error"`
		FailedStoreID string `bson:"failed_store_id"`
	}{}

	filter := bson.M{"_id": id}
	projection := bson.M{"status": 1, "error": 1, "failed_store_id": 1} // fields to include

	// Find the document by ID with projection
	err := collection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", "", "", mongo.ErrNoDocuments
		}
		return "", "", "", err
	}

	return result.Status, result.Error, result.FailedStoreID, nil
}
