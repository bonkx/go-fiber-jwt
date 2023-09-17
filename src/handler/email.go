package handler

import (
	"myapp/pkg/configs"
	"myapp/pkg/helpers"
	"myapp/src/models"

	"github.com/gofiber/fiber/v2"
)

type EmailHandler struct {
	userUsecase models.UserUsecase
}

func NewEmailHandler(r *fiber.App, userUsecase models.UserUsecase) {
	handler := &EmailHandler{
		userUsecase: userUsecase,
	}

	r.Get("/verify-email/:verificationCode", handler.VerifyEmail)

	// test view email
	r.Get("/view/register-email", handler.ViewRegisterEmail)
	r.Get("/view/verify-success", handler.ViewVerifySuccess)
	r.Get("/view/otp", handler.ViewOtpEmail)

}

func (h *EmailHandler) VerifyEmail(c *fiber.Ctx) error {
	code := c.Params("verificationCode")

	if err := h.userUsecase.VerificationEmail(c.Context(), code); err != nil {
		errD := fiber.NewError(err.Code, err.Message)
		return c.Status(errD.Code).JSON(errD)
	}

	siteData, _ := configs.GetSiteData(".")
	emailData := helpers.EmailData{
		Subject:  "Account verification success",
		SiteData: siteData,
	}

	return c.Render("emails/verification_done", emailData)
}

func (h *EmailHandler) ViewRegisterEmail(c *fiber.Ctx) error {
	siteData, _ := configs.GetSiteData(".")

	emailData := helpers.EmailData{
		URL:          siteData.ClientOrigin + "/verify-email/" + "QdkGUPVhjqu7sy7hGQqsGmg2YOOx9OIcyZQveNPljRpmWuE9NKMQ1pz6x49mEGfm",
		FirstName:    "Admin",
		Subject:      "Your account verification",
		TypeOfAction: "Register",
		SiteData:     siteData,
	}

	return c.Render("emails/verification_code", emailData)
}

func (h *EmailHandler) ViewVerifySuccess(c *fiber.Ctx) error {
	siteData, _ := configs.GetSiteData(".")

	emailData := helpers.EmailData{
		Subject:  "Account verification success",
		SiteData: siteData,
	}

	return c.Render("emails/verification_done", emailData)
}

func (h *EmailHandler) ViewOtpEmail(c *fiber.Ctx) error {
	siteData, _ := configs.GetSiteData(".")

	emailData := helpers.EmailData{
		URL:          "123456",
		FirstName:    "farrid",
		Subject:      "Your OTP code",
		TypeOfAction: "Forgot Password",
		SiteData:     siteData,
	}

	return c.Render("emails/otp_code", emailData)
}
