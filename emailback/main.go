/* Copyright (C) Intel Corporation
 *
 * All Rights Reserved
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 *
 * Written by Ying Xia <ying.xia@intel.com>, 2019
 */

package main

func main() {

	sender := NewSender("cnbtest2019@gmail.com", "test2019")

	//The receiver needs to be in slice as the receive supports multiple receiver
	Receiver := []string{"ying.xia@intel.com", "rockcicak@yahoo.com"}

	Subject := "Testing HTLML Email from golang"
	message := `
	<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
	<html>
	<head>
	<meta http-equiv="content-type" content="text/html"; charset=ISO-8859-1">
	</head>
	<body>This is the body<br>
	<div class="moz-signature"><i><br>
	<br>
	Regards<br>
	Alex<br>
	<i></div>
	</body>
	</html>
	`
	bodyMessage := sender.WriteHTMLEmail(Receiver, Subject, message)

	sender.SendMail(Receiver, Subject, bodyMessage)
}
