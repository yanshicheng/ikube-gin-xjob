package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 对密码进行加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash 验证密码与哈希是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	password := "123456"

	// 加密密码
	hashedPassword, err := HashPassword(password)
	fmt.Println("Original password:", hashedPassword)
	fmt.Println("Original password:", hashedPassword)
	hashedPassword = "$2a$10$bIiJUsIlWBg3dXX2vZ0gHO0czaFoqXhDHQA25I1PbV1.L7uAgpbUy"
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed password:", hashedPassword)

	// 验证密码
	match := CheckPasswordHash(password, hashedPassword)
	fmt.Println("Password match:", match)
}
