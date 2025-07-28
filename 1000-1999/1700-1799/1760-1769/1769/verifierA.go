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

func expected(a []int) []int {
	n := len(a)
	res := make([]int, n)
	if n > 0 {
		res[0] = a[0] - 1
	}
	for i := 1; i < n; i++ {
		cand := a[i] - (i + 1)
		if cand < res[i-1]+1 {
			cand = res[i-1] + 1
		}
		res[i] = cand
	}
	return res
}

func runCase(bin string, a []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(a)))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteString("\n")
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	want := expected(a)
	if len(fields) != len(want) {
		return fmt.Errorf("expected %d numbers got %d", len(want), len(fields))
	}
	for i, f := range fields {
		var x int
		if _, err := fmt.Sscan(f, &x); err != nil {
			return fmt.Errorf("invalid int output: %v", err)
		}
		if x != want[i] {
			return fmt.Errorf("mismatch at position %d: expected %d got %d", i, want[i], x)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(100) + 1
	a := make([]int, n)
	cur := rng.Intn(5) + 1
	for i := 0; i < n; i++ {
		a[i] = cur
		cur += rng.Intn(10) + 1
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
