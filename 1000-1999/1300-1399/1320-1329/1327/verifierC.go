package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n      int
	m      int
	k      int
	starts [][2]int
	ends   [][2]int
}

func validate(tc testCase, out string) error {
	toks := strings.Fields(out)
	if len(toks) == 0 {
		return fmt.Errorf("empty output")
	}
	if toks[0] == "-1" {
		// A construction always exists within 2*n*m (standard full-grid sweep).
		return fmt.Errorf("reported impossible, but a valid sequence always exists")
	}
	l, err := strconv.Atoi(toks[0])
	if err != nil {
		return fmt.Errorf("invalid first token (operations count)")
	}
	if l < 0 || l > 2*tc.n*tc.m {
		return fmt.Errorf("invalid operations count %d", l)
	}
	path := ""
	if l > 0 {
		if len(toks) != 2 {
			return fmt.Errorf("expected count and path only")
		}
		path = toks[1]
		if len(path) != l {
			return fmt.Errorf("count/path length mismatch: %d vs %d", l, len(path))
		}
		for i := 0; i < len(path); i++ {
			c := path[i]
			if c != 'U' && c != 'D' && c != 'L' && c != 'R' {
				return fmt.Errorf("invalid move character %q", c)
			}
		}
	} else if len(toks) != 1 {
		return fmt.Errorf("unexpected extra output tokens")
	}

	pos := make([][2]int, tc.k)
	seen := make([]bool, tc.k)
	for i := 0; i < tc.k; i++ {
		pos[i] = tc.starts[i]
		if pos[i] == tc.ends[i] {
			seen[i] = true
		}
	}

	for i := 0; i < len(path); i++ {
		mv := path[i]
		for j := 0; j < tc.k; j++ {
			x, y := pos[j][0], pos[j][1]
			switch mv {
			case 'U':
				if x > 1 {
					x--
				}
			case 'D':
				if x < tc.n {
					x++
				}
			case 'L':
				if y > 1 {
					y--
				}
			case 'R':
				if y < tc.m {
					y++
				}
			}
			pos[j] = [2]int{x, y}
			if pos[j] == tc.ends[j] {
				seen[j] = true
			}
		}
	}

	for i := 0; i < tc.k; i++ {
		if !seen[i] {
			return fmt.Errorf("chip %d never visits target", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(3)
	t := 100
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		k := rand.Intn(5) + 1
		st := make([][2]int, k)
		ed := make([][2]int, k)
		for j := 0; j < k; j++ {
			st[j] = [2]int{rand.Intn(n) + 1, rand.Intn(m) + 1}
		}
		for j := 0; j < k; j++ {
			ed[j] = [2]int{rand.Intn(n) + 1, rand.Intn(m) + 1}
		}
		tests[i] = testCase{n: n, m: m, k: k, starts: st, ends: ed}
	}

	for idx := 0; idx < t; idx++ {
		tc := tests[idx]
		n, m, k := tc.n, tc.m, tc.k
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for j := 0; j < k; j++ {
			input.WriteString(fmt.Sprintf("%d %d\n", tc.starts[j][0], tc.starts[j][1]))
		}
		for j := 0; j < k; j++ {
			input.WriteString(fmt.Sprintf("%d %d\n", tc.ends[j][0], tc.ends[j][1]))
		}
		in := input.String()

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(in)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Runtime error: %v\n%s", err, out.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if err := validate(tc, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nGot:\n%s\n", idx+1, err, in, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
