package database

import (
	mgo "gopkg.in/mgo.v2"
	
	"log"
	"fmt"
)

var (
	host = "mongodb://localhost:27017"
	dbName = "project"
)

type Db struct {
	Session *mgo.Session
}

func (db *Db) StartConnection(){

	var err error
	if db.Session, err = mgo.Dial(host); err != nil {
		log.Fatal(err)
	} 

	fmt.Println("Session started...")
}

func (db *Db) GetSession() *mgo.Session {

	if db.Session == nil {
		fmt.Println("Session doesn't exist...")
		db.StartConnection()
	}

	return db.Session.Clone()
}

func (db *Db)CloseSession(){
	db.Session.Close()
}