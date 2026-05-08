package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/config"
	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	LoginPIN(req dto.LoginPINRequest) (*dto.LoginResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return s.generateTokens(user.ID.String(), user.Name, user.Email, string(user.Role), user.BranchID.String())
}

func (s *authService) LoginPIN(req dto.LoginPINRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByPIN(req.PIN)
	if err != nil {
		return nil, errors.New("invalid PIN")
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	return s.generateTokens(user.ID.String(), user.Name, user.Email, string(user.Role), user.BranchID.String())
}

func (s *authService) generateTokens(userID, name, email, role, branchID string) (*dto.LoginResponse, error) {
	accessToken, err := utils.GenerateJWT(userID, role, branchID, s.cfg.JWTSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(userID, role, branchID, s.cfg.JWTSecret, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       userID,
			Name:     name,
			Email:    email,
			Role:     role,
			BranchID: branchID,
		},
	}, nil
}

