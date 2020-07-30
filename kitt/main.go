package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pborman/getopt/v2"
)

type Config struct {
	passPath    string
	composePath string
	consulPath  string
	vaultPath   string
	pass        map[string]string
	env         map[string]string
}

// list of requred kitt vars
// must be added to the kitt directory in pass
func kittVars() []string {
	return []string{
		"CF_API_EMAIL",
		"CF_DNS_API_TOKEN",
		"CF_ZONE_API_TOKEN",
		"COMPOSE_FILE",
		"CONSUL_HTTP_TOKEN",
		"KITT_IP",
		"KITT_DOMAIN",
		"KITT_DIRECTORY",
		"KITT_TUNNEL_HOSTNAME",
		"KITT_TUNNEL_URL",
	}
}

//func init() {
//}

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

	vault, err := exec.LookPath("vault")
	if err != nil {
		fmt.Printf("cannot find vault in $PATH: %s\n", err)
		os.Exit(1)
	}

	// populate our config
	mypass := passEnv(pass)
	myenv := osEnv()
	conf := &Config{
		passPath:    pass,
		composePath: compose,
		consulPath:  consul,
		vaultPath:   vault,
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
	helpFlag := getopt.BoolLong("help", 'h', "display the help message")
	initFlag := getopt.BoolLong("init", 'a', "initialize kitt")
	startFlag := getopt.BoolLong("start", 's', "start kitt")
	stopFlag := getopt.BoolLong("stop", 'z', "stop kitt")
	getopt.Parse()

	switch {
	case *initFlag:
		// add vault init stuff here
		bootConsul(conf)
		os.Exit(0)
	case *startFlag:
		os.Exit(0)
	case *stopFlag:
		os.Exit(0)
	case *helpFlag:
		getopt.PrintUsage(os.Stderr)
	default:
		getopt.PrintUsage(os.Stderr)
	}

}
