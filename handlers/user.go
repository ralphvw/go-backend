package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Hash      []byte
}

func CreateUser(c *gin.Context, db *sql.DB) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plaintextPassword := user.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {

		log.Fatal(err)
		return
	}

	user.Hash = hashedPassword
	query := "INSERT INTO users (first_name, last_name, email, hash) VALUES($1, $2, $3, $4) RETURNING id"

	var id int
	err = db.QueryRow(query, user.FirstName, user.LastName, user.Email, user.Hash).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}

	user.ID = id

	c.JSON(http.StatusCreated, gin.H{"message": "Entity created successfully"})
}

func GetAllUsers(c *gin.Context, db *sql.DB) {
	query := "SELECT id, first_name, FROM users"
	countQuery := "SELECT COUNT(*) FROM users"
	var id int
	var name string
	var age int
	var args []interface{}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	GetDataHandler(c, db, pageSize, page, query, countQuery, args, &id, &name, &age)
}
