package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

const appleKeysURL = "https://appleid.apple.com/auth/keys"

type AppleKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type AppleKeysResponse struct {
	Keys []AppleKey `json:"keys"`
}

// GetAppleKeys 获取Apple的公钥
func GetAppleKeys() (*AppleKeysResponse, error) {
	resp, err := http.Get(appleKeysURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get apple keys: status code %d", resp.StatusCode)
	}

	var result AppleKeysResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyAppleToken 使用Apple的公钥验证身份令牌
func VerifyAppleToken(tokenStr string) (*jwt.Token, error) {
	fmt.Printf("Token: %s\n", tokenStr)
	appleKeys, err := GetAppleKeys()
	if err != nil {
		return nil, err
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid token: kid is not string")
		}

		for _, key := range appleKeys.Keys {
			if key.Kid == kid {
				// 解析e和n来生成RSA公钥
				e, err := base64.RawURLEncoding.DecodeString(key.E)
				fmt.Println(e)
				fmt.Println(key.E)

				buf := make([]byte, 4) // 4 bytes buffer for Uint32
				copy(buf[4-len(e):], e)

				if err != nil {
					return nil, fmt.Errorf("failed to decode exponent: %v", err)
				}
				n, err := base64.RawURLEncoding.DecodeString(key.N)
				if err != nil {
					return nil, fmt.Errorf("failed to decode modulus: %v", err)
				}

				exponent := big.NewInt(int64(binary.BigEndian.Uint32(buf)))
				modulus := new(big.Int).SetBytes(n)

				return &rsa.PublicKey{
					N: modulus,
					E: int(exponent.Int64()),
				}, nil
			}
		}

		return nil, fmt.Errorf("key not found: %s", kid)
	}
	fmt.Println(jwt.Parse(tokenStr, keyFunc))
	fmt.Println(95)
	return jwt.Parse(tokenStr, keyFunc)
}
