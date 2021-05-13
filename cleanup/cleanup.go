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
	"os"
	"os/exec"
	"strings"
)

func main() {
	// call kubespray reset
	fmt.Println("Running kubespray reset ......")
	cmd := exec.Command("ansible-playbook", "-i", "./kubespray/inventory/cnb-cluster/hosts.ini",
		"--become", "--become-user=root", "./kubespray/reset.yml")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	// remove docker and images
	fmt.Println("Remove Docker and its images ......")
	cmd = exec.Command("apt-get", "purge", "-y", "docker-ce")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = exec.Command("rm", "-rf", "/var/lib/docker").Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = exec.Command("rm", "-rf", "/etc/docker").Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	err = exec.Command("rm", "-rf", "/var/lib/dockershim").Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	// If we backup hostname file before, copy it back
	var hostBack = "./hostname.back"
	if _, err = os.Stat(hostBack); err == nil {
		fmt.Println("Restore system hostname ......")
		err := exec.Command("cp", hostBack, "/etc/hostname").Run()
		if err != nil {
			log.Fatal(err.Error())
		}
		content, err := ioutil.ReadFile(hostBack)
		if err != nil {
			log.Fatalf("Error reading all.yaml file, %s", err)
		}
		err = exec.Command("hostname", strings.TrimSpace(string(content))).Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// If kube-manifests directory are created, remove it
	if _, err = os.Stat("/root/kube-manifests"); err == nil {
		fmt.Println("Remove /root/kube-manifests directory ......")
		err = exec.Command("rm", "-rf", "/root/kube-manifests").Run()
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
