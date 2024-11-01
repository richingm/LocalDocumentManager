package domain

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"os/exec"
	"richingm/LocalDocumentManager/configs"
	"time"
)

type BackupBiz struct {
}

func NewBackupBiz() *BackupBiz {
	return &BackupBiz{}
}

func (b *BackupBiz) Save() error {
	backupPath := fmt.Sprintf("/tmp/backup_%s.sql", time.Now().Format("20060102150405"))

	err := b.backDb(backupPath)
	if err != nil {
		return err
	}

	//mysqlConf := configs.ConfigXx.MysqlConfig
	//
	//conf := configs.ConfigXx.EmailConfig
	//err = b.sendEmail([]string{conf.Sender}, "网站备份", "备份数据库:"+mysqlConf.Dbname, backupPath)
	//if err != nil {
	//	return err
	//}

	//// 删除备份文件
	//err = os.Remove(backupPath)
	//if err != nil {
	//	return err
	//}

	return nil
}

func (b *BackupBiz) backDb(backupPath string) error {
	conf := configs.ConfigXx.MysqlConfig
	cmd := exec.Command("/usr/local/bin/mysqldump", "-u"+conf.Username, "-p"+conf.Password, conf.Dbname)
	outfile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run mysqldump: %v", err)
	}

	return nil
}

func (b *BackupBiz) sendEmail(mailTo []string, mailTitle, mailBody, attach string) error {
	conf := configs.ConfigXx.EmailConfig
	m := gomail.NewMessage()
	m.SetHeader("From", conf.Sender)  //发送者腾讯邮箱账号
	m.SetHeader("To", mailTo...)      //接收者邮箱列表
	m.SetHeader("Subject", mailTitle) //邮件标题
	m.SetBody("text/html", mailBody)  //邮件内容,可以是html
	m.Attach(attach)                  // 添加附件
	d := gomail.NewDialer(conf.Smtp, conf.SmtpCode, conf.Sender, conf.AuthCode)
	err := d.DialAndSend(m)
	return err
}
