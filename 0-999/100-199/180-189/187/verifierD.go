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

func simulate(segs []int64, g, r int64, start int64) int64 {
	t := start
	cycle := g + r
	for i := 0; i < len(segs)-1; i++ {
		t += segs[i]
		mod := t % cycle
		if mod >= g {
			t += cycle - mod
		}
	}
	t += segs[len(segs)-1]
	return t
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(5) + 1
	g := int64(rng.Intn(5) + 1)
	r := int64(rng.Intn(5) + 1)
	segs := make([]int64, n+1)
	for i := 0; i <= n; i++ {
		segs[i] = int64(rng.Intn(5) + 1)
	}
	q := rng.Intn(3) + 1
	starts := make([]int64, q)
	for i := 0; i < q; i++ {
		starts[i] = int64(rng.Intn(10))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, g, r)
	for i, v := range segs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", q)
	for i := 0; i < q; i++ {
		fmt.Fprintf(&sb, "%d\n", starts[i])
	}
	answers := make([]int64, q)
	for i := 0; i < q; i++ {
		answers[i] = simulate(segs, g, r, starts[i])
	}
	return sb.String(), answers
}

func runCase(bin, input string, exp []int64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) < len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", len(exp), len(fields))
	}
	for i, v := range exp {
		var got int64
		fmt.Sscan(fields[i], &got)
		if got != v {
			return fmt.Errorf("answer %d: expected %d got %d", i+1, v, got)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
