package main

import (
	"FirstProject/model/database"
	"FirstProject/model/api"
)

var (
	Api				api.Api
	Db				database.Db
)

func main(){
	Db.StartConnection()
	Api.Start(Db.Session)
}