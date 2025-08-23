package usecase

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/androsyz/nexus-user-svc/config"
	"github.com/androsyz/nexus-user-svc/constant"
	"github.com/androsyz/nexus-user-svc/graph/model"
	modelDB "github.com/androsyz/nexus-user-svc/model"
	"github.com/androsyz/nexus-user-svc/pkg/authctx"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type repoUserInterface interface {
	Create(ctx context.Context, user *modelDB.UserDB) (*string, error)
	GetByID(ctx context.Context, id string) (*modelDB.UserDB, error)
	GetByEmail(ctx context.Context, email string) (*modelDB.UserDB, error)
}

type UcUser struct {
	cfg      *config.Config
	repoUser repoUserInterface
	zlog     zerolog.Logger
}

func NewUserUsecase(cfg *config.Config, repoUser repoUserInterface, zlog zerolog.Logger) *UcUser {
	return &UcUser{
		cfg:      cfg,
		repoUser: repoUser,
		zlog:     zlog,
	}
}

func (uc *UcUser) Register(ctx context.Context, request model.RegisterRequest) (*model.AuthResponse, error) {
	if request.Email == "" || request.Password == "" {
		return nil, constant.ErrMissingCredentials
	}

	existingUser, err := uc.repoUser.GetByEmail(ctx, request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existingUser != nil {
		return nil, constant.ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrHashingPassword, err)
	}

	payload := &modelDB.UserDB{
		Email:    request.Email,
		Name:     request.Name,
		Password: string(hashedPassword),
	}

	userID, err := uc.repoUser.Create(ctx, payload)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrCreatingField("user"), err)
	}

	return uc.generateAuthResponse(*userID)
}

func (uc *UcUser) Login(ctx context.Context, request model.LoginRequest) (*model.AuthResponse, error) {
	if request.Email == "" || request.Password == "" {
		return nil, constant.ErrMissingCredentials
	}

	user, err := uc.repoUser.GetByEmail(ctx, request.Email)
	if err != nil {
		return nil, constant.ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, constant.ErrInvalidCredentials
	}

	userID := user.ID.String()

	return uc.generateAuthResponse(userID)
}

func (uc *UcUser) RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(uc.cfg.Settings.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, constant.ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, constant.ErrInvalidClaims
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return nil, constant.ErrInvalidSubject
	}

	return uc.generateAuthResponse(userID)
}

func (uc *UcUser) User(ctx context.Context) (*model.User, error) {
	userID, err := authctx.GetAuthUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := uc.repoUser.GetByID(ctx, userID)
	if err != nil {
		return nil, constant.ErrUserNotFound
	}

	resp := &model.User{
		ID:    user.ID.String(),
		Email: user.Email,
		Name:  user.Name,
	}

	return resp, nil
}

func (uc *UcUser) generateAuthResponse(userID string) (*model.AuthResponse, error) {
	accessToken, err := uc.generateJWT(userID, false)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGeneratingJWT, err)
	}

	refreshToken, err := uc.generateJWT(userID, true)
	if err != nil {
		return nil, constant.ErrWithMsg(constant.ErrGeneratingJWT, err)
	}

	return &model.AuthResponse{
		Token:        accessToken,
		RefreshToken: &refreshToken,
	}, nil
}

func (uc *UcUser) generateJWT(userID string, isRefreshToken bool) (string, error) {
	tokenDuration := uc.cfg.Settings.TokenDuration
	if isRefreshToken {
		tokenDuration = uc.cfg.Settings.RefreshTokenDuration
	}

	duration := time.Hour * time.Duration(tokenDuration)

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(duration).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.cfg.Settings.JWTSecret))
}
