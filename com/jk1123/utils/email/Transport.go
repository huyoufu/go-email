package email

import (
	"log"
	"net/smtp"
)

var conf *Config

func init() {
	//初始化参数
	conf = parseYaml()
}
func SendEmail(to, subject, message string) {
	SendEmailWithAttachments(to, subject, message, nil)

}
func SendEmailWithAttachments(to, subject, message string, atts []string) {
	//初始化参数
	auth := conf.auth()
	newMail := NewMail()
	newMail.SetAuth(conf.Account, conf.Password, conf.Host)
	newMail.SetFrom(conf.From)
	newMail.SetTO(to)
	newMail.SetSubject(subject)
	newMail.SetBody(message)
	if atts != nil {
		for _, att := range atts {
			newMail.SetAttachments(att)
		}
	}
	//newMail.SetAttachments("C:\\Users\\Administrator\\Desktop\\ceshi.png")
	//newMail.SetAttachments("C:\\Users\\Administrator\\Desktop\\后台开发API文档.md")
	bf := newMail.toBuffer()
	e := smtp.SendMail(conf.Host+":"+conf.Port, auth, conf.Account, []string{to}, bf.Bytes())
	if e != nil {
		log.Println(e)
	}
}
