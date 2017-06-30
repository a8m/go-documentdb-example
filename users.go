package main

import (
	. "fmt"
	"github.com/a8m/documentdb-go"
)

// User document
type User struct {
	documentdb.Document
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// DB interface
type DB interface {
	Get(id string) *User
	GetAll() []*User
	Add(u *User) *User
	Update(u *User) *User
	Remove(id string) error
}

// UsersDB implement DB interface
type UsersDB struct {
	Database   string
	Collection string
	db         *documentdb.Database
	coll       *documentdb.Collection
	client     *documentdb.DocumentDB
}

// Return new UserDB
// Test if database and collection exist. if not, create them.
func NewDB(db, coll string, config *Config) (users UsersDB) {
	users.Database = db
	users.Collection = coll
	users.client = documentdb.New(config.Url, documentdb.Config{config.MasterKey})
	// Find or create `test` db and `users` collection
	if err := users.findOrDatabase(db); err != nil {
		panic(err)
	}
	if err := users.findOrCreateCollection(coll); err != nil {
		panic(err)
	}
	return
}

// Get user by given id
func (u *UsersDB) Get(id string) (user User, err error) {
	var users []User
	query := Sprintf("SELECT * FROM ROOT r WHERE r.id='%s'", id)
	if err = u.client.QueryDocuments(u.coll.Self, query, &users); err != nil || len(users) == 0 {
		return user, err
	}
	user = users[0]
	return user, nil
}

// Get all users
func (u *UsersDB) GetAll() (users []User, err error) {
	err = u.client.ReadDocuments(u.coll.Self, &users)
	return
}

// Create user
func (u *UsersDB) Add(user *User) (err error) {
	return u.client.CreateDocument(u.coll.Self, user)
}

// Update user by id
func (u *UsersDB) Update(id string, user *User) (err error) {
	var users []User
	if err = u.client.QueryDocuments(u.coll.Self,
		Sprintf("SELECT * FROM ROOT r WHERE r.id='%s'", id), &users); err != nil || len(users) == 0 {
		return
	}
	nUser := users[0]
	user.Id = nUser.Id
	if err = u.client.ReplaceDocument(nUser.Self, &user); err != nil {
		return
	}
	return
}

// Remove user by id
func (u *UsersDB) Remove(id string) (err error) {
	var users []User
	if err = u.client.QueryDocuments(u.coll.Self,
		Sprintf("SELECT * FROM ROOT r WHERE r.id='%s'", id), &users); err != nil || len(users) == 0 {
		return
	}
	nUser := users[0]
	if err = u.client.DeleteDocument(nUser.Self); err != nil {
		return
	}
	return
}

// Find or create collection by id
func (u *UsersDB) findOrCreateCollection(name string) (err error) {
	if colls, err := u.client.QueryCollections(u.db.Self, Sprintf("SELECT * FROM ROOT r WHERE r.id='%s'", name)); err != nil {
		return err
	} else if len(colls) == 0 {
		if coll, err := u.client.CreateCollection(u.db.Self, Sprintf(`{ "id": "%s" }`, name)); err != nil {
			return err
		} else {
			u.coll = coll
		}
	} else {
		u.coll = &colls[0]
	}
	return
}

// Find or create database by id
func (u *UsersDB) findOrDatabase(name string) (err error) {
	if dbs, err := u.client.QueryDatabases(Sprintf("SELECT * FROM ROOT r WHERE r.id='%s'", name)); err != nil {
		return err
	} else if len(dbs) == 0 {
		if db, err := u.client.CreateDatabase(Sprintf(`{ "id": "%s" }`, name)); err != nil {
			return err
		} else {
			u.db = db
		}
	} else {
		u.db = &dbs[0]
	}
	return
}
