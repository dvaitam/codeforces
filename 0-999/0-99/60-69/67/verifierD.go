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

func solveD(n int, x, y []int) string {
	exitPos := make([]int, n+1)
	for idx, v := range y {
		pos := idx + 1
		if v >= 1 && v <= n {
			exitPos[v] = pos
		}
	}
	tails := make([]int, 0, n)
	for i := 0; i < n; i++ {
		e := exitPos[x[i]]
		b := n + 1 - e
		lo, hi := 0, len(tails)
		for lo < hi {
			mid := (lo + hi) >> 1
			if tails[mid] < b {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo == len(tails) {
			tails = append(tails, b)
		} else {
			tails[lo] = b
		}
	}
	return fmt.Sprintf("%d", len(tails))
}

func generatePermutation(rng *rand.Rand, n int) []int {
	p := rng.Perm(n)
	for i := 0; i < n; i++ {
		p[i]++
	}
	return p
}

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	x := generatePermutation(rng, n)
	y := generatePermutation(rng, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", x[i])
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", y[i])
	}
	sb.WriteByte('\n')
	expected := solveD(n, x, y)
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
