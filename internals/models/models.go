package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type VisitInfo struct {
	StoreID    string   `bson:"store_id" json:"store_id"`
	VisitTime  string   `bson:"visit_time" json:"visit_time"`
	ImageURLs  []string `bson:"image_urls" json:"image_url"`
	ImageUUIDs []string `bson:"image_uuids"`
	Perimeters []int64  `bson:"perimeters"`
}

type StoresVisit struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Status        string             `bson:"status"`
	Error         string             `bson:"error"`
	FailedStoreID string             `bson:"failed_store_id"`
	Count         int                `bson:"count" json:"count"`
	Visits        []VisitInfo        `bson:"visits" json:"visits"`
}
