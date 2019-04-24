package collection

import (
	mgo "gopkg.in/mgo.v2"
	"log"
)

type Collection struct {
	Db 		*mgo.Database
}

func NewCollection(s *mgo.Session) Collection {
	return Collection {
		Db: s.DB("project"),
	}
}


func (c *Collection) IsEmpty(collection string) (isEmpty bool) {

	isEmpty = true
	noRegisters, err := c.Db.C(collection).Count()

	if err != nil {
		log.Print("error counting no. of documents")
	}

	if noRegisters > 0 {
		isEmpty = false
	}
	return
}

