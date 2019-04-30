package main

import (
	"FirstProject/Model/database"
	"FirstProject/Model/api"
)

var (
	Api				api.Api
	Db				database.Db
)

func main(){
	Db.StartConnection()
	Api.Start(Db.Session)
}

