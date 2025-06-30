package service

import (
	"aidanwoods.dev/go-paseto"
	"context"
	"fmt"
	"log"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/repo"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type AuthService interface {
	// RegisterUser attempts to register a user with the provided credentials
	//
	// Returns the auth token string if succeeded
	RegisterUser(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (status_codes.UserRegister, string, error)

	// LoginUser checks if the login credentials are valid and returns an auth token
	LoginUser(
		ctx context.Context,
		credentials entities.UserCredentials,
	) (status_codes.UserLogin, string, error)

	// GetUserFromToken attempts to get a user from a token string
	GetUserFromToken(ctx context.Context, token string) (*entities.User, error)
}

type authService struct {
	publicKey  paseto.V4AsymmetricPublicKey
	privateKey paseto.V4AsymmetricSecretKey
	userRepo   repo.UserRepository
}

func NewAuthService(config entities.Config, userRepo repo.UserRepository) AuthService {
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(config.Auth.PublicKey)
	if err != nil {
		log.Printf("[NewV4AsymmetricPublicKeyFromHex] | %v", err)
		panic(err)
	}

	privateKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(config.Auth.PrivateKey)
	if err != nil {
		log.Printf("[NewV4AsymmetricSecretKeyFromHex] | %v", err)
		panic(err)
	}

	return authService{
		publicKey:  publicKey,
		privateKey: privateKey,
		userRepo:   userRepo,
	}
}

func (s authService) RegisterUser(
	ctx context.Context,
	credentials entities.UserCredentials,
) (status_codes.UserRegister, string, error) {
	// Clean fields
	credentials.Name = strings.TrimSpace(credentials.Name)

	// Validate fields
	if !rules.IsValidUserName(credentials.Name) {
		return status_codes.UserRegisterInvalidName, "", nil
	}
	if !rules.IsValidUserPassword(credentials.Password) {
		return status_codes.UserRegisterInvalidPassword, "", nil
	}

	// Check if a user with this name already exists
	user, err := s.userRepo.GetUserByName(ctx, credentials.Name)
	if err != nil {
		log.Printf("[GetUserByName] | %v", err)
		return -1, "", fmt.Errorf("[GetUserByName] | %v", err)
	}

	if user != nil {
		return status_codes.UserRegisterInvalidAlreadyRegistered, "", nil
	}

	// Hash the password
	credentials.Password, err = util.HashPassword(credentials.Password)
	if err != nil {
		log.Printf("[HashPassword] | %v", err)
		return -1, "", fmt.Errorf("[HashPassword] | %v", err)
	}

	newUser, err := s.userRepo.RegisterUser(ctx, credentials)
	if err != nil {
		log.Printf("[RegisterUser] | %v", err)
		return -1, "", fmt.Errorf("[RegisterUser] | %v", err)
	}

	// Generate auth token
	token, err := util.GenerateAuthToken(newUser.ID, s.privateKey)
	if err != nil {
		log.Printf("[GenerateAuthToken] | %v", err)
		return -1, "", fmt.Errorf("[GenerateAuthToken] | %v", err)
	}

	return status_codes.UserRegisterSuccess, token, nil
}

func (s authService) LoginUser(
	ctx context.Context,
	credentials entities.UserCredentials,
) (status_codes.UserLogin, string, error) {
	// Clean fields
	credentials.Name = strings.TrimSpace(credentials.Name)

	user, err := s.userRepo.GetUserByName(ctx, credentials.Name)
	if err != nil {
		log.Printf("[GetUserByName] | %v", err)
		return -1, "", fmt.Errorf("[GetUserByName] | %v", err)
	}

	// Check if user exists
	if user == nil {
		return status_codes.UserLoginNotFound, "", nil
	}

	// Check if the password matches
	if !util.CheckPasswordHash(credentials.Password, user.Password) {
		return status_codes.UserLoginWrongPassword, "", nil
	}

	// Generate auth token
	token, err := util.GenerateAuthToken(user.ID, s.privateKey)
	if err != nil {
		log.Printf("[GenerateAuthToken] | %v", err)
		return -1, "", fmt.Errorf("[GenerateAuthToken] | %v", err)
	}

	return status_codes.UserLoginSuccess, token, nil
}

func (s authService) GetUserFromToken(
	ctx context.Context,
	token string,
) (*entities.User, error) {
	// Get ID in token
	userID, err := util.GetIDFromToken(token, s.publicKey)
	if err != nil {
		return nil, fmt.Errorf("[GetUserIDFromToken] | %v", err)
	}

	// Get user from ID
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[GetUserByID] | %v", err)
	}

	return user, nil
}
