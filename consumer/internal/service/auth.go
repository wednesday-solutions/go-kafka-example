package service

import (
	"crypto/sha1"
	"os"

	"github.com/wednesday-solutions/go-template-consumer/internal/config"
	"github.com/wednesday-solutions/go-template-consumer/internal/jwt"
	"github.com/wednesday-solutions/go-template-consumer/pkg/utl/secure"
)

// Secure returns new secure service
func Secure(cfg *config.Configuration) *secure.Service {
	return secure.New(cfg.App.MinPasswordStr, sha1.New())
}

// JWT returns new JWT service
func JWT(cfg *config.Configuration) (jwt.Service, error) {
	return jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
}
