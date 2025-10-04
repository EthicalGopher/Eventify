// Package type contains the all the types
package datatype

type User struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	Interests []string `bson:"interests,omitempty"`
	Role      string   `bson:"role"`
}
