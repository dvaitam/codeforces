package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct{ n int }

func solveA(n int) string {
	if n == 1 {
		return "9"
	}
	if n == 2 {
		return "98"
	}
	if n == 3 {
		return "989"
	}
	res := make([]byte, n)
	res[0] = '9'
	res[1] = '8'
	res[2] = '9'
	for i := 3; i < n; i++ {
		res[i] = byte('0' + (i-3)%10)
	}
	return string(res)
}

func buildInputA(n int) string {
	return fmt.Sprintf("1\n%d\n", n)
}

func runCaseA(bin string, tc testCaseA) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputA(tc.n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := solveA(tc.n)
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func generateCasesA() []testCaseA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseA, 0, 100)
	for _, n := range []int{1, 2, 3, 4, 5, 10, 11, 20, 50, 100} {
		cases = append(cases, testCaseA{n})
	}
	for len(cases) < 100 {
		cases = append(cases, testCaseA{rng.Intn(100) + 1})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesA()
	for i, tc := range cases {
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
