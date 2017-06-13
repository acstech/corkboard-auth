package corkboardauth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/couchbase/gocb"
	uuid "github.com/satori/go.uuid"
)

var (
	errNoUserFound = errors.New("No user found")
)

//User is how a user is stored in couchbase
type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Sites    []string `json:"sites"`
}

//FakeUser is used during marshaling so that I can add a field to the struct
type FakeUser User

//MarshalJSON turns the user struct into the appropriate json format
func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type string `json:"_type"`
		FakeUser
	}{
		Type:     "User",
		FakeUser: FakeUser(*user),
	})
}

func getUserKey(id uuid.UUID) string {
	return fmt.Sprintf("user:%s", id.String())
}

func (cba *CorkboardAuth) addUser(user *User) error {
	newID := uuid.NewV4()
	user.ID = newID.String()
	_, err := cba.bucket.Insert(getUserKey(newID), user, 0)
	return err
}

func (cba *CorkboardAuth) updateUser(user *User) error {
	id, err := uuid.FromString(user.ID)
	if err != nil {
		return err
	}
	_, err = cba.bucket.Upsert(getUserKey(id), user, 0)
	if err != nil {
		return err
	}
	return nil
}

func (cba *CorkboardAuth) findUser(userEmail string) (*User, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT `email`, `password`, `sites`, `id` FROM `%s` WHERE _type = 'User' AND email = $1 LIMIT 1", cba.bucket.Name())).AdHoc(true)
	res, err := cba.bucket.ExecuteN1qlQuery(query, []interface{}{userEmail})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	user := new(User)
	for res.Next(user) {
		return user, nil
	}
	return nil, errNoUserFound
}

func (cba *CorkboardAuth) findUserFromSite(userEmail string, siteID uuid.UUID) (*User, error) {
	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT `email`, `password`, `sites`, `id` FROM `%s` WHERE _type = 'User' AND email = $1 AND ANY site IN sites SATISFIES site = $2 END LIMIT 1", cba.bucket.Name())).AdHoc(true)
	res, err := cba.bucket.ExecuteN1qlQuery(query, []interface{}{userEmail, siteID.String()})
	if err != nil {
		return nil, err
	}
	defer res.Close()
	user := new(User)
	for res.Next(user) {
		return user, nil
	}
	return nil, errNoUserFound
}
