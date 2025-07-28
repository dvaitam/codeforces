package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type pair struct{ x, y int }

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

func runCase(bin string, p pair) error {
	input := fmt.Sprintf("1\n%d %d\n", p.x, p.y)
	expectedA, expectedB := p.x, p.y
	if expectedA > expectedB {
		expectedA, expectedB = expectedB, expectedA
	}
	got, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	var a, b int
	if _, err := fmt.Sscan(got, &a, &b); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if a != expectedA || b != expectedB {
		return fmt.Errorf("expected %d %d got %d %d", expectedA, expectedB, a, b)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]pair, 0, 100)
	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			cases = append(cases, pair{x, y})
		}
	}
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
