package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	ref := "./refB.bin"
	if err := exec.Command("go", "build", "-o", ref, "1007B.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}

	want, err := run(ref, data)
	if err != nil {
		fmt.Println("reference runtime error:", err)
		os.Exit(1)
	}
	got, err := run(bin, data)
	if err != nil {
		fmt.Println("candidate runtime error:", err)
		os.Exit(1)
	}
	if want != got {
		fmt.Println("outputs differ")
		fmt.Println("expected:\n", want)
		fmt.Println("got:\n", got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
