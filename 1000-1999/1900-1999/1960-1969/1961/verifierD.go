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

func pow(n, m int) int {
	res := 1
	for m > 0 {
		res *= n
		m--
	}
	return res
}

func runCase(bin string, n, m int) error {
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	var v int
	if _, err := fmt.Sscan(got, &v); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if v != pow(n, m) {
		return fmt.Errorf("expected %d got %d", pow(n, m), v)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	idx := 0
	for n := 0; n < 10; n++ {
		for m := 0; m < 10; m++ {
			idx++
			if err := runCase(bin, n, m); err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
