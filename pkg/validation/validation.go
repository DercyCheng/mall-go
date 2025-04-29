package validation

import (
	"html"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"mall-go/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate

	// 正则表达式
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{3,20}$`)
	// 移除复杂的密码正则表达式，改用函数验证
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
)

func init() {
	validate = validator.New()

	// 注册结构体标签
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// 注册自定义验证
	_ = validate.RegisterValidation("username", validateUsername)
	_ = validate.RegisterValidation("password", validatePassword)
}

// ValidateStruct 验证结构体
func ValidateStruct(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			// 如果只需要返回第一个错误
			if len(ve) > 0 {
				e := ve[0]
				// 获取字段名和错误类型
				field := e.Field()
				tag := e.Tag()

				// 构建错误消息
				var message string
				switch tag {
				case "required":
					message = field + "不能为空"
				case "email":
					message = field + "必须是有效的电子邮件地址"
				case "min":
					message = field + "长度不能小于" + e.Param()
				case "max":
					message = field + "长度不能大于" + e.Param()
				case "username":
					message = field + "只能包含字母、数字、下划线、短横线和点，长度为3-20个字符"
				case "password":
					message = field + "必须包含数字、大小写字母和特殊字符，长度至少为8位"
				default:
					message = field + "格式不正确"
				}

				return errors.ValidationError(errors.CodeValidationError, message)
			}
		}
		return errors.ValidationError(errors.CodeValidationError, "输入验证失败")
	}
	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

// validateUsername 验证用户名（内部函数）
func validateUsername(fl validator.FieldLevel) bool {
	return ValidateUsername(fl.Field().String())
}

// ValidatePassword 验证密码
// 密码必须包含数字、大小写字母和特殊字符，长度至少为8位
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

// validatePassword 验证密码（内部函数）
func validatePassword(fl validator.FieldLevel) bool {
	return ValidatePassword(fl.Field().String())
}

// ValidateEmail 验证电子邮件
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePhone 验证电话号码
func ValidatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

// SanitizeString 清理字符串，防止XSS攻击
func SanitizeString(input string) string {
	// 转义HTML标签
	return html.EscapeString(strings.TrimSpace(input))
}

// SanitizeMap 清理map中所有字符串值
func SanitizeMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range input {
		switch v := value.(type) {
		case string:
			result[key] = SanitizeString(v)
		case map[string]interface{}:
			result[key] = SanitizeMap(v)
		default:
			result[key] = v
		}
	}
	return result
}
