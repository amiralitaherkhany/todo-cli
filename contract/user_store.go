package contract

import "github.com/amiralitaherkhany/todo-cli/entity"

type UserWriteStore interface {
	Save(u entity.User)
}

type UserReadStore interface {
	Load() []entity.User
}
