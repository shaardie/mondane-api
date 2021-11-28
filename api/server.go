package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	router    *gin.Engine
	k8sclient k8sClient
}

func (s *server) route() error {
	s.router = gin.Default()
	s.router.POST("/api/v1/user", s.CreateUser())
	s.router.DELETE("/api/v1/user", s.DeleteUser())
	s.router.GET("/api/v1/user", s.ReadUser())
	s.router.PUT("/api/v1/user", s.UpdateUser())

	// Health Checks
	s.router.GET("/ready", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	s.router.GET("/health", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	return nil
}

func Run() error {
	s := server{
		k8sclient: &k8sClientImpl{},
	}

	k, err := initK8sClient("mondane-user")
	if err != nil {
		return fmt.Errorf("unable to initialize k8s client, %w", err)
	}
	s.k8sclient = k

	if err := s.route(); err != nil {
		return fmt.Errorf("unable to add routes to server, %w", err)
	}
	if err := s.router.Run(); err != nil {
		return fmt.Errorf("unable to run server, %w", err)
	}
	return nil
}
