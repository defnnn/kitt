package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func bootConsul(c *Config) {
	var arg []string
	var env []string
	var out string

	arg = []string{"up", "-d", "consul"}
	err, out := cli(c, c.composePath, arg, flatAll(c), strings.NewReader(""))
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}

	fmt.Println("waiting for consul to start...")
	// suppress stdout while we are going through the info loop
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		fmt.Println("os.Pipe redirect Error: ", err)
		w.Close()
		os.Stdout = old
	} else {
		os.Stdout = w
	}
	arg = []string{"info"}
	env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	err, out = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	for err != nil {
		time.Sleep(5 * time.Second)
		arg = []string{"info"}
		env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	        err, out = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	}
	w.Close()

	// redirect stdout to a var
	r, w, err = os.Pipe() // there's probably a better way to flush the previous pipe
	if err != nil {
		fmt.Println("os.Pipe redirect Error: ", err)
		w.Close()
		os.Stdout = old
	} else {
		os.Stdout = w
	}
	arg = []string{"acl", "bootstrap", "-format=json"}
	env = append(flatOs(c), "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	err, out = cli(c, c.consulPath, arg, env, strings.NewReader(""))
	if err != nil {
		fmt.Println(out+" Error: ", err)
		fmt.Println("looks like consul may already be bootstrapped")
	}
	w.Close()
	os.Stdout = old
	stdout, _ := ioutil.ReadAll(r)

	if json.Valid(stdout) {
		fmt.Println("acl bootstrap success!")
		var acl map[string]interface{}
		json.Unmarshal([]byte(stdout), &acl)
		secret := acl["SecretID"].(interface{})
		str := fmt.Sprintf("%v", secret)
		arg = []string{"insert", "-e", "kitt/CONSUL_HTTP_TOKEN"}
		err, out := cli(c, c.passPath, arg, flatOs(c), strings.NewReader(str))
		if err != nil {
			fmt.Println(out+" Error: ", err)
			fmt.Println("unable to insert consul secret id acl token into pass")
			fmt.Println("please manually run:")
			fmt.Println("echo " + str + " | pass insert -e kitt/CONSUL_HTTP_TOKEN")
		} else {
			arg = []string{"git", "push"}
		        err, out := cli(c, c.passPath, arg, flatOs(c), strings.NewReader(""))
			if err != nil {
				fmt.Println(out+" Error: ", err)
				fmt.Println("please manually run: pass git push")
			}
		}
	} else {
		fmt.Printf("%s", stdout) // print any error message from consul acl bootstrap cmd
	}

	arg = []string{"down"}
	err, out = cli(c, c.composePath, arg, flatAll(c), strings.NewReader(""))
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}
}
