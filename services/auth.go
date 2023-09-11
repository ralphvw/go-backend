package services

import (
    "fmt"
    "net/http"
    "time"
	"database/sql"

    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


var jwtSecret = []byte("your-secret-key")

type User struct {
    ID       int    `json:"id"`
    Password string `json:"password"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Hash []byte
}

type Claims struct {
    UserID   int    `json:"userId"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`

    jwt.StandardClaims
}


func LoginHandler(c *gin.Context, db *sql.DB) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    authenticatedUser, err := authenticateUser(db, user.Email, user.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := createToken(authenticatedUser.ID, authenticatedUser.FirstName, authenticatedUser.LastName,authenticatedUser.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}


func authenticateUser(db *sql.DB, email string, password string) (*User, error) {
    query := "SELECT id, email, hash FROM users WHERE username = $1"
    row := db.QueryRow(query, email)

    var user User
    err := row.Scan(&user.ID, &user.Email, &user.Hash)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    if !comparePasswords(password, string(user.Hash)) {
        return nil, fmt.Errorf("authentication failed")
    }

    return &user, nil
}

func createToken(userID int, firstName string, lastName string, email string) (string, error) {
    claims := &Claims{
        UserID:   userID,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
}


func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
            c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("firstName", claims.FirstName)
			c.Set("lastName",claims.LastName)
            c.Next()
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
        }
    }
}

func protectedHandler(c *gin.Context) {
    userID := c.GetInt("userID")
    username := c.GetString("username")
    c.JSON(http.StatusOK, gin.H{"message": "Authenticated User", "userID": userID, "username": username})
}

func comparePasswords(plaintextPassword, hashedPassword string) bool {

    hashedPasswordBytes := []byte(hashedPassword)

    err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, []byte(plaintextPassword))

    return err == nil
}