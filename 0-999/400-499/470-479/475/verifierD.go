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

func gcd475(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type pair475 struct {
	g int
	c int64
}

func solveD(input string) string {
	idx := 0
	data := []byte(input)
	nextInt := func() int {
		for idx < len(data) && (data[idx] < '0' || data[idx] > '9') {
			idx++
		}
		v := 0
		for idx < len(data) && data[idx] >= '0' && data[idx] <= '9' {
			v = v*10 + int(data[idx]-'0')
			idx++
		}
		return v
	}

	n := nextInt()
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = nextInt()
	}

	total := make(map[int]int64)
	var current []pair475
	for i := 0; i < n; i++ {
		next := make([]pair475, 0, len(current)+1)
		next = append(next, pair475{a[i], 1})
		for _, p := range current {
			next = append(next, pair475{gcd475(p.g, a[i]), p.c})
		}
		current = current[:0]
		for _, p := range next {
			if len(current) > 0 && current[len(current)-1].g == p.g {
				current[len(current)-1].c += p.c
			} else {
				current = append(current, p)
			}
		}
		for _, p := range current {
			total[p.g] += p.c
		}
	}

	q := nextInt()
	var sb strings.Builder
	for i := 0; i < q; i++ {
		x := nextInt()
		fmt.Fprintf(&sb, "%d\n", total[x])
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(100) + 1
	}
	q := rng.Intn(5) + 1
	queries := make([]int, q)
	for i := range queries {
		queries[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d\n", q)
	for _, v := range queries {
		fmt.Fprintf(&sb, "%d\n", v)
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := solveD(input)
		got, err := runProg(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
