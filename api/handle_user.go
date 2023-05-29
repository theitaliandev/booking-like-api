package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/theitaliandev/booking-like-api/store"
	"github.com/theitaliandev/booking-like-api/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userStore store.UserStore
}

func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	user, err := h.userStore.GetUserByID(objID)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers()
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var userParams types.CreateUserParams
	err := c.BodyParser(&userParams)
	if err != nil {
		return err
	}
	newUser, validationErrors, err := types.NewUserFromParams(&userParams)
	if err != nil {
		return err
	}
	if len(validationErrors) > 0 {
		errors := map[string]string{}
		for key, val := range validationErrors {
			errors[key] = val.Error()
		}
		return c.JSON(errors)
	}
	user, err := h.userStore.CreateUser(newUser)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	err = h.userStore.DeleteUser(objID)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"id": id})
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var updateUserParams *types.UpdateUserParams
	err = c.BodyParser(&updateUserParams)
	if err != nil {
		return err
	}
	validationErrors := updateUserParams.ValidateUpdateUser()
	if len(validationErrors) > 0 {
		errors := map[string]string{}
		for key, val := range validationErrors {
			errors[key] = val.Error()
		}
		return c.JSON(errors)
	}
	updatedUser, err := h.userStore.UpdateUser(objID, updateUserParams)
	if err != nil {
		return err
	}
	return c.JSON(updatedUser)
}
