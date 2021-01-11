package router

import (
	"os"
	"gin-todo/controllers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(dbConn *gorm.DB, UserDB *gorm.DB) {
	todoHandler := controllers.TodoHandler{
		Db:     dbConn,
		UserDb: UserDB,
	}
	r := gin.Default()
	//セッションの設定
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.GET("/", todoHandler.AuthTest)
	r.POST("login/:account", todoHandler.Login)
	r.GET("/register", todoHandler.Register)
	r.POST("/register/user", todoHandler.RegisterPOST)
	r.POST("/logout", todoHandler.Logout)
	r.GET("/todo", todoHandler.GetAll)
	r.GET("/alluser", todoHandler.GetAllUser)
	r.POST("/todo/create", todoHandler.CreateTask)
	r.GET("/todo/:id", todoHandler.EditTask)
	r.POST("/todo/edit/:id", todoHandler.UpdateTask)
	r.POST("/todo/delete/:id", todoHandler.DeleteTask)
	r.POST("/user/delete/:id", todoHandler.DeleteUser)

	port := os.Getenv("PORT")
	if port != ""{
		r.Run(":" + port)
	}else{
		r.Run(":8080")
	}
}
