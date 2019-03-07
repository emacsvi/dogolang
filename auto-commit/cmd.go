package main

import (
	"fmt"
	"os/exec"
	"time"
)

func autoCommit() (err error) {
	var (
		timeString string
	)
	// add 所有
	if err = execOneCmdAndPrint("git", "add", "--all"); err != nil {
		return err
	}

	// commit 所有
	timeString = time.Now().Format("2006-01-02 15:04:05")
	if err = execOneCmdAndPrint("git", "commit", "-a", "-m", timeString); err != nil {
		return err
	}

	// push到服务器
	if err = execOneCmdAndPrint("git", "push", "-u", "origin", "master"); err != nil {
		return err
	}
	return nil
}

// 执行单条命令
func execOneCmdAndPrint(name string, arg ...string) error {
	var (
		err error
		cmd *exec.Cmd
		out []byte
	)
	cmd = exec.Command(name, arg...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(out))

	return err
}
