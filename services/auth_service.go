package services

import (
	"errors"
	"fmt"
	"gin-freemarket/models"
	"gin-freemarket/repositories"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Signup(email string, password string) error
	Login(email string, password string) (*string, error)
	GetUserFromToken(tokenString string) (*models.User, error)
}

type AuthService struct {
	repository repositories.IAuthRepository
}

func NewAuthService(repository repositories.IAuthRepository) IAuthService {
	return &AuthService{repository: repository}
}

func (s *AuthService) Signup(email string, password string) error {
	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}
	return s.repository.CreateUser(user)
}

func (s *AuthService) Login(email string, password string) (*string, error) {
	foundUser, err := s.repository.FindUser(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	// Tokenの生成
	token, err := CreateToken(foundUser.ID, foundUser.Email)
	if err != nil {
		return nil, err
	}

	return token, nil
}

// プライベートメソッド（ポインタレシーバいらない）
func CreateToken(userId uint, email string) (*string, error) {

	// token生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   userId, //ユーザー識別子
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(), //Tokenの有効期限
	})

	tokenString, err := token.SignedString([]byte(os.Getenv(("SECRET_KEY"))))
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}

func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {

	// jwtトークンを解析するために、jwt.Parseの第一引数に渡ってきたTokenStringを第二引数は無名関数。
	// これはjwt暗号化アルゴリズムが正しいかどうかをチェックしてあっていたら暗号化時に使用したsecret_keyを返すというもの
	// 最終的にTokenStringとsecret_keyを使ってパースする
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// token.Method.(*jwt.SigningMethodHMAC);の部分は型アサーションという
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}

	var user *models.User
	// tokensからClaimsを取得し、ok変数にtrueが返ってきたら期限をチェック
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 現在時刻がClaims内のExpireを超過していたら、期限切れを返す
		if exp, ok := claims["exp"].(float64); ok {
			if float64(time.Now().Unix()) > exp {
				return nil, jwt.ErrTokenExpired
			}
		} else {
			return nil, errors.New("invalid 'exp' claim type")
		}

		// ここで上で宣言したvar user *models.UserにDB検索結果を入れる。
		// 宣言した変数に格納するので:=ではなく、=。errも使いまわせる
		if email, ok := claims["email"].(string); ok {
			user, err = s.repository.FindUser(email)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("invalid token: email claim missing or not a string")
		}
	}
	return user, nil
}
