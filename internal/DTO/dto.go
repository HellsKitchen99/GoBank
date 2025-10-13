package dto

import "GoBank/internal/domain"

func ConvertFromUserToUserDetails(user *domain.User) *domain.UserDetails {
	var userDetails domain.UserDetails = domain.UserDetails{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Roles: user.Roles,
	}
	return &userDetails
}
