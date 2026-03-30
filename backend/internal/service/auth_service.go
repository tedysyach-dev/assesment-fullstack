package service

import (
	"context"
	"wms/core/base"
	"wms/core/errors"
	"wms/core/utils"
	"wms/internal/entity"
	"wms/internal/model"
	"wms/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
)

type AuthService struct {
	DB              *bun.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	UsersRepository *repository.UsersRepository
	TokenUtil       *utils.TokenUtil
}

func NewAuthService(
	db *bun.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	usersRepository *repository.UsersRepository,
	tokenUtil *utils.TokenUtil,
) *AuthService {
	return &AuthService{
		DB:              db,
		Log:             logger,
		Validate:        validate,
		UsersRepository: usersRepository,
		TokenUtil:       tokenUtil,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, req *model.RegisterRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return errors.NewValidationError(err)
	}

	user := entity.Users{}
	s.UsersRepository.FindOne(ctx, &user, base.WithWhere("email = ?", req.Email))

	if user.ID != "" {
		return errors.NewConflictError("user already exits")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return errors.NewInternalError(err)
	}

	user.ID = uuid.New().String()
	user.Email = req.Email
	user.PasswordHash = hashedPassword

	if err := s.UsersRepository.Create(ctx, &user); err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, errors.NewValidationError(err)
	}

	user := entity.Users{}
	s.UsersRepository.FindOne(ctx, &user, base.WithWhere("email = ?", req.Email))

	if user.ID == "" {
		return nil, errors.NewUnauthorizedError("The email or password you entered is incorrect")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.NewUnauthorizedError("The email or password you entered is incorrect")
	}

	accessToken, err := s.TokenUtil.CreateAccessToken(ctx, &model.Auth{
		UID: user.ID,
	})
	if err != nil {
		s.Log.WithError(err).Error("failed to generate access token")
		return nil, errors.NewInternalError(err)
	}

	return &model.LoginResponse{
		AccessToken: accessToken,
	}, nil
}
