package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/yanshicheng/ikube-gin-xjob/global"
	"strings"
	"time"
)

type ApplicationRole struct {
	Application string   `json:"application"`
	Role        []string `json:"role"`
}

type JWTClaims struct {
	Account            string          `json:"account"`
	Application        ApplicationRole `json:"application"`
	jwt.StandardClaims                 // 内嵌标准的声明
}

type JWTResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// ParseJWT 解析并验证 JWT
func ParseJWT(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	// 提取 JWT
	JwtToken, err := extractTokenFromBearerString(tokenString)
	if err != nil {
		return nil, err
	}
	// 解析 JWT
	token, err := jwt.ParseWithClaims(JwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return global.JwtKey, nil
	})

	if err != nil {
		return nil, err
	}
	// 验证 JWT
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// extractTokenFromBearerString 从 Bearer 令牌字符串中提取 JWT
func extractTokenFromBearerString(bearerToken string) (string, error) {
	if len(bearerToken) < 7 || !strings.HasPrefix(bearerToken, "Bearer ") {
		return "", fmt.Errorf("invalid bearer token")
	}

	return bearerToken[7:], nil
}

func GenerateToken(account string, aRole ApplicationRole) (*JWTResponse, error) {
	now := time.Now()
	jti, err := GenerateRandomID()
	if err != nil {
		return nil, err
	}

	claims := JWTClaims{
		Account:     account,
		Application: aRole,
		StandardClaims: jwt.StandardClaims{
			Issuer:    "www.ikubeops.com",
			NotBefore: now.Unix(),                       // 生效时间：Unix时间戳，token在此时间之前不可用
			IssuedAt:  now.Unix(),                       // 发行时间：Unix时间戳，指明token何时被发行
			ExpiresAt: now.Add(15 * time.Minute).Unix(), // 过期时间：Unix时间戳，指明token何时过期
			Id:        jti,
		},
	}

	// 创建访问令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := token.SignedString(global.JwtKey)
	if err != nil {
		return nil, err
	}
	// 创建刷新令牌，通常具有更长的有效期
	claims.ExpiresAt = now.Add(24 * time.Hour).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err := refreshToken.SignedString(global.JwtKey)
	if err != nil {
		return nil, err
	}
	return &JWTResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}
