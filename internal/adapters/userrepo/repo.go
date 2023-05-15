package userrepo

import (
	"context"
	"fmt"
	"homework9/internal/app"
	"homework9/internal/users"
	"sync"
)

func New() app.UserRepository {
	return &userRepo{make(map[int64]users.User, 0), 0, sync.Mutex{}}
}

type userRepo struct {
	users map[int64]users.User
	idx   int64
	mutex sync.Mutex
}

func (r *userRepo) DeleteUser(ctx context.Context, ID int64) error {
	_, ok := r.users[ID]
	if !ok {
		return fmt.Errorf("user not found")
	}
	delete(r.users, ID)
	return nil
}

func (r *userRepo) CreateUser(ctx context.Context, Nickname string, Email string) (users.User, error) {
	for _, user := range r.users {
		if user.Email == Email {
			return users.User{}, fmt.Errorf("user already exists")
		}
	}
	r.mutex.Lock()
	newUser := users.User{ID: r.idx, Nickname: Nickname, Email: Email}
	r.users[r.idx] = newUser
	r.idx++
	r.mutex.Unlock()
	return newUser, nil
}

func (r *userRepo) GetUser(ctx context.Context, ID int64) (users.User, error) {
	user, ok := r.users[ID]
	if !ok {
		return users.User{}, fmt.Errorf("user not found")
	}
	return user, nil
}
func (r *userRepo) GetUsers(ctx context.Context) map[int64]users.User {
	return r.users
}
