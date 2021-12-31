package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserAllData struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email           string             `json:"email" bson:"email,omitempty"`
	Password        string             `json:"password" bson:"password,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Sessions        []Session          `json:"sessions" bson:"sessions,"`
	Username        string             `json:"username" bson:"username,"`
	Story           string             `json:"story" bson:"story,"`
	Subject         string             `json:"subject" bson:"subject,"`
	State           string             `json:"state" bson:"state,"`
	ProfilePhotoUrl string             `json:"profilephotourl" bson:"profilephotourl,"`
}

type UserAllDataPublic struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	//Email           string             `json:"email" bson:"email,omitempty"`
	//Password        string             `json:"password" bson:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	//UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	//Sessions        []Session          `json:"sessions" bson:"sessions,"`
	Username        string `json:"username" bson:"username,"`
	Story           string `json:"story" bson:"story,"`
	Subject         string `json:"subject" bson:"subject,"`
	State           string `json:"state" bson:"state,"`
	ProfilePhotoUrl string `json:"profilephotourl" bson:"profilephotourl,"`
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
	Uuid     string   `json:"uuid" bson:"uuid,omitempty"`
	Username string   `json:"username" bson:"username,"`
	Story    string   `json:"story" bson:"story,"`
	Subject  string   `json:"subject" bson:"subject,"`
	State    string   `json:"state" bson:"state,"`
	Posts    []string `json:"post" bson:"postid,"`
}

type Session struct {
	Uuid     string `json:"uuid" bson:"uuid,omitempty"`
	Device   string `json:"device" bson:"device,omitempty"`
	Location string `json:"location" bson:"location,omitempty"`
}

type Login struct {
	Email    string `json:"email" bson:"email,"`
	Password string `json:"password" bson:"password,"`
	Device   string `json:"device" bson:"device,omitempty"`
	Location string `json:"location" bson:"location,omitempty"`
}

type UsernameCheck struct {
	Username string `json:"username" bson:"username,omitempty"`
}

type Post struct {
	Uuid         string             `json:"uuid" bson:"uuid,omitempty"`
	Username     string             `json:"username" bson:"username,"`
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Type         string             `json:"type" bson:"type,"`
	Subject      string             `json:"subject" bson:"subject,omitempty"`
	Date         time.Time          `json:"time" bson:"time,omitempty"`
	Interactions int64              `json:"interactions" bson:"interactions,omitempty"`
	Brief        string             `json:"brief" bson:"brief,omitempty"`
	MediaURL     string             `json:"mediaurl" bson:"mediaurl,"`
	Comments     []Comment          `json:"comments" bson:"comments,omitempty"`
}

type Comment struct {
	Username        string             `json:"username" bson:"username,"`
	ProfilePhotoUrl string             `json:"profilephotourl" bson:"profilephotourl,omitempty"`
	Id              primitive.ObjectID `json:"id" bson:"id,"`
	Date            time.Time          `json:"time" bson:"time,"`
	Interactions    int64              `json:"interactions" bson:"interactions,omitempty"`
	Comment         string             `json:"comment" bson:"comment,omitempty"`
}

type Request struct {
	Uuid string `json:"uuid" bson:"uuid,omitempty"`
}

type DeletePost struct {
	Uuid   string `json:"uuid" bson:"uuid,omitempty"`
	PostID string `json:"postid" bson:"postid,omitempty"`
}

type GetUser struct {
	Uuid     string `json:"uuid" bson:"uuid,omitempty"`
	UserName string `json:"username" bson:"username,omitempty"`
}
