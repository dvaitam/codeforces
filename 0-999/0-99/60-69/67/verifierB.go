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

func solveB(n, k int, B []int) string {
	A := make([]int, 0, n)
	for j := n; j >= 1; j-- {
		need := B[j]
		cnt := 0
		idx := 0
		for idx = 0; idx < len(A); idx++ {
			if cnt == need {
				break
			}
			if A[idx] >= j+k {
				cnt++
			}
		}
		A = append(A, 0)
		copy(A[idx+1:], A[idx:])
		A[idx] = j
	}
	var sb strings.Builder
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) (string, string) {
	n := rng.Intn(7) + 1
	k := rng.Intn(n) + 1
	A := make([]int, 0, n)
	Bvals := make([]int, n+1)
	for j := n; j >= 1; j-- {
		pos := rng.Intn(len(A) + 1)
		cnt := 0
		for i := 0; i < pos; i++ {
			if A[i] >= j+k {
				cnt++
			}
		}
		Bvals[j] = cnt
		A = append(A, 0)
		copy(A[pos+1:], A[pos:])
		A[pos] = j
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", Bvals[i])
	}
	sb.WriteByte('\n')
	expected := solveB(n, k, Bvals)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseB(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
