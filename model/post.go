package model

type Post struct {
	ID          string `bson:"_id" json:"id"`
	UserID      string `bson:"userId" json:"userId"`
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Hashtag     string `bson:"hashtag" json:"hashtag"`
}
