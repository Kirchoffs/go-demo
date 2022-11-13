package controllers

import (
    "mongodb-api-demo/models"
    "mongodb-api-demo/services"
    "net/http"

    "github.com/gin-gonic/gin"
)

type UserController struct {
    UserService services.UserService
}

func New(userService services.UserService) UserController {
    return UserController{
        UserService: userService,
    }
}

func (userController *UserController) CreateUser(ctx *gin.Context) {
    var user models.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    if err := userController.UserService.CreateUser(&user); err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (userController *UserController) GetUser(ctx *gin.Context) {
    username := ctx.Param("name")
    user, err := userController.UserService.GetUser(&username)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, user)
}

func (userController *UserController) GetAll(ctx *gin.Context) {
    users, err := userController.UserService.GetAll()
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, users)
}

func (userController *UserController) UpdateUser(ctx *gin.Context) {
    var user models.User
    if err := ctx.ShouldBindJSON(&user); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    if err := userController.UserService.UpdateUser(&user); err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (userController *UserController) DeleteUser(ctx *gin.Context) {
    username := ctx.Param("name")
    if err := userController.UserService.DeleteUser(&username); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (userController *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
    userRoute := rg.Group("/user")
    userRoute.POST("/create", userController.CreateUser)
    userRoute.GET("/get/:name", userController.GetUser)
    userRoute.GET("/get-all", userController.GetAll)
    userRoute.PATCH("/update", userController.UpdateUser)
    userRoute.DELETE("/delete/:name", userController.DeleteUser)
}
