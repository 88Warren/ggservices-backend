package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	email "github.com/lmwcode/ggservices-backend/utils"
	"github.com/lmwcode/ggservices-backend/utils/emailtemplates"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

type ContactForm struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Subject string `json:"subject" binding:"required"`
	Message string `json:"message" binding:"required"`
	Token   string `json:"token"`
}

func (hc *HomeController) GetHome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Home Page",
	})
}

func (hc *HomeController) HandleContactForm(ctx *gin.Context) {
	var form ContactForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		log.Printf("Form binding error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":  err.Error(),
			"status": "error",
		})
		return
	}

	isDev := os.Getenv("GO_ENV") == "development"
	if !isDev && !verifyRecaptcha(form.Token) {
		log.Printf("reCAPTCHA verification failed for token: %s", form.Token)
		ctx.JSON(http.StatusForbidden, gin.H{"error": "reCAPTCHA verification failed"})
		return
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")

	emailBody := emailtemplates.GenerateContactFormEmailBody(form.Name, form.Email, form.Subject, form.Message)

	err := email.SendEmail(
		os.Getenv("SMTP_FROM"),
		os.Getenv("SMTP_TO"),
		"New Contact Form Submission: "+form.Subject,
		emailBody,
		form.Email,
		smtpPassword,
	)
	if err != nil {
		log.Printf("Email sending failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Message received!"})
}

func verifyRecaptcha(token string) bool {
	secret := os.Getenv("RECAPTCHA_SECRET")
	verifyURL := "https://www.google.com/recaptcha/api/siteverify"

	resp, err := http.PostForm(verifyURL, url.Values{
		"secret":   {secret},
		"response": {token},
	})
	if err != nil {
		log.Printf("reCAPTCHA verification error: %v", err)
		return false
	}
	defer resp.Body.Close()

	var result struct {
		Success    bool     `json:"success"`
		ErrorCodes []string `json:"error-codes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("reCAPTCHA response decode error: %v", err)
		return false
	}

	return result.Success
}
