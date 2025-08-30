package services

import (
	"fmt"

	"github.com/manuelrojas19/go-oauth2-server/api"
	"github.com/manuelrojas19/go-oauth2-server/store/repositories"
	"go.uber.org/zap"
)

type GetUserinfoCommand struct {
	AccessToken string
}

type UserinfoResponse struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

type userinfoService struct {
	accessTokenRepository repositories.AccessTokenRepository
	userRepository        repositories.UserRepository
	logger                *zap.Logger
}

func NewUserinfoService(accessTokenRepository repositories.AccessTokenRepository, userRepository repositories.UserRepository, logger *zap.Logger) UserinfoService {
	return &userinfoService{
		accessTokenRepository: accessTokenRepository,
		userRepository:        userRepository,
		logger:                logger,
	}
}

func (s *userinfoService) GetUserinfo(command *GetUserinfoCommand) (*UserinfoResponse, error) {
	s.logger.Info("Attempting to retrieve user info", zap.String("accessToken", command.AccessToken))

	accessTokenEntity, err := s.accessTokenRepository.FindByAccessToken(command.AccessToken)
	if err != nil {
		s.logger.Error("Access token not found or invalid", zap.Error(err))
		return nil, fmt.Errorf(api.ErrInvalidToken.Error())
	}

	userEntity, err := s.userRepository.FindById(accessTokenEntity.UserId)
	if err != nil {
		s.logger.Error("User not found for accessToken", zap.String("userId", accessTokenEntity.UserId), zap.Error(err))
		return nil, fmt.Errorf(api.ErrInvalidToken.Error())
	}

	userinfoResponse := &UserinfoResponse{
		Subject: userEntity.Id,
		Email:   userEntity.Email,
		Name:    userEntity.Name,
	}

	s.logger.Info("User info retrieved successfully", zap.Any("userinfo", userinfoResponse))
	return userinfoResponse, nil
}
