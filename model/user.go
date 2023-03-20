package model

type User struct {
	ID          string  `bson:"_id" json:"id"`
	DisplayName *string `bson:"displayName" json:"displayName"`
	Username    string  `bson:"username" json:"username"`
	Password    string  `bson:"password" json:"-"`
}
