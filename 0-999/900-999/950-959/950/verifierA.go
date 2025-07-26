package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	l, r, a int
}

func expected(l, r, a int) int {
	total := l + r + a
	half := total / 2
	if half > l+a {
		half = l + a
	}
	if half > r+a {
		half = r + a
	}
	return 2 * half
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d\n", tc.l, tc.r, tc.a)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(tc.l, tc.r, tc.a)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := make([]testCase, 0, 100)
	for l := 0; l < 10; l++ {
		for r := 0; r < 10; r++ {
			a := (l * r) % 10
			cases = append(cases, testCase{l, r, a})
		}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
