package config

import (
	mgo "gopkg.in/mgo.v2"
	
	"log"
	"fmt"
)

var host = "mongodb://localhost:27017"
var dbName = "project"
var session *mgo.Session

func InitializeConnection(){ // first time goes so slow

	var err error
	if session, err = mgo.Dial(host); err != nil {
		log.Fatal(err)
	} 

	fmt.Println("Session started...")

}

func GetSession() (*mgo.Database, error){

	if session == nil {
		fmt.Println("Session not exist...")
		InitializeConnection()

	}

	sess := session.Clone()
	db := sess.DB(dbName)

	return db, nil

}