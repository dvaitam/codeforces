package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		cmd := exec.Command(bin)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("run %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if cmd.ProcessState.ExitCode() != 0 {
			fmt.Printf("run %d exited with %d\n", i+1, cmd.ProcessState.ExitCode())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
