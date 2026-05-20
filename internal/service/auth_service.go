package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
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
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	return s.generateTokens(user.ID.String(), user.Name, user.Email, string(user.Role), user.BranchID.String())
}

func (s *authService) LoginPIN(req dto.LoginPINRequest) (*dto.LoginResponse, error) {
	// FIX: Kita tidak tahu branchID sebelum login, jadi coba tanpa filter dulu
	// Atau: PIN login harus include branch_id di request
	// Untuk fix cepat: scan semua tapi dengan cost=10 hash = lebih cepat

	// Coba cari user dengan PIN di semua branch (cost=10 = ~50ms per user)
	// Di production, tambahkan branch_id di PIN login request

	// Workaround: coba cari user yang PIN-nya match
	users, err := s.userRepo.FindAll()
	if err != nil {
		return nil, errors.New("invalid PIN")
	}

	var matchedUser *models.User
	for i := range users {
		if users[i].IsActive && users[i].PIN != "" {
			if utils.CheckPasswordHash(req.PIN, users[i].PIN) {
				matchedUser = &users[i]
				break
			}
		}
	}

	if matchedUser == nil {
		return nil, errors.New("invalid PIN")
	}
	if !matchedUser.IsActive {
		return nil, errors.New("account is disabled")
	}
	return s.generateTokens(matchedUser.ID.String(), matchedUser.Name, matchedUser.Email, string(matchedUser.Role), matchedUser.BranchID.String())
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
			ID: userID, Name: name, Email: email, Role: role, BranchID: branchID,
		},
	}, nil
}
