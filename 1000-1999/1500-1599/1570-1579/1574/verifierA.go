package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TestCase struct {
	input string
	nvals []int
}

func genCase(rng *rand.Rand) TestCase {
	t := rng.Intn(4) + 1 // up to 4 testcases per run
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t) + "\n")
	nvals := make([]int, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(50) + 1
		nvals[i] = n
		sb.WriteString(fmt.Sprintf("%d\n", n))
	}
	return TestCase{sb.String(), nvals}
}

func runProg(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func checkOutput(tc TestCase, out string) error {
	out = strings.TrimSpace(out)
	lines := strings.Split(out, "\n")
	idx := 0
	for caseNum, n := range tc.nvals {
		if idx+n > len(lines) {
			return fmt.Errorf("not enough lines for case %d", caseNum+1)
		}
		seen := make(map[string]bool)
		for i := 0; i < n; i++ {
			seq := strings.TrimSpace(lines[idx])
			idx++
			if len(seq) != 2*n {
				return fmt.Errorf("case %d line %d: expected length %d got %d", caseNum+1, i+1, 2*n, len(seq))
			}
			bal := 0
			for _, ch := range seq {
				if ch == '(' {
					bal++
				} else if ch == ')' {
					bal--
				} else {
					return fmt.Errorf("case %d line %d: invalid char %q", caseNum+1, i+1, ch)
				}
				if bal < 0 {
					return fmt.Errorf("case %d line %d: not balanced", caseNum+1, i+1)
				}
			}
			if bal != 0 {
				return fmt.Errorf("case %d line %d: not balanced", caseNum+1, i+1)
			}
			if seen[seq] {
				return fmt.Errorf("case %d line %d: duplicate sequence", caseNum+1, i+1)
			}
			seen[seq] = true
		}
	}
	if idx != len(lines) {
		return fmt.Errorf("extra output lines")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		out, err := runProg(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
