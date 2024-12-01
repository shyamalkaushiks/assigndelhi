package auth

import (
	"errors"
	"time"

	"otppro/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	JWTClaimAccountId     int
	JWTClaimAccountRoleId int
	JWTClaimUserId        int
	JWTClaimPhoneno       string
)

var jwtKey = []byte(config.Config.API_SECRET)

// JWTClaim :
type JWTClaim struct {
	//AccountId int    `json:"account_id"`
	Id      int    `json:"id"`
	Phoneno string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT :
func GenerateJWT(userID int, Phoneno string) (tokenString string, err error) {
	// expirationTime := time.Now().Add(24 * time.Hour)
	expirationTime := time.Now().Add(time.Hour * 24)
	claims := &JWTClaim{
		//AccountId: accountId,
		// AccountRoleId: accountRoleId,
		Id:      userID,
		Phoneno: Phoneno,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

// Auth :
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		err := ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}
}

// ValidateToken :
func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	//
	//	JWTClaimAccountId = claims.AccountId
	// JWTClaimAccountRoleId = claims.AccountRoleId
	JWTClaimUserId = claims.Id
	JWTClaimPhoneno = claims.Phoneno

	if JWTClaimUserId == 0 || JWTClaimPhoneno == "" {
		err = errors.New("token not valid")
		return
	}
	// if JWTClaimAccountId == 0 ||  JWTClaimUserId == 0 || JWTClaimEmail == "" {
	// 	err = errors.New("token not valid")
	// 	return
	// }
	return
}
