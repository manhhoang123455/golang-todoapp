package controller

import (
	"github.com/gin-gonic/gin"
	"mini_project/models"
	"mini_project/request"
	"mini_project/services"
	"mini_project/utils/responses"
	"net/http"
	"strings"
	"time"
)

type UserController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type userController struct {
	// inject user service
	userService services.UserService

	// inject jwt service
	jwtService services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		// inject user service
		userService: userService,

		// inject jwt service
		jwtService: jwtService,
	}
}

func ParseToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "No token found"
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "Bearer token not in proper format"
	}

	authHeader = strings.TrimSpace(splitToken[1])

	return authHeader
}

// Login is a function for user login
// @Summary "User Login"
// @Tags	Auth
// @Version 1.0
// @Produce application/json
// @Param	* body request.LoginRequest true "User Login"
// @Success 200 object responses.Response{errors=string,data=string} "Login successfully"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 401 object responses.Response{errors=string,data=string} "Failed to process request"
// @Failure 500 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/login [post]
func (h *userController) Login(c *gin.Context) {
	// create new instance of LoginRequest
	var input request.LoginRequest
	jwtTTL := services.GetTokenTTL()

	// bind the input with the request body
	err := c.ShouldBindJSON(&input)
	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		// response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), responses.EmptyObject{})
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Check if the email and password is valid
	loginResult := h.userService.VerifyCredential(input.Email, input.Password)
	if v, ok := loginResult.(models.User); ok {
		generatedToken := h.jwtService.GenerateToken(v.ID, time.Now().Add(time.Duration(jwtTTL)*time.Second))
		if len(generatedToken) < 1 {
			response := responses.ErrorsResponseByCode(http.StatusInternalServerError, "Failed to process request", responses.SignatureFailed, nil)
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}

		v.Token = generatedToken
		response := responses.SuccessResponse(http.StatusOK, "Login successfully", v)
		c.JSON(http.StatusOK, response)
		return
	}

	// If the email and password is not valid
	response := responses.ErrorsResponseByCode(http.StatusUnauthorized, "Failed to process request", responses.InvalidCredential, nil)
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	return
}

// Register is a function for user register
// @Summary "User Register"
// @Tags	Auth
// @Version 1.0
// @Produce application/json
// @Param	* body request.RegisterRequest true "User Register"
// @Success 201 object responses.Response{errors=string,data=string} "Register Success"
// @Failure 400 object responses.Response{errors=string,data=string} "Failed to process request"
// @Router	/auth/register [post]
func (h *userController) Register(c *gin.Context) {
	// create new instance of RegisterRequest
	var input request.RegisterRequest

	// bind the register with the request body
	err := c.ShouldBind(&input)
	// Check if there is any error in binding
	if err != nil {
		response := responses.ErrorsResponse(http.StatusBadRequest, "Failed to process request", err.Error(), nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// if the email is valid and unique in the database then register the user
	// create new user
	createdUser := h.userService.CreateUser(input)
	// Check if then email exists
	if createdUser.ID == 0 {
		response := responses.ErrorsResponseByCode(http.StatusBadRequest, "Failed to process request", responses.EmailAlreadyExists, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// response with the user data and token
	response := responses.SuccessResponse(http.StatusCreated, "Register Success", createdUser)
	// return the response
	c.JSON(http.StatusCreated, response)
	return
}
