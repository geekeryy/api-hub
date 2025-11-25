package xerror

import (
	coreError "github.com/geekeryy/api-hub/core/error"
)

// 错误码参考：https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Reference/Status
// 内置系统错误 (100~599)
var (
	// 400 Bad Request
	InvalidParameterErr = newerr(400, "INVALID_PARAMETER")

	// 401 Unauthorized 通知前端刷新或请求Token
	UnauthorizedErr = newerr(401, "UNAUTHORIZED")

	// 403 Forbidden
	ForbiddenErr = newerr(403, "FORBIDDEN")

	// 404 Not Found
	NotFoundErr = newerr(404, "NOT_FOUND")

	// 409 Conflict
	ConflictErr = newerr(409, "CONFLICT")

	// 429 Too Many Requests
	TooManyRequestsErr = newerr(429, "TOO_MANY_REQUESTS")

	// 500 Internal Server Error
	InternalServerErr = newerr(500, "INTERNAL_SERVER_ERROR")

	// 504 Gateway Timeout
	GatewayTimeoutErr = newerr(504, "GATEWAY_TIMEOUT")
)

// 框架错误(1000~9999)
var (
	// 未知错误
	UnknownErr = newerr(1000, "UNKNOWN_ERROR")

	// 数据库错误
	DBErr = newerr(1001, "DB_ERROR")

	// 请求频率过高
	RequestRateLimitErr = newerr(1002, "REQUEST_RATE_LIMIT_ERR")

	// 请求失败
	RequestFailedErr = newerr(1003, "REQUEST_FAILED_ERR")
)

// 自定义错误（10000~99999）
var (
	// 会员注册失败
	MemberRegisterErr = newerr(10000, "MEMBER_REGISTER")

	// 邮箱已被注册
	EmailHasRegisteredErr = newerr(10001, "EMAIL_HAS_REGISTERED")

	// 邮箱未注册
	EmailNotRegisteredErr = newerr(10002, "EMAIL_NOT_REGISTERED")

	// 手机号未注册
	PhoneNotRegisteredErr = newerr(10003, "PHONE_NOT_REGISTERED")

	// 手机号已注册
	PhoneHasRegisteredErr = newerr(10004, "PHONE_HAS_REGISTERED")

	// 密码错误
	PasswordErr = newerr(10005, "PASSWORD_ERR")

	// 用户不存在
	MemberNotFoundErr = newerr(10006, "MEMBER_NOT_FOUND")

	// 邮件发送失败
	MailSendErr = newerr(10007, "MAIL_SEND_ERR")

	// 短信发送失败
	SmsSendErr = newerr(10008, "SMS_SEND_ERR")

	// 验证码错误
	VerificationCodeErr = newerr(10009, "VERIFICATION_CODE_ERR")

	// 管理员不存在
	AdminNotFoundErr = newerr(10010, "ADMIN_NOT_FOUND")

	// 邮箱验证码错误
	EmailVerificationCodeErr = newerr(10011, "EMAIL_VERIFICATION_CODE_ERR")

	// 手机验证码错误
	PhoneVerificationCodeErr = newerr(10012, "PHONE_VERIFICATION_CODE_ERR")

	// 翻译错误
	TranslationErr = newerr(10013, "TRANSLATION_ERR")

	// 管理员已存在
	AdminExistsErr = newerr(10014, "ADMIN_EXISTS")

	// 密码错误次数超过限制，请等待 {{.num}} 分钟后再尝试!
	LoginLimitError = newerr(10015, "LOGIN_LIMIT_ERROR")

	// 登录错误，还可以尝试 {{.num}} 次
	LoginError = newerr(10016, "LOGIN_ERROR")

	// 登录错误，还可以尝试 {{.num}} 次
	LoginRecaptchaError = newerr(10017, "LOGIN_RECAPTCHA_ERROR")

	// 请求频率过高
	RequestRateLimitError = newerr(10018, "REQUEST_RATE_LIMIT_ERROR")

	// 数据库错误
	DBError = newerr(10019, "DB_ERROR")

	// 文件上传限制，最大{.limit}MB
	FileUploadLimitError = newerr(10046, "FILE_UPLOAD_LIMIT_ERROR")

	// 邮箱格式错误
	EmailFormatError = newerr(10047, "EMAIL_FORMAT_ERROR")
)

func newerr(c int64, msg string) *coreError.Error {
	return &coreError.Error{
		Code:         c,
		MessageId:    msg,
		Plural:       0,
		TemplateDate: make(map[string]string),
	}
}

func New(err error, e *coreError.Error) *coreError.Error {
	slacks := coreError.Callers()
	if len(slacks) > 1 {
		slacks = slacks[len(slacks)-1:]
	}
	return &coreError.Error{
		Code:          e.Code,
		MessageId:     e.MessageId,
		Plural:        0,
		TemplateDate:  make(map[string]string),
		OriginalError: err.Error(),
		Slacks:        slacks,
	}
}
