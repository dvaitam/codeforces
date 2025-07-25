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

type caseA struct {
	nA, nB int
	k, m   int
	A, B   []int
}

func genCase(rng *rand.Rand) caseA {
	nA := rng.Intn(20) + 1
	nB := rng.Intn(20) + 1
	k := rng.Intn(nA) + 1
	m := rng.Intn(nB) + 1
	A := make([]int, nA)
	cur := rng.Intn(41) - 20
	for i := 0; i < nA; i++ {
		cur += rng.Intn(5)
		A[i] = cur
	}
	B := make([]int, nB)
	cur = rng.Intn(41) - 20
	for i := 0; i < nB; i++ {
		cur += rng.Intn(5)
		B[i] = cur
	}
	return caseA{nA, nB, k, m, A, B}
}

func expected(tc caseA) string {
	if tc.A[tc.k-1] < tc.B[tc.nB-tc.m] {
		return "YES"
	}
	return "NO"
}

func formatInput(tc caseA) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.nA, tc.nB))
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.k, tc.m))
	for i, v := range tc.A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.B {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, tc caseA) error {
	input := formatInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	ans := strings.ToUpper(fields[0])
	if ans != "YES" && ans != "NO" {
		return fmt.Errorf("expected YES or NO, got %s", fields[0])
	}
	if ans != expected(tc) {
		return fmt.Errorf("expected %s got %s", expected(tc), ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, formatInput(tc))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
