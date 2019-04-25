package collection

import (
	mgo "gopkg.in/mgo.v2"
)

type Collection struct {
	Db 		*mgo.Database
}

func NewCollection(s *mgo.Session) Collection {
	return Collection {
		Db: s.DB("project"),
	}
}

