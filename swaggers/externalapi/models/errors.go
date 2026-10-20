package models

import "go/types"

// Описание ошибок по статус кодам из pkg/models/errors/errors.go.
// errorText формируется как auth.errors.<текст ошибки в camelCase>,
// additionalErrors.reason - публичная причина ошибки (Cause) либо текст ошибки.

type ErrAdditional struct {
	// @tg desc=`Причина ошибки`
	// @tg example=`email is empty`
	Reason string `json:"reason"`
}

type Err400 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.badRequest, auth.errors.invalidRequest, auth.errors.badMeta, auth.errors.valueNotFoundInCache, auth.errors.invalidCurrencyCode, auth.errors.currencyCodeIsEmpty, auth.errors.optionIdsAreEmpty, auth.errors.titleIsEmpty, auth.errors.invalidOptionIdsOptionDoesNotExistOrWasDeleted, auth.errors.invalidType, auth.errors.isActiveFlagAndActivatedAtAreBothSet, auth.errors.invalidTimestampFormat`
	// @tg example=`auth.errors.badRequest`
	ErrorText string `json:"errorText"`
	// @tg desc=`Некорректный запрос: невалидные параметры или тело запроса`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err401 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.invalidToken`
	// @tg example=`auth.errors.invalidToken`
	ErrorText string `json:"errorText"`
	// @tg desc=`Пользователь не авторизован: токен отсутствует, невалиден или истек`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err402 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.notEnoughMoney`
	// @tg example=`auth.errors.notEnoughMoney`
	ErrorText string `json:"errorText"`
	// @tg desc=`На балансе недостаточно средств`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err403 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.accessDenied, auth.errors.forbidden, auth.errors.statusCode403`
	// @tg example=`auth.errors.accessDenied`
	ErrorText string `json:"errorText"`
	// @tg desc=`Доступ запрещен: недостаточно прав для выполнения операции`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err404 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.notFound, auth.errors.canTGetSubscription, auth.errors.promotionDoesntExist, auth.errors.promotionIsNotActive`
	// @tg example=`auth.errors.notFound`
	ErrorText string `json:"errorText"`
	// @tg desc=`Запрашиваемый объект не найден`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err405 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.methodNotAllowed`
	// @tg example=`auth.errors.methodNotAllowed`
	ErrorText string `json:"errorText"`
	// @tg desc=`HTTP метод не поддерживается для данного пути`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err409 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.keyAlreadyExists, auth.errors.conflictCreatingOperationIdAlreadyExsist, auth.errors.objectInUse, auth.errors.countShouldBeGreaterThanZero, auth.errors.couldNotGetSubscription, auth.errors.couldNotUpdateSubscription, auth.errors.subscriptionNotFound, auth.errors.eventsNotFound, auth.errors.subscriptionIsAlreadyActivatedByPremium, auth.errors.subscriptionIsAlreadyActivated, auth.errors.alreadyHasActiveSubscription, auth.errors.cannotDeleteOptionBecauseItIsInUse`
	// @tg example=`auth.errors.keyAlreadyExists`
	ErrorText string `json:"errorText"`
	// @tg desc=`Конфликт состояния: объект уже существует или используется`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err422 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.unprocessableEntity, auth.errors.promotionValidateError`
	// @tg example=`auth.errors.unprocessableEntity`
	ErrorText string `json:"errorText"`
	// @tg desc=`Запрос корректен, но не может быть обработан: данные не прошли валидацию`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err429 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.limitExceeded`
	// @tg example=`auth.errors.limitExceeded`
	ErrorText string `json:"errorText"`
	// @tg desc=`Превышен лимит запросов, повторите позже`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}

type Err500 struct {
	// @tg example=``
	Data types.Nil `json:"data,omitempty"`
	// @tg desc=`Флаг показывающий, что ответ пришел с ошибкой`
	// @tg example=`true`
	Error bool `json:"error"`
	// @tg desc=`Заголовок ошибки. Возможные значения: auth.errors.internalServerError, а также ошибки вида auth.errors.failedTo...`
	// @tg example=`auth.errors.internalServerError`
	ErrorText string `json:"errorText"`
	// @tg desc=`Внутренняя ошибка сервера`
	AdditionalErrors ErrAdditional `json:"additionalErrors"`
}
