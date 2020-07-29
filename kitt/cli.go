package main

import (
	"io"
	"os"
	"os/exec"
)

func cli(c *Config, path string, args []string, envs []string, in io.Reader) (error, string) {
	cli := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   cli,
		Env:    envs,
		Dir:    c.pass["KITT_DIRECTORY"],
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
