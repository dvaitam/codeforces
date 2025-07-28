package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read testcasesC.txt:", err)
		os.Exit(1)
	}
	input := string(data)

	ref := "refC_bin"
	if err := exec.Command("go", "build", "-o", ref, "1886C.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	want, err := run("./"+ref, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "reference error:", err)
		os.Exit(1)
	}
	got, err := run(cand, input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if want != got {
		fmt.Println("expected:\n" + want)
		fmt.Println("got:\n" + got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
