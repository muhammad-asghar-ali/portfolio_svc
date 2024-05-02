package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/0xbase-Corp/portfolio_svc/internal/models"
	"github.com/0xbase-Corp/portfolio_svc/shared/errors"
	"github.com/0xbase-Corp/portfolio_svc/shared/utils"
)

type (
	GenerateHashRequest struct {
		PublicKey string `json:"public_key" binding:"required"`
	}

	GenerateHashResonse struct {
		Type            string `json:"type"`
		HashedPublicKey string `json:"hashed_public_key"`
	}

	VerifyRequest struct {
		PublicKey       string `json:"public_key" binding:"required"`
		HashedPublicKey string `json:"hashed_public_key" binding:"required"`
	}

	TokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)

//	@BasePath	/api/v1

// AuthGenerateHash godoc
//
// @Summary      Generate the hash for authorization
// @Description  Generate the hash for authorization
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body GenerateHashRequest true  "GenerateHashRequest object"
// @Success      200 {object} GenerateHashResonse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /generate-hash [post]
func AuthGenerateHash(c *gin.Context, db *gorm.DB) {
	request := &GenerateHashRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}

	var response GenerateHashResonse
	// check in user table to find the record with this key is exist pass other wise create
	exist, _ := models.GetUserByPublicKey(db, request.PublicKey)

	if exist != nil {
		response = GenerateHashResonse{
			Type:            "key", // decide
			HashedPublicKey: exist.HashedPublicKey,
		}
	} else {
		user := &models.User{
			PublicKey:       request.PublicKey,
			Email:           "test@local.com",       // TODO: make it optional or not
			Username:        "xbase",                // TODO: make it optional or not or generate a random user name
			HashedPublicKey: utils.RandomString(64), // generate the HashedPublicKey
		}

		err := models.CreateUser(db, user)
		if err != nil {
			errors.HandleHttpError(c, errors.NewBadRequestError("Error while creating user: "+err.Error()))
			return
		}

		response = GenerateHashResonse{
			Type:            "key", // TODO: decide
			HashedPublicKey: user.HashedPublicKey,
		}
	}

	c.JSON(http.StatusOK, response)
}

//	@BasePath	/api/v1

// AuthVerifyHashKey godoc
//
// @Summary      Verify the hash for authorization
// @Description  Verify the hash for authorization
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body VerifyRequest true  "VerifyRequest object"
// @Success      200 {object} TokenResponse
// @Failure      400 {object} errors.APIError
// @Failure      404 {object} errors.APIError
// @Failure      500 {object} errors.APIError
// @Router       /verify-hash [post]
func AuthVerifyHashKey(c *gin.Context, db *gorm.DB) {
	request := &VerifyRequest{}

	if err := c.ShouldBindJSON(request); err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}

	user, _ := models.GetUserByPublicKey(db, request.PublicKey)
	if user == nil {
		errors.HandleHttpError(c, errors.NewNotFoundError("User not found"))
		return
	}

	matches := compareHash(request.HashedPublicKey, user.HashedPublicKey)
	if !matches {
		errors.HandleHttpError(c, errors.NewForbiddenError("Invalid hash value"))
		return
	}

	// TODO: may update the hashed value in database

	// generate jwt token
	accessToken, err := utils.GenerateAccessToken(user.UserId, user.Email, user.PublicKey)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.UserId, user.Email, user.PublicKey)
	if err != nil {
		errors.HandleHttpError(c, errors.NewBadRequestError(err.Error()))
		return
	}

	resp := &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	// return the response
	c.JSON(http.StatusOK, resp)
}

func compareHash(input, hash string) bool {
	return input == hash
}
