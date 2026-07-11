package handlers

import (
	"errors"
	"net/http"

	"gotask-api/constants"
	"gotask-api/datatransfers"
	"gotask-api/domains"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	domain domains.UserDomain
}

func NewUserHandler(domain domains.UserDomain) *UserHandler {
	return &UserHandler{domain: domain}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req datatransfers.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		datatransfers.ErrorRes(c, http.StatusBadRequest, constants.BAD_REQUEST, err.Error())
		return
	}

	user, err := h.domain.Register(req)
	if err != nil {
		if errors.Is(err, domains.ErrEmailAlreadyExists) {
			datatransfers.ErrorRes(c, http.StatusConflict, "EMAIL_EXISTS", err.Error())
			return
		}
		datatransfers.ErrorRes(c, http.StatusInternalServerError, constants.INTERNAL_ERROR, err.Error())
		return
	}

	datatransfers.SuccessRes(c, http.StatusCreated, constants.SUCCESS, user)
}
