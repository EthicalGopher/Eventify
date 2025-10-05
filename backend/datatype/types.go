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

// Event type to store event data
type Event struct {
	ID               string   `bson:"id"`
	Title            string   `bson:"title"`
	Discription      string   `bson:"discription"`
	Category         []string `bson:"category,omitempty"`
	Location         string   `bson:"location"`
	Date             string   `bson:"date"`
	Price            float32  `bson:"price"`
	TicketsAvailable int64    `bson:"tickets_available"`
	OrganizerID      string   `bson:"organizer_id"`
	ImageURL         string   `bson:"image_url,omitempty"`
	Participants     []string `bson:"participants,omitempty"`
}
