package fixture

import (
	"FirstProject/Collection"
	"FirstProject/Fixtures/user"
)

type Fixture struct {
	Collection		collection.Collection
}

func NewFixture(c collection.Collection) *Fixture {
	return &Fixture {
		Collection: c,
	}
}

func (f *Fixture) LoadFixtures(){

	if f.Collection.IsEmpty("users") {
		UserFixture := user.UserFixture{}
		UserFixture.LoadUsersFixture()
	}
	
}