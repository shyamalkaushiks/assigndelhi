package services

import (
	"otppro/auth"

	"github.com/gin-gonic/gin"
)

// HandlerService :
type HandlerService struct{}

// Bootstrap :
func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	// r.GET("/oauth/callback", callbackHandler)
	public := r.Group("/api/users")

	public.POST("/register", RegisterUser)
	public.POST("/login", LoginUser)
	public.POST("/Verifyotp", VerifyOTP)
	public.POST("/resend-otp", ResendOTP)
	//public.POST("/register", RegisterUser)
	r.Use(auth.Auth())
	groupRoute := r.Group("/api/users")
	//course api
	groupRoute.GET("/Profile", GetUserDetails)

}
