package services

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"otppro/auth"
	"otppro/model"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	type RegisterRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.DBConn
	var user model.UserRegisterwe
	db.AutoMigrate(&model.UserRegisterwe{})
	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already registered"})
		return
	}

	token, err := auth.GenerateJWT(1, req.PhoneNumber)
	if err != nil {
		fmt.Println(err)
	}
	//user.Token = token
	otp := generateOTP()
	newUser := model.UserRegisterwe{
		Token:        token,
		PhoneNumber:  req.PhoneNumber,
		OTP:          otp,
		OTPExpiresAt: time.Now().Add(5 * time.Minute),
		DeviceID:     req.DeviceID,
	}

	if err := db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	sendSMS(req.PhoneNumber, otp)
	fmt.Println("otp", otp)
	fmt.Println("token", token)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func generateOTP() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%06d", int(b[0])%1000000)
}

func sendSMS(phone, otp string) {
	fmt.Printf("Sending OTP %s to phone %s\n", otp, phone)
}

func LoginUser(c *gin.Context) {
	type LoginRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		DeviceID    string `json:"device_id" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.DBConn
	var user model.UserRegisterwe

	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	otp := generateOTP()
	user.OTP = otp
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)
	db.Save(&user)

	sendSMS(req.PhoneNumber, otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent for login"})
}
func VerifyOTP(c *gin.Context) {
	type VerifyOTPRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		OTP         string `json:"otp" binding:"required"`
	}

	var req VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.DBConn
	var user model.UserRegisterwe

	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != req.OTP || time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	token, err := auth.GenerateJWT(1, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	user.Token = token
	db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}

func ResendOTP(c *gin.Context) {
	type ResendRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var req ResendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.DBConn
	var user model.UserRegisterwe

	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	otp := generateOTP()
	user.OTP = otp
	user.OTPExpiresAt = time.Now().Add(5 * time.Minute)

	db.Save(&user)

	sendSMS(req.PhoneNumber, otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP resent successfully"})
}
func GetUserDetails(c *gin.Context) {
	type GetUserRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}

	var req GetUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := model.DBConn
	var user model.UserRegisterwe

	if err := db.Where("phone_number = ?", req.PhoneNumber).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
