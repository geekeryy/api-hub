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

	// 取消策略不存在
	CancelPolicyNotExistError = newerr(10020, "CANCEL_POLICY_NOT_EXIST")

	// 取消策略批量更新失败
	CancelPolicyBatchUpsertError = newerr(10021, "CANCEL_POLICY_BATCH_UPDATE_ERROR")

	// 取消策略已存在
	CancelPolicyExistsError = newerr(10022, "CANCEL_POLICY_EXISTS_ERROR")

	// 取消策略不能删除
	CancelPolicyCannotDeleteError = newerr(10023, "CANCEL_POLICY_CANNOT_DELETE_ERROR")

	// 取消策略不能编辑
	CancelPolicyCannotEditError = newerr(10024, "CANCEL_POLICY_CANNOT_EDIT_ERROR")

	// 取消策略不能添加
	CancelPolicyCannotAddError = newerr(10025, "CANCEL_POLICY_CANNOT_ADD_ERROR")

	// 取消策略创建失败
	CancelPolicyCreateError = newerr(10026, "CANCEL_POLICY_CREATE_ERROR")

	// 取消策略删除失败
	CancelPolicyDeleteError = newerr(10027, "CANCEL_POLICY_DELETE_ERROR")

	// 取消策略为空
	CancelPolicyIsNullError = newerr(10028, "CANCEL_POLICY_IS_NULL_ERROR")

	// 取消策略更新失败
	CancelPolicyUpdateError = newerr(10029, "CANCEL_POLICY_UPDATE_ERROR")

	// 取消策略编辑详情失败
	CancelPolicyEditDetailsError = newerr(10030, "CANCEL_POLICY_EDIT_DETAILS_ERROR")

	// 字段未翻译
	FieldNotTranslateError = newerr(10031, "FIELD_NOT_TRANSLATE_ERROR")

	// 类别名重复
	CategoryNameRepeatError = newerr(10032, "CATEGORY_NAME_REPEAT_ERROR")

	// 商品名重复
	CommodityNameRepeatError = newerr(10033, "COMMODITY_NAME_REPEAT_ERROR")

	// 无法删除已上架的商品
	CommodityUpCannotDeleteError = newerr(10034, "COMMODITY_UP_CANNOT_DELETE_ERROR")

	// 商品分类已绑定商品
	CategoryHasCommodityError = newerr(10035, "CATEGORY_HAS_COMMODITY_ERROR")

	// 套餐库存不足
	ComboStockNotEnoughError = newerr(10036, "COMBO_STOCK_NOT_ENOUGH_ERROR")

	// 商品已下架
	CommodityDownError = newerr(10037, "COMMODITY_DOWN_ERROR")

	// 商品已经删除
	CommodityAlreadyDeleteError = newerr(10038, "COMMODITY_ALREADY_DELETE_ERROR")

	// 套餐已经删除
	ComboAlreadyDeleteError = newerr(10039, "COMBO_ALREADY_DELETE_ERROR")

	// 套餐已经过期
	ComboExpiredError = newerr(10040, "COMBO_EXPIRED_ERROR")

	// 租赁结束时间不在售卖期内
	ComboLeaseEndTimeError = newerr(10041, "COMBO_LEASE_END_TIME_ERROR")

	// 罚金不能大于押金
	PenaltyTooBigError = newerr(10042, "PENALTY_TOO_BIG_ERROR")

	// 每个订单最多100种商品
	OrderCommodityLimitError = newerr(10043, "ORDER_COMMODITY_LIMIT_ERROR")

	// 当前时间已超过可取消时间，请联系商家
	OrderCantCancelForTimeError = newerr(10044, "ORDER_CANT_CANCEL_FOR_TIME_ERROR")

	// 已有进行中的子订单，无法取消
	OrderCantCancelForStatusError = newerr(10045, "ORDER_CANT_CANCEL_FOR_STATUS_ERROR")

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
