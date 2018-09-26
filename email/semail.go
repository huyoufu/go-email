package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewMail() *mimeMail {
	return &mimeMail{}
}

type mimeMail struct {
	from, to, subject, body string
	auth                    smtp.Auth
	hasAttach               bool
	boundary                string
	attachments             []string
}

func (mm *mimeMail) SetAuth(account, password, host string) {
	mm.auth = smtp.PlainAuth("", account, password, host)
}
func (mm *mimeMail) SetFrom(fromStr string) {
	mm.from = fromStr
}
func (mm *mimeMail) SetTO(toStr string) {
	mm.to = toStr
}
func (mm *mimeMail) SetSubject(subjectStr string) {
	mm.subject = subjectStr
}
func (mm *mimeMail) SetBody(bodyStr string) {
	mm.body = bodyStr
}
func (mm *mimeMail) SetAttachments(att string) {
	if mm.hasAttach {
		mm.attachments = append(mm.attachments, att)
	} else {
		mm.attachments = []string{att}
		mm.hasAttach = true
		mm.boundary = mm.Boundary()
	}
}

func (mm *mimeMail) Boundary() string {
	return strconv.Itoa(int(time.Now().UnixNano()))
}

func (mm *mimeMail) toBuffer() *bytes.Buffer {
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(fmt.Sprintf("TO: %s\r\n", mm.to))
	buffer.WriteString(fmt.Sprintf("From: %s\r\n", mm.from))
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", mm.subject))
	buffer.WriteString("MIME-Version: 1.0\r\n")
	if !mm.hasAttach {
		//不带附件的
		buffer.WriteString("Content-Type: text/html;charset=utf-8\r\n\r\n")
		buffer.WriteString(mm.body)
	} else {
		//开始正文

		buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", mm.Boundary()))
		//html文本正文
		buffer.WriteString(fmt.Sprintf("--%s\r\n", mm.boundary))
		buffer.WriteString("Content-Type: text/html;charset=utf-8\r\n\r\n")
		buffer.WriteString(mm.body)
		//html内容结束
		for _, att := range mm.attachments {
			file, _ := ioutil.ReadFile(att)
			attachFileName := Substr(att, strings.LastIndex(att, string(os.PathSeparator))+1)
			attachFileName = base64AttachName(attachFileName)
			buffer.WriteString(fmt.Sprintf("\r\n--%s\r\n", mm.boundary))
			buffer.WriteString("Content-Type: application/octet-stream\r\n")
			buffer.WriteString("Content-Transfer-Encoding: base64\r\n")
			buffer.WriteString("Content-Description: 附件\r\n")
			buffer.WriteString("Charset: utf-8\r\n")
			buffer.WriteString("Content-Disposition: attachment; filename=\"" + attachFileName + "\"\r\n\r\n")
			b := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
			base64.StdEncoding.Encode(b, file)
			buffer.Write(b)
		}
		buffer.WriteString("\r\n--" + mm.boundary + "--\r\n\r\n")
	}
	log.Println(buffer.String())
	//fmt.Println(buffer.String())
	return buffer
}

func base64AttachName(attaFileName string) string {
	b := make([]byte, base64.StdEncoding.EncodedLen(len([]byte(attaFileName))))
	base64.StdEncoding.Encode(b, []byte(attaFileName))
	return "=?utf-8?B?" + string(b) + "?="
}

func Substr(str string, start int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}
	return string(rs[start:length])
}
