package main

import (
	"bytes"
	"io"
	"os/exec"
)

func cli(c *Config, path string, args []string, envs []string, in io.Reader) (string, string, error) {
	var stdout bytes.Buffer
	cli := append([]string{path}, args...)

	cmd := &exec.Cmd{
		Path:   path,
		Args:   cli,
		Env:    envs,
		Dir:    c.pass["KITT_DIRECTORY"],
		Stdin:  in,
		Stdout: &stdout,
		Stderr: &stdout,
	}

	err := cmd.Run()
	if err != nil {
		return cmd.String(), string(stdout.Bytes()), err
	}

	return cmd.String(), string(stdout.Bytes()), nil

}
