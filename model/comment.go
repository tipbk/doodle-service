package model

type Comment struct {
	ID      string  `bson:"_id" json:"id"`
	UserID  string  `bson:"userId" json:"userId"`
	Comment string  `bson:"comment" json:"comment"`
	PostId  string  `bson:"postId" json:"postId"`
	ReplyOn *string `bson:"replyOn" json:"replyOn"`
}
