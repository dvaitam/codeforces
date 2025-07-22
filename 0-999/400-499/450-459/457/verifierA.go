package main

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var s, t string
	fmt.Fscan(rdr, &s, &t)
	n, m := len(s), len(t)
	L := n
	if m > L {
		L = m
	}
	pos := -1
	for i := 0; i < L; i++ {
		var cs, ct byte
		if i < L-n {
			cs = '0'
		} else {
			cs = s[i-(L-n)]
		}
		if i < L-m {
			ct = '0'
		} else {
			ct = t[i-(L-m)]
		}
		if cs != ct {
			pos = i
			break
		}
	}
	if pos == -1 {
		return "="
	}
	var cs, ct byte
	if pos < L-n {
		cs = '0'
	} else {
		cs = s[pos-(L-n)]
	}
	if pos < L-m {
		ct = '0'
	} else {
		ct = t[pos-(L-m)]
	}
	var sign0 float64
	if cs == '1' {
		sign0 = 1.0
	} else {
		sign0 = -1.0
	}
	invQ := 2.0 / (1.0 + math.Sqrt(5.0))
	tval := sign0
	power := 1.0
	for j := pos + 1; j < L; j++ {
		power *= invQ
		if power < 1e-18 {
			break
		}
		if j < L-n {
			cs = '0'
		} else {
			cs = s[j-(L-n)]
		}
		if j < L-m {
			ct = '0'
		} else {
			ct = t[j-(L-m)]
		}
		if cs == ct {
			continue
		}
		if cs == '1' {
			tval += power
		} else {
			tval -= power
		}
	}
	const eps = 1e-9
	if tval > eps {
		return ">"
	} else if tval < -eps {
		return "<"
	}
	return "="
}

func generateCases() []testCase {
	rand.Seed(1)
	cases := []testCase{}
	fixed := []struct{ a, b string }{
		{"0", "0"},
		{"1", "0"},
		{"0", "1"},
		{"101", "101"},
		{"1111", "10000"},
		{"0", "11111"},
	}
	for _, f := range fixed {
		inp := fmt.Sprintf("%s\n%s\n", f.a, f.b)
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(20) + 1
		m := rand.Intn(20) + 1
		a := make([]byte, n)
		b := make([]byte, m)
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				a[i] = '0'
			} else {
				a[i] = '1'
			}
		}
		for i := 0; i < m; i++ {
			if rand.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		inp := fmt.Sprintf("%s\n%s\n", string(a), string(b))
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierA.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
