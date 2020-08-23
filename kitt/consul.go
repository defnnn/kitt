package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

func bootConsul(c *Config) {
	var arg []string
	var env []string
	var out string

	arg = []string{"up", "-d", "consul"}
	out, stdout, err := cli(c, c.composePath, arg, flatAll(c), strings.NewReader(""))
	fmt.Println(stdout)
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}

	fmt.Println("waiting for consul to start...")
	arg = []string{"info"}
	env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	out, stdout, err = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	for err != nil { // possible infinite loop?
		time.Sleep(5 * time.Second)
		arg = []string{"info"}
		env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
		out, stdout, err = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	}

	arg = []string{"acl", "bootstrap", "-format=json"}
	env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	out, stdout, err = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	fmt.Println(stdout)
	if err != nil {
		fmt.Println(out+" Error: ", err)
		fmt.Println("looks like consul may already be bootstrapped")
	}

	if json.Valid([]byte(stdout)) {
		fmt.Println("acl bootstrap success!")
		var acl map[string]interface{}
		json.Unmarshal([]byte(stdout), &acl)
		secret := acl["SecretID"].(interface{})
		str := fmt.Sprintf("%v", secret)
		arg = []string{"insert", "-e", "kitt/CONSUL_HTTP_TOKEN"}
		out, stdout, err := cli(c, c.passPath, arg, flatOs(c), strings.NewReader(str))
		fmt.Println(stdout)
		if err != nil {
			fmt.Println(out+" Error: ", err)
			fmt.Println("unable to insert consul secret id acl token into pass")
			fmt.Println("please manually run:")
			fmt.Println("echo " + str + " | pass insert -e kitt/CONSUL_HTTP_TOKEN")
		} else {
			arg = []string{"git", "push"}
			out, stdout, err := cli(c, c.passPath, arg, flatOs(c), strings.NewReader(""))
			fmt.Println(stdout)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				fmt.Println("please manually run: pass git push")
			}
		}
	}

	arg = []string{"down"}
	out, stdout, err = cli(c, c.composePath, arg, flatAll(c), strings.NewReader(""))
	fmt.Println(stdout)
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}
}
