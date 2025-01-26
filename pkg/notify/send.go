package notify

import (
	"bytes"
	"ferry/models/system"
	"ferry/pkg/logger"
	"ferry/pkg/notify/dingtalk"
	"ferry/pkg/notify/email"
	"ferry/pkg/sms"
	"fmt"
	"text/template"

	"github.com/spf13/viper"
)

/*
  @Author : lanyulei
  @同时发送多种通知方式
*/

type BodyData struct {
	SendTo        interface{} // 接受人
	EmailCcTo     []string    // 抄送人邮箱列表
	Subject       string      // 标题
	Classify      []int       // 通知类型
	Id            int         // 工单ID
	Title         string      // 工单标题
	Creator       string      // 工单创建人
	Priority      int         // 工单优先级
	PriorityValue string      // 工单优先级
	CreatedAt     string      // 工单创建时间
	Content       string      // 通知的内容
	Description   string      // 表格上面的描述信息
	ProcessId     int         // 流程ID
	Domain        string      // 域名地址
}

func (b *BodyData) ParsingTemplate() (err error) {
	// 读取模版数据
	var (
		buf bytes.Buffer
	)

	tmpl, err := template.ParseFiles("./static/template/email.html")
	if err != nil {
		return
	}

	b.Domain = viper.GetString("settings.domain.url")
	err = tmpl.Execute(&buf, b)
	if err != nil {
		return
	}

	b.Content = buf.String()

	return
}

func (b *BodyData) SendNotify() (err error) {
	var (
		emailList []string
		phoneList []string
	)

	switch b.Priority {
	case 1:
		b.PriorityValue = "正常"
	case 2:
		b.PriorityValue = "紧急"
	case 3:
		b.PriorityValue = "非常紧急"
	}

	for _, c := range b.Classify {
		switch c {
		case 1: // 邮件
			users := b.SendTo.(map[string]interface{})["userList"].([]system.SysUser)
			if len(users) > 0 {
				for _, user := range users {
					emailList = append(emailList, user.Email)
					phoneList = append(phoneList, user.Phone)
				}
				err = b.ParsingTemplate()
				if err != nil {
					logger.Errorf("模版内容解析失败，%v", err.Error())
					return
				}
				go email.SendMail(emailList, b.EmailCcTo, b.Subject, b.Content)
				dingtalkEnable := viper.GetBool("settings.dingtalk.enable")
				if dingtalkEnable {
					url := fmt.Sprintf("%s/#/process/handle-ticket?workOrderId=%d&processId=%d", b.Domain, b.Id, b.ProcessId)
					go dingtalk.SendDingMsg(phoneList, url, b.Title, b.Creator, b.PriorityValue, b.CreatedAt)
				}
				phoneEnable := viper.GetBool("settings.phone.enable")
				if phoneEnable {
					for index, phone := range phoneList {
						fmt.Printf("Phone #%d: %s\n", index+1, phone)
						go sms.SendSMS(phone, b.Title, b.Creator)
					}
				}
			}
		}
	}
	return
}
