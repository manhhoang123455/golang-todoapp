package services

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"mini_project/entity"
	"reflect"

	"os"
	"strconv"
	"time"
)

type JWTService interface {
	GenerateToken(userID uint64, t time.Time) string

	ValidateToken(token string) (*jwt.Token, error)

	RefreshToken(token string) string

	Logout(authHeader string) bool

	AuthJWT(authHeader string) string

	GoogleGenerateToken(data interface{}) string
}

type jwtCustomClaim struct {
	UserID uint64 `json:"user_id"`

	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey   string
	issuer      string
	redisEntity entity.RedisEntity
	userEntity  entity.UserEntity
}

func NewJWTService(redisEntity entity.RedisEntity, userEntity entity.UserEntity) JWTService {
	return &jwtService{
		secretKey:   getSecretKey(),
		issuer:      "gojwt",
		redisEntity: redisEntity,
		userEntity:  userEntity,
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "learnGolangJWTToken"
	}
	return secretKey
}

// getUserDataByToken Get user data by token
func getUserDataByToken(authHeader string, claims jwt.Claims, secretKey string) (interface{}, error) {
	token, err := jwt.ParseWithClaims(authHeader, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return token, err
}

func GetTokenTTL() int {
	stringTTL := os.Getenv("JWT_TTL")
	if stringTTL == "" {
		stringTTL = "900"
	}
	intTTL, _ := strconv.Atoi(stringTTL)
	return intTTL
}

func (s *jwtService) GenerateToken(userID uint64, t time.Time) string {
	jwtTTL := GetTokenTTL()
	claims := &jwtCustomClaim{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(t),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
		},
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generateToken.SignedString([]byte(s.secretKey))
	if err != nil {
		panic("Failed to process request : Signature failed")
		return ""
	}

	_, setRedisErr := s.redisEntity.Set("token"+strconv.FormatUint(userID, 10), token, time.Duration(jwtTTL)*time.Second)
	if setRedisErr != nil {
		panic("Failed to set the token in redis : " + setRedisErr.Error())
		return ""
	}
	return token
}

func (s *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			panic("Failed to process request : Unexpected signing method")
			return nil, fmt.Errorf("Unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) RefreshToken(authHeader string) string {
	jwtTTL := GetTokenTTL()
	claims := &jwtCustomClaim{}
	_, err := getUserDataByToken(authHeader, claims, s.secretKey)
	if err != nil {
		panic("Failed to get user data (refresh token) : " + err.Error())
		return ""
	}

	token := s.GenerateToken(claims.UserID, time.Now().Add(time.Duration(jwtTTL)*time.Second))

	return token
}

// Logout User logout and remove token from redis
func (s *jwtService) Logout(authHeader string) bool {
	claims := &jwtCustomClaim{}
	_, erro := getUserDataByToken(authHeader, claims, s.secretKey)
	if erro != nil {
		panic("Failed to get user data (Logout) : " + erro.Error())
		return false
	}

	get, err := s.redisEntity.Del("token" + strconv.Itoa(int(claims.UserID)))
	if get == int64(1) {
		return true
	}
	if err != nil {
		panic("Failed to set the token in redis (Logout) : ")
		return false
	}

	return false
}

// AuthJWT Get the redis token and return the middleware
func (s *jwtService) AuthJWT(authHeader string) string {
	claims := &jwtCustomClaim{}
	_, erro := getUserDataByToken(authHeader, claims, s.secretKey)
	if erro != nil {
		panic("Failed to get user data (AuthJWT) : " + erro.Error())
		return ""
	}

	get, err := s.redisEntity.Get("token" + strconv.Itoa(int(claims.UserID)))
	if err != nil {
		panic("Failed to get the token in redis (AuthJWT) : ")
		return ""
	}

	return fmt.Sprintf("%v", get)
}

func (s *jwtService) GoogleGenerateToken(data interface{}) string {
	jwtTTL := GetTokenTTL()
	googleInfo := reflect.ValueOf(data).Elem()
	findByEmail := s.userEntity.FindByEmail(googleInfo.FieldByName("Email").String())
	if findByEmail.ID == 0 {
		return ""
	}

	token := s.GenerateToken(googleInfo.FieldByName("Id").Uint(), time.Now().Add(time.Duration(jwtTTL)*time.Second))
	return token
}
