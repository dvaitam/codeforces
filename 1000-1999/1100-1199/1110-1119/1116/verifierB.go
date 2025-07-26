package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		cmd := exec.Command(bin)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "run %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		valStr := strings.TrimSpace(out.String())
		val, err := strconv.Atoi(valStr)
		if err != nil || val < 0 || val > 2 {
			fmt.Fprintf(os.Stderr, "test %d: invalid output %q\n", i+1, valStr)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
