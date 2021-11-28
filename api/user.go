package api

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const xAuthUserHeader = "X-Auth-Request-User"

func (s *server) getUserID(c *gin.Context) (string, error) {
	userHeader := c.GetHeader(xAuthUserHeader)
	if userHeader == "" {
		return "", fmt.Errorf("header %v is missing", xAuthUserHeader)
	}
	return base64.StdEncoding.EncodeToString([]byte(userHeader)), nil

}

func (s *server) CreateUser() gin.HandlerFunc {

	type input struct {
		Username string   `json:"username" binding:"required"`
		URLs     []string `json:"urls" binding:"required"`
		Email    string   `json:"email" binding:"required"`
	}

	return func(c *gin.Context) {
		i := &input{}

		id, err := s.getUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		if err := c.ShouldBindJSON(i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &User{
			ID:       id,
			Username: i.Username,
			URLs:     i.URLs,
			Email:    i.Email,
		}

		if err := s.k8sclient.Create(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, user)
	}
}

func (s *server) ReadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := s.getUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user, err := s.k8sclient.Read(c, id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func (s *server) UpdateUser() gin.HandlerFunc {

	type input struct {
		Username string   `json:"username" binding:"required"`
		URLs     []string `json:"urls" binding:"required"`
		Email    string   `json:"email" binding:"required"`
	}

	return func(c *gin.Context) {
		id, err := s.getUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		i := &input{}
		if err := c.ShouldBindJSON(i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := &User{
			ID:       id,
			Username: i.Username,
			URLs:     i.URLs,
			Email:    i.Email,
		}

		if err := s.k8sclient.Update(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (s *server) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := s.getUserID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		user := &User{ID: id}
		if err := s.k8sclient.Delete(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
	}
}
