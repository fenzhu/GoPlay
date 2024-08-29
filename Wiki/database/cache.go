package database

import "sync"

type Cache struct {
	Data *sync.Map
}
