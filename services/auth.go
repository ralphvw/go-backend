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

// Define a secret key for JWT token signing (should be kept secure).
var jwtSecret = []byte("your-secret-key")

// User represents a user in the application.
type User struct {
    ID       int    `json:"id"`
    Password string `json:"password"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	Email string `json:"email"`
	Hash []byte
}

// Create a sample user database (replace with a real user database).
// var users = []User{
//     {ID: 1, Username: "user1", Password: "password1"},
//     {ID: 2, Username: "user2", Password: "password2"},
// }

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

// authenticateUser verifies user credentials (replace with your actual authentication logic).
func authenticateUser(db *sql.DB, email string, password string) (*User, error) {
    // Query the database to retrieve the user based on the username.
    query := "SELECT id, email, hash FROM users WHERE username = $1"
    row := db.QueryRow(query, email)

    var user User
    err := row.Scan(&user.ID, &user.Email, &user.Hash)
    if err != nil {
        // Handle the error, such as "no rows found" or a database error.
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, err
    }

    // Compare the provided password with the hashed password stored in the database.
    if !comparePasswords(password, string(user.Hash)) {
        return nil, fmt.Errorf("authentication failed")
    }

    return &user, nil
}

// createToken generates a JWT token with custom claims.
func createToken(userID int, firstName string, lastName string, email string) (string, error) {
    claims := &Claims{
        UserID:   userID,
		FirstName: firstName,
		LastName: lastName,
		Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Token expires in 1 hour
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

// authMiddleware is a middleware function that verifies JWT tokens.
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Parse and validate the JWT token.
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        if claims, ok := token.Claims.(*Claims); ok && token.Valid {
            // Authentication successful, set user information in the context.
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

// protectedHandler is a protected route that requires JWT authentication.
func protectedHandler(c *gin.Context) {
    userID := c.GetInt("userID")
    username := c.GetString("username")
    c.JSON(http.StatusOK, gin.H{"message": "Authenticated User", "userID": userID, "username": username})
}


// ComparePasswords securely compares a plaintext password with a hashed password.
func comparePasswords(plaintextPassword, hashedPassword string) bool {

    hashedPasswordBytes := []byte(hashedPassword)

    err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, []byte(plaintextPassword))

    return err == nil
}