package authenticate

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sync"
	"time"
)

var secretKey = []byte(uuid.New().String())

type TokenBlockList struct {
	sync.Mutex
	Tokens []string `json:"token"`
}

func CreateToken(username []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyTokenStructure(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil
}

func VerifyToken(tokenString string, blockList *TokenBlockList) error {
	err := VerifyTokenStructure(tokenString)
	if err != nil {
		return err
	}

	// check to make sure the token is not in the blocklist
	if blockList.CheckToken(tokenString) {
		return fmt.Errorf("token is in blocklist")
	}
	return nil
}

// AddToken function to add token to blocklist
func (t *TokenBlockList) AddToken(token string) {
	t.Lock()
	defer t.Unlock()
	t.Tokens = append(t.Tokens, token)

}

// CheckToken function to check if token is in blocklist
func (t *TokenBlockList) CheckToken(token string) bool {
	t.Lock()
	defer t.Unlock()
	for _, v := range t.Tokens {
		if v == token {
			return true
		}
	}
	return false
}

// RemoveExpiredTokens function to remove expired tokens from blocklist
func (t *TokenBlockList) RemoveExpiredTokens() {
	t.Lock()
	defer t.Unlock()
	temp := make([]string, 0)
	for _, v := range t.Tokens {
		err := VerifyTokenStructure(v)
		if err == nil {
			temp = append(temp, v)
		}
	}
	t.Tokens = temp
}

// NewTokenBlockList function to create new token blocklist
func NewTokenBlockList() *TokenBlockList {
	return &TokenBlockList{
		Tokens: make([]string, 0),
	}
}
