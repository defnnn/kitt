package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Config struct {
	passPath    string
	composePath string
	consulPath  string
	pass        map[string]string
	env         map[string]string
}

func kittVars() []string {
	return []string{
		"CF_API_EMAIL",
		"CF_DNS_API_TOKEN",
		"CF_ZONE_API_TOKEN",
		"KITT_IP",
		"KITT_DOMAIN",
		"KITT_DIRECTORY",
		"KITT_TUNNEL_HOSTNAME",
		"KITT_TUNNEL_URL",
	}
}

func main() {

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

	// populate our config
	mypass := passEnv(pass)
	myenv := osEnv()
	conf := &Config{
		passPath:    pass,
		composePath: compose,
		consulPath:  consul,
		pass:        mypass,
		env:         myenv,
	}

	// check for missing kitt vars
	miss := missEnv(conf, kittVars())
	if len(miss) > 0 {
		str := strings.Join(miss, " ")
		fmt.Println("Error: missing essential kitt vars in pass.", str)
		os.Exit(1)
	}

	// shall we begin?

	// next step.... add getopts

	bootConsul(conf)

}
