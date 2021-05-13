/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package email

import (
	"bytes"
	"context"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
	"time"
)

/**
	Modified from https://gist.github.com/jpillora/cb46d183eca0710d909a
**/

const (
	/**
		Gmail SMTP Server
	**/
	SMTPServer = "smtp.gmail.com"
)

const (
	messageTop = `
        <!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
        <html>
        <head>
        <meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
        </head>
        <body>`

	messageTail = `</body>
        </html>`
)

type Sender struct {
	User     string
	Password string
}

func NewSender(Username, Password string) Sender {

	return Sender{Username, Password}
}

// Sometimes email could get stuck there so need cancel after 3 seconds
func (sender Sender) SendEmailWithTimeout(receivers []string, subject string, bodyMessage string,
	ctx context.Context) bool {
	// Create a channel for signal handling
	c := make(chan bool)

	// Define a cancellation after 5s in the context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Run SendMail via a goroutine
	go func() {
		sender.sendMail(receivers, subject, bodyMessage, c)
	}()

	// Listening to signals
	select {
	case <-ctx.Done():
		return false
	case <-c:
		return true
	}
}

func (sender Sender) sendMail(Dest []string, Subject, bodyMessage string, c chan bool) {

	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, Dest, []byte(msg))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return
	}
	c <- true
}

func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)

	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	var encodedMessage bytes.Buffer
	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	for key, value := range header {
		finalMessage.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}

	finalMessage.Write([]byte("\r\n"))
	finalMessage.Write([]byte(messageTop))
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Write([]byte(messageTail))
	finalMessage.Close()

	return encodedMessage.String()
}

func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/html", subject, bodyMessage)
}

func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/plain", subject, bodyMessage)
}
