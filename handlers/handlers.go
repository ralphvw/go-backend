package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetDataHandler(c *gin.Context, db *sql.DB, pageSize, page int, query string, countQuery string, args []interface{}, destinations ...interface{}) {
	offset := (page - 1) * pageSize

	queryString := fmt.Sprintf("%s LIMIT %d OFFSET %d", query, pageSize, offset)

	rows, err := db.Query(queryString, args...)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		if err := rows.Scan(destinations...); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		result := make(map[string]interface{})
		for i := 0; i < len(destinations); i += 2 {
			key, ok := destinations[i].(string)
			if !ok {
				log.Fatal("Invalid key type")
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}
			value := destinations[i+1]
			result[key] = value
		}

		results = append(results, result)
	}

	totalRows := 0
	err = db.QueryRow(countQuery).Scan(&totalRows)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	totalPages := (totalRows + pageSize - 1) / pageSize

	responseData := gin.H{
		"data":       results,
		"totalPages": totalPages,
	}
	c.JSON(http.StatusOK, responseData)
}


func GetSingleDataHandler(c *gin.Context, db *sql.DB, query string, args []interface{}, destinations ...interface{}) {
    row := db.QueryRow(query, args...)

    if err := row.Scan(destinations...); err != nil {
        log.Fatal(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
        return
    }

    responseData := gin.H{}
    for i := 0; i < len(destinations); i += 2 {
        key, ok := destinations[i].(string)
        if !ok {
            log.Fatal("Invalid key type")
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
            return
        }

        value := destinations[i+1]
        responseData[key] = value
    }

    c.JSON(http.StatusOK, responseData)
}

//Example call
    /* Call GetSingleDataHandler to retrieve user data with multiple query arguments
    var user User
    query := "SELECT id, name, age FROM users WHERE id = $1 AND department = $2"
    userID := 123     // Replace with the desired user's ID
    userDepartment := "HR" // Replace with the desired department

    // Call GetSingleDataHandler with the query, query arguments, and destination variables
    GetSingleDataHandler(c, db, query, []interface{}{userID, userDepartment}, &user.ID, &user.Name, &user.Age)

    // Now 'user' variable contains the retrieved user data
    fmt.Println("User ID:", user.ID)
    fmt.Println("User Name:", user.Name)
    fmt.Println("User Age:", user.Age)
	*/
