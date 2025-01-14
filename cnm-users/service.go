package users

import (
	"errors"
	e "shopping-cart/cnm-users/entities"
	 "shopping-cart/cnm-users/repository"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

// Service is the interface that provides booking methods.
type  Service interface {
	add(request userRequest) (userResponse , error)
	list() (userListResponse, error)
}

type service struct {
}


func (s *service) add(request userRequest) (userResponse, error) {
	instance:= repository.UsersRepositoryInstance()
	usr, err:=instance.Add(e.UserDetails{
		FirstName: request.FirstName,
		MiddleName: request.MiddleName,
		LastName: request.LastName,
		UserEmail: request.UserEmail,
		Password: request.Password,
		PhoneNumber: request.PhoneNumber,
			})

	return userResponse{UserId: usr.UserId}, err
}

func (s *service) list() (userListResponse, error) {
	instance:= repository.UsersRepositoryInstance()
	users, err:=instance.List();
	return userListResponse{UserDetails:users , Err: err}, err
}

// NewService creates a auth service with necessary dependencies.
func NewService() Service {
	return &service{}
}

