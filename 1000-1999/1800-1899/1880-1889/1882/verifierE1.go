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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// applyOp applies the pivot operation: given permutation p of length n,
// pivot at 1-indexed position i: result is p[i+1..n], p[i], p[1..i-1]
func applyOp(p []int, i int) []int {
	n := len(p)
	result := make([]int, n)
	idx := 0
	for j := i; j < n; j++ {
		result[idx] = p[j]
		idx++
	}
	result[idx] = p[i-1]
	idx++
	for j := 0; j < i-1; j++ {
		result[idx] = p[j]
		idx++
	}
	return result
}

func isIdentity(p []int) bool {
	for i, v := range p {
		if v != i+1 {
			return false
		}
	}
	return true
}

func validateOutput(n, m int, a, b []int, output string) error {
	output = strings.TrimSpace(output)
	if output == "-1" {
		// Check if -1 is actually valid by trying the reference solver
		// For the easy version, we trust that if candidate says -1, we verify
		// by checking our own solver too
		// Actually, let's just try all (s,r) combos to see if solution exists
		canSolve := false
		for s := 0; s <= 1; s++ {
			for r := 0; r <= 1; r++ {
				ansP := solveSingle(a, n, s)
				ansQ := solveSingle(b, m, r)
				if len(ansP)%2 == len(ansQ)%2 {
					canSolve = true
				}
			}
		}
		if canSolve {
			return fmt.Errorf("candidate said -1 but solution exists")
		}
		return nil
	}

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	k, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("bad k: %v", err)
	}
	if k < 0 || k > 10000 {
		return fmt.Errorf("k=%d out of range [0,10000]", k)
	}
	if len(lines) != k+1 {
		return fmt.Errorf("expected %d operation lines, got %d", k, len(lines)-1)
	}

	p := make([]int, n)
	copy(p, a)
	q := make([]int, m)
	copy(q, b)

	for op := 0; op < k; op++ {
		fields := strings.Fields(lines[op+1])
		if len(fields) != 2 {
			return fmt.Errorf("op %d: expected 2 fields, got %d", op+1, len(fields))
		}
		pi, err1 := strconv.Atoi(fields[0])
		qj, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("op %d: bad numbers", op+1)
		}
		if pi < 1 || pi > n {
			return fmt.Errorf("op %d: i=%d out of range [1,%d]", op+1, pi, n)
		}
		if qj < 1 || qj > m {
			return fmt.Errorf("op %d: j=%d out of range [1,%d]", op+1, qj, m)
		}
		p = applyOp(p, pi)
		q = applyOp(q, qj)
	}

	if !isIdentity(p) {
		return fmt.Errorf("p is not identity after operations")
	}
	if !isIdentity(q) {
		return fmt.Errorf("q is not identity after operations")
	}
	return nil
}

func solveSingle(p []int, n int, s int) []int {
	A := make([]int, n+1)
	A[0] = 0
	for i := 1; i <= n; i++ {
		A[i] = p[i-1]
	}
	T := make([]int, n+1)
	for i := 0; i <= n; i++ {
		T[i] = (i + s) % (n + 1)
	}

	ans := make([]int, 0)
	for {
		match := true
		for i := 0; i <= n; i++ {
			if A[i] != T[i] {
				match = false
				break
			}
		}
		if match {
			break
		}

		z := -1
		for i := 0; i <= n; i++ {
			if A[i] == 0 {
				z = i
				break
			}
		}

		var x int
		if T[z] == 0 {
			for i := 0; i <= n; i++ {
				if A[i] != T[i] {
					x = A[i]
					break
				}
			}
		} else {
			x = T[z]
		}

		y := -1
		for i := 0; i <= n; i++ {
			if A[i] == x {
				y = i
				break
			}
		}

		dist := (y - z) % (n + 1)
		if dist < 0 {
			dist += n + 1
		}
		ans = append(ans, dist)
		A[z], A[y] = A[y], A[z]
	}
	return ans
}

func genCase(rng *rand.Rand) (string, int, int, []int, []int) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	pp := rng.Perm(n)
	qq := rng.Perm(m)
	a := make([]int, n)
	b := make([]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range pp {
		a[i] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	for i, v := range qq {
		b[i] = v + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String(), n, m, a, b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	// Fixed case
	{
		input := "1 1\n1\n1\n"
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fixed case 1 failed: %v\ninput:\n%s", err, input)
			os.Exit(1)
		}
		if err := validateOutput(1, 1, []int{1}, []int{1}, out); err != nil {
			fmt.Fprintf(os.Stderr, "fixed case 1 failed: %v\ninput:\n%s\noutput:\n%s\n", err, input, out)
			os.Exit(1)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, n, m, a, b := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := validateOutput(n, m, a, b, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
