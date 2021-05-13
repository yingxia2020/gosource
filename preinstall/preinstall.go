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

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var (
	httpProxy  = "http://proxy-chain.intel.com:911"
	httpsProxy = "http://proxy-chain.intel.com:911"

	aptConf = `Acquire::http::proxy "%s";
Acquire::https::proxy "%s";`

	etcEnv = `http_proxy="%s"
https_proxy="%s"
no_proxy="localhost,127.0.0.1,.intel.com,%s,%s:6443,10.233.0.0/16"
`

	allYaml = `http_proxy: "%s"
https_proxy: "%s"
no_proxy: "localhost,127.0.0.1,.intel.com,%s,%s:6443,10.233.0.0/16"`

	allHttpProxy  = `http_proxy: "%s"`
	allHttpsProxy = `https_proxy: "%s"`
	allNoProxy    = `no_proxy: "localhost,127.0.0.1,.intel.com,%s,%s:6443,10.233.0.0/16"`
)

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[:idx]
}

func checkSystem() bool {
	out, err := exec.Command("lsb_release", "-d").CombinedOutput()
	if err != nil {
		log.Fatal(err.Error())
	}
	return strings.Contains(string(out), "Ubuntu 18")
}

func createAptConf(aptConfContent string) {
	// Remove the file if it exists already, best efforts
	os.Remove("/etc/apt/apt.conf")

	err := ioutil.WriteFile("/etc/apt/apt.conf", []byte(aptConfContent), 644)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func appendEtcEnv(input string) {
	// check if proxy has been set in this file, if so, do nothing
	content, err := ioutil.ReadFile("/etc/environment")
	if err != nil {
		log.Fatalf("Error reading /etc/environment file, %s", err)
	}

	if strings.Contains(string(content), input) {
		fmt.Println("You have set proxy in /etc/environment file before.")
		return
	}

	f, err := os.OpenFile("/etc/environment", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	if _, err := f.WriteString(input); err != nil {
		log.Fatal(err.Error())
	}
}

func replaceAllYaml(input []string) {
	var filepath = "./kubespray/inventory/cnb-cluster/group_vars/all/all.yml"

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading all.yaml file, %s", err)
	}

	// Only need replace no_proxy line in this file
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "http_proxy") || strings.HasPrefix(line, "#http_proxy") {
			lines[i] = input[0]
		}
		if strings.HasPrefix(line, "https_proxy") || strings.HasPrefix(line, "#https_proxy") {
			lines[i] = input[1]
		}
		if strings.HasPrefix(line, "no_proxy") || strings.HasPrefix(line, "#no_proxy") {
			lines[i] = input[2]
			break
		}
	}

	err = ioutil.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func removeProxyAllYaml() {
	var filepath = "./kubespray/inventory/cnb-cluster/group_vars/all/all.yml"

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatalf("Error reading all.yaml file, %s", err)
	}

	// Only need replace no_proxy line in this file
	lines := strings.Split(string(content), "\n")

	for i, line := range lines {
		if strings.HasPrefix(line, "http_proxy") {
			lines[i] = `#http_proxy: ""`
		}
		if strings.HasPrefix(line, "https_proxy") {
			lines[i] = `#https_proxy: ""`
		}
		if strings.HasPrefix(line, "no_proxy") {
			lines[i] = `#no_proxy: ""`
			break
		}
	}

	err = ioutil.WriteFile(filepath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func backupHostname() {
	// If the file is already there, do nothing
	if _, err := os.Stat("./hostname.back"); err == nil {
		return
	}
	err := exec.Command("cp", "/etc/hostname", "./hostname.back").Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func setDateTime() {
	_, err := exec.LookPath("date")
	if err != nil {
		log.Fatalf("Date binary not found, cannot set system date: %s\n", err.Error())
	}

	formatTime := `date -s "$(wget -qSO- --max-redirect=0 google.com 2>&1 | grep Date: | cut -d' ' -f5-8)Z"`
	err = exec.Command("bash", "-c", formatTime).Run()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	if !checkSystem() {
		fmt.Println("CNB needs to be run on Ubuntu 18.")
		os.Exit(0)
	}

	// backup /etc/hostname file
	fmt.Println("Backup hostname file\n")
	backupHostname()

	var input string
	fmt.Println("Set up proxy before run CNB installations")
	fmt.Printf(`Your test machine(s) under INTEL proxy and the proxy is "%s"? (Y/N): `, httpProxy)
	fmt.Scanln(&input)

	if strings.ToLower(input) == "n" {
		var hproxy, hsproxy string
		fmt.Print("\nPlease enter your http proxy (Enter for do not need proxy): ")
		fmt.Scanln(&hproxy)
		if len(hproxy) == 0 {
			fmt.Println("Remove proxy setups from kubespray all.yml file.")
			removeProxyAllYaml()
			os.Exit(0)
		}
		fmt.Print("\nPlease enter your https proxy: ")
		fmt.Scanln(&hsproxy)
		if len(hsproxy) == 0 {
			fmt.Println("HTTPS proxy value cannot be empty!")
			os.Exit(0)
		}
		httpProxy = hproxy
		httpsProxy = hsproxy
	}

	fmt.Println("The following changes will be made on your machine:")
	ip := GetOutboundIP()
	fmt.Printf("\n[/etc/apt/apt.conf] CREATE\n%s\n", fmt.Sprintf(aptConf, httpProxy, httpsProxy))
	fmt.Printf("\n[/etc/environment] APPEND\n%s\n", fmt.Sprintf(etcEnv, httpProxy, httpsProxy,
		ip, ip))
	fmt.Printf("\n[kubespray config file: all.yml] REPLACE\n%s\n", fmt.Sprintf(allYaml, httpProxy,
		httpsProxy, ip, ip))

	fmt.Print("\nAre you OK with the changes? (Y/N): ")
	fmt.Scanln(&input)

	if strings.ToLower(input) == "n" {
		fmt.Println("Please set up proxy manually before run CNB installations!")
		os.Exit(0)
	}

	createAptConf(fmt.Sprintf(aptConf, httpProxy, httpsProxy))

	appendEtcEnv(fmt.Sprintf(etcEnv, httpProxy, httpsProxy, ip, ip))

	replaceAllYaml([]string{fmt.Sprintf(allHttpProxy, httpProxy), fmt.Sprintf(allHttpsProxy, httpsProxy),
		fmt.Sprintf(allNoProxy, ip, ip)})

	fmt.Println("Please reboot your system for your changes to take effects before run CNB installations!")
}
