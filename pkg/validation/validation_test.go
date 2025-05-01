package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用结构体
type TestUser struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required,password"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone" validate:"required"`
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		want     bool
	}{
		{"ValidUsername", "user123", true},
		{"ValidWithUnderscore", "user_123", true},
		{"ValidWithDash", "user-123", true},
		{"ValidWithDot", "user.123", true},
		{"TooShort", "us", false},
		{"TooLong", "usernameistoolongforvalidation", false},
		{"WithSpace", "user 123", false},
		{"WithSpecialChars", "user@123", false},
		{"WithChinese", "用户123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUsername(tt.username)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{"ValidPassword", "Test123!", true},
		{"ValidComplex", "P@ssw0rd123", true},
		{"TooShort", "Abc123!", false},
		{"NoUppercase", "test123!", false},
		{"NoLowercase", "TEST123!", false},
		{"NoDigit", "Testabcd!", false},
		{"NoSpecialChar", "Test1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePassword(tt.password)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"ValidEmail", "test@example.com", true},
		{"ValidWithDots", "test.user@example.com", true},
		{"ValidWithPlus", "test+user@example.com", true},
		{"ValidSubdomain", "test@sub.example.com", true},
		{"InvalidNoAt", "testexample.com", false},
		{"InvalidNoDomain", "test@", false},
		{"InvalidNoUsername", "@example.com", false},
		{"InvalidFormat", "test@example", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEmail(tt.email)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name  string
		phone string
		want  bool
	}{
		{"ValidPhone", "1234567890", true},
		{"ValidInternational", "+12345678901", true},
		{"ValidMinLength", "1234567", true},
		{"TooShort", "123456", false},
		{"WithLetters", "123456abc", false},
		{"WithSpaces", "123 456 7890", false},
		{"WithHyphens", "123-456-7890", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidatePhone(tt.phone)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateStruct(t *testing.T) {
	t.Run("ValidStruct", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Password: "Test123!",
			Email:    "test@example.com",
			Phone:    "1234567890",
		}
		err := ValidateStruct(user)
		assert.NoError(t, err)
	})

	t.Run("InvalidUsername", func(t *testing.T) {
		user := TestUser{
			Username: "t",  // 太短
			Password: "Test123!",
			Email:    "test@example.com",
			Phone:    "1234567890",
		}
		err := ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Password: "password",  // 不符合复杂度要求
			Email:    "test@example.com",
			Phone:    "1234567890",
		}
		err := ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		user := TestUser{
			Username: "testuser",
			Password: "Test123!",
			Email:    "invalidemail",  // 非法邮箱格式
			Phone:    "1234567890",
		}
		err := ValidateStruct(user)
		assert.Error(t, err)
	})

	t.Run("EmptyField", func(t *testing.T) {
		user := TestUser{
			Username: "",  // 空字段
			Password: "Test123!",
			Email:    "test@example.com",
			Phone:    "1234567890",
		}
		err := ValidateStruct(user)
		assert.Error(t, err)
	})
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"NormalString", "Hello World", "Hello World"},
		{"StringWithSpaces", "  Hello  ", "Hello"},
		{"ScriptTag", "<script>alert('XSS')</script>", "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;"},
		{"HTMLTags", "<b>Bold</b>", "&lt;b&gt;Bold&lt;/b&gt;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSanitizeMap(t *testing.T) {
	input := map[string]interface{}{
		"name":  "  John  ",
		"email": "<script>alert('XSS')</script>john@example.com",
		"nested": map[string]interface{}{
			"address": "<b>New York</b>",
		},
		"age": 30, // 非字符串类型保持不变
	}

	expected := map[string]interface{}{
		"name":  "John",
		"email": "&lt;script&gt;alert(&#39;XSS&#39;)&lt;/script&gt;john@example.com",
		"nested": map[string]interface{}{
			"address": "&lt;b&gt;New York&lt;/b&gt;",
		},
		"age": 30,
	}

	result := SanitizeMap(input)
	assert.Equal(t, expected, result)
}