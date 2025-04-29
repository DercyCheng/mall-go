package assembler

import (
	"time"

	"mall-go/services/user-service/application/dto"
	"mall-go/services/user-service/domain/model"
)

// ToUserDTO converts a domain User model to a UserDTO
func ToUserDTO(user *model.User) dto.UserDTO {
	roleIds := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleIds[i] = role.ID
	}

	return dto.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		NickName:  user.NickName,
		Phone:     user.Phone,
		Icon:      user.Icon,
		Status:    int(user.Status),
		Note:      user.Note,
		RoleIds:   roleIds,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		LastLogin: user.LastLogin.Format(time.RFC3339),
	}
}

// ToUserDTOList converts a list of domain User models to a list of UserDTOs
func ToUserDTOList(users []*model.User) []dto.UserDTO {
	dtos := make([]dto.UserDTO, len(users))
	for i, user := range users {
		dtos[i] = ToUserDTO(user)
	}
	return dtos
}

// ToRoleDTO converts a domain Role model to a RoleDTO
func ToRoleDTO(role model.Role) dto.RoleDTO {
	return dto.RoleDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format(time.RFC3339),
	}
}

// ToRoleDTOList converts a list of domain Role models to a list of RoleDTOs
func ToRoleDTOList(roles []*model.Role) []dto.RoleDTO {
	dtos := make([]dto.RoleDTO, len(roles))
	for i, role := range roles {
		dtos[i] = ToRoleDTO(*role)
	}
	return dtos
}

// ToUserModel converts a UserCreateRequest to a domain User model
func ToUserModel(req dto.UserCreateRequest) (*model.User, error) {
	return model.NewUser(req.Username, req.Password, req.Email, req.NickName)
}

// ToRoleModel converts a RoleCreateRequest to a domain Role model
func ToRoleModel(req dto.RoleCreateRequest) model.Role {
	now := time.Now()
	return model.Role{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
