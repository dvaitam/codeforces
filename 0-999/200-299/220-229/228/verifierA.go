package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(a, b, c, d int) string {
	m := map[int]bool{a: true, b: true, c: true, d: true}
	return fmt.Sprintf("%d", 4-len(m))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	type T [4]int
	tests := []T{}
	for a := 1; len(tests) < 100 && a <= 4; a++ {
		for b := 1; len(tests) < 100 && b <= 4; b++ {
			for c := 1; len(tests) < 100 && c <= 4; c++ {
				for d := 1; len(tests) < 100 && d <= 4; d++ {
					tests = append(tests, T{a, b, c, d})
				}
			}
		}
	}
	for i, t := range tests {
		in := fmt.Sprintf("%d %d %d %d\n", t[0], t[1], t[2], t[3])
		exp := expected(t[0], t[1], t[2], t[3])
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: input=%q expected=%s got=%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("ok %d tests\n", len(tests))
}
