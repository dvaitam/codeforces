package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func buildRef() string {
	ref := "refD_bin"
	cmd := exec.Command("go", "build", "-o", ref, "575D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		panic(fmt.Sprintf("failed to build reference: %v\n%s", err, string(out)))
	}
	return ref
}

func run(bin string) (string, error) {
	c := exec.Command(bin)
	var out bytes.Buffer
	c.Stdout = &out
	c.Stderr = &out
	err := c.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref := buildRef()
	defer os.Remove(ref)
	exp, err := run(ref)
	if err != nil {
		fmt.Println("reference failed:", err)
		return
	}
	for i := 0; i < 100; i++ {
		got, err := run(bin)
		if err != nil {
			fmt.Printf("binary failed on iteration %d: %v\n", i, err)
			return
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("mismatch on iteration %d\nexpected:%s\nactual:%s\n", i, exp, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
