package customhandlers

import (
	interfaces "github.com/ARM-auto-place/auth/pkg/interfaces/externalapi"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const authAPIServiceName = "auth-api"

func UserRegister(ctx *fiber.Ctx, svc interfaces.ExternalAPI, email string, password string) error {
	return HandelResponse(ctx, "post", "/v1/auth/register", "UserRegister", authAPIServiceName, map[string]interface{}{
		"email":    email,
		"password": password,
	}, func() (interface{}, error) {
		return svc.UserRegister(ctx.UserContext(), email, password)
	})

}

func CheckCode(ctx *fiber.Ctx, svc interfaces.ExternalAPI, userID uuid.UUID, code string) error {
	return HandelResponse(ctx, "GET", "/v1/auth/checkCode", "CheckCode", authAPIServiceName, map[string]interface{}{
		"userID": userID,
		"code":   code,
	}, func() (interface{}, error) {
		return nil, svc.CheckCode(ctx.UserContext(), userID, code)
	})
}
