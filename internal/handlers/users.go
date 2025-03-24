package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api-gateway-golang/internal/model"
	"api-gateway-golang/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetUsers godoc
//
//	@Summary		Get all users
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object} []model.GetUserResponse
//
// @Router			/api/users [get]
func (h *Handlers) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.Storage.GetUsers(ctx)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	userResponse := []model.GetUserResponse{}
	for _, user := range users {
		userResponse = append(userResponse, model.GetUserResponse{
			ID:          user.ID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
		})

	}

	err = h.Sender.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("Error when requesting /users")

		panic(err)
	}
}

// GetUser godoc
//
//	@Summary		Get a specific user
//	@Tags			Users
//	@Produce		json
//	@Param			id path int	true "User ID"
//	@Success		200	{object} model.GetUserResponse
//
// @Router			/api/user/{id} [get]
func (h *Handlers) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Lấy user từ storage
	user, err := h.Storage.GetUser(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Kiểm tra nếu user không tồn tại
	if user.ID == 0 {
		h.Sender.JSON(w, http.StatusBadRequest, fmt.Sprintf("User with id=%d not found", id))
		return
	}

	// Xây dựng response cho user
	userResponse := model.GetUserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	// Trả về response
	err = h.Sender.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		logger.OutputLog.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal(fmt.Sprintf("Error when requesting /user/%d", user.ID))
		panic(err)
	}
}

// AddUser godoc
//
//		@Summary		Add a specific user
//		@Tags			Users
//		@Produce		json
//	 	@Accept			json
//		@Param			firstName body string true "User firstName"
//		@Param			lastName body string true "User lastName"
//		@Param			email body string true "User email"
//		@Param			phoneNumber body string true "User phoneNumber"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/user/add [post]
func (h *Handlers) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.AddUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(user)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	id, err := h.Storage.AddUser(ctx, user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// UpdateUser godoc
//
//		@Summary		Update a specific user
//		@Tags			Users
//		@Produce		json
//	 	@Accept			json
//		@Param			firstName body string false "User firstName"
//		@Param			lastName body string false "User lastName"
//		@Param			email body string false "User email"
//		@Param			phoneNumber body string false "User phoneNumber"
//		@Success		200	{object} model.IDResponse
//
// @Router			/api/user/update [patch]
func (h *Handlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var user model.UpdateUserRequest

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = Validate.Struct(user)
	if err != nil {
		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, err.Field()+" "+err.Tag())
		}
		h.Sender.JSON(w, http.StatusBadRequest, strings.Join(errs, ", "))
		return
	}

	exists, err := h.Storage.VerifyUserExists(ctx, user.ID)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		h.Sender.JSON(w, http.StatusBadRequest, "User with id="+fmt.Sprint(user.ID)+" not found")
		return
	}

	id, err := h.Storage.UpdateUser(ctx, user)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := model.IDResponse{ID: id}
	err = h.Sender.JSON(w, http.StatusOK, response)
	if err != nil {
		panic(err)
	}
}

// DeleteUser godoc
//
//	@Summary		Delete a specific user
//	@Tags			Users
//	@Produce		json
//	@Param			id path int	true "User ID"
//	@Success		200	{object} model.GetUserResponse
//
// @Router			/api/user/delete/{id} [delete]
func (h *Handlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	exists, err := h.Storage.VerifyUserExists(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		err = h.Sender.JSON(w, http.StatusBadRequest, "User with id="+fmt.Sprint(id)+" not found")
		if err != nil {
			panic(err)
		}
		return
	}

	err = h.Storage.DeleteUser(ctx, id)
	if err != nil {
		h.Sender.JSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.Sender.JSON(w, http.StatusOK, map[string]bool{"success": true})
	if err != nil {
		panic(err)
	}
}

// cSpell:ignore godoc logrus
