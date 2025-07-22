package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	l, r int64
}

func solveA(l, r int64) string {
	if r-l+1 < 3 || (r-l+1 == 3 && l%2 != 0) {
		return "-1\n"
	}
	if l%2 == 0 {
		return fmt.Sprintf("%d %d %d\n", l, l+1, l+2)
	}
	return fmt.Sprintf("%d %d %d\n", l+1, l+2, l+3)
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func generateTests() []testCaseA {
	tests := make([]testCaseA, 0, 100)
	for i := 0; i < 90; i++ {
		l := int64(i + 1)
		var diff int64
		switch i % 3 {
		case 0:
			diff = 0
		case 1:
			diff = 1
		default:
			diff = int64(i%49 + 2)
		}
		r := l + diff
		if diff > 50 {
			r = l + 50
		}
		tests = append(tests, testCaseA{l: l, r: r})
	}
	base := int64(900000000000000000)
	for i := int64(0); i < 10; i++ {
		l := base + i
		r := l + 50
		tests = append(tests, testCaseA{l: l, r: r})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d\n", t.l, t.r)
		expected := solveA(t.l, t.r)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expected) {
			fmt.Printf("case %d failed: l=%d r=%d expected %q got %q\n", i+1, t.l, t.r, strings.TrimSpace(expected), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
