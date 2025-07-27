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

type testCaseA struct {
	n int
}

func solveCaseA(n int) (int, int, int, bool) {
	for x := 0; x <= n/3; x++ {
		left := n - 3*x
		z := 0
		for left > 0 && left%5 != 0 {
			z++
			left -= 7
		}
		if left >= 0 && left%5 == 0 {
			y := left / 5
			if 3*x+5*y+7*z == n {
				return x, y, z, true
			}
		}
	}
	return 0, 0, 0, false
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
	fields := strings.Fields(out.String())
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if fields[0] == "-1" {
		if _, _, _, ok := solveCaseA(tc.n); ok {
			return fmt.Errorf("expected solution but got -1")
		}
		if len(fields) != 1 {
			return fmt.Errorf("unexpected extra output")
		}
		return nil
	}
	if len(fields) != 3 {
		return fmt.Errorf("expected three numbers got %d", len(fields))
	}
	var x, y, z int
	if _, err := fmt.Sscan(fields[0], &x); err != nil {
		return fmt.Errorf("bad x: %v", err)
	}
	if _, err := fmt.Sscan(fields[1], &y); err != nil {
		return fmt.Errorf("bad y: %v", err)
	}
	if _, err := fmt.Sscan(fields[2], &z); err != nil {
		return fmt.Errorf("bad z: %v", err)
	}
	if 3*x+5*y+7*z != tc.n || x < 0 || y < 0 || z < 0 {
		return fmt.Errorf("incorrect triple")
	}
	return nil
}

func generateCasesA() []testCaseA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseA, 0, 103)
	// fixed small cases
	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 30, 67, 100, 999, 1000} {
		cases = append(cases, testCaseA{n})
	}
	for len(cases) < 100 {
		cases = append(cases, testCaseA{rng.Intn(1000) + 1})
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
