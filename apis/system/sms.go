package system

import (
	"ferry/pkg/sms"
	"ferry/tools/app"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GenerateSMSHandler(c *gin.Context) {
	// 从请求参数中获取手机号
	phone := c.Query("phone")
	if phone == "" {
		app.Error(c, -1, nil, "手机号不能为空")
		return
	}

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
