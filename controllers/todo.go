package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todo-list-api/utils"
	_ "modernc.org/sqlite"
)

type todo struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	IsComplete bool   `json:"is_complete"`
}

type todoList struct {
	todoStore []todo
	db        *sql.DB
}

func Service(db *sql.DB) *todoList {
	return &todoList{
		todoStore: []todo{},
		db:        db,
	}
}

// -----------------------Create a new todo list.---------------------------------------------
func (t *todoList) CreateTodo(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.POST) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only POST method is allowed."})
		return
	}

	//read and store
	var newTodo todo
	err := json.NewDecoder(c.Request.Body).Decode(&newTodo)
	if err != nil {
		panic(err)
	}

	stmt, err := t.db.Prepare("INSERT INTO todo(title,body) VALUES(?,?);")
	if err != nil {
		panic(err)
	}
	row, err := stmt.Exec(newTodo.Title, newTodo.Body)
	if err != nil {
		panic(err)
	}
	id, err := row.LastInsertId()
	if err != nil {
		panic(err)
	}

	newTodo.Id = int(id)

}

// -------------------------------Get all todo.----------------------------------------------
func (t *todoList) GetTodo(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.GET) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only GET method is allowed."})
		return
	}

	rows, err := t.db.Query("SELECT * FROM todo;")
	if err != nil {
		panic(err)
	}
	var data = []todo{}
	for rows.Next() {
		var newTodo todo
		rows.Scan(&newTodo.Id, &newTodo.Title, &newTodo.Body, &newTodo.IsComplete)
		data = append(data, newTodo)

	}
	if len(data) > 0 {
		c.JSON(http.StatusOK, gin.H{"data": data})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "No todo found."})
	}
}

// 	rows, err := t.db.Query("SELECT * FROM todo")
// 	if err != nil {
// 		panic(err)
// 	}

// 	data := []todo{}

// 	for rows.Next() {
// 		var newTodo todo
// 		rows.Scan(&newTodo.Id, &newTodo.Title, &newTodo.Body, &newTodo.IsComplete)
// 		data = append(data, newTodo)
// 	}
// 	if len(data) > 0 {
// 		c.JSON(200, data)
// 	} else {
// 		c.JSON(200, nil)
// 	}
// }

// delete a todo.
// func (t *todoList) DeleteTodo(c *gin.Context) {
// 	id := c.Param("id")
// 	stmt, err := t.db.Prepare("DELETE FROM todo WHERE id =?")
// 	if err != nil {
// 		panic(err)
// 	}
// 	row, err2 := stmt.Exec(id)
// 	if err2 != nil {
// 		c.JSON(200, gin.H{
// 			"status": 200,
// 			"error":  err2.Error(),
// 		})

//		}
//	}
//
// ----------------------------------Delete a todo----------------------------------------------
func (t *todoList) DeleteTodo(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.DELETE) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only DELETE method is allowed."})
		return
	}
	//getting if from url
	id := c.Param("id")
	stmt, err := t.db.Prepare("DELETE FROM todo WHERE id =?;")
	if err != nil {
		panic(err)
	}

	rows, err := stmt.Exec(id)
	if err != nil {
		panic(err)
	}

	affectedRows, err := rows.RowsAffected()
	if err != nil {
		panic(err)
	}
	if affectedRows > 0 {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "todo deleted with id " + id,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "todo not found with id " + id,
		})
	}
}

// -----------------------------------Update the todo------------------------------------------
func (t *todoList) UpdateTodo(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.PATCH) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only PATCH method is allowed."})
		return
	}

	id := c.Param("id")
	stmt, err := t.db.Prepare("UPDATE todo SET title =?, body =? WHERE id =?;")
	if err != nil {
		panic(err)
	}

	var newTodo todo
	rows, err := stmt.Exec(&newTodo.Title, &newTodo.Body, id)
	if err != nil {
		panic(err)
	}

	affectedRows, err := rows.RowsAffected()
	if err != nil {
		panic(err)
	}
	if affectedRows > 0 {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "todo updated with id " + id,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "todo not found with id " + id,
		})
	}
}

//----------------------------------Get a single todo----------------------------------------

func (t *todoList) GetTodobyID(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.GET) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only GET method is allowed."})
		return
	}
	id := c.Param("id")

	rows := t.db.QueryRow("SELECT * FROM todo WHERE id=?;", id)

	var newTodo todo

	err := rows.Scan(&newTodo.Id, &newTodo.Title, &newTodo.Body, &newTodo.IsComplete)
	if err != nil {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "todo not found with id " + id,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  200,
			"blog":    newTodo,
			"message": "todo found with id " + id,
		})
	}
}

// ----------------------------------Change the status of a todo----------------------------------------
func (t *todoList) ChangeTodoStatus(c *gin.Context) {
	if !utils.Checkmethod(c.Request.Method, utils.PATCH) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Only PATCH method is allowed."})
		return
	}
	id := c.Param("id")
	stmt, err := t.db.Prepare("UPDATE todo SET isComplete =? WHERE id =?;")
	if err != nil {
		panic(err)
	}
	row, err := stmt.Exec(true, id)
	if err != nil {
		panic(err)
	}
	affectedRows, err := row.RowsAffected()
	if err != nil {
		panic(err)
	}
	if affectedRows > 0 {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Blog status updated for id " + id,
		})
	} else {
		c.JSON(200, gin.H{
			"status":  200,
			"message": "Todo not found with id " + id,
		})
	}

}
