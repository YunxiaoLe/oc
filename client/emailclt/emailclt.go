package emailclt

import (
	"bytes"
	"github.com/OverClockTeam/OC-WordsInSongs-BE/conf"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"strings"
)

const (
	REGISTER_EMAIL = iota
	RESET_PASSWORD
)

var (
	EmailSettings conf.EmailSenderSettings
)

func InitEmailCtl(settings conf.EmailSenderSettings) {
	EmailSettings = settings
}

// ParseEmailTemplate TODO subject应该根据产品名称来确定 所以到时候需要改 and 需要前端支持html的编写
func ParseEmailTemplate(emailType int, verifyCode string) (subject, content string) {
	var t *template.Template
	var err error
	switch emailType {
	case RESET_PASSWORD:
		subject = "重置密码"
		t, err = template.ParseFiles("email_template/resetPassword.html")
	case REGISTER_EMAIL:
		subject = "邮箱验证"
		t, err = template.ParseFiles("email_template/hust-mail.html")
	}
	if err != nil {
		log.Println(err)
		return
	}
	buffer := new(bytes.Buffer)
	var data interface{}
	if err = t.Execute(buffer, data); err != nil {
		log.Println(err)
		return
	}

	content = strings.Replace(buffer.String(), "VerifyCodePlace", verifyCode, 1)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	content += mime
	return
}

//TODO 交给HY和LYX好了 不用GRPC先实现一个能用的即可 邮件服务你们熟悉

func RequestSendEmail(receiver, title, content string) {
	emailServerAddr := EmailSettings.ServerAddress
	mailHeader := map[string][]string{
		"From":    {emailServerAddr},
		"To":      {receiver},
		"Subject": {title},
	}
	m := gomail.NewMessage()
	m.SetHeaders(mailHeader)
	m.SetBody("text/html", content)
	d := gomail.NewDialer(EmailSettings.ServerHost, EmailSettings.ServerPort, EmailSettings.ServerAddress, EmailSettings.ServerPassword)
	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
	// Set up a connection to the server.
	/*emailServerAddr := EmailSettings.ServerAddress + ":" + strconv.Itoa(EmailSettings.ServerPort)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		conn, err := grpc.DialContext(ctx, emailServerAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			logrus.Error(err)
			return
		}
		defer conn.Close()
		c := email.NewEmailServiceClient(conn)

		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		r, err := c.SendEmail(ctx, &email.SendEmailInfo{ReceiveEmail: receiver, Title: title,
			Content: content}) //buffer.String()
		if err != nil {
			log.Println("could not greet: ", err.Error())
		}
		log.Println("Greeting: ", r.GetMessage())
	}*/
}
