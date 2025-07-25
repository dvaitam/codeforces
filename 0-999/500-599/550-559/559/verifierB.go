package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	a, b string
}

func canonical(s string) string {
	if len(s)%2 == 1 {
		return s
	}
	mid := len(s) / 2
	left := canonical(s[:mid])
	right := canonical(s[mid:])
	if left < right {
		return left + right
	}
	return right + left
}

func solveCase(a, b string) string {
	if canonical(a) == canonical(b) {
		return "YES"
	}
	return "NO"
}

func runCase(bin string, tc testCase) error {
	input := tc.a + "\n" + tc.b + "\n"
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := solveCase(tc.a, tc.b)
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(2)
	cases := make([]testCase, 100)
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := range cases {
		n := rand.Intn(20) + 1
		var sb1, sb2 strings.Builder
		for j := 0; j < n; j++ {
			sb1.WriteRune(letters[rand.Intn(len(letters))])
			sb2.WriteRune(letters[rand.Intn(len(letters))])
		}
		cases[i] = testCase{sb1.String(), sb2.String()}
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
