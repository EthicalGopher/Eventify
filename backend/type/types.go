// Package datatype contains the all the types
package datatype

// User for storing user data
type User struct {
	ID        string   `bson:"_id,omitempty"`
	Name      string   `bson:"name"`
	Email     string   `bson:"email"`
	Password  string   `bson:"password"`
	Interests []string `bson:"interests,omitempty"`
	Role      string   `bson:"role"`
}

// Event for storing Event data
type Event struct {
	ID               string   `bson:"_id,omitempty"`
	Title            string   `bson:"title"`
	Description      string   `bson:"description"`
	Category         []string `bson:"Category"`
	Location         string   `bson:"location"`
	Price            float64  `bson:"price"`
	TicketsAvailable int      `bson:"ticketsAvailable"`
	OrganizerID      string   `bson:"organizerID,omitempty"`
	ImageURL         string   `bson:"imageURL,omitempty"`
}
