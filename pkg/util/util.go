package util

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GetCurrentTimeInHex() string {
	// Get the current time in the desired format
	currentTime := time.Now().Unix()

	// Convert the formatted time to hexadecimal
	return Int64oHex(currentTime)
}

func Int64oHex(num int64) string {
	// Convert the integer to a hexadecimal string with a max length of 4 characters
	return strconv.FormatInt(num, 16)
}

func IntToHex(input int) string {
	// Convert the integer to a hexadecimal string with a max length of 4 characters
	return fmt.Sprintf("%04X", input)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
