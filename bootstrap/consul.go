package main

import (
	"bufio"
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

// need to establish kitt directory

func dotenv() map[string]string {
	var envs = make(map[string]string)

	if _, err := os.Stat(".env"); err == nil {
		file, err := os.Open(".env")
		if err != nil {
			fmt.Println("Error: .env ", err)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			text := strings.Split(scanner.Text(), "=")
			if len(text) == 2 {
				envs[text[0]] = text[1]
			} else {
				fmt.Println("Error: parsing .env")
				fmt.Println("please check .env file format")
				os.Exit(1)
			}
		}
		file.Close()
	} else {
		fmt.Println("Error: .env ", err)
		fmt.Println("did you run make env?")
		os.Exit(1)
	}

	return envs
}

func pass(path string) map[string]string {
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
		pass[three] = res.String()
	}

	return pass
}

func myenv(path string) []string {
	var env []string

	pass := pass(path)
	fmt.Printf("%v", pass)

	//dotenv := dotenv()
	//fmt.Printf("%v", dotenv)

	val, present := os.LookupEnv("HOME")
	if present {
		env = append([]string{"HOME=" + val}, env...)
	} else {
		fmt.Println("Error: variable $HOME not set")
		os.Exit(1)
	}

	val, present = os.LookupEnv("KITT_DOMAIN")
	if present {
		env = append([]string{"KITT_DOMAIN=" + val}, env...)
	} else {
		fmt.Println("Error: variable $KITT_DOMAIN not set")
		os.Exit(1)
	}

	return env
}

func cli(path string, args []string, envs []string, in io.Reader) (error, string) {
	cli := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   cli,
		Env:    envs,
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
	var out string
	in := strings.NewReader("")

	pass, err := exec.LookPath("pass")
	if err != nil {
		fmt.Printf("cannot find pass in $PATH: %s\n", err)
		os.Exit(1)
	}
	myenv := myenv(pass)

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

	arg = []string{"up", "-d", "consul"}
	err, out = cli(compose, arg, myenv, in)
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
	env = append(myenv, "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	err, out = cli(consul, arg, env, in)
	for err != nil {
		time.Sleep(5 * time.Second)
		arg = []string{"info"}
		env = append(myenv, "CONSUL_HTTP_ADDR=169.254.32.1:8500")
		err, out = cli(consul, arg, env, in)
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
	env = append(myenv, "CONSUL_HTTP_ADDR=169.254.32.1:8500")
	err, out = cli(consul, arg, env, in)
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
		err, out = cli(pass, arg, myenv, strings.NewReader(str))
		if err != nil {
			fmt.Println(out+" Error: ", err)
			fmt.Println("unable to insert consul secret id acl token into pass")
			fmt.Println("please manually run:")
			fmt.Println("echo " + str + " | pass insert -e kitt/CONSUL_HTTP_TOKEN")
		} else {
			arg = []string{"git", "push"}
			err, out = cli(pass, arg, myenv, in)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				fmt.Println("please manually run: pass git push")
			}
		}
	} else {
		fmt.Printf("%s", stdout) // print any error message from consul acl bootstrap
	}

	arg = []string{"down"}
	err, out = cli(compose, arg, myenv, in)
	if err != nil {
		fmt.Println(out+" Error: ", err)
		os.Exit(1)
	}

}
