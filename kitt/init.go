package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func backupDir(c *Config) {
	root := c.pass["KITT_DIRECTORY"]
	_, err := os.Stat(root + "/backup")
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(root+"/backup", 0755)
			if err != nil {
				fmt.Println("make backup dir Error: ", err)
				os.Exit(1)
			} else {
				fmt.Println("kitt backup dir created")
			}
		} else {
			fmt.Println("stat backup dir Error: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("kitt backup dir exists")
	}
}

func dockerNet(c *Config) {
	arg := []string{"network", "create", "kitt"}
	out, stdout, err := cli(c, c.dockerPath, arg, flatAll(c), strings.NewReader(""))
	if err != nil {
		if strings.Contains(string(stdout), "network with name kitt already exists") {
			fmt.Println("kitt docker network exists")
		} else {
			fmt.Println(out+" Error: ", err)
			fmt.Println(stdout)
			os.Exit(1)
		}
	} else {
		fmt.Println("kitt docker network created")
	}
}

func dummyNet(c *Config) {
	sudo, err := exec.LookPath("sudo")
	if err != nil {
		fmt.Printf("Error: cannot find sudo in $PATH: %s\n", err)
		os.Exit(1)
	}
	switch runtime.GOOS {
	case "linux":
		if !dummyCheck("dummy0") {
			arg := []string{"nmcli", "connection", "add", "con-name", "dummy0", "ifname", "dummy0", "type", "dummy", "ipv4.method", "manual", "ipv4.address", "169.254.32.1/32"}
			out, stdout, err := cli(c, sudo, arg, flatAll(c), strings.NewReader(""))
			fmt.Println(stdout)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				os.Exit(1)
			}
		}
	case "darwin":
		if !dummyCheck("lo0") {
			arg := []string{"ifconfig", "lo0", "alias", "169.254.32.1", "netmask", "255.255.255.255"}
			out, stdout, err := cli(c, sudo, arg, flatAll(c), strings.NewReader(""))
			fmt.Println(stdout)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				os.Exit(1)
			}
		}
	default:
		fmt.Println("Error: cannot add kitt dummy0 interface, your operating system is not supported")
		os.Exit(1)
	}
}

func dummyCheck(intface string) bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error interfaces: ", err)
		os.Exit(1)
	}
	for _, iface := range ifaces {
		if iface.Name == intface {
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error interface addresses: ", err)
				os.Exit(1)
			}
			for _, addr := range addrs {
				if addr.String() == "169.254.32.1/32" {
					fmt.Println("kitt dummy0 interface exists")
					return true
				}
			}
		}
	}
	return false
}
