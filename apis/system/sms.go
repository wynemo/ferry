package system

import (
	"ferry/global/orm"
	"ferry/models/system"
	"ferry/pkg/sms"
	"ferry/tools/app"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

const (
	regular = "^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\\d{8}$"
)

func validate(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func GenerateSMSHandler(c *gin.Context) {
	// 从请求参数中获取手机号
	phone := c.Query("phone")
	if phone == "" {
		app.Error(c, -1, fmt.Errorf("手机号不能为空"), "手机号不能为空")
		return
	}
	var authUserCount int
	var err = orm.Eloquent.Model(&system.SysUser{}).
		Where("phone = ?", phone).
		Count(&authUserCount).Error
	if err != nil || authUserCount == 0 {
		app.Error(c, -1, fmt.Errorf("手机号不存在"), "手机号不存在")
		return
	}

	// if !validate(phone) {
	// 	app.Error(c, -1, fmt.Errorf("手机号格式不正确"), "手机号格式不正确")
	// 	return
	// }

	// 调用 GenerateSMSCode 方法生成验证码
	code, err := sms.GenerateSMSCode(phone)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("验证码获取失败, %v", err.Error()))
		return
	}
	fmt.Println("code is", code)

	err = sms.SendSMS(phone, "测试", code)
	if err != nil {
		app.Error(c, -1, err, fmt.Sprintf("验证码获取失败, %v", err.Error()))
		return
	}

	// 返回成功响应
	app.Custum(c, gin.H{
		"code":  200,
		"data":  nil,
		"phone": phone,
		"msg":   "success",
	})
}
