package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func expected(n int) (int64, []int) {
	p := make([]int, n+1)
	for i := range p {
		p[i] = -1
	}
	for i := n; i >= 0; i-- {
		if p[i] != -1 {
			continue
		}
		if i == 0 {
			p[0] = 0
			break
		}
		k := bits.Len(uint(i)) - 1
		b := (1 << (k + 1)) - 1
		j := b - i
		if j >= 0 && j <= n && p[j] == -1 {
			p[i] = j
			p[j] = i
		} else {
			p[i] = i
		}
	}
	var sum int64
	for i := 0; i <= n; i++ {
		sum += int64(i ^ p[i])
	}
	return sum, p
}

func parseOutput(out string, n int) (int64, []int, error) {
	lines := strings.Fields(out)
	if len(lines) < n+2 {
		return 0, nil, fmt.Errorf("expected at least %d numbers", n+2)
	}
	var m int64
	if _, err := fmt.Sscan(lines[0], &m); err != nil {
		return 0, nil, fmt.Errorf("cannot parse m: %v", err)
	}
	perm := make([]int, n+1)
	for i := 0; i <= n; i++ {
		if _, err := fmt.Sscan(lines[i+1], &perm[i]); err != nil {
			return 0, nil, fmt.Errorf("cannot parse permutation: %v", err)
		}
	}
	return m, perm, nil
}

func verifyPerm(p []int) bool {
	seen := make([]bool, len(p))
	for _, v := range p {
		if v < 0 || v >= len(p) || seen[v] {
			return false
		}
		seen[v] = true
	}
	for _, ok := range seen {
		if !ok {
			return false
		}
	}
	return true
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotM, gotPerm, err := parseOutput(out.String(), n)
	if err != nil {
		return err
	}
	expM, expPerm := expected(n)
	if gotM != expM {
		return fmt.Errorf("expected beauty %d got %d", expM, gotM)
	}
	if len(gotPerm) != len(expPerm) || !verifyPerm(gotPerm) {
		return fmt.Errorf("invalid permutation")
	}
	var sum int64
	for i := 0; i <= n; i++ {
		sum += int64(i ^ gotPerm[i])
	}
	if sum != expM {
		return fmt.Errorf("permutation produces beauty %d want %d", sum, expM)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := []int{0, 1, 2, 5, 10}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 100 {
		cases = append(cases, rng.Intn(50))
	}
	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d (n=%d) failed: %v\n", i+1, n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
