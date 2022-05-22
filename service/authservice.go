package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/PudgeKim/go-holdem/config"
	"github.com/PudgeKim/go-holdem/domain/entity"
	"github.com/PudgeKim/go-holdem/domain/repository"
	"github.com/PudgeKim/go-holdem/errors/autherror"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (a *AuthService) SignUp(ctx context.Context, email, password, nickname string) error {
	user, err := a.userRepo.FindByEmail(ctx, email); if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}
	if user != nil {
		return errors.New("user already exists")
	}
	
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost); if err != nil {
		return err 
	}
	newUser := entity.NewUser(nickname, email, string(hashedPW))
	if err := a.userRepo.Save(ctx, newUser); err != nil {
		return err 
	}
	
	return nil 
}

func (a *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := a.userRepo.FindByEmail(ctx, email); if user != nil {
		return "", err 
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err 
	}

	token, err := a.createToken(user.Id); if err != nil {
		return "", err 
	}

	return token, nil 
}

func (a *AuthService) FindUser(ctx context.Context, userId int64) (*entity.User, error) {
	user, err := a.userRepo.FindOne(ctx, userId); if err != nil {
		return nil, err 
	}
	return user, nil 
}

func (a *AuthService) UpdateBalance(ctx context.Context, userId int64, balance uint64) (uint64, error) {
	totalBalance, err := a.userRepo.UpdateBalance(ctx, userId, balance); if err != nil {
		return 0, err 
	}
	return totalBalance, nil 
}

func (a *AuthService) ValidateToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return config.TokenSecret, nil 
	})
	if err != nil {
		return 0, err 
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId, ok := claims["id"].(int64); if !ok {
			return 0, errors.New("Conversion Error")
		}
		return userId ,nil 
	}
	return 0, autherror.InvalidToken
}

type CustomClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func (a *AuthService) createToken(userId int64) (string, error) {
	claims := CustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: 30000,
			Issuer: "goholdem",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.TokenSecret); if err != nil {
		return "", err 
	}

	return tokenString, nil 
}