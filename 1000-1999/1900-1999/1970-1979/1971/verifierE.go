package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func timeAt(n int, a, b []int, d int) int {
	if d == 0 {
		return 0
	}
	prevA, prevB := 0, 0
	for i := 0; i < len(a); i++ {
		if d <= a[i] {
			num := (d - prevA) * (b[i] - prevB)
			den := a[i] - prevA
			return prevB + num/den
		}
		prevA, prevB = a[i], b[i]
	}
	return b[len(b)-1]
}

type caseE struct {
	n, k, q int
	a, b    []int
	queries []int
}

func genCase(rng *rand.Rand) caseE {
	n := rng.Intn(100) + 1
	k := rng.Intn(5) + 1
	if k > n {
		k = n
	}
	q := rng.Intn(5) + 1
	a := make([]int, k)
	b := make([]int, k)
	lastA := 0
	lastB := 0
	for i := 0; i < k; i++ {
		lastA += rng.Intn(10) + 1
		if i == k-1 && lastA < n {
			lastA = n
		}
		if lastA > n {
			lastA = n
		}
		a[i] = lastA
		lastB += rng.Intn(10) + 1
		b[i] = lastB
	}
	qs := make([]int, q)
	for i := 0; i < q; i++ {
		qs[i] = rng.Intn(n + 1)
	}
	return caseE{n, k, q, a, b, qs}
}

func runCase(bin string, tc caseE) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("1\n%d %d %d\n", tc.n, tc.k, tc.q))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, v := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	expected := make([]int, len(tc.queries))
	for i, d := range tc.queries {
		expected[i] = timeAt(tc.n, tc.a, tc.b, d)
	}
	for i, exp := range expected {
		if !scanner.Scan() {
			return fmt.Errorf("missing output line %d", i+1)
		}
		var got int
		if _, err := fmt.Sscan(scanner.Text(), &got); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if got != exp {
			return fmt.Errorf("query %d expected %d got %d", i+1, exp, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
