package main

import (
	"github.com/martini-contrib/sessionauth"
)

// MyUserModel can be any struct that represents a user in my system
type PosterUserModel struct {
	Id            int64  `form:"id" db:"id"`
	Username      string `form:"name" db:"username"`
	Password      string `form:"password" db:"password"`
	authenticated bool   `form:"-" db:"-"`
}

// GetAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() sessionauth.User {
	return &PosterUserModel{}
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *PosterUserModel) Login() {
	// Update last login time
	// Add to logged-in user's list
	// etc ...
	u.authenticated = true
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *PosterUserModel) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.authenticated = false
}

func (u *PosterUserModel) IsAuthenticated() bool {
	return u.authenticated
}

func (u *PosterUserModel) UniqueId() interface{} {
	return u.Id
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *PosterUserModel) GetById(id interface{}) error {
	err := dbmap.SelectOne(u, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
