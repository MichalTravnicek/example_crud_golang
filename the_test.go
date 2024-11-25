package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB sets up an in-memory SQLite database for testing.
func SetupTestDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  true,        // Enable color
		},
	)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect to the test database")
	}
	db.AutoMigrate(&DbUser{})
	return db
}

func CreateDbUser(uuidString string, name string, email string) DbUser {
	var uuidDb uuid.UUID
	var err error
	if uuidString == "" {
		uuidDb = uuid.New()
	} else {
		uuidDb, err = uuid.Parse(uuidString)
		if err != nil {
			log.Fatal("Error - Invalid UUID:", uuidString)
		}
	}
	return DbUser{UUID: uuidDb, Name: name, Email: email, Birth: time.Now()}
}

func CreateRestUser(uuidString string, name string, email string, timeString string) RestUser {
	if uuidString == "" {
		uuidString = uuid.New().String()
	}
	if timeString == "" {
		timeString = time.Now().Format(timeFormat)
	}
	return RestUser{ID: uuidString, Name: name, Email: email, Birth: timeString}
}

func TestCreateUser(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Create a new user.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	user := CreateRestUser(uuid, "John Doe", "john@example.com", "2020-01-01T12:12:35+00:00")
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdUser RestUser
	json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "john@example.com", createdUser.Email)
}

func TestCreateDuplicatedUser(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Create a new user.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	user := CreateRestUser(uuid, "John Doe", "john@example.com", "2020-01-01T12:12:35+00:00")
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	var createdUser RestUser
	json.Unmarshal(w.Body.Bytes(), &createdUser)
	assert.Equal(t, "John Doe", createdUser.Name)
	assert.Equal(t, "john@example.com", createdUser.Email)
	assert.Contains(t, w2.Body.String(),"UNIQUE constraint failed")
}

func TestCreateUserBadRequest(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Create a new user.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	user := CreateRestUser(uuid, "John Doe", "john@example.com", "xxxx-01T12:12:35+00:00")
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUser(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Insert a user into the in-memory database.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	user := CreateDbUser(uuid, "Jane Doe", "jane@example.com")
	db.Create(&user)

	// Make a GET request.
	req, _ := http.NewRequest("GET", "/users/"+uuid, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedUser RestUser
	json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.Equal(t, "Jane Doe", fetchedUser.Name)
	assert.Equal(t, "jane@example.com", fetchedUser.Email)
}

func TestGetUserNotFound(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Make a GET request for a non-existent user.
	req, _ := http.NewRequest("GET", "/users/f7993007-8c03-4672-bb79-c8bee775e387", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUserBadRequest(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Make a GET request for a non-existent UUID.
	req, _ := http.NewRequest("GET", "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUser(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Insert a user into the in-memory database.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	dbUser := CreateDbUser(uuid, "Jane Doe", "jane@example.com")
	db.Create(&dbUser)

	// Update user details.
	updatedUser := CreateRestUser(uuid, "Jane Smith", "jane.smith@example.com", "")
	jsonValue, _ := json.Marshal(updatedUser)
	req, _ := http.NewRequest("PUT", "/users/"+uuid, bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var fetchedUser RestUser
	json.Unmarshal(w.Body.Bytes(), &fetchedUser)
	assert.Equal(t, "Jane Smith", fetchedUser.Name)
	assert.Equal(t, "jane.smith@example.com", fetchedUser.Email)
}

func TestDeleteUser(t *testing.T) {
	db := SetupTestDB()
	r := SetupRouter(db)

	// Insert a user into the in-memory database.
	uuid := "d95cc5a3-62d7-49ce-a094-f65a82caac5f"
	dbUser := CreateDbUser(uuid, "Jane Doe", "jane@example.com")
	db.Create(&dbUser)

	// Delete the user.
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%s", "/users/", uuid), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Verify that the user is deleted.
	var fetchedUser DbUser
	err := db.Where("UUID = ?", uuid).First(&fetchedUser).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}
