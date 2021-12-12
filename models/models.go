package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAllData struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Sessions  []Session          `json:"sessions" bson:"sessions,"`
	Username  string             `json:"username" bson:"username,"`
	Story     string             `json:"story" bson:"story,"`
	Subject   string             `json:"subject" bson:"subject,"`
	State     string             `json:"state" bson:"state,"`
}

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email     string             `json:"email" bson:"email,omitempty"`
	Password  string             `json:"password" bson:"password,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Sessions  []Session          `json:"sessions" bson:"sessions,"`
}

type UpdateUser struct {
	Uuid     string `json:"uuid" bson:"uuid,omitempty"`
	Username string `json:"username" bson:"username,"`
	Story    string `json:"story" bson:"stsory,"`
	Subject  string `json:"subject" bson:"subject,"`
	State    string `json:"state" bson:"state,"`
}

type Session struct {
	Uuid     string `json:"uuid" bson:"uuid,omitempty"`
	Device   string `json:"device" bson:"device,omitempty"`
	Location string `json:"location" bson:"location,omitempty"`
}

type Login struct {
	Email    string `json:"email" bson:"email,omitempty"`
	Password string `json:"password" bson:"password,omitempty"`
	Device   string `json:"device" bson:"device,omitempty"`
	Location string `json:"location" bson:"location,omitempty"`
}
