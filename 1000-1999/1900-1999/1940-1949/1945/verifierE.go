package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// generatePerm returns a random permutation of size n (1-indexed values).
func generatePerm(rng *rand.Rand, n int) []int {
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return perm
}

// generateCase creates a single test case and also returns the permutation and target x.
func generateCase(rng *rand.Rand) (string, []int, int) {
	n := rng.Intn(20) + 1
	x := rng.Intn(n) + 1
	p := generatePerm(rng, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, x)
	for i, v := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), p, x
}

// checkOutput validates the candidate output against the problem rules.
func checkOutput(n, x int, p []int, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("missing number of operations")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("invalid number of operations: %v", err)
	}
	if k < 0 || k > 2 {
		return fmt.Errorf("number of operations out of range: %d", k)
	}
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing index i for swap %d", i+1)
		}
		a, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid index i for swap %d: %v", i+1, err)
		}
		if !scanner.Scan() {
			return fmt.Errorf("missing index j for swap %d", i+1)
		}
		b, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("invalid index j for swap %d: %v", i+1, err)
		}
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("swap indices out of range: %d %d", a, b)
		}
		p[a-1], p[b-1] = p[b-1], p[a-1]
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output detected: %s", scanner.Text())
	}
	l, r := 1, n+1
	for r-l > 1 {
		m := (l + r) / 2
		if p[m-1] <= x {
			l = m
		} else {
			r = m
		}
	}
	if p[l-1] != x {
		return fmt.Errorf("after operations binary search ends at value %d instead of %d", p[l-1], x)
	}
	return nil
}

func runCase(exe string, input string, p []int, x int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if err := checkOutput(len(p), x, append([]int(nil), p...), strings.TrimSpace(out.String())); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, p, x := generateCase(rng)
		if err := runCase(exe, input, p, x); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
