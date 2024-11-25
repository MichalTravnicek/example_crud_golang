package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"net/http"
)

const timeFormat = "2006-01-02T15:04:05+00:00"

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

func (DbUser) TableName() string {
    return "users"
}

// DbUser represents a simple user model.
type DbUser struct {
    ID    uint   `gorm:"primaryKey"`
    UUID  uuid.UUID   `gorm:"type:uuid;unique_index"`
    Name  string `gorm:"type:varchar(150)"`
    Email string `gorm:"type:varchar(150);unique_index" form:"email"`
    Birth time.Time
}

type RestUser struct {
    ID  string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
    Birth string `json:"date_of_birth"`
}

func MapToDb(user RestUser) (DbUser, error) {
    log.Printf("%+v\n", user)
    uuid, err := uuid.Parse(user.ID)
    if (err != nil) {
        log.Println("UUID parse failed:",user.ID)
        return DbUser{}, err
    }
    birth, err2 :=  time.Parse(timeFormat, user.Birth)
    if (err2 != nil) {
        log.Println("Time parse failed:", user.Birth," required:",timeFormat)
        return DbUser{}, err2
    }
    return DbUser{UUID:uuid, Name: user.Name, Email:user.Email, Birth: birth}, nil
}

func MapToRest(dbUser DbUser) RestUser{
    log.Printf("%+v\n", dbUser)
    return RestUser{ID:dbUser.UUID.String(), Name: dbUser.Name, Email:dbUser.Email, Birth: dbUser.Birth.String()}
}

// SetupRouter initializes the Gin engine with routes.
func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    // Inject the database into the handler
    r.POST("/users", func(c *gin.Context) {
        var user RestUser

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        dbUser, err := MapToDb(user)

        if(err!=nil){
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return 
        }
        if err := db.Create(&dbUser).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, MapToRest(dbUser))
    })

    r.GET("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        dbUuid, err := uuid.Parse(id)
        if (err != nil){
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        var user = DbUser{}

        if err := db.Where("UUID = ?", dbUuid).First(&user).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, MapToRest(user))
    })

    r.PUT("/users/:id", func(c *gin.Context) {
        var user RestUser
        id := c.Param("id")        
        dbUuid, err := uuid.Parse(id)
        if (err != nil){
            c.JSON(http.StatusBadRequest, err)
            return
        }
        var dbUser = DbUser{}

        if err := db.Where("UUID = ?", dbUuid).First(&dbUser).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        dbUser, err = MapToDb(user)

        if (err != nil){
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Save(&dbUser).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
		c.JSON(http.StatusOK, MapToRest(dbUser))
	})

    r.DELETE("/users/:id", func(c *gin.Context) {
        id := c.Param("id")
        dbUuid, err := uuid.Parse(id)
        if (err != nil){
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        if err := db.Where("UUID = ?", dbUuid).Delete(&DbUser{}).Error; err != nil {
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

    dbHandle.AutoMigrate(&DbUser{})

    r := SetupRouter(dbHandle)
    r.Run(":8080")
}