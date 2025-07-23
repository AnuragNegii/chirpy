package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string)(string, error){
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "Chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject: userID.String(),
	})
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil{
		return "", err
	}
	return ss, nil	
}

func ValidateJWT(tokenString, tokenSecret string)(uuid.UUID, error){
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{},func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}
	
	strignId, err := token.Claims.GetSubject()
	if err != nil{
		return uuid.Nil, err
	}
	
	// parseId, err  := uuid.Parse(strignId)
	// if err != nil {
	// 	return uuid.Nil, err
	// }
	return uuid.Parse(strignId) 
	// claims, ok := token.Claims.(*jwt.RegisteredClaims)
	// if !ok || !token.Valid{
	// 	return uuid.Nil, fmt.Errorf("not valid token")
	// }
	//
	// return uuid.Parse(claims.Subject)
}

func GetBearerToken(headers http.Header)(string, error){
	atString := headers.Get("Authorization")
	if atString == ""{
		return "", fmt.Errorf("no authorization header found")
	}
	atString = strings.TrimPrefix(atString, "Bearer ")
	return atString,nil 
}
