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
	total := 0
	for _, v := range a {
		total += v
	}
	match := make([]bool, 101)
	match[0] = true
	copied := 0
	for _, ai := range a {
		for x := 1; x <= ai; x++ {
			p1 := 100 * x / ai
			p2 := 100 * (copied + x) / total
			if p1 == p2 {
				match[p1] = true
			}
		}
		copied += ai
	}
	res := make([]int, 0)
	for i := 0; i <= 100; i++ {
		if match[i] {
			res = append(res, i)
		}
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
			return fmt.Errorf("mismatch at pos %d: expected %d got %d", i, want[i], x)
		}
	}
	return nil
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(20) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(1000) + 1
	}
	return a
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB1.go /path/to/binary")
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
