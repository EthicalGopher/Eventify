// Package datatype type contains the all the types
package datatype

// User type to store user data
type User struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	Interests []string `bson:"interests,omitempty"`
	Role      string   `bson:"role"`
}
