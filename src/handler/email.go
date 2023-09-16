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
	r.Get("/view/register-email", handler.registerEmail)
	r.Get("/view/verify-success", handler.verifySuccess)

}

func (h *EmailHandler) VerifyEmail(c *fiber.Ctx) error {
	code := c.Params("verificationCode")

	if err := h.userUsecase.VerificationEmail(c.Context(), code); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	siteData, _ := configs.GetSiteData(".")
	emailData := helpers.EmailData{
		Subject:  "Account verification success",
		SiteData: siteData,
	}

	return c.Render("emails/verificationDone", emailData)
}

func (h *EmailHandler) registerEmail(c *fiber.Ctx) error {
	siteData, _ := configs.GetSiteData(".")

	emailData := helpers.EmailData{
		URL:          siteData.ClientOrigin + "/verify-email/" + "QdkGUPVhjqu7sy7hGQqsGmg2YOOx9OIcyZQveNPljRpmWuE9NKMQ1pz6x49mEGfm",
		FirstName:    "Admin",
		Subject:      "Your account verification",
		TypeOfAction: "Register",
		SiteData:     siteData,
	}

	return c.Render("emails/verificationCode", emailData)
}

func (h *EmailHandler) verifySuccess(c *fiber.Ctx) error {
	siteData, _ := configs.GetSiteData(".")

	emailData := helpers.EmailData{
		Subject:  "Account verification success",
		SiteData: siteData,
	}

	return c.Render("emails/verificationDone", emailData)
}
