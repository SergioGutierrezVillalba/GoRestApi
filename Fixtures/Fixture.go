package fixture

import (

	"FirstProject/Collection"
	usercol "FirstProject/Collection/user"

	userfix "FirstProject/Fixtures/user"
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
	f.LoadUsersFixture()
	// and so on...
	
	
}

func (f *Fixture) LoadUsersFixture(){
	// 1. Collection
	UserCollection := usercol.NewUserCollection(f.Collection.Db)
	// 2. Fixture
	UserFixture := userfix.NewUserFixture(UserCollection)
	// 3. Load
	UserFixture.LoadFixture()
}