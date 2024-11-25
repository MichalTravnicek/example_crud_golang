package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB sets up an in-memory SQLite database for testing.
func SetupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect to the test database")
    }
    db.AutoMigrate(&User{})
    return db
}

func TestCreateUser(t *testing.T) {
    db := SetupTestDB()
    r := SetupRouter(db)

    // Create a new user.
    user := User{Name: "John Doe", Email: "john@example.com"}
    jsonValue, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)

    var createdUser User
    json.Unmarshal(w.Body.Bytes(), &createdUser)
    assert.Equal(t, "John Doe", createdUser.Name)
    assert.Equal(t, "john@example.com", createdUser.Email)
}

func TestGetUser(t *testing.T) {
    db := SetupTestDB()
    r := SetupRouter(db)

    // Insert a user into the in-memory database.
    user := User{Name: "Jane Doe", Email: "jane@example.com"}
    db.Create(&user)

    // Make a GET request.
    req, _ := http.NewRequest("GET", "/users/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var fetchedUser User
    json.Unmarshal(w.Body.Bytes(), &fetchedUser)
    assert.Equal(t, "Jane Doe", fetchedUser.Name)
    assert.Equal(t, "jane@example.com", fetchedUser.Email)
}

func TestGetUserNotFound(t *testing.T) {
    db := SetupTestDB()
    r := SetupRouter(db)

    // Make a GET request for a non-existent user.
    req, _ := http.NewRequest("GET", "/users/999", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateUser(t *testing.T) {
    db := SetupTestDB()
    r := SetupRouter(db)

    // Insert a user into the in-memory database.
    user := User{Name: "Jane Doe", Email: "jane@example.com"}
    db.Create(&user)
	log.Printf("User:%v",user)
    // Update user details.
    updatedUser := User{Name: "Jane Smith", Email: "jane.smith@example.com"}
    jsonValue, _ := json.Marshal(updatedUser)
    req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(jsonValue))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    var fetchedUser User
    json.Unmarshal(w.Body.Bytes(), &fetchedUser)
    assert.Equal(t, "Jane Smith", fetchedUser.Name)
    assert.Equal(t, "jane.smith@example.com", fetchedUser.Email)
}

func TestDeleteUser(t *testing.T) {
    db := SetupTestDB()
    r := SetupRouter(db)

    // Insert a user into the in-memory database.
    user := User{Name: "Jane Doe", Email: "jane@example.com"}
    db.Create(&user)
	
    // Delete the user.
    req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s%d","/users/",user.ID), nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)

    // Verify that the user is deleted.
    var fetchedUser User
    err := db.First(&fetchedUser, user.ID).Error
    assert.Error(t, err)
    assert.Equal(t, gorm.ErrRecordNotFound, err)
}