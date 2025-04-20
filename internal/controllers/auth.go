package controllers

import (
	"Ecadr/internal/app/models"
	"Ecadr/internal/app/service"
	"Ecadr/pkg/errs"
	"Ecadr/pkg/logger"
	utils2 "Ecadr/pkg/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

const emptyInt = 0

// SignUp godoc
// @Summary Register a new user
// @Description This endpoint registers a new user with a username, email, and password.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.UserRequest true "User information"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/sign-up [post]
func SignUp(c *gin.Context) {
	var user models.UserRequest

	if err := c.BindJSON(&user); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if user.Password == "" {
		HandleError(c, errs.ErrPasswordIsEmpty)
		return
	}

	if user.Email == "" {
		HandleError(c, errs.ErrEmailIsEmpty)
		return
	}

	if user.Username == "" {
		HandleError(c, errs.ErrUsernameIsEmpty)
		return
	}

	if user.RoleID == emptyInt {
		HandleError(c, errs.ErrRoleIsRequired)
		return
	}

	if user.RoleID >= 4 {
		HandleError(c, errs.ErrWrongRoleID)
		return
	}

	userID, err := service.CreateUser(models.User{
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    int(user.RoleID),
	})
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			HandleError(c, errs.ErrIncorrectUsernameOrPassword)
			return
		}
		HandleError(c, err)
		return
	}

	accessToken, refreshToken, err := utils2.GenerateToken(userID, user.RoleID, user.Username)
	if err != nil {
		logger.Error.Printf("Error generating access token: %s", err)
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       userID,
	})
}

// SignIn godoc
// @Summary User login
// @Description This endpoint logs in an existing user using their username, email, and password.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.UserLogin true "User login information"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/sign-in [post]
func SignIn(c *gin.Context) {
	var user models.UserLogin

	if err := c.BindJSON(&user); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	if user.Username == "" {
		HandleError(c, errs.ErrUsernameIsEmpty)
		return
	}

	if user.Password == "" {
		HandleError(c, errs.ErrPasswordIsEmpty)
		return
	}

	user.Password = utils2.GenerateHash(user.Password)

	signedUser, accessToken, refreshToken, err := service.SignIn(user.Username, user.Password)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			HandleError(c, errs.ErrIncorrectUsernameOrPassword)
			return
		}

		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       signedUser.ID,
	})
}

// RefreshToken godoc
// @Summary Refresh Token
// @Description This endpoint refreshes the access token.
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body models.RefreshRequest true "User login information"
// @Success 200 {object} models.RefreshTokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		HandleError(c, errs.ErrValidationFailed)
		return
	}

	// Проверка валидности refresh_token
	token, err := jwt.ParseWithClaims(requestBody.RefreshToken, &utils2.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		HandleError(c, errs.ErrInvalidToken)
		return
	}

	// Генерация нового access_token
	claims, ok := token.Claims.(*utils2.CustomClaims)
	if !ok || claims.ExpiresAt < time.Now().Unix() {
		HandleError(c, errs.ErrRefreshTokenExpired)
		return
	}

	accessToken, refreshToken, err := utils2.GenerateToken(claims.RoleID, claims.UserID, claims.Username)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
