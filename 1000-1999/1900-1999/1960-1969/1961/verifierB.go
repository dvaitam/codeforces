package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
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

func expected(n int) int {
	return n * (n + 1) / 2
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	var v int
	if _, err := fmt.Sscan(got, &v); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if v != expected(n) {
		return fmt.Errorf("expected %d got %d", expected(n), v)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		if err := runCase(bin, i); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
