package corkboardauth

import (
	"encoding/json"
	"errors"
)

var (
	errNoUserFound = errors.New("No user found")
)

//TODO: modify this struct to match Joel's user table

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	//Sites    []string `json:"sites"`
}

//FakeUser is used during marshaling so that I can add a field to the struct
//type FakeUser User

//TODO: rewrite MarshalJSON to comply with what we need for new struct and bindings
//Do we even need this method with a MySQL implementation?

//MarshalJSON turns the user struct into the appropriate json format
func (user *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
	// 	Type string `json:"_type"`
	// 	//FakeUser
	// }{
	// 	Type:     "User",
		//FakeUser: FakeUser(*user),
	})
}

// func getUserKey(id uuid.UUID) string {
// 	return fmt.Sprintf("user:%s", id.String())
// }

//TODO: rewrite this function
func (cba *CorkboardAuth) addUser(user *User) error {
	// newID := uuid.NewV4()
	// user.ID = newID.String()
	// _, err := cba.bucket.Insert(getUserKey(newID), user, 0)
	// return err
}

//TODO: rewrite this function
func (cba *CorkboardAuth) updateUser(user *User) error {
	// id, err := uuid.FromString(user.ID)
	// if err != nil {
	// 	return err
	// }
	// _, err = cba.bucket.Upsert(getUserKey(id), user, 0)
	// return err
}

//TODO: rewrite this method
func (cba *CorkboardAuth) findUser(userEmail string) (*User, error) {
	// query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT `email`, `password`, `sites`, `id` FROM `%s` WHERE _type = 'User' AND email = $1 LIMIT 1", cba.bucket.Name())).AdHoc(true) //nolint: gas
	// res, err := cba.bucket.ExecuteN1qlQuery(query, []interface{}{userEmail})
	// if err != nil {
	// 	return nil, err
	// }
	// defer res.Close() //nolint: errcheck
	// user := new(User)
	// for res.Next(user) {
	// 	return user, nil
	// }
	// return nil, errNoUserFound
}

//Don't think we need this method anymore
// func (cba *CorkboardAuth) findUserFromSite(userEmail string, siteID uuid.UUID) (*User, error) {
// 	query := gocb.NewN1qlQuery(fmt.Sprintf("SELECT `email`, `password`, `sites`, `id` FROM `%s` WHERE _type = 'User' AND email = $1 AND ANY site IN sites SATISFIES site = $2 END LIMIT 1", cba.bucket.Name())).AdHoc(true) //nolint: gas
// 	res, err := cba.bucket.ExecuteN1qlQuery(query, []interface{}{userEmail, siteID.String()})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Close() //nolint: errcheck
// 	user := new(User)
// 	for res.Next(user) {
// 		return user, nil
// 	}
// 	return nil, errNoUserFound
// }
