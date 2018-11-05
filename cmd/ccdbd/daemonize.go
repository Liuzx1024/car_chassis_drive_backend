package main

import (
	"errors"
	"os"
	"syscall"
)

var errForkFail = errors.New("syscall fork fail")

func createDaemon() {
	if err := createChild(); err != nil {
		panic(err)
	}
	if _, err := syscall.Setsid(); err != nil {
		panic(err)
	}
	if err := createChild(); err != nil {
		panic(err)
	}
	if err := syscall.Chdir("/"); err != nil {
		panic(err)
	}
	syscall.Umask(0)
	for i := 0; i < 3; i++ {
		if err := syscall.Close(i); err != nil {
			panic(err)
		}
	}
}

func callFork() (uintptr, error) {
	r1, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	return r1, error(err)
}

func createChild() error {
	if pid, err := callFork(); err != nil || pid < 0 {
		return errForkFail
	} else if pid == 0 {
		os.Exit(0)
	}
	return nil
}
