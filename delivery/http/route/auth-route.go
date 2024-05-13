package route

import (
	"database/sql"
	"os"
	"profiln-be/delivery/http"
	email "profiln-be/libs/email"
	"profiln-be/package/auth"
	repository "profiln-be/package/auth/repository"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewAuthRoute(app *gin.RouterGroup, db *sql.DB, log *logrus.Logger) {
	smtpSender := os.Getenv("SENDER_NAME")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	authEmail := os.Getenv("AUTH_EMAIL")
	authPassword := os.Getenv("AUTH_PASSWORD")

	email := email.NewEmail(smtpPort, smtpSender, smtpHost, authEmail, authPassword, log)
	authRepository := repository.NewAuthRepository(db)
	authUsecase := auth.NewAuthUsecase(authRepository, email, log)
	authController := http.NewAuthController(authUsecase)

	app.POST("/login", authController.Login)
	app.POST("/reset-password", authController.ResetPassword)
	app.POST("/register", authController.Register)
	app.POST("/user-otp", authController.VerifiedEmail)
	app.POST("/email/reset-password", authController.SendResetPasswordEmail)
	app.GET("/user-otp/:email", authController.GetUserOtpByEmail)
	app.POST("resend-otp", authController.ResendOTP)
}
