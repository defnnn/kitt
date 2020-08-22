package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type Config struct {
	passPath    string
	dockerPath  string
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

func main() {

	// ensure our os binaries exist and are in our $PATH
	pass, err := exec.LookPath("pass")
	if err != nil {
		fmt.Printf("cannot find pass in $PATH: %s\n", err)
		os.Exit(1)
	}

	docker, err := exec.LookPath("docker")
	if err != nil {
		fmt.Printf("cannot find docker in $PATH: %s\n", err)
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
		dockerPath:  docker,
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
	var arg []string

	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "initialize kitt",
		Long:  `init will bootstrap a brand new kitt instance.`,
		Run: func(cmd *cobra.Command, args []string) {
			backupDir(conf)
			dockerNet(conf)
			// add dummy0 interface
			bootConsul(conf)
			// initialize vault
		},
	}

	var cmdStart = &cobra.Command{
		Use:   "start",
		Short: "start kitt",
		Long:  `start will run all kitt services configured in COMPOSE_FILE.`,
		Run: func(cmd *cobra.Command, args []string) {
			arg = []string{"up", "-d"}
			err, out, stdout := cli(conf, compose, arg, flatAll(conf), strings.NewReader(""))
			fmt.Println(stdout)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				os.Exit(1)
			}

		},
	}

	var cmdStop = &cobra.Command{
		Use:   "stop",
		Short: "stop kitt",
		Long:  `stop will shutdown a running kitt service.`,
		Run: func(cmd *cobra.Command, args []string) {
			arg = []string{"down"}
			err, out, stdout := cli(conf, compose, arg, flatAll(conf), strings.NewReader(""))
			fmt.Println(stdout)
			if err != nil {
				fmt.Println(out+" Error: ", err)
				os.Exit(1)
			}

		},
	}

	var rootCmd = &cobra.Command{Use: "kitt"}
	rootCmd.AddCommand(cmdInit, cmdStart, cmdStop)
	rootCmd.Execute()

}
