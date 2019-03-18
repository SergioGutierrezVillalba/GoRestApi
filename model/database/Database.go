package database

import (
	mgo "gopkg.in/mgo.v2"
	
	"log"
	"fmt"
)

var (
	host = "mongodb://localhost:27017"
	dbName = "project"
	session *mgo.Session
)

// var host = "mongodb://localhost:27017"
// var dbName = "project"
// var session *mgo.Session

func StartConnection(){

	var err error
	if session, err = mgo.Dial(host); err != nil {
		log.Fatal(err)
	} 

	fmt.Println("Session started...")

}

func GetSession() (*mgo.Database, error){

	if session == nil {
		fmt.Println("Session doesn't exist...")
		StartConnection()
	}

	sess := session.Clone()
	db := sess.DB(dbName)

	return db, nil
}

func CloseSession(db *mgo.Database){
	db.Session.Close()
}