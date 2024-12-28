package authentication

import (
	"context"
	"example.com/go-web-base/internal/application"
	"github.com/google/uuid"
)

type contextKey int

const userContextKey contextKey = iota

func WithUserContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) (User, bool) {
	user := ctx.Value(userContextKey)
	if user == nil {
		return User{}, false
	}

	user, ok := user.(User)
	if !ok {
		return User{}, false
	}

	return user.(User), true
}

type User struct {
	Id    string
	Name  string
	Email string
}

type NewUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUser(ctx context.Context, app application.Application, request NewUserRequest) (User, error) {
	hashedPassword, err := hashPassword(request.Password)
	if err != nil {
		app.Log.Error(ctx, "unable to hash users password: "+err.Error())
		return User{}, err
	}

	userId, err := uuid.NewV7()
	if err != nil {
		app.Log.Error(ctx, "unable to generate uuid for user: "+err.Error())
		return User{}, err
	}

	_, err = app.DB.Exec("INSERT INTO users (id, name, email, password) VALUES ($1, $2, $3, $4)", userId.String(), request.Name, request.Email, hashedPassword)
	if err != nil {
		app.Log.Error(ctx, "unable to save user: "+err.Error())
		return User{}, err
	}

	return User{
		Id:    userId.String(),
		Name:  request.Name,
		Email: request.Email,
	}, nil
}
