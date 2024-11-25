package main

import (
	"fmt"
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"net/http"
)

var (
    dbHandle *gorm.DB
	user string
	password string
	db string
	host string
	port string
	ssl string
)

func init() {
	user = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	db = os.Getenv("POSTGRES_DB")
	host = os.Getenv("POSTGRES_HOST")
	port = os.Getenv("POSTGRES_PORT")
	ssl = os.Getenv("POSTGRES_SSL")
}

func initDb() *gorm.DB{
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, db, port, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if (err != nil){
		log.Println("Database init error:", err)
	}
	return db
}

// User represents a simple user model.
type User struct {
    ID    uint   `gorm:"primaryKey"`
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}

// SetupRouter initializes the Gin engine with routes.
func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // Inject the database into the handler
    r.POST("/users", func(c *gin.Context) {
        var user User

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Create(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, user)
    })

    r.GET("/users/:id", func(c *gin.Context) {
        var user User
        id := c.Param("id")

        if err := db.First(&user, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, user)
    })

    r.PUT("/users/:id", func(c *gin.Context) {
        var user User
        id := c.Param("id")

        if err := db.First(&user, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Save(&user).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
		c.JSON(http.StatusOK, user)
	})

    r.DELETE("/users/:id", func(c *gin.Context) {
        id := c.Param("id")

        if err := db.Delete(&User{}, id).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
    })

    return r
}

func main() {
	dbHandle = initDb()
    log.Println(dbHandle)
    fmt.Println("App started")

    dbHandle.AutoMigrate(&User{})

    r := SetupRouter(dbHandle)
    r.Run(":8080")
}