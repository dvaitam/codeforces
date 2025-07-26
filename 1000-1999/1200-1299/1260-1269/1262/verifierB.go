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

func construct(q []int) ([]int, bool) {
	n := len(q)
	used := make([]bool, n+1)
	res := make([]int, n)
	next := 1
	for i := 0; i < n; i++ {
		if i == 0 || q[i] > q[i-1] {
			if q[i] < 1 || q[i] > n || used[q[i]] {
				return nil, false
			}
			res[i] = q[i]
			used[q[i]] = true
		} else {
			for next <= n && (used[next] || next >= q[i]) {
				next++
			}
			if next >= q[i] || next > n {
				return nil, false
			}
			res[i] = next
			used[next] = true
		}
		for next <= n && used[next] {
			next++
		}
	}
	return res, true
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string, n int) ([]int, bool, error) {
	fields := strings.Fields(out)
	if len(fields) == 1 && fields[0] == "-1" {
		return nil, false, nil
	}
	if len(fields) != n {
		return nil, false, fmt.Errorf("expected %d numbers, got %d", n, len(fields))
	}
	res := make([]int, n)
	used := make(map[int]bool)
	for i, f := range fields {
		var v int
		if _, err := fmt.Sscan(f, &v); err != nil {
			return nil, false, fmt.Errorf("bad number %q", f)
		}
		if v < 1 || v > n || used[v] {
			return nil, false, fmt.Errorf("invalid permutation value %d", v)
		}
		used[v] = true
		res[i] = v
	}
	return res, true, nil
}

func runCase(bin string, q []int) error {
	n := len(q)
	input := fmt.Sprintf("1\n%d\n", n)
	for i, v := range q {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	_, ok := construct(q)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	got, valid, err := parseOutput(out, n)
	if err != nil {
		return err
	}
	if !ok {
		if valid {
			return fmt.Errorf("expected -1 but got permutation")
		}
		return nil
	}
	if !valid {
		return fmt.Errorf("expected permutation but got -1")
	}
	// verify got matches q
	pm := 0
	for i := 0; i < n; i++ {
		if got[i] > pm {
			pm = got[i]
		}
		if pm != q[i] {
			return fmt.Errorf("prefix max mismatch at %d", i)
		}
	}
	return nil
}

func genValidCase(rng *rand.Rand, n int) []int {
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	q := make([]int, n)
	mx := 0
	for i := 0; i < n; i++ {
		if p[i] > mx {
			mx = p[i]
		}
		q[i] = mx
	}
	return q
}

func genInvalidCase(rng *rand.Rand, n int) []int {
	q := make([]int, n)
	mx := 0
	for i := 0; i < n; i++ {
		inc := rng.Intn(2)
		if inc == 1 {
			mx += rng.Intn(2)
		}
		q[i] = mx + 1
	}
	return q
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// some edge cases from problem
	cases := [][]int{
		{1, 3, 4, 5, 5},
		{1, 1, 3, 4},
		{2, 2},
		{1},
	}
	for idx, q := range cases {
		if err := runCase(bin, q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	for i := len(cases); i < 100; i++ {
		n := rng.Intn(8) + 1
		var q []int
		if rng.Intn(2) == 0 {
			q = genValidCase(rng, n)
		} else {
			q = genInvalidCase(rng, n)
		}
		if err := runCase(bin, q); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
