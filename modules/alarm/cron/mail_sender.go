// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cron

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/open-falcon/falcon-plus/modules/alarm/g"
	"github.com/open-falcon/falcon-plus/modules/alarm/model"
	"github.com/open-falcon/falcon-plus/modules/alarm/redi"
	"net/http"
	"time"
	"strings"
	"fmt"
	"encoding/json"
)

type EmailMessage struct {
	Subject string   `json:"subject"`
	Content string   `json:"content"`
	Tos     []string `json:"tos"`
}

func ConsumeMail() {
	for {
		L := redi.PopAllMail()
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendMailList(L)
	}
}

func SendMailList(L []*model.Mail) {
	for _, mail := range L {
		MailWorkerChan <- 1
		go SendMail(mail)
	}
}

func SendMail(mail *model.Mail) {
	defer func() {
		<-MailWorkerChan
	}()

	url := g.Config().Api.Mail
	var emailMessage EmailMessage
	emailMessage.Tos = strings.Split(mail.Tos, ",")
	emailMessage.Content = mail.Content
	emailMessage.Subject = mail.Subject

	jsonStr, err:= json.Marshal(emailMessage)
	if err != nil {
		fmt.Println(err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = 10 * time.Second
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("send mail fail, receiver:%s, subject:%s, content:%s, error:%v", mail.Tos, mail.Subject, mail.Content, err)
	}
	defer resp.Body.Close()

	log.Debugf("send mail:%v, resp:%v, url:%s", mail, resp, url)
}
