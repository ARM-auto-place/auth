package externalapi

import (
	"context"

	"github.com/google/uuid"
)

// @tg version=0.0.1
// @tg backend=auth-api
// @tg title=`auth API`
// @tg servers=
//
//go:generate tg transport --services . --out ../../../internal/transport/jsonRPC/externalapi --outSwagger ../../../swaggers/externalapi/swagger.yaml

// ExternalAPI
// @tg http-server metrics log
// @tg http-prefix=/api/auth
// @tg 200=github.com/ARM-auto-place/auth/swaggers/externalapi/models:Resp200
// @tg 402=github.com/ARM-auto-place/auth/swaggers/externalapi/models:Err402
// @tg 403=github.com/ARM-auto-place/auth/swaggers/externalapi/models:Err403
// @tg 405=github.com/ARM-auto-place/auth/swaggers/externalapi/models:Err405
// @tg 500=github.com/ARM-auto-place/auth/swaggers/externalapi/models:Err500
type ExternalAPI interface {
	// UserRegister ...
	// @tg http-method=POST
	// @tg http-path=/v1/register
	// @tg http-args=email|email
	// @tg http-args=password|password
	// @tg uuidPackage=github.com/google/uuid
	// @tg http-response=github.com/ARM-auto-place/auth/internal/transport/custom_handlers:GetPayments
	// @tg summary=`Регистрация пользователя`
	// @tg desc=`Регистрация пользователя`
	UserRegister(ctx context.Context, email string, password string) (userID uuid.UUID, err error)
	// UserRegister ...
	// @tg http-method=GET
	// @tg http-path=/v1/checkCode
	// @tg http-cookies=userID|x-user-id
	// @tg http-args=code|code
	// @tg uuidPackage=github.com/google/uuid
	// @tg http-response=github.com/ARM-auto-place/auth/internal/transport/custom_handlers:CheckCode
	// @tg summary=`Проверка кода верефикации`
	// @tg desc=`Проверка кода верефикации`
	CheckCode(ctx context.Context, userID uuid.UUID, code string) (err error)

	// TODO: arm-3 this task 3
	// RenewVerefyCode(ctx context.Context, userID uuid.UUID)(err error)
	// Login(ctx context.Context, email string, password string)(userID uuid.UUID,err error)
	// TODO: Next step: change password,change email
}
