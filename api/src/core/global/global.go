package global

import (
	"dse/src/core/models"
	"sync"
)

// ------------------------------------------------------------
// : Globals
// ------------------------------------------------------------
const (
	VERSION = "3.0.6"
)

var (
	mutex = &sync.Mutex{}
	users = map[string]*models.User{}
)

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func AddUser(user *models.User) {
	mutex.Lock()
	defer mutex.Unlock()
	users[user.Token] = user
}

func RemoveUser(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(users, token)
}

func GetUser(token string) *models.User {
	mutex.Lock()
	defer mutex.Unlock()
	return users[token]
}

func GetUsers() map[string]*models.User {
	mutex.Lock()
	defer mutex.Unlock()
	return users
}

func HasUser(token string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := users[token]
	return ok
}
