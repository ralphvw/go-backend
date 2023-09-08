package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type User struct {
	id int
	name string
	age int
}

func CreateUser(c *gin.Context, db *sql.DB) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	query := "INSERT INTO users (name, age) VALUES($1, $2) RETURNING id"

	var id int
	err := db.QueryRow(query, user.name, user.age).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
        return
	}

	user.id = id

	c.JSON(http.StatusCreated, gin.H{"message": "Entity created successfully"})
}

func GetAllUsers(c *gin.Context, db *sql.DB) {
	query := "SELECT id, name, age FROM users"
	countQuery := "SELECT COUNT(*) FROM users"
	var id int
	var name string
	var age int
	var args []interface{}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	GetDataHandler(c, db, pageSize, page, query, countQuery, args, &id, &name, &age)
}