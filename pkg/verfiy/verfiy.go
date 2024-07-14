package verfiy

import (
	"errors"
	"regexp"
	"strings"
)

// IsUsername 验证用户名合法性，用户名必须由5-12位(数字、字母、_、-)组成，且必须以字母开头。
func IsUsername(username string) error {
	if len(username) == 0 {
		return errors.New("请输入用户名")
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z_-]{5,12}$", username)
	if err != nil || !matched {
		return errors.New("用户名必须由5-12位(数字、字母、_、-)组成")
	}
	return nil
}

// IsEmail 验证是否是合法的邮箱
func IsEmail(email string) (err error) {
	if len(email) == 0 {
		err = errors.New("邮箱格式不符合规范")
		return
	}
	pattern := `^([A-Za-z0-9_\-\.])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,4})$`
	matched, _ := regexp.MatchString(pattern, email)
	if !matched {
		err = errors.New("邮箱格式不符合规范")
	}
	return
}

// IsValidPassword 是否是合法的密码
func IsValidPassword(password, rePassword string) error {
	if err := IsPassword(password); err != nil {
		return err
	}
	if password != rePassword {
		return errors.New("两次输入密码不匹配")
	}
	return nil
}

func IsPassword(password string) error {
	if len(password) == 0 {
		return errors.New("请输入密码")
	}
	if len(password) < 6 {
		return errors.New("密码过于简单")
	}
	if len(password) > 1024 {
		return errors.New("密码长度不能超过128")
	}
	return nil
}

// IsURL 是否是合法的URL
func IsURL(url string) error {
	if len(url) == 0 {
		return errors.New("URL格式错误")
	}
	indexOfHttp := strings.Index(url, "http://")
	indexOfHttps := strings.Index(url, "https://")
	if indexOfHttp == 0 || indexOfHttps == 0 {
		return nil
	}
	return errors.New("URL格式错误")
}

// func CategoryVer(gc *manage.MallGoodsCategory) bool {
// 	return contains([]int{1,2,3,4},gc.CategoryLevel)
// }

func Contains(arr []int, elem int) bool {
    for _, v := range arr {
        if v == elem {
            return true
        }
    }
    return false
}