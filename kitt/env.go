package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func passEnv(path string) map[string]string {
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

func osEnv() map[string]string {
	var env = make(map[string]string)

	path, err := exec.LookPath("env")
	if err != nil {
		fmt.Printf("cannot find env in $PATH: %s\n", err)
		os.Exit(1)
	}

	var out bytes.Buffer
	cmd := exec.Command(path)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: env ", err)
		os.Exit(1)
	}
	one := strings.Split(out.String(), "\n")
	for _, field := range one {
		if field == "" {
			continue
		}
		two := strings.Split(field, "=")
		env[two[0]] = two[1]
	}

	return env
}

func flatOs(c *Config) []string {
	var flat []string

	for k, v := range c.env {
		flat = append(flat, fmt.Sprintf("%s=%s", k, v))
	}

	return flat
}

func flatPass(c *Config) []string {
	var flat []string

	for k, v := range c.pass {
		flat = append(flat, fmt.Sprintf("%s=%s", k, v))
	}

	return flat
}

func flatAll(c *Config) []string {
	envs := append(flatOs(c), flatPass(c)...)

	return envs
}

func missEnv(c *Config, kitt []string) []string {
	var missing []string
	var ya bool

	for _, v := range kitt {
		for k, _ := range c.pass {
			if k == v {
				ya = true
				break
			} else {
				ya = false
			}
		}
		if !ya {
			missing = append(missing, v)
		}
	}

	return missing
}
