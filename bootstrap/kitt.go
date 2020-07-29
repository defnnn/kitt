package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func passenv(path string) map[string]string {
	var pass = make(map[string]string)

	var out bytes.Buffer
	cmd := exec.Command(path, "kitt")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: pass kitt ", err)
		fmt.Println("did you enter your kitt variables into pass?")
		os.Exit(1)
	}
	one := strings.Split(out.String(), "\n")
	for _, field := range one {
		if field == "kitt" {
			continue
		}
		if field == "" {
			continue
		}
		// get rid of that pesky └── & ├── pass likes to output
		two := strings.Split(field, "─")
		three := strings.TrimSpace(two[len(two)-1])
		// get the value of our pass key
		var res bytes.Buffer
		value := exec.Command(path, fmt.Sprintf("kitt/%s", three))
		value.Stdout = &res
		err := value.Run()
		if err != nil {
			fmt.Println("Error: pass kitt/"+three, err)
			os.Exit(1)
		}
		pass[three] = strings.TrimSpace(res.String())
	}

	return pass
}

func myenv(pass map[string]string, env string, opt bool) string {
	var out string

	if val, ok := pass[env]; ok {
		out = fmt.Sprintf("%s=%s", env, val)
		return out
	}

	val, present := os.LookupEnv(env)
	if present {
		out = fmt.Sprintf("%s=%s", env, val)
		return out
	} else {
		if !opt {
			fmt.Println("Error: variable " + env + " not set")
			os.Exit(1)
		}
	}

	return ""
}

func cli(path string, args []string, envs []string, dir string, in io.Reader) (error, string) {
	cli := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   cli,
		Env:    envs,
		Dir:    dir,
		Stdin:  in,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err := cmd.Run()
	if err != nil {
		return err, cmd.String()
	}

	return nil, cmd.String()

}

func main() {
	var arg []string
	var env []string
	var xtrenv []string
	var dir string
	var out string
	in := strings.NewReader("")

	// placeholder: until we have packaging solidified,
	// we'll use an os env var to find kitt's root directory
	val, present := os.LookupEnv("KITT_DIRECTORY")
	if present {
		dir = val
	} else {
		fmt.Println("Error: os variable KITT_DIRECTORY not set")
		os.Exit(1)
	}

	// ensure our os binaries exist and are in our $PATH
	pass, err := exec.LookPath("pass")
	if err != nil {
		fmt.Printf("cannot find pass in $PATH: %s\n", err)
		os.Exit(1)
	}

	compose, err := exec.LookPath("docker-compose")
	if err != nil {
		fmt.Printf("cannot find docker-compose in $PATH: %s\n", err)
		os.Exit(1)
	}

	consul, err := exec.LookPath("consul")
	if err != nil {
		fmt.Printf("cannot find consul in $PATH: %s\n", err)
		os.Exit(1)
	}

	// load our environment variables
	kenvs := passenv(pass)
	home := myenv(kenvs, "HOME", false)
	sshauth := myenv(kenvs, "SSH_AUTH_SOCK", true)
	sshpid := myenv(kenvs, "SSH_AGENT_PID", true)
	email := myenv(kenvs, "CF_API_EMAIL", false)
	dns := myenv(kenvs, "CF_DNS_API_TOKEN", false)
	api := myenv(kenvs, "CF_ZONE_API_TOKEN", true)
	http := myenv(kenvs, "CONSUL_HTTP_TOKEN", true)
	ip := myenv(kenvs, "KITT_IP", false)
	domain := myenv(kenvs, "KITT_DOMAIN", false)
	host := myenv(kenvs, "KITT_TUNNEL_HOSTNAME", false)
	url := myenv(kenvs, "KITT_TUNNEL_URL", false)
	xtrenv = append(xtrenv, home, email, api, dns, http, ip, domain, host, url)

	arg = []string{"up", "-d", "consul"}
	err, out = cli(compose, arg, xtrenv, dir, in)
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
	env = []string{home, "CONSUL_HTTP_ADDR=169.254.32.1:8500"}
	err, out = cli(consul, arg, env, dir, in)
	for err != nil {
		time.Sleep(5 * time.Second)
		arg = []string{"info"}
		env = []string{home, "CONSUL_HTTP_ADDR=169.254.32.1:8500"}
		err, out = cli(consul, arg, env, dir, in)
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
	env = []string{home, "CONSUL_HTTP_ADDR=169.254.32.1:8500"}
	err, out = cli(consul, arg, env, dir, in)
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
		err, out = cli(pass, arg, []string{home}, dir, strings.NewReader(str))
		if err != nil {
			fmt.Println(out+" Error: ", err)
			fmt.Println("unable to insert consul secret id acl token into pass")
			fmt.Println("please manually run:")
			fmt.Println("echo " + str + " | pass insert -e kitt/CONSUL_HTTP_TOKEN")
		} else {
			arg = []string{"git", "push"}
			err, out = cli(pass, arg, []string{home, sshauth, sshpid}, dir, in)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				fmt.Println("please manually run: pass git push")
			}
		}
	} else {
		fmt.Printf("%s", stdout) // print any error message from consul acl bootstrap cmd
	}

	arg = []string{"down"}
	err, out = cli(compose, arg, xtrenv, dir, in)
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}

}
