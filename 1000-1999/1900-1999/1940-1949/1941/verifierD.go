package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseD struct {
	n, m, x int
	r       []int
	c       []byte
	exp     string
}

func solveD(n, m, x int, r []int, c []byte) string {
	x--
	cur := make([]bool, n)
	cur[x] = true
	for i := 0; i < m; i++ {
		rr := r[i]
		dir := c[i]
		next := make([]bool, n)
		for pos := 0; pos < n; pos++ {
			if !cur[pos] {
				continue
			}
			if dir == '0' || dir == '?' {
				to := (pos + rr) % n
				next[to] = true
			}
			if dir == '1' || dir == '?' {
				to := pos - rr
				to %= n
				if to < 0 {
					to += n
				}
				next[to] = true
			}
		}
		cur = next
	}
	var res []int
	for i := 0; i < n; i++ {
		if cur[i] {
			res = append(res, i+1)
		}
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, len(res))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseD {
	rng := rand.New(rand.NewSource(4))
	cases := make([]testCaseD, 100)
	for i := range cases {
		n := rng.Intn(6) + 2
		m := rng.Intn(6) + 1
		x := rng.Intn(n) + 1
		r := make([]int, m)
		c := make([]byte, m)
		for j := 0; j < m; j++ {
			r[j] = rng.Intn(n-1) + 1
			d := rng.Intn(3)
			if d == 0 {
				c[j] = '0'
			} else if d == 1 {
				c[j] = '1'
			} else {
				c[j] = '?'
			}
		}
		cases[i] = testCaseD{n: n, m: m, x: x, r: r, c: c, exp: solveD(n, m, x, r, c)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.x)
		for j := 0; j < tc.m; j++ {
			fmt.Fprintf(&sb, "%d %c\n", tc.r[j], tc.c[j])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
