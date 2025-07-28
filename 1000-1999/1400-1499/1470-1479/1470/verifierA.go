package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveCase(n, m int, k []int, c []int) int64 {
	sort.Slice(k, func(i, j int) bool { return k[i] > k[j] })
	var cost int64
	ptr := 0
	for _, idx := range k {
		if ptr < m && ptr <= idx {
			cost += int64(c[ptr])
			ptr++
		} else {
			cost += int64(c[idx])
		}
	}
	return cost
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := n + rng.Intn(5) + 1
	k := make([]int, n)
	for i := 0; i < n; i++ {
		k[i] = rng.Intn(m)
	}
	c := make([]int, m)
	cur := rng.Intn(5) + 1
	c[0] = cur
	for i := 1; i < m; i++ {
		cur += rng.Intn(5) + 1
		c[i] = cur
	}
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", k[i]+1)
	}
	input += "\n"
	for i := 0; i < m; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", c[i])
	}
	input += "\n"
	ans := solveCase(n, m, k, c)
	return input, fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
