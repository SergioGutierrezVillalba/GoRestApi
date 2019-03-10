package config

import (
	mgo "gopkg.in/mgo.v2"
	"log"
	"fmt"
)

var host = "mongodb://localhost:27017"
var dbName = "project"
var session *mgo.Session

func initialize(){ // first time goes so slow

	var err error
	if session, err = mgo.Dial(host); err != nil {
		log.Fatal(err)
	} 

	fmt.Println("Session started...")

}

func Connect() (*mgo.Database, error){

	if session == nil {
		fmt.Println("Session not exist...")
		initialize()

	}

	sess := session.Clone()
	// defer sess.Close() pending...
	db := sess.DB(dbName)

	return db, nil

	// session, err := mgo.Dial(host)

	// if err != nil {
	// 	return nil, err
	// } else {
	// 	db := session.DB(dbName)
	// 	return db, err
	// }

}