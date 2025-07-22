package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func solveC(n int64, weights []int64) int64 {
	m := int64(len(weights))
	lim := n - 1
	s := int64(math.Sqrt(float64(1 + 8*lim)))
	k := (1 + s) / 2
	for k*(k-1)/2 > lim {
		k--
	}
	if m < k {
		k = m
	}
	sort.Slice(weights, func(i, j int) bool { return weights[i] > weights[j] })
	var total int64
	for i := int64(0); i < k; i++ {
		total += weights[i]
	}
	return total
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(100) + 1)
	m := rng.Intn(6) + 1
	weights := make([]int64, m)
	qs := make([]int, m)
	used := map[int]bool{}
	for i := 0; i < m; i++ {
		w := int64(rng.Intn(100) + 1)
		weights[i] = w
		q := rng.Intn(1000) + 1
		for used[q] {
			q = rng.Intn(1000) + 1
		}
		used[q] = true
		qs[i] = q
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	for i := 0; i < m; i++ {
		input += fmt.Sprintf("%d %d\n", qs[i], weights[i])
	}
	ans := solveC(n, weights)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
