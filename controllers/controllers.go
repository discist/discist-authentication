package controllers

import (
	"authentication/db"
	"authentication/models"
	"authentication/utils"
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//db inits

var rdb = db.RedisInit()
var mongoclient = db.MongoInit()
var userCollection = mongoclient.Database("discistuserdb").Collection("users")
var ctx = context.Background()

func Save(user *models.User) error { //   save to db

	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ADDED NEW USER ", user.Email)
	return err
}

func GetByEmail(email string) (models.User, error) { // get by email
	var result models.User
	//var userlogin models.User

	err := userCollection.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(result, " this is the result bb")
	return result, err
}

func UpdatePassword(key string, value string, user models.User) {
	filter := bson.D{{key, value}}

	update := bson.D{{"$set", bson.D{{"password", user.Password}}}}

	_, e := userCollection.UpdateOne(context.Background(), filter, update)
	utils.CheckErorr(e)

}

func AddNewKey(objid string, addkey string, addvalue string) error {
	_id, err := primitive.ObjectIDFromHex(objid)
	if err != nil {

		return err

	}
	filter := bson.D{{"_id", _id}}

	update := bson.D{{"$set", bson.D{{addkey, addvalue}}}}

	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logrus.Info(err)
	}
	fmt.Println("updated ", addkey, addvalue)

	return err

}

func UpdateSessions(key string, value string, empty []models.Session) error {

	id, err := primitive.ObjectIDFromHex(value)
	if err != nil {
		panic(err)
	}
	filter := bson.D{{key, id}}
	fmt.Println(id)

	update := bson.D{{"$set", bson.D{{"sessions", empty}}}}

	_, err = userCollection.UpdateOne(context.Background(), filter, update)
	utils.CheckErorr(err)
	fmt.Println("update sucesss")
	return err
}

func GetByKey(key string, value string) (models.User, error) {

	filter := bson.D{{key, value}}
	var res models.User

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}
func UGetByKey(key string, value string) (models.UpdateUser, error) {

	filter := bson.D{{key, value}}
	var res models.UpdateUser

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}

func UGetByID(value string) (models.UpdateUser, error) {
	_id, err1 := primitive.ObjectIDFromHex(value)
	utils.CheckErorr(err1)
	filter := bson.D{{"_id", _id}}
	var res models.UpdateUser

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}

func GetByID(value string) (models.User, error) {
	_id, err1 := primitive.ObjectIDFromHex(value)
	utils.CheckErorr(err1)
	filter := bson.D{{"_id", _id}}
	var res models.User

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}

func Delete(id string) (*mongo.DeleteResult, error) {

	_id, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		panic(err1)
	}

	opts := options.Delete().SetCollation(&options.Collation{})

	res, err := userCollection.DeleteOne(context.Background(), bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Panic(err)
	}

	return res, err
}

func GetFullDoc(value string) (models.UserAllData, error) {
	_id, err1 := primitive.ObjectIDFromHex(value)
	utils.CheckErorr(err1)
	filter := bson.D{{"_id", _id}}
	var res models.UserAllData

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}

func GetAll() ([]models.UserAllDataPublic, error) {
	cursor, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Panic(err)
	}
	var docs []models.UserAllDataPublic

	for cursor.Next(context.Background()) {
		var single models.UserAllDataPublic
		err := cursor.Decode(&single)
		if err != nil {
			log.Panic(err)
		}
		docs = append(docs, single)
	}

	return docs, err

}

func Close() error {
	err := mongoclient.Disconnect(context.Background())
	fmt.Println("db closed")
	utils.CheckErorr(err)
	return err
}

func RedisAddKey(key string, value string) error {

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		logrus.Fatalln(err)
	}

	return err

}

func RedisGetKey(key string) (string, error) {

	value, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		fmt.Printf("%s doesnt exist \n", key)
		return "", err

	} else if err != nil {
		logrus.Info(err)
		return "", err

	}

	return value, nil

}

func RedisDelKey(key string) error {
	value, err := rdb.Del(ctx, key).Result()

	if err == redis.Nil {
		fmt.Printf("%s doesnt exist \n", key)
		fmt.Println(value)
		return err

	} else if err != nil {
		logrus.Info(err)
		return err

	}

	return nil

}
