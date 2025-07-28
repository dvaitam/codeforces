package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin string) (string, error) {
	cmd := exec.Command(bin)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	expect := "Seriously"
	out, err := run(bin)
	if err != nil {
		fmt.Println(err)
		return
	}
	if out != expect {
		fmt.Printf("expected %q got %q\n", expect, out)
		return
	}
	fmt.Println("All tests passed")
}
