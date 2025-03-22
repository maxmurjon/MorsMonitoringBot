package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
)

type queryParams struct {
	Key string
	Val interface{}
}

type TokenInfo struct {
	UserID     string `json:"user_id"`
	ClientType string `json:"client_type"`
}

func ReplaceQueryParams(namedQuery string, params map[string]interface{}) (string, []interface{}) {
	var (
		i    int = 1
		args []interface{}
		arr  []queryParams
	)

	for k, v := range params {
		arr = append(arr, queryParams{
			Key: k,
			Val: v,
		})
	}

	sort.Slice(arr, func(i, j int) bool {
		return len(arr[i].Key) > len(arr[j].Key)
	})

	for _, v := range arr {
		if v.Key != "" && strings.Contains(namedQuery, ":"+v.Key) {
			namedQuery = strings.ReplaceAll(namedQuery, ":"+v.Key, "$"+strconv.Itoa(i))
			args = append(args, v.Val)
			i++
		}
	}

	return namedQuery, args
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

// GenerateJWT ...
func GenerateJWT(m map[string]interface{}, tokenExpireTime time.Duration, tokenSecretKey string) (tokenString string, err error) {
	var token *jwt.Token

	token = jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for key, value := range m {
		claims[key] = value
	}

	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(tokenExpireTime).Unix()

	tokenString, err = token.SignedString([]byte(tokenSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseClaims(token string, secretKey string) (result TokenInfo, err error) {
	var claims jwt.MapClaims
	fmt.Println("Hello")
	claims, err = ExtractClaims(token, secretKey)
	if err != nil {
		return result, err
	}
	fmt.Println(claims)
	result.UserID = cast.ToString(claims["UserId"])
	if len(result.UserID) <= 0 {
		err = errors.New("cannot parse 'user_id' field")
		return result, err
	}

	result.ClientType = cast.ToString(claims["client_type"])

	return
}

// ExtractClaims extracts claims from given token
func ExtractClaims(tokenString string, tokenSecretKey string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return []byte(tokenSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// ExtractToken checks and returns token part of input string
func ExtractToken(bearer string) (token string, err error) {
	strArr := strings.Split(bearer, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return token, errors.New("wrong token format")
}
