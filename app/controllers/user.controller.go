package controllers

import (
	"github.com/gofiber/fiber/v2"

	"note-api/app/dto"
	"note-api/app/services"
	"note-api/core/utils"

)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctrl *UserController) Create(c *fiber.Ctx) error {
	var user dto.UserRequest
	utils.MustBindAndValidate(c, &user, user.Validate)
	createdUser, err := ctrl.service.CreateUser(&user)
	if err != nil {
		panic(fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}
	return utils.GlobalResponse(c, fiber.StatusCreated, "User created successfully", createdUser)
}
