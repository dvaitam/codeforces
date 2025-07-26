package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n int
}

func expectedB(n int) int64 {
	var count int64
	for a := 1; a <= n; a++ {
		for b := a; b <= n; b++ {
			c := a ^ b
			if c < b || c > n {
				continue
			}
			if a+b <= c {
				continue
			}
			count++
		}
	}
	return count
}

func genTestsB() []testCaseB {
	rand.Seed(2)
	tests := make([]testCaseB, 0, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(2500) + 1 // 1..2500
		tests = append(tests, testCaseB{n: n})
	}
	return tests
}

func runCase(bin string, tc testCaseB) error {
	input := fmt.Sprintf("%d\n", tc.n)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprint(expectedB(tc.n))
	if got != expect {
		return fmt.Errorf("n=%d expected %s got %s", tc.n, expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsB()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
