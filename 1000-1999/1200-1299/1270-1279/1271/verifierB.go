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

func tryOps(orig []byte, target byte) ([]int, bool) {
	n := len(orig)
	a := make([]byte, n)
	copy(a, orig)
	ops := make([]int, 0, n)
	for i := 0; i < n-1; i++ {
		if a[i] != target {
			ops = append(ops, i+1)
			a[i] = target
			if a[i+1] == 'W' {
				a[i+1] = 'B'
			} else {
				a[i+1] = 'W'
			}
		}
	}
	for i := 0; i < n; i++ {
		if a[i] != target {
			return nil, false
		}
	}
	return ops, true
}

func expected(n int, s string) string {
	orig := []byte(s)
	if ops, ok := tryOps(orig, 'W'); ok {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(ops))
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		if len(ops) > 0 {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
		return strings.TrimRight(sb.String(), "\n")
	}
	if ops, ok := tryOps(orig, 'B'); ok {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(ops))
		for i, v := range ops {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		if len(ops) > 0 {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
		return strings.TrimRight(sb.String(), "\n")
	}
	return "-1"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(199) + 2 // 2..200
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 'W'
		} else {
			b[i] = 'B'
		}
	}
	s := string(b)
	input := fmt.Sprintf("%d\n%s\n", n, s)
	return input, expected(n, s)
}

func validate(input, output string) error {
	// parse input
	var n int
	var s string
	fmt.Sscan(input, &n, &s)
	orig := []byte(s)

	output = strings.TrimSpace(output)
	// check for impossible
	if output == "-1" {
		// verify truly impossible: neither all-W nor all-B reachable
		if _, ok := tryOps(orig, 'W'); ok {
			return fmt.Errorf("solution said -1 but a solution exists")
		}
		if _, ok := tryOps(orig, 'B'); ok {
			return fmt.Errorf("solution said -1 but a solution exists")
		}
		return nil
	}

	lines := strings.SplitN(output, "\n", 2)
	var k int
	if _, err := fmt.Sscan(lines[0], &k); err != nil {
		return fmt.Errorf("cannot parse k: %v", err)
	}
	if k < 0 || k > 3*n {
		return fmt.Errorf("k=%d out of range [0, 3*%d]", k, n)
	}

	// simulate operations
	a := make([]byte, n)
	copy(a, orig)
	if k > 0 {
		if len(lines) < 2 {
			return fmt.Errorf("expected ops line after k=%d", k)
		}
		fields := strings.Fields(lines[1])
		if len(fields) != k {
			return fmt.Errorf("expected %d ops, got %d", k, len(fields))
		}
		for _, f := range fields {
			var p int
			if _, err := fmt.Sscan(f, &p); err != nil {
				return fmt.Errorf("cannot parse op %q: %v", f, err)
			}
			if p < 1 || p >= n {
				return fmt.Errorf("op %d out of range [1, %d]", p, n-1)
			}
			// invert a[p-1] and a[p]
			if a[p-1] == 'W' {
				a[p-1] = 'B'
			} else {
				a[p-1] = 'W'
			}
			if a[p] == 'W' {
				a[p] = 'B'
			} else {
				a[p] = 'W'
			}
		}
	}
	// check all same
	for i := 1; i < n; i++ {
		if a[i] != a[0] {
			return fmt.Errorf("after operations blocks are not all the same color: %s", string(a))
		}
	}
	return nil
}

func runCase(exe, input, _ string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return validate(input, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
