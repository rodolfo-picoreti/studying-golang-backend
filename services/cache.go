package services

import (
	"time"

	"github.com/gin-contrib/cache/persistence"
)

var store *persistence.InMemoryStore

func GetCacheStore() *persistence.InMemoryStore {
	if store == nil {
		store = persistence.NewInMemoryStore(time.Second)
	}
	return store
}
