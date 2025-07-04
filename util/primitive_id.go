package util

import "go.mongodb.org/mongo-driver/v2/bson"

func GetPrimitiveID(id string) (*bson.ObjectID, error) {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	return &objId, nil
}
