package main

import (
	"fmt"
        "io/ioutil"
	"os"
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
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		w.Close()
		os.Stdout = old
		fmt.Println("os.Pipe redirect Error: ", err)
	} else {
		os.Stdout = w
	}
	arg := []string{"network", "create", "kitt"}
	err, out := cli(c, c.dockerPath, arg, flatAll(c), strings.NewReader(""))
        w.Close()
        os.Stdout = old
        stdout, _ := ioutil.ReadAll(r)
	if err != nil {
		if strings.Contains(string(stdout), "network with name kitt already exists") {
			fmt.Println("kitt docker network exists")
		} else {
			fmt.Println(out+" Error: ", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("kitt docker network created")
	}
}

func dummyNet() {

}
