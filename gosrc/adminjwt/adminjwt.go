package adminjwt

import (
	"fmt"
	"time"

	"os"

	"github.com/dgrijalva/jwt-go"
)

// KeyStore는 비밀키를 관리합니다.
// type KeyStore struct {
// 	keys map[int64]string // 시간 기반 키 저장
// }

// // NewKeyStore는 새로운 KeyStore를 생성합니다.
// func NewKeyStore() *KeyStore {
// 	return &KeyStore{
// 		keys: make(map[int64]string),
// 	}
// }

// // SetKey는 주어진 시간에 대한 비밀키를 설정합니다.
// func (ks *KeyStore) SetKey(t time.Time, key string) {
// 	ks.keys[t.Unix()] = key
// }

// // GetKey는 주어진 시간에 가장 가까운 비밀키를 반환합니다.
// func (ks *KeyStore) GetKey(t time.Time) string {
// 	// 여기서는 단순화를 위해 직접 구현하지 않고 개념만 설명합니다.
// 	// 실제로는 주어진 시간에 가장 가까운 키를 찾아 반환해야 합니다.
// 	return ks.keys[t.Unix()]
// }

var mySigningKey = []byte(os.Getenv("JWT_SIGNKEY"))

// 토큰 생성
func GenerateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", fmt.Errorf("토큰생성실패: %s", err.Error())
	}

	return tokenString, nil
}

// 토큰 검증
func ParseToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return mySigningKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// exp 클레임이 현재 시간보다 이전인지 확인
		if exp, ok := claims["exp"].(float64); ok {
			// exp는 Unix 타임스탬프로 제공됩니다.
			expTime := time.Unix(int64(exp), 0)
			if expTime.Before(time.Now()) {
				return nil, fmt.Errorf("토큰이 만료되었습니다")
			}
			// else {
			// 	return nil, fmt.Errorf("토큰이 유효합니다.")
			// }
		} else {
			return nil, fmt.Errorf("exp 클레임을 찾을 수 없습니다")
		}
	} else {
		return nil, fmt.Errorf("토큰 검증 실패")
	}

	return token, err
}
