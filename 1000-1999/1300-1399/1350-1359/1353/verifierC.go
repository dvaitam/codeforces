package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testC struct {
	n int64
}

func genTestsC() []testC {
	rand.Seed(3353)
	tests := make([]testC, 0, 100)
	tests = append(tests, testC{1}, testC{3}, testC{5})
	for len(tests) < 100 {
		// n must be odd and <= 500000
		n := int64(rand.Intn(250000))*2 + 1
		tests = append(tests, testC{n})
	}
	return tests[:100]
}

func expectedC(n int64) int64 {
	return n * (n*n - 1) / 3
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String() + errBuf.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTestsC()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		exp := fmt.Sprintf("%d", expectedC(tc.n))
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
