package usersservice

import (
	"bank_app/internal/storage/cache"
	"bank_app/internal/storage/repos/users"
)

type UsersService struct {
	usersRepo    *users.Repo
	cacheService cache.CacheInterface
}

func NewUsersService(usersRepo *users.Repo, cacheService cache.CacheInterface) *UsersService {
	return &UsersService{
		usersRepo:    usersRepo,
		cacheService: cacheService,
	}
}
