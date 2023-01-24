package main

import (
	"github.com/gin-gonic/gin"
	"github.com/todo-list-api/controllers"
	"github.com/todo-list-api/initializers"
)

func main() {

	db := initializers.InitializingDB()
	defer db.Close()

	newServ := controllers.Service(db)
	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.POST("/todo/create", newServ.CreateTodo)
	r.GET("/todo/getall", newServ.GetTodo)
	r.GET("/todo/getbyid/:id", newServ.GetTodobyID)
	r.DELETE("/todo/delete/:id", newServ.DeleteTodo)
	r.PATCH("/todo/update/:id", newServ.UpdateTodo)
	r.PATCH("/todo/update/statusbyid/:id", newServ.ChangeTodoStatus)

	r.Run(":3000")

}
