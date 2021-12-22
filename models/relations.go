package models

type Relation struct {
	UserID         string `bson:"user_id" json:"user_id"`
	UserRelationID string `bson:"user_relation_id" json:"user_relation_id"`
}

type ResponseRelation struct {
	Status bool `json:"status"`
}
