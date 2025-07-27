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

type testCaseD struct {
	s string
}

func solveCaseD(s string) int {
	runs := []int{}
	last := byte(0)
	for i := 0; i < len(s); i++ {
		if i == 0 || s[i] != last {
			runs = append(runs, 1)
			last = s[i]
		} else {
			runs[len(runs)-1]++
		}
	}
	ops := 0
	idx := 0
	m := len(runs)
	for idx < m {
		if idx == m-1 {
			ops++
			break
		}
		if runs[idx] > 1 {
			ops++
			idx++
		} else {
			ops++
			runs[idx+1]--
			idx++
			if idx < m && runs[idx] == 0 {
				idx++
			}
		}
	}
	return ops
}

func buildInputD(s string) string {
	return fmt.Sprintf("1\n%d\n%s\n", len(s), s)
}

func runCaseD(bin string, tc testCaseD) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputD(tc.s))
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
	expect := solveCaseD(tc.s)
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesD() []testCaseD {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseD, 0, 100)
	cases = append(cases, testCaseD{s: "0"}, testCaseD{s: "1"}, testCaseD{s: "01"}, testCaseD{s: "10"}, testCaseD{s: "1111"})
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		cases = append(cases, testCaseD{s: sb.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesD()
	for i, tc := range cases {
		if err := runCaseD(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (s=%s)\n", i+1, err, tc.s)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
